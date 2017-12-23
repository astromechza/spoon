package sink

import (
	"sync"

	"github.com/AstromechZA/spoon/conf"
)

type StatsdSink struct {
	lock sync.Mutex
}

func NewStatsdSink(cfg *conf.SpoonConfigSink) (*StatsdSink, error) {
	return &StatsdSink{}, nil
}

// Put writes a path/value pair to the log
func (s *StatsdSink) Put(path string, value float64) error {
	return nil
}

func (s *StatsdSink) PutBatch(batch []Metric) error {
	return nil
}
