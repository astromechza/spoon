package agents

import (
    "github.com/shirou/gopsutil/host"

    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/sink"
)

type uptimeAgent struct {
    config conf.SpoonConfigAgent
}

func NewUpTimeAgent(config *conf.SpoonConfigAgent) (interface{}, error) {
    return &uptimeAgent{config: (*config)}, nil
}

func (a *uptimeAgent) GetConfig() conf.SpoonConfigAgent {
    return a.config
}

func (a *uptimeAgent) Tick(sink sink.Sink) error {
    ut, err := host.Uptime()
    if err != nil { return err}
    return sink.Put(a.config.Path, float64(ut))
}
