package agents

import (
    "time"

    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/sink"
)

type timeAgent struct {
    config conf.SpoonConfigAgent
}

func NewTimeAgent(config *conf.SpoonConfigAgent) (interface{}, error) {
    return &timeAgent{config: (*config)}, nil
}

func (a *timeAgent) GetConfig() conf.SpoonConfigAgent {
    return a.config
}

func (a *timeAgent) Tick(sink sink.Sink) error {
    return sink.Put(a.config.Path, float64(time.Now().UnixNano()))
}
