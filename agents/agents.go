package agents

import (
    "strings"
    "time"
    "math"
    "fmt"

    "github.com/op/go-logging"

    "github.com/AstromechZA/spoon/conf"
)

var log = logging.MustGetLogger("spoon.agents")

// The Agent type is an object that can gather and return a result.
type Agent interface {

    // return the config object for this agent
    GetConfig() conf.SpoonConfigAgent

    // fetch a result for this agent
    Tick() (float64, error)
}

// BuildAgent will return a pointer to a constructed object that follows
// the Agent interface.
// This method will return an error if there is no constructor for the
// agent type or if an error occurs while constructing the object.
func BuildAgent(agentConfig *conf.SpoonConfigAgent) (interface{}, error) {
    switch strings.ToLower(agentConfig.Type) {
    case "time":
        return NewTimeAgent(agentConfig)
    default:
        return nil, fmt.Errorf("Unrecognised agent type '%v'", agentConfig.Type)
    }
}

// SpawnAgent will begin running the given agent in a loop based on the
// interval for that agent.
func SpawnAgent(agent interface{}) error {
    agentO, ok := agent.(Agent)
    if ok != true {
        return fmt.Errorf("Failed to cast interface %v to Agent", agent)
    }
    go func(agent Agent) {
        conf := agent.GetConfig()
        log.Infof("Starting %s agent %s with interval %.2f seconds", conf.Type, conf.Path, conf.Interval)

        intervalNanos := float64(conf.Interval) * float64(time.Second)
        spawnTime := time.Now().UnixNano()
        var tickStart time.Time
        var tickElapsed time.Duration
        for {
            // do call to agent
            tickStart = time.Now()
            value, err := agent.Tick()
            tickElapsed = time.Since(tickStart)
            if err != nil {
                log.Errorf("tick for agent %v returned an error after %v: %v", conf.Path, tickElapsed.String(), err.Error())
            } else {
                log.Debugf("tick for agent %v returned a value after %v: %v", conf.Path, tickElapsed.String(), value)
            }

            // now calculate time to sleep to meet the next tick time
            sleep := intervalNanos * (1.0 - math.Mod(float64(time.Now().UnixNano() - spawnTime) / intervalNanos, float64(1)))
            time.Sleep(time.Duration(sleep))
        }
    }(agentO)
    return nil
}
