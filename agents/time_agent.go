package agents

import (
    "time"

    "github.com/AstromechZA/spoon/conf"
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

func (self *timeAgent) Tick() (float64, error) {
    return float64(time.Now().UnixNano()), nil
}
