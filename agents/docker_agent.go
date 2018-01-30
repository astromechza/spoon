package agents

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/sink"
)

type dockerAgent struct {
	dockerAgentSettings
	config conf.SpoonConfigAgent
}

type dockerAgentSettings struct {
	ContainerFilters map[string]string
}

func NewDockerAgent(config *conf.SpoonConfigAgent) (Agent, error) {
	s := dockerAgentSettings{}
	if err := json.Unmarshal(config.SettingsRaw, &s); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %s", err)
	}

	agent := &dockerAgent{
		dockerAgentSettings: s,
		config:              (*config),
	}
	return agent, nil
}

func (a *dockerAgent) GetConfig() conf.SpoonConfigAgent {
	return a.config
}

func (a *dockerAgent) Tick(s sink.Sink) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("failed to setup docker client: %s", err)
	}
	cli.NegotiateAPIVersion(ctx)
	defer cli.Close()

	filters := filters.NewArgs()
	filters.Add("status", "running")
	for k, v := range a.ContainerFilters {
		filters.Add(k, v)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters})
	if err != nil {
		return fmt.Errorf("failed to list containers: %s", err)
	}

	wg := sync.WaitGroup{}
	for _, c := range containers {

		// name comes from container name itself
		name := strings.Replace(strings.Trim(c.Names[0], "/"), ".", "_", -1)
		id := c.ID
		uptime := time.Now().Sub(time.Unix(c.Created, 0))
		wg.Add(1)
		go func() {
			a.doStatsForContainer(s, cli, id, name, uptime)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func (a *dockerAgent) doStatsForContainer(s sink.Sink, cli *client.Client, cid, cname string, uptime time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Interval)*time.Second)
	data, err := cli.ContainerStats(ctx, cid, false)
	if err != nil {
		log.Printf("unable to pull stats for container %s: %s", cid, err)
		return
	}
	defer data.Body.Close()
	defer cancel()
	stats := new(types.StatsJSON)
	if err = json.NewDecoder(data.Body).Decode(stats); err != nil {
		log.Printf("failed to parse stats from container %s: %s", cid, err)
		return
	}

	s.Gauge(fmt.Sprintf("%s.%s.uptime_seconds", a.config.Path, cname), uptime.Seconds())
	s.Gauge(fmt.Sprintf("%s.%s.cpus.usage_percent", a.config.Path, cname), calculateCPUPercent(stats))
	s.Gauge(fmt.Sprintf("%s.%s.memory.usage_bytes", a.config.Path, cname), calculateMemoryBytes(stats))
	s.Gauge(fmt.Sprintf("%s.%s.memory.usage_percent", a.config.Path, cname), calculateMemoryUsage(stats))

	for iface, nstats := range stats.Networks {
		s.Gauge(fmt.Sprintf("%s.%s.networks.%s.rx_bytes", a.config.Path, cname, iface), float64(nstats.RxBytes))
		s.Gauge(fmt.Sprintf("%s.%s.networks.%s.tx_bytes", a.config.Path, cname, iface), float64(nstats.TxBytes))
		// TODO: can do a lot more here with packets/dropped/errors
	}
}

func calculateCPUPercent(
	stats *types.StatsJSON,
) float64 {
	usageDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systmDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	if usageDelta > 0 && systmDelta > 0 {
		return (usageDelta / systmDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return 0.0
}

func calculateMemoryBytes(stats *types.StatsJSON) int64 {
	return int64(stats.MemoryStats.Usage)
}

func calculateMemoryUsage(stats *types.StatsJSON) float64 {
	if stats.MemoryStats.Limit != 0 {
		return float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit) * 100.0
	}
	return 0
}
