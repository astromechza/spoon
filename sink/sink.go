package sink

import (
    "sync"

    "github.com/op/go-logging"
)

var log = logging.MustGetLogger("spoon.sink")

// A Sink is an object that acts as the destination for results from the
// agents.
type Sink interface {
    Put(path string, value float64) error
}

// LoggingSink is an implementation of sync which just logs data points to the
// standard logging output.
type LoggingSink struct {
    Lock sync.Mutex
}

// Put writes a path/value pair to the log
func (s *LoggingSink) Put(path string, value float64) error {
    s.Lock.Lock()
    defer s.Lock.Unlock()

    log.Infof("Value for '%v' = %v", path, value)

    return nil
}

func NewLoggingSink() LoggingSink {
    return LoggingSink{Lock: sync.Mutex{}}
}
