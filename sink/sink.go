package sink

import (
    "fmt"

    "github.com/op/go-logging"

    "github.com/AstromechZA/spoon/conf"
)

var log = logging.MustGetLogger("spoon.sink")

type Metric struct {
    Path string
    Value float64
    Timestamp int64
}

// A Sink is an object that acts as the destination for results from the
// agents.
type Sink interface {
    Put(path string, value float64) error
    PutBatch(batch []Metric) error
}

func BuildSink(cfg *conf.SpoonConfigSink) (interface{}, error) {
    switch cfg.Type {
    case "log":
        return NewLoggingSink(), nil
    case "carbon":
        return NewRobustCarbonSink(cfg)
    default:
        return nil, fmt.Errorf("Unrecognised sink type '%v'", cfg.Type)
    }
}
