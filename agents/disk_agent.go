package agents

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/disk"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/sink"
)

type diskAgent struct {
	diskAgentSettings
	config   conf.SpoonConfigAgent
	settings map[string]string
}

type diskAgentSettings struct {
	DeviceRegex string `json:"device_regex"`
}

func NewDiskAgent(config *conf.SpoonConfigAgent) (Agent, error) {
	s := diskAgentSettings{}
	if err := json.Unmarshal(config.SettingsRaw, &s); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %s", err)
	}
	return &diskAgent{
		diskAgentSettings: s,
		config:            (*config),
	}, nil
}

func (a *diskAgent) GetConfig() conf.SpoonConfigAgent {
	return a.config
}

func (a *diskAgent) Tick(s sink.Sink) error {

	// fetch all the physical disk partitions. the boolean indicates whether
	// non-physical partitions should be returned too.
	parts, err := disk.Partitions(true)
	if err == nil {
		// loop through all the partitions returned
		for _, p := range parts {

			// check against regex if provided
			if m, _ := regexp.MatchString(a.DeviceRegex, p.Device); m == false {
				continue
			}

			usage, uerr := disk.Usage(p.Mountpoint)
			if uerr == nil {
				log.Printf("Outputting Usage for %v because it matched device_regex", p.Device)
				prefixPath := fmt.Sprintf("%s.%s", a.config.Path, a.formatDeviceName(p.Device))

				s.Gauge(fmt.Sprintf("%s.total_bytes", prefixPath), float64(usage.Total))
				s.Gauge(fmt.Sprintf("%s.free_bytes", prefixPath), float64(usage.Free))
				s.Gauge(fmt.Sprintf("%s.used_bytes", prefixPath), float64(usage.Used))
				s.Gauge(fmt.Sprintf("%s.used_percent", prefixPath), float64(usage.UsedPercent))
				s.Gauge(fmt.Sprintf("%s.inode_free_count", prefixPath), float64(usage.InodesFree))
				s.Gauge(fmt.Sprintf("%s.inode_used_count", prefixPath), float64(usage.InodesUsed))
				s.Gauge(fmt.Sprintf("%s.inode_used_percent", prefixPath), float64(usage.InodesUsedPercent))

			} else {
				log.Printf("Fetching usage for disk %v failed: %v", p.Mountpoint, uerr.Error())
			}
		}
	} else {
		// just log this error, we can continue
		log.Printf("Fetching list of physical disk partitions failed: %v", err.Error())
	}

	iocounters, err := disk.IOCounters()
	if err == nil {

		for path, iostat := range iocounters {
			deviceName := "/dev/" + path

			// check against regex if provided
			if m, _ := regexp.MatchString(a.DeviceRegex, deviceName); m == false {
				continue
			}

			log.Printf("Outputting IO Counters for %v because it matched device_regex", deviceName)
			prefixPath := fmt.Sprintf("%s.%s", a.config.Path, a.formatDeviceName(deviceName))

			s.Gauge(fmt.Sprintf("%s.read_count", prefixPath), float64(iostat.ReadCount))
			s.Gauge(fmt.Sprintf("%s.write_count", prefixPath), float64(iostat.WriteCount))
			s.Gauge(fmt.Sprintf("%s.read_bytes", prefixPath), float64(iostat.ReadBytes))
			s.Gauge(fmt.Sprintf("%s.write_bytes", prefixPath), float64(iostat.WriteBytes))
			s.Gauge(fmt.Sprintf("%s.read_count", prefixPath), float64(iostat.ReadCount))
		}

	} else {
		log.Printf("Fetching iocounters for system failed: %v", err.Error())
	}

	return nil
}

func (a *diskAgent) formatDeviceName(device string) string {
	// first replace all forward slashes with -
	device = strings.Replace(device, "/", "_", -1)
	// then trim them off
	return strings.Trim(device, "_")
}
