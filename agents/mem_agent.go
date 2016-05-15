package agents

import (
    "fmt"

    "github.com/shirou/gopsutil/mem"
    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/sink"
)

type memAgent struct {
    config conf.SpoonConfigAgent
}

func NewMemAgent(config *conf.SpoonConfigAgent) (interface{}, error) {
    return &memAgent{
        config: (*config),
    }, nil
}

func (a *memAgent) GetConfig() conf.SpoonConfigAgent {
    return a.config
}

func (a *memAgent) Tick(sinkBatcher *sink.Batcher) error {
    sinkBatcher.Clear()

    vmemInfo, err := mem.VirtualMemory()
    if err != nil { return err }

    err = sinkBatcher.Put(fmt.Sprintf("%s.total", a.config.Path), float64(vmemInfo.Total))
    if err != nil { return err }

    err = sinkBatcher.Put(fmt.Sprintf("%s.used", a.config.Path), float64(vmemInfo.Used))
    if err != nil { return err }

    err = sinkBatcher.Put(fmt.Sprintf("%s.used_percent", a.config.Path), float64(vmemInfo.UsedPercent))
    if err != nil { return err }

    err = sinkBatcher.Put(fmt.Sprintf("%s.available", a.config.Path), float64(vmemInfo.Available))
    if err != nil { return err }

    smemInfo, err := mem.SwapMemory()
    if err != nil { return err }

    err = sinkBatcher.Put(fmt.Sprintf("%s.swap.total", a.config.Path), float64(smemInfo.Total))
    if err != nil { return err }

    err = sinkBatcher.Put(fmt.Sprintf("%s.swap.used", a.config.Path), float64(smemInfo.Used))
    if err != nil { return err }

    err = sinkBatcher.Put(fmt.Sprintf("%s.swap.used_percent", a.config.Path), float64(smemInfo.UsedPercent))
    if err != nil { return err }

    err = sinkBatcher.Put(fmt.Sprintf("%s.swap.free", a.config.Path), float64(smemInfo.Free))
    if err != nil { return err }

    return sinkBatcher.Flush()
}
