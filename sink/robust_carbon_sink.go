package sink

import (
    "fmt"
    "time"
    "net"
    "bytes"
    "sync"
    "strconv"
    "errors"

    "github.com/AstromechZA/spoon/conf"
)

/* Robust graphite sink for metrics

when Put is called, it attempts to send a metric, if the metric fails to send
if checks the connection, if the connection is broken, it attempts to reconnect
if the reconnection fails, it returns an error and the metric is not sent
if the reconnection succeeds, it sends the metric.

alot of inspiration taken from github.com/marpaia/graphite-golang

*/

type RobustCarbonSink struct {
    Host string
    Port int
    Protocol string
    ConnTimeout time.Duration

    lock sync.Mutex
    connection net.Conn
}

const defaultConnTimeout = 5

func NewRobustCarbonSink(cfg *conf.SpoonConfigSink) (*RobustCarbonSink, error) {

    carbonHost, found := cfg.Settings["carbon_host"]
    if found == false { return nil, errors.New("Sink settings did not contain 'carbon_host' key") }
    carbonHostString, ok := carbonHost.(string)
    if ok == false {return nil, fmt.Errorf("Error casting %v to string", carbonHost)}

    carbonPort, found := cfg.Settings["carbon_port"]
    if found == false { return nil, errors.New("Sink settings did not contain 'carbon_port' key") }
    carbonPortNum, ok := carbonPort.(float64)
    if ok == false {return nil, fmt.Errorf("Error casting %v to number", carbonPort)}

    return &RobustCarbonSink{
        Host: carbonHostString,
        Port: int(carbonPortNum),
        Protocol: "tcp",
    }, nil
}

func (s *RobustCarbonSink) Reconnect() error {
    if s.connection != nil {
        // try close, ignore errors
        s.connection.Close()
    }

    // these should probably be done only once and stored
    address := fmt.Sprintf("%v:%v", s.Host, s.Port)
    if s.ConnTimeout <= 0 {
        s.ConnTimeout = defaultConnTimeout * time.Second
    }

    // connect with timeout
    log.Infof("Attempting to connect graphite %v socket to %v with timeout %v", s.Protocol, address, s.ConnTimeout)
    conn, err := net.DialTimeout(s.Protocol, address, s.ConnTimeout)
    if err != nil { return err }

    log.Info("Connection successful.")
    s.connection = conn
    return nil
}

func (s *RobustCarbonSink) Disconnect() error {
    log.Infof("Disconnecting graphite socket")
    err := s.connection.Close()
    s.connection = nil
    if err != nil {
        log.Errorf("Error while disconnecting: %v", err.Error())
    }
    return err
}

func (s *RobustCarbonSink) Put(path string, value float64) error {
    // get the timestamp BEFORE we lock
    timestamp := time.Now().Unix()

    // now get lock
    s.lock.Lock()
    defer s.lock.Unlock()

    // construct output buffer
    buf := bytes.NewBufferString("")
    buf.WriteString(fmt.Sprintf(
        "%s %s %d\n",
        path,
        strconv.FormatFloat(value, 'f', -1, 64),
        timestamp,
    ))

    // if no connection, try reconnect
    if s.connection == nil {
        err := s.Reconnect()
        if err != nil { return err }
    }

    // now try send the data
    _, err := s.connection.Write(buf.Bytes())

    // if err is not temporary, disconnect and we can try redo the connection
    if err != nil {
        netError, ok := err.(net.Error)
        if ok {
            if netError.Timeout() {
                log.Errorf("Graphite connection timed out while sending")
                return err
            }
            if netError.Temporary() {
                log.Errorf("Graphite connection hit a temporary error")
                return err
            }

        }

        log.Errorf("Graphite connection hit a more permanent error")
        s.Disconnect()
        return err
    }

    return nil
}
