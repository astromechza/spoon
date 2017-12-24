package sink

import (
	"fmt"

	"github.com/AstromechZA/spoon/conf"
)

type Metric struct {
	Path      string
	Value     float64
	Timestamp int64
}

// A Sink is an object that acts as the destination for results from the
// agents.
type Sink interface {
	Gauge(bucket string, value interface{})
}

func BuildSink(cfg *conf.SpoonConfigSink) (interface{}, error) {
	switch cfg.Type {
	case "log":
		return NewLoggingSink(), nil
	case "statsd":
		return NewStatsdSink(cfg)
	default:
		return nil, fmt.Errorf("Unrecognised sink type '%v'", cfg.Type)
	}
}
