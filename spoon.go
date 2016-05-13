package main

import (
    "flag"
    "fmt"
    "encoding/json"
    "os"
    "os/signal"

    "github.com/op/go-logging"

    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/agents"
    "github.com/AstromechZA/spoon/sink"
    slogging "github.com/AstromechZA/spoon/logging"
)

const usageString =
`Spoon is a simple metric gatherer for Linux systems. Like the popular Diamond
daemon, it runs a configurable number of gathering agents and forwards the
results to a Carbon Aggregator or Cache.

By default, it looks for a config file at /etc/spoon.json but this path can be
specified at the command line.

Spoon does not require root permissions to run, but might need them depending on
what agents are configured.

`

var log = logging.MustGetLogger("spoon")

func main() {

    // first set up config flag options
    configFlag := flag.String("config", "", "Path to a Spoon config file.")
    generateFlag := flag.Bool("generate", false, "Generate a new example config and print it to stdout.")
    validateFlag := flag.Bool("validate", false, "Validate the config passed in via '-config'.")

    // set a more verbose usage message.
    flag.Usage = func() {
        os.Stderr.WriteString(usageString)
        flag.PrintDefaults()
    }
    // parse them
    flag.Parse()

    // first do arg checking
    if (*generateFlag) {
        if (*validateFlag) {
            os.Stderr.WriteString("Cannot use both -validate and -generate.\n")
            os.Exit(1)
        }
        if (*configFlag) != "" {
            os.Stderr.WriteString("Cannot use both -generate and -config.\n")
            os.Exit(1)
        }
    }

    // first handle generate and validate
    if (*generateFlag) {
        bytes, err := json.MarshalIndent(GenerateExampleConfig(), "", "    ")
        if err != nil {
            fmt.Printf("Failed to serialise config: %v\n", err.Error())
            os.Exit(1)
        }
        fmt.Println(string(bytes))
        os.Exit(0)
    }

    // load the config file
    configPath := (*configFlag)
    if configPath == "" { configPath = "/etc/spoon.json" }

    // quick validate the config if user asked
    if (*validateFlag) {
        err := LoadAndValidateConfig(configPath)
        if err != nil {
            fmt.Printf("Config failed validation: %v\n", err.Error())
            os.Exit(1)
        }
        fmt.Printf("No problems found in %v config. Looks good to me!\n", configPath)
        os.Exit(0)
    }

    // set up initial logging to stdout
    slogging.Initial()

    // load config
    log.Infof("Loading config from %s", configPath)
    cfg, err := conf.Load(&configPath)
    if err != nil {
        fmt.Printf("Failed to load config: %v\n", err.Error())
        os.Exit(1)
    }

    // now that we have the config, we can reconfigure the logger according to
    // the config
    slogging.Reconfigure(&cfg.Logging)

    // build sink
    activeSink := sink.NewLoggingSink()

    // build the list of real agents
    agentList := make([]interface{}, len(cfg.Agents))
    for i, c := range cfg.Agents {
        agent, err := agents.BuildAgent(&c)
        if err != nil {
            os.Exit(1)
        }
        agentList[i] = agent
    }

    hn, err := os.Hostname()
    if err == nil {
        log.Infof("Hostname: %v", hn)
    }

    // now spawn each of the agents
    for _, a := range agentList {
        err = agents.SpawnAgent(a, &activeSink)
        if err != nil {
            log.Errorf("Failed to spawn agent %v: %v", a, err.Error())
            os.Exit(1)
        }
    }

    // instead of sitting in a for loop or something, we wait for sigint
    signalChannel := make(chan os.Signal, 1)
    // notify that we are going to handle interrupts
    signal.Notify(signalChannel, os.Interrupt)
    for sig := range signalChannel {
        log.Infof("Received %v signal. Stopping.", sig)
        os.Exit(0)
    }
}

func LoadAndValidateConfig(path string) error {
    cfg, err := conf.Load(&path)
    if err != nil { return err }

    // check logging config

    // check Sink config

    for _, c := range cfg.Agents {
        _, err := agents.BuildAgent(&c)
        if err != nil { return err }
    }
    return nil
}

func GenerateExampleConfig() *conf.SpoonConfig {
    return &conf.SpoonConfig{
        Logging: conf.SpoonConfigLog{
            Path: "-",
        },
        Agents: []conf.SpoonConfigAgent{
            conf.SpoonConfigAgent{
                Type: "cpu",
                Interval: float32(60),
                Path: "example.cpu",
            },
            conf.SpoonConfigAgent{
                Type: "disk",
                Interval: float32(60),
                Path: "example.disk",
            },
            conf.SpoonConfigAgent{
                Type: "mem",
                Interval: float32(60),
                Path: "example.mem",
            },
            conf.SpoonConfigAgent{
                Type: "meta",
                Interval: float32(30),
                Path: "example.spoon",
                Enabled: true,
            },
            conf.SpoonConfigAgent{
                Type: "time",
                Interval: float32(10),
                Path: "example.time",
                Enabled: true,
            },
            conf.SpoonConfigAgent{
                Type: "uptime",
                Interval: float32(60),
                Path: "example.uptime",
                Enabled: true,
            },
            conf.SpoonConfigAgent{
                Type: "net",
                Interval: float32(60),
                Path: "example.net",
                Settings: map[string]interface{}{
                    "nic_regex": "e(th|n)\\d",
                },
                Enabled: true,
            },
            conf.SpoonConfigAgent{
                Type: "cmd",
                Interval: float32(30),
                Path: "example.cmd",
                Settings: map[string]interface{}{
                    "cmd": []string{
                        "python",
                        "-c",
                        "import random; print 'test.path', random.randint(-100, 100)",
                    },
                },
                Enabled: true,
            },
        },
    }
}
