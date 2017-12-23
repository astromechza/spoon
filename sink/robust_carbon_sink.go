package sink

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/AstromechZA/spoon/conf"
)

type RobustCarbonSink struct {
	Host        string
	Port        int
	Protocol    string
	ConnTimeout time.Duration

	lock       sync.Mutex
	connection net.Conn
}

const defaultConnTimeout = 5
const connectionAttemptWait = 10

func NewRobustCarbonSink(cfg *conf.SpoonConfigSink) (*RobustCarbonSink, error) {

	carbonHost, found := cfg.Settings["carbon_host"]
	if found == false {
		return nil, errors.New("Sink settings did not contain 'carbon_host' key")
	}
	carbonHostString, ok := carbonHost.(string)
	if ok == false {
		return nil, fmt.Errorf("Error casting %v to string", carbonHost)
	}

	carbonPort, found := cfg.Settings["carbon_port"]
	if found == false {
		return nil, errors.New("Sink settings did not contain 'carbon_port' key")
	}
	carbonPortNum, ok := carbonPort.(float64)
	if ok == false {
		return nil, fmt.Errorf("Error casting %v to number", carbonPort)
	}

	return &RobustCarbonSink{
		Host:     carbonHostString,
		Port:     int(carbonPortNum),
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
	log.Printf("Attempting to connect graphite %v socket to %v with timeout %v", s.Protocol, address, s.ConnTimeout)
	conn, err := net.DialTimeout(s.Protocol, address, s.ConnTimeout)
	if err != nil {
		// sleep here since we could not connect
		log.Printf("Failed to connect, sleeping %v seconds until next attempt", connectionAttemptWait)
		time.Sleep(time.Duration(connectionAttemptWait) * time.Second)
		return err
	}

	log.Printf("Connection successful.")
	s.connection = conn
	return nil
}

func (s *RobustCarbonSink) Disconnect() error {
	log.Printf("Disconnecting graphite socket")
	err := s.connection.Close()
	s.connection = nil
	if err != nil {
		log.Printf("Error while disconnecting: %v", err.Error())
	}
	return err
}

func (s *RobustCarbonSink) Put(path string, value float64) error {
	// get the timestamp BEFORE we lock
	m := Metric{
		Path:      path,
		Value:     value,
		Timestamp: time.Now().Unix(),
	}
	return s.PutBatch([]Metric{m})
}

func (s *RobustCarbonSink) PutBatch(batch []Metric) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// construct output buffer
	buf := bytes.NewBufferString("")

	for _, m := range batch {
		buf.WriteString(fmt.Sprintf(
			"%s %s %d\n",
			m.Path,
			strconv.FormatFloat(m.Value, 'f', -1, 64),
			m.Timestamp,
		))
	}

	// if no connection, try reconnect
	if s.connection == nil {
		err := s.Reconnect()
		if err != nil {
			return err
		}
	}

	// now try send the data
	_, err := s.connection.Write(buf.Bytes())

	// if err is not temporary, disconnect and we can try redo the connection
	if err != nil {
		netError, ok := err.(net.Error)
		if ok {
			if netError.Timeout() {
				log.Printf("Graphite connection timed out while sending")
				return err
			}
			if netError.Temporary() {
				log.Printf("Graphite connection hit a temporary error")
				return err
			}

		}

		log.Printf("Graphite connection hit a more permanent error")
		s.Disconnect()
		return err
	}

	log.Printf("RobustCarbonSink sent batch of %v metrics", len(batch))
	return nil
}
