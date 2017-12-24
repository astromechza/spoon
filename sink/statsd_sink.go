package sink

import (
	"fmt"
	"log"

	"github.com/AstromechZA/go-statsd"
	"github.com/AstromechZA/spoon/conf"
)

type StatsdSink struct {
	client *statsd.Client
}

func NewStatsdSink(cfg *conf.SpoonConfigSink) (*StatsdSink, error) {
	addr, ok := cfg.Settings["address"]
	if !ok {
		return nil, fmt.Errorf("statsd sink settings missing 'address'")
	}

	client, err := statsd.New(
		statsd.Address(addr.(string)),
		statsd.ErrorHandler(func(e error) {
			log.Printf("Statsd sink error: %s", e)
		}),
		statsd.LazyConnect(),
		statsd.FlushesBetweenReconnect(10*60*5),
	)
	if err != nil {
		return nil, err
	}

	return &StatsdSink{client: client}, nil
}

// Gauge writes a path/value pair to the log
func (s *StatsdSink) Gauge(path string, value interface{}) {
	s.client.Gauge(path, value)
}
