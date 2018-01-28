package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/AstromechZA/spoon/agents"
	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/constants"
	"github.com/AstromechZA/spoon/sink"
)

const usageString = `Spoon is a simple metric gatherer for Linux systems. Like the popular Diamond
daemon, it runs a configurable number of gathering agents and forwards the
results to Statsd.

By default, it looks for a config file at /etc/spoon.json but this path can be
specified at the command line using '-config'.

Spoon does not require root permissions to run, but might need them depending on
which agents are configured.
`

// SpoonVersion is the version string
// format should be 'X.YZ'
// Set this at build time using the -ldflags="-X main.SpoonVersion=X.YZ"
var SpoonVersion = "<unofficial build>"

func mainInner() error {

	// first set up config flag options
	configFlag := flag.String("config", "", "Path to a Spoon config file.")
	generateFlag := flag.Bool("generate", false, "Generate a new example config and print it to stdout.")
	validateFlag := flag.Bool("validate", false, "Validate the config passed in via '-config'.")
	versionFlag := flag.Bool("version", false, "Print the version string.")
	onceFlag := flag.Bool("once", false, "Run each agent once, immediately, and then exit.")

	// set a more verbose usage message.
	flag.Usage = func() {
		fmt.Println(usageString)
		flag.PrintDefaults()
	}
	// parse them
	flag.Parse()

	// first do arg checking
	if *versionFlag {
		fmt.Println("Spoon Version " + SpoonVersion)
		fmt.Println(constants.SpoonImage)
		fmt.Println("Project: https://github.com/AstromechZA/spoon")
		return nil
	}

	if *generateFlag {
		if *validateFlag {
			return fmt.Errorf("Cannot use both -validate and -generate.")
		}
		if *configFlag != "" {
			return fmt.Errorf("Cannot use both -generate and -config.")
		}
	}

	// first handle generate and validate
	if *generateFlag {
		bytes, err := json.MarshalIndent(GenerateExampleConfig(), "", "    ")
		if err != nil {
			return fmt.Errorf("Failed to serialise config: %v", err)
		}
		fmt.Println(string(bytes))
		return nil
	}

	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)

	// load the config file
	configPath := (*configFlag)
	if configPath == "" {
		configPath = "/etc/spoon.json"
	}
	configPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("Failed to identify config path: %s", err)
	}

	// load config
	log.Printf("Loading config from %s", configPath)
	cfg, err := Load(&configPath)
	if err != nil {
		return fmt.Errorf("Failed to load config: %s", err)
	}

	err = CleanAndValidate(cfg)
	if err != nil {
		return fmt.Errorf("Invalid configuration: %s", err)
	}

	if *validateFlag {
		fmt.Printf("No problems found in config from %s. Looks good to me!\n", configPath)
		return nil
	}

	// build sink
	activeSink, err := sink.BuildSink(&cfg.Sink)
	if err != nil {
		return fmt.Errorf("Failed to setup metric sink: %s", err)
	}

	// build the list of real agents
	agentList := make([]agents.Agent, len(cfg.Agents))
	for i, c := range cfg.Agents {

		if len(c.Path) > 0 && c.Path[0] == '.' {
			c.Path = cfg.BasePath + c.Path
		}

		agent, aerr := agents.BuildAgent(&c)
		if aerr != nil {
			os.Exit(1)
		}
		agentList[i] = agent
	}

	// run each agent once
	if *onceFlag {
		hasErrors := false
		group := sync.WaitGroup{}
		for _, a := range agentList {
			group.Add(1)
			go func(current agents.Agent) {
				if aerr := current.Tick(activeSink.(sink.Sink)); aerr != nil {
					log.Printf("Error: %T: %s", current, aerr)
					hasErrors = true
				}
				group.Done()
			}(a)
		}
		group.Wait()
		if hasErrors {
			return fmt.Errorf("some agents had errors")
		}
		return nil
	}

	// now spawn each of the agents
	for _, a := range agentList {
		err = agents.SpawnAgent(a, activeSink.(sink.Sink))
		if err != nil {
			return fmt.Errorf("Failed to spawn agent %v: %s", a, err)
		}
	}

	// instead of sitting in a for loop or something, we wait for sigint
	signalChannel := make(chan os.Signal, 1)
	// notify that we are going to handle interrupts
	signal.Notify(signalChannel, os.Interrupt)
	for sig := range signalChannel {
		log.Printf("Received %v signal. Stopping.", sig)
		break
	}
	return nil
}

func main() {
	if err := mainInner(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func GenerateExampleConfig() *conf.SpoonConfig {
	return &conf.SpoonConfig{
		BasePath: "example.%(hostname)",
		Agents: []conf.SpoonConfigAgent{
			conf.SpoonConfigAgent{
				Type:     "cpu",
				Interval: float32(60),
				Path:     ".cpu",
				Enabled:  true,
			},
			conf.SpoonConfigAgent{
				Type:     "disk",
				Interval: float32(60),
				Path:     ".disk",
				Enabled:  true,
				Settings: map[string]interface{}{
					"device_regex": "da\\d$|disk\\d$",
				},
			},
			conf.SpoonConfigAgent{
				Type:     "mem",
				Interval: float32(60),
				Path:     ".mem",
				Enabled:  true,
			},
			conf.SpoonConfigAgent{
				Type:     "meta",
				Interval: float32(30),
				Path:     ".spoon_process",
				Enabled:  true,
			},
			conf.SpoonConfigAgent{
				Type:     "time",
				Interval: float32(10),
				Path:     ".time",
				Enabled:  true,
			},
			conf.SpoonConfigAgent{
				Type:     "uptime",
				Interval: float32(60),
				Path:     ".uptime",
				Enabled:  true,
			},
			conf.SpoonConfigAgent{
				Type:     "net",
				Interval: float32(60),
				Path:     ".net",
				Settings: map[string]interface{}{
					"nic_regex": "^e(th|n|m)\\d$",
				},
				Enabled: true,
			},
			conf.SpoonConfigAgent{
				Type:     "cmd",
				Interval: float32(30),
				Path:     ".cmd",
				Settings: map[string]interface{}{
					"cmd": []string{
						"python",
						"-c",
						"import random; print '.test.path', random.randint(-100, 100)",
					},
				},
				Enabled: true,
			},
		},
		Sink: conf.SpoonConfigSink{
			Type: "log",
		},
	}
}
