package sink

import (
    "sync"
    "fmt"
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

// Put writes a path/value pair to the log
func (s *LoggingSink) Put(path string, value float64) error {
    s.lock.Lock()
    defer s.lock.Unlock()

    log.Infof("Value for '%v' = %v", path, value)

    return nil
}

func (s *LoggingSink) PutBatch(batch []Metric) error {
    s.lock.Lock()
    defer s.lock.Unlock()

    output := fmt.Sprintf("LoggingSink received batch of %v metrics:", len(batch))
    for _, m := range batch {
        output += fmt.Sprintf("\nValue for '%v' = %v", m.Path, m.Value)
    }

    log.Info(output)

    return nil
}
