package agents

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/shirou/gopsutil/net"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/sink"
)

type netAgent struct {
	netAgentSettings
	config conf.SpoonConfigAgent
}

type netAgentSettings struct {
	NicRegex string `json:"nic_regex"`
}

func NewNetAgent(config *conf.SpoonConfigAgent) (Agent, error) {
	s := netAgentSettings{}
	if err := json.Unmarshal(config.SettingsRaw, &s); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %s", err)
	}
	return &netAgent{
		netAgentSettings: s,
		config:           (*config),
	}, nil
}

func (a *netAgent) GetConfig() conf.SpoonConfigAgent {
	return a.config
}

func (a *netAgent) Tick(s sink.Sink) error {

	iocounters, err := net.IOCounters(true)
	if err != nil {
		return err
	}

	for _, nicio := range iocounters {
		if a.NicRegex != "" {
			m, _ := regexp.MatchString(a.NicRegex, nicio.Name)
			if m == false {
				continue
			}
		}
		log.Printf("Outputting metrics for %v because it matched nic_regex", nicio.Name)
		prefixPath := fmt.Sprintf("%s.%s", a.config.Path, nicio.Name)

		s.Gauge(fmt.Sprintf("%s.tx_bytes", prefixPath), float64(nicio.BytesSent))
		s.Gauge(fmt.Sprintf("%s.rx_bytes", prefixPath), float64(nicio.BytesRecv))
		s.Gauge(fmt.Sprintf("%s.tx_packets", prefixPath), float64(nicio.PacketsSent))
		s.Gauge(fmt.Sprintf("%s.rx_packets", prefixPath), float64(nicio.PacketsRecv))

		// TODO do we need the error and dropped counts?
	}

	// TODO protocol stats from gopsutil
	// would be useful to track udp/tcp
	// conntrack stats?

	return nil
}
