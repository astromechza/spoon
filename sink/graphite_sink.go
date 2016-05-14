package sink

import (
    "sync"
    "time"
    "strconv"
    "github.com/marpaia/graphite-golang"
)

/*
This graphite-golang library is nice, but probably needs to be a bit more robust.
Connection should be able to be broken, and recreated later. We should be able
to go through periods where no metrics can be sent. As long as these are logged
appropriately. Otherwise I think a long running process just becomes stuck.

So we need a proper buffered managed graphite connection. So it sends bunches
of metrics at a time and makes sure the connection doesn't die.

If I drop the docker container while Spoon is still running and sending, Put
begins failing with 'broken pipe' and even after bringing up the container again
it continues to fail. So we definitely need to intelligently bring the pipe back.
*/

// A GraphiteSink is a sink that pushes each metric to a carbon
// cache/relay/aggregator port over a tcp connection
type GraphiteSink struct {
    lock sync.Mutex
    graphite graphite.Graphite
}

// NewGraphiteSink creates a new Graphite Sink object
func NewGraphiteSink(graphiteHost string, graphitePort int) (*GraphiteSink, error) {
    g, err := graphite.NewGraphite(graphiteHost, graphitePort)
    if err != nil { return nil, err }

    return &GraphiteSink{
        lock: sync.Mutex{},
        graphite: *g,
    }, nil
}

func (s *GraphiteSink) Put(path string, value float64) error {
    s.lock.Lock()
    defer s.lock.Unlock()

    return s.graphite.SendMetric(graphite.Metric{
        Name: path,
        Value: strconv.FormatFloat(value, 'f', -1, 64),
        Timestamp: time.Now().Unix(),
    })
}
