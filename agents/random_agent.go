package agents

import (
    "fmt"
    "math/rand"
    "time"
    "strconv"
    "errors"

    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/sink"
)

type randomAgent struct {
    config conf.SpoonConfigAgent
    minValue float64
    maxValue float64
    random *rand.Rand
}

func NewRandomAgent(config *conf.SpoonConfigAgent) (interface{}, error) {
    minVal := float64(0)
    maxVal := float64(100)
    var err error
    v, found := config.Settings["min"]
    if found == true {
        minVal, err = strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
        if err != nil {
            return nil, errors.New("Failed to parse min for randomAgent as float64")
        }
    }

    v, found = config.Settings["max"]
    if found == true {
        maxVal, err = strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
        if err != nil {
            return nil, errors.New("Failed to parse max for randomAgent as float64")
        }
    }

    if maxVal <= minVal {
        return nil, fmt.Errorf("max (%f) for randomAgent must be greater than min (%f)", maxVal, minVal)
    }

    return &randomAgent{
        config: (*config),
        minValue: minVal,
        maxValue: maxVal,
        random: rand.New(rand.NewSource(time.Now().UnixNano())),
    }, nil
}

func (a *randomAgent) GetConfig() conf.SpoonConfigAgent {
    return a.config
}

func (a *randomAgent) Tick(sinkBatcher *sink.Batcher) error {
    rng := a.maxValue - a.minValue
    v := a.random.Float64() * rng + a.minValue
    return sinkBatcher.PutAndFlush(a.config.Path, v)
}
