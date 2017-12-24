package sink

import (
	"log"
	"sync"
)

// LoggingSink is an implementation of sync which just logs data points to the
// standard logging output.
type LoggingSink struct {
	lock sync.Mutex
}

// NewLoggingSink constructs a new Logging metric sink.
// A logging sink does not send metrics anywhere but simply logs them at
// info level to the module logger.
func NewLoggingSink() *LoggingSink {
	return &LoggingSink{lock: sync.Mutex{}}
}

// Gauge writes a path/value pair to the log
func (s *LoggingSink) Gauge(path string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	log.Printf("Value for '%v' = %v", path, value)
}
