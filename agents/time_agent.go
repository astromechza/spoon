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

func (self *timeAgent) GetConfig() conf.SpoonConfigAgent {
    return self.config
}

func (self *timeAgent) Tick(sink sink.Sink) error {
    return sink.Put(self.config.Path, float64(time.Now().UnixNano()))
}
