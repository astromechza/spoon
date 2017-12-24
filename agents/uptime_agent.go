package agents

import (
	"github.com/shirou/gopsutil/host"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/sink"
)

type uptimeAgent struct {
	config conf.SpoonConfigAgent
}

func NewUpTimeAgent(config *conf.SpoonConfigAgent) (Agent, error) {
	return &uptimeAgent{config: (*config)}, nil
}

func (a *uptimeAgent) GetConfig() conf.SpoonConfigAgent {
	return a.config
}

func (a *uptimeAgent) Tick(s sink.Sink) error {
	ut, err := host.Uptime()
	if err != nil {
		return err
	}
	s.Gauge(a.config.Path, float64(ut))
	return nil
}
