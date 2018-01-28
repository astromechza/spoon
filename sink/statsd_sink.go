package sink

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AstromechZA/go-statsd"
	"github.com/AstromechZA/spoon/conf"
)

type StatsdSink struct {
	client *statsd.Client
}

type StatsdSinkSettings struct {
	Address string `json:"address"`
}

func NewStatsdSink(cfg *conf.SpoonConfigSink) (*StatsdSink, error) {
	s := &StatsdSinkSettings{}
	if err := json.Unmarshal(cfg.SettingsRaw, s); err != nil {
		return nil, fmt.Errorf("failed to parse statsd settings: %s", err)
	}
	if s.Address == "" {
		return nil, fmt.Errorf("statsd sink settings missing 'address'")
	}

	client, err := statsd.New(
		statsd.Address(s.Address),
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
