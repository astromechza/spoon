package sink

import (
    "fmt"

    "github.com/op/go-logging"

    "github.com/AstromechZA/spoon/conf"
)

var log = logging.MustGetLogger("spoon.sink")

// A Sink is an object that acts as the destination for results from the
// agents.
type Sink interface {
    Put(path string, value float64) error
}

func BuildSink(cfg *conf.SpoonConfigSink) (interface{}, error) {

    switch cfg.Type {
    case "log":
        return NewLoggingSink(), nil
    default:
        return nil, fmt.Errorf("Unrecognised sink type '%v'", cfg.Type)
    }

}
