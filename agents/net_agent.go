package agents

import (
	"fmt"
	"log"
	"regexp"

	"github.com/shirou/gopsutil/net"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/sink"
)

type netAgent struct {
	config   conf.SpoonConfigAgent
	settings map[string]string
}

func NewNetAgent(config *conf.SpoonConfigAgent) (Agent, error) {

	settings := make(map[string]string, 0)
	for k, v := range config.Settings {
		vs, ok := v.(string)
		if ok == false {
			return nil, fmt.Errorf("Error casting settings value %v to string", v)
		}
		settings[k] = vs
	}

	return &netAgent{config: (*config), settings: settings}, nil
}

func (a *netAgent) GetConfig() conf.SpoonConfigAgent {
	return a.config
}

func (a *netAgent) Tick(s sink.Sink) error {

	iocounters, err := net.IOCounters(true)
	if err != nil {
		return err
	}

	nicre := a.settings["nic_regex"]
	for _, nicio := range iocounters {
		if nicre != "" {
			m, _ := regexp.MatchString(nicre, nicio.Name)
			if m == false {
				continue
			}
		}
		log.Printf("Outputting metrics for %v because it matched nic_regex", nicio.Name)
		prefixPath := fmt.Sprintf("%s.%s", a.config.Path, nicio.Name)

		s.Gauge(fmt.Sprintf("%s.bytes_sent", prefixPath), float64(nicio.BytesSent))
		s.Gauge(fmt.Sprintf("%s.bytes_recv", prefixPath), float64(nicio.BytesRecv))
		s.Gauge(fmt.Sprintf("%s.packets_sent", prefixPath), float64(nicio.PacketsSent))
		s.Gauge(fmt.Sprintf("%s.packets_recv", prefixPath), float64(nicio.PacketsRecv))

		// TODO do we need the error and dropped counts?
	}

	// TODO protocol stats from gopsutil
	// would be useful to track udp/tcp
	// conntrack stats?

	return nil
}
