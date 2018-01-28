package agents

import (
	"fmt"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/sink"
	"github.com/shirou/gopsutil/mem"
)

type memAgent struct {
	config conf.SpoonConfigAgent
}

func NewMemAgent(config *conf.SpoonConfigAgent) (Agent, error) {
	return &memAgent{
		config: (*config),
	}, nil
}

func (a *memAgent) GetConfig() conf.SpoonConfigAgent {
	return a.config
}

func (a *memAgent) Tick(s sink.Sink) error {

	vmemInfo, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	s.Gauge(fmt.Sprintf("%s.total_bytes", a.config.Path), float64(vmemInfo.Total))
	s.Gauge(fmt.Sprintf("%s.used_bytes", a.config.Path), float64(vmemInfo.Used))
	s.Gauge(fmt.Sprintf("%s.used_percent", a.config.Path), float64(vmemInfo.UsedPercent))
	s.Gauge(fmt.Sprintf("%s.available_bytes", a.config.Path), float64(vmemInfo.Available))

	smemInfo, err := mem.SwapMemory()
	if err != nil {
		return err
	}

	s.Gauge(fmt.Sprintf("%s.swap.total_bytes", a.config.Path), float64(smemInfo.Total))
	s.Gauge(fmt.Sprintf("%s.swap.used_bytes", a.config.Path), float64(smemInfo.Used))
	s.Gauge(fmt.Sprintf("%s.swap.used_percent", a.config.Path), float64(smemInfo.UsedPercent))
	s.Gauge(fmt.Sprintf("%s.swap.free_bytes", a.config.Path), float64(smemInfo.Free))

	return nil
}
