package sink

import (
    "time"
)

type Batcher struct {
    Metrics []Metric
    MaxSize int

    currentSize int
    sink Sink
}

func NewBatch(sink Sink, maxSize int) *Batcher {
    return &Batcher{
        Metrics: make([]Metric, maxSize + 1),
        MaxSize: maxSize,
        currentSize: 0,
        sink: sink,
    }
}

func (s *Batcher) Put(path string, value float64) error {
    // create metric object
    m := Metric{
        Path: path,
        Value: value,
        Timestamp: time.Now().Unix(),
    }

    // add it
    s.Metrics[s.currentSize] = m
    s.currentSize++
    if s.currentSize >= s.MaxSize {
        // if we have too many metrics, flush
        return s.Flush()
    }
    return nil
}

func (s *Batcher) Flush() error {
    if s.currentSize == 0 { return nil }
    err := s.sink.PutBatch(s.Metrics[:s.currentSize])
    s.currentSize = 0
    s.Metrics = make([]Metric, s.MaxSize + 1)
    return err
}
