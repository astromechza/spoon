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
            fmt.Printf("Failed to serialise config: %v", err.Error())
            os.Exit(1)
        }
        fmt.Println(string(bytes))
        os.Exit(0)
    }

    // TODO - implement validate

    // set up initial logging to stdout
    slogging.Initial()
    // load the config file
    configPath := (*configFlag)
    if configPath == "" { configPath = "/etc/spoon.json" }
    log.Infof("Loading config from %s", configPath)
    cfg, err := conf.Load(&configPath)
    if err != nil {
        log.Criticalf("Failed to load config: %s", err.Error())
        os.Exit(1)
    }

    // now that we have the config, we can reconfigure the logger according to
    // the config
    slogging.Reconfigure(&cfg.Logging)

    // build sink
    activeSink := sink.NewLoggingSink()

    // now spawn each of the agents
    for i, c := range cfg.Agents {
        if c.Enabled == false {
            log.Infof("Skipping agent %v: %v because it is disabled.", i, c)
            continue
        }
        log.Debugf("Building agent %v: %v", i, c)
        agent, err := agents.BuildAgent(&c)
        if err != nil {
            log.Errorf("Failed to build agent %v: %v", &c, err.Error())
            continue
        }
        err = agents.SpawnAgent(agent, &activeSink)
        if err != nil {
            log.Errorf("Failed to spawn agent %v: %v", &c, err.Error())
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
        },
    }
}
