package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"

	"github.com/AstromechZA/spoon/agents"
	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/constants"
	"github.com/AstromechZA/spoon/sink"
)

const usageString = `Spoon is a simple metric gatherer for Linux systems. Like the popular Diamond
daemon, it runs a configurable number of gathering agents and forwards the
results to a Carbon Aggregator or Cache.

By default, it looks for a config file at /etc/spoon.json but this path can be
specified at the command line using '-config'.

Spoon does not require root permissions to run, but might need them depending on
which agents are configured.

`

// SpoonVersion is the version string
// format should be 'X.YZ'
// Set this at build time using the -ldflags="-X main.SpoonVersion=X.YZ"
var SpoonVersion = "<unofficial build>"

func main() {

	// first set up config flag options
	configFlag := flag.String("config", "", "Path to a Spoon config file.")
	generateFlag := flag.Bool("generate", false, "Generate a new example config and print it to stdout.")
	validateFlag := flag.Bool("validate", false, "Validate the config passed in via '-config'.")
	versionFlag := flag.Bool("version", false, "Print the version string.")

	// set a more verbose usage message.
	flag.Usage = func() {
		os.Stderr.WriteString(usageString)
		flag.PrintDefaults()
	}
	// parse them
	flag.Parse()

	// first do arg checking
	if *versionFlag {
		fmt.Println("Spoon Version " + SpoonVersion)
		fmt.Println(constants.SpoonImage)
		fmt.Println("Project: https://github.com/AstromechZA/spoon")
		os.Exit(0)
	}

	if *generateFlag {
		if *validateFlag {
			os.Stderr.WriteString("Cannot use both -validate and -generate.\n")
			os.Exit(1)
		}
		if (*configFlag) != "" {
			os.Stderr.WriteString("Cannot use both -generate and -config.\n")
			os.Exit(1)
		}
	}

	// first handle generate and validate
	if *generateFlag {
		bytes, err := json.MarshalIndent(GenerateExampleConfig(), "", "    ")
		if err != nil {
			fmt.Printf("Failed to serialise config: %v\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(bytes))
		os.Exit(0)
	}

	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lshortfile)

	// load the config file
	configPath := (*configFlag)
	if configPath == "" {
		configPath = "/etc/spoon.json"
	}
	configPath, err := filepath.Abs(configPath)
	if err != nil {
		fmt.Printf("Failed to identify config path: %v\n", err.Error())
		os.Exit(1)
	}

	// quick validate the config
	err = LoadAndValidateConfig(configPath)
	if err != nil {
		fmt.Printf("Config failed validation: %v\n", err.Error())
		os.Exit(1)
	}

	if *validateFlag {
		fmt.Printf("No problems found in config from %s. Looks good to me!\n", configPath)
		os.Exit(0)
	}

	// load config
	log.Printf("Loading config from %s", configPath)
	cfg, err := conf.Load(&configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err.Error())
		os.Exit(1)
	}

	// build sink
	activeSink, err := sink.BuildSink(&cfg.Sink)
	if err != nil {
		log.Printf("Failed to setup metric sink: %v\n", err.Error())
		os.Exit(1)
	}

	// build the list of real agents
	agentList := make([]agents.Agent, len(cfg.Agents))
	for i, c := range cfg.Agents {

		if len(c.Path) > 0 && c.Path[0] == '.' {
			c.Path = cfg.BasePath + c.Path
		}

		agent, err := agents.BuildAgent(&c)
		if err != nil {
			os.Exit(1)
		}
		agentList[i] = agent
	}

	// now spawn each of the agents
	for _, a := range agentList {
		err = agents.SpawnAgent(a, activeSink.(sink.Sink))
		if err != nil {
			log.Printf("Failed to spawn agent %v: %v", a, err.Error())
			os.Exit(1)
		}
	}

	// instead of sitting in a for loop or something, we wait for sigint
	signalChannel := make(chan os.Signal, 1)
	// notify that we are going to handle interrupts
	signal.Notify(signalChannel, os.Interrupt)
	for sig := range signalChannel {
		log.Printf("Received %v signal. Stopping.", sig)
		os.Exit(0)
	}
}

func LoadAndValidateConfig(path string) error {
	cfg, err := conf.Load(&path)
	if err != nil {
		return err
	}

	// check base path
	if cfg.BasePath != "" {
		m, err := regexp.MatchString(constants.ValidBasePathRegexStrict, cfg.BasePath)
		if err != nil {
			return err
		}
		if m == false {
			return fmt.Errorf("Base path %s does not match required format", cfg.BasePath)
		}
	}

	// check Sink config
	_, err = sink.BuildSink(&cfg.Sink)
	if err != nil {
		return err
	}

	for _, c := range cfg.Agents {

		// validate agent path
		m, err := regexp.MatchString(constants.ValidAgentPathRegexStrict, c.Path)
		if err != nil {
			return err
		}
		if m == false {
			return fmt.Errorf("%s agent path %s does not match required format", c.Type, c.Path)
		}

		if len(c.Path) > 0 && c.Path[0] == '.' {

			if cfg.BasePath == "" {
				return fmt.Errorf("%s agent path %s is relative, but no base path was specified in config", c.Type, c.Path)
			}

			c.Path = cfg.BasePath + c.Path
		}

		if c.Interval <= 0 {
			return fmt.Errorf("%s agent interval cannot be <= 0", c.Type)
		}

		_, err = agents.BuildAgent(&c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateExampleConfig() *conf.SpoonConfig {
	return &conf.SpoonConfig{
		BasePath: "example",
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
