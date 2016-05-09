package main

import (
    "flag"
    "os"
    "os/signal"

    "github.com/op/go-logging"

    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/agents"

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
    configFlag := flag.String("config", "/etc/spoon.json", "Path to a Spoon config file.")
    // set a more verbose usage message.
    flag.Usage = func() {
        os.Stderr.WriteString(usageString)
        flag.PrintDefaults()
    }
    // parse them
    flag.Parse()

    // set up initial logging to stdout
    slogging.Initial()
    // load the config file
    log.Infof("Loading config from %s", *configFlag)
    cfg, err := conf.Load(configFlag)
    if err != nil {
        log.Criticalf("Failed to load config: %s", err.Error())
        os.Exit(1)
    }

    // now that we have the config, we can reconfigure the logger according to
    // the config
    slogging.Reconfigure(&cfg.Logging)

    // now spawn each of the agents
    for i, c := range cfg.Agents {
        log.Debugf("Building agent %v: %v", i, c)
        agent, err := agents.BuildAgent(&c)
        if err != nil {
            log.Errorf("Failed to build agent %v: %v", &c, err.Error())
            continue
        }
        err = agents.SpawnAgent(agent)
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
