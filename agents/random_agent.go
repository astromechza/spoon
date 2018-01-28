package agents

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/sink"
)

type randomAgent struct {
	randomAgentSettings
	config conf.SpoonConfigAgent
	random *rand.Rand
}

type randomAgentSettings struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

func NewRandomAgent(config *conf.SpoonConfigAgent) (Agent, error) {
	s := randomAgentSettings{
		Min: 0,
		Max: 100,
	}
	if err := json.Unmarshal(config.SettingsRaw, &s); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %s", err)
	}

	if s.Max <= s.Min {
		return nil, fmt.Errorf("max (%f) for randomAgent must be greater than min (%f)", s.Max, s.Min)
	}

	return &randomAgent{
		randomAgentSettings: s,
		config:              (*config),
		random:              rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (a *randomAgent) GetConfig() conf.SpoonConfigAgent {
	return a.config
}

func (a *randomAgent) Tick(s sink.Sink) error {
	rng := a.Max - a.Min
	v := a.random.Float64()*rng + a.Min
	s.Gauge(a.config.Path, v)
	return nil
}
