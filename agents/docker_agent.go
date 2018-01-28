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
	config         conf.SpoonConfigAgent
	containerLabel string
}

func NewDockerAgent(config *conf.SpoonConfigAgent) (Agent, error) {
	agent := &dockerAgent{
		config: (*config),
	}
	if v, ok := config.Settings["container_label"]; ok {
		agent.containerLabel = v.(string)
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

	filters := filters.NewArgs()
	filters.ExactMatch("state", "running")
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters})
	if err != nil {
		return fmt.Errorf("failed to list containers: %s", err)
	}

	wg := sync.WaitGroup{}
	for _, c := range containers {

		// name comes from container name itself
		name := strings.Replace(strings.Trim(c.Names[0], "/"), ".", "_", -1)
		id := c.ID

		// if setting is defined, only run on labeled containers
		if a.containerLabel != "" {
			if _, ok := c.Labels[a.containerLabel]; !ok {
				continue
			}
		}

		wg.Add(1)
		go func() {
			a.doStatsForContainer(s, cli, id, name)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func (a *dockerAgent) doStatsForContainer(s sink.Sink, cli *client.Client, cid, cname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Interval)*time.Second)
	data, err := cli.ContainerStats(ctx, cid, false)
	if err != nil {
		log.Printf("unable to pull stats for container %s: %s", cid, err)
		return
	}
	defer data.Body.Close()
	cancel()
	stats := new(types.StatsJSON)
	if err = json.NewDecoder(data.Body).Decode(stats); err != nil {
		log.Printf("failed to parse stats from container %s: %s", cid, err)
		return
	}

	s.Gauge(fmt.Sprintf("%s.%s.cpus.usage.percent", a.config.Path, cname), calculateCPUPercent(stats))
	s.Gauge(fmt.Sprintf("%s.%s.memory.usage.bytes", a.config.Path, cname), calculateMemoryBytes(stats))
	s.Gauge(fmt.Sprintf("%s.%s.memory.usage.percent", a.config.Path, cname), calculateMemoryUsage(stats))

	for iface, nstats := range stats.Networks {
		s.Gauge(fmt.Sprintf("%s.%s.networks.%s.rx.bytes", a.config.Path, cname, iface), float64(nstats.RxBytes))
		s.Gauge(fmt.Sprintf("%s.%s.networks.%s.tx.bytes", a.config.Path, cname, iface), float64(nstats.TxBytes))
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
