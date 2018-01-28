package agents

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/sink"
)

// The Agent type is an object that can gather and return a result.
type Agent interface {
	// return the config object for this agent
	GetConfig() conf.SpoonConfigAgent

	// fetch a result for this agent
	Tick(sink.Sink) error
}

// BuildAgent will return a pointer to a constructed object that follows
// the Agent interface.
// This method will return an error if there is no constructor for the
// agent type or if an error occurs while constructing the object.
func BuildAgent(agentConfig *conf.SpoonConfigAgent) (Agent, error) {
	switch strings.ToLower(agentConfig.Type) {
	case "disk":
		return NewDiskAgent(agentConfig)
	case "cpu":
		return NewCPUAgent(agentConfig)
	case "mem":
		return NewMemAgent(agentConfig)
	case "time":
		return NewTimeAgent(agentConfig)
	case "uptime":
		return NewUpTimeAgent(agentConfig)
	case "meta":
		return NewMetaAgent(agentConfig)
	case "net":
		return NewNetAgent(agentConfig)
	case "cmd":
		return NewCMDAgent(agentConfig)
	case "random":
		return NewRandomAgent(agentConfig)
	case "timed":
		return NewTimedCMDAgent(agentConfig)
	case "docker":
		return NewDockerAgent(agentConfig)
	default:
		return nil, fmt.Errorf("Unrecognised agent type '%v'", agentConfig.Type)
	}
}

// SpawnAgent will begin running the given agent in a loop based on the
// interval for that agent.
func SpawnAgent(agent Agent, s sink.Sink) error {

	if agent.GetConfig().Enabled == false {
		log.Printf("Skipping agent %v because it is disabled.", agent.GetConfig())
		return nil
	}

	go func(agent Agent) {
		conf := agent.GetConfig()
		log.Printf("Starting %s agent %s with interval %.2f seconds", conf.Type, conf.Path, conf.Interval)

		// calculate random delay using half of the interval
		delay := rand.Float64() * float64(conf.Interval) * 0.5
		log.Printf("Delaying %s agent by %.2f seconds in order to spread the agents out and reduce spikes", conf.Type, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		intervalNanos := float64(conf.Interval) * float64(time.Second)
		spawnTime := time.Now().UnixNano()
		var tickStart time.Time
		var tickElapsed time.Duration
		for {
			// do call to agent
			tickStart = time.Now()
			err := agent.Tick(s)
			tickElapsed = time.Since(tickStart)
			if err != nil {
				log.Printf("Agent %v@%v returned an error after %v: %v", conf.Type, conf.Path, tickElapsed.String(), err.Error())
			} else {
				log.Printf("Tick for agent %v@%v returned after %v", conf.Type, conf.Path, tickElapsed.String())
			}

			// now calculate time to sleep to meet the next tick time
			sleep := intervalNanos * (1.0 - math.Mod(float64(time.Now().UnixNano()-spawnTime)/intervalNanos, float64(1)))
			time.Sleep(time.Duration(sleep))
		}
	}(agent)
	return nil
}
