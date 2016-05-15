package agents

import (
    "os"
    "runtime"
    "time"

    "github.com/shirou/gopsutil/process"

    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/sink"
)

type metaAgent struct {
    pid int32
    config conf.SpoonConfigAgent
    process process.Process

    // some cpu vars to track cpu change
    numCPU int
    hasPrevCPU bool
    prevCPUTime time.Time
    prevCPUTotal float64
}

func NewMetaAgent(config *conf.SpoonConfigAgent) (interface{}, error) {
    pid := int32(os.Getpid())
    procInfo, err := process.NewProcess(pid)
    if err != nil { return nil, err }

    return &metaAgent{
        pid: pid,
        config: (*config),
        process: *procInfo,
        numCPU: runtime.NumCPU(),
        hasPrevCPU: false,
    }, nil
}

func (a *metaAgent) GetConfig() conf.SpoonConfigAgent {
    return a.config
}

func (a *metaAgent) Tick(sinkBatcher *sink.Batcher) error {
    sinkBatcher.Clear()

    err1 := a.doCPU(sinkBatcher)
    err2 := a.doMem(sinkBatcher)
    if err1 != nil { return err1 }
    if err2 != nil { return err2 }

    return sinkBatcher.Flush()
}

func (a *metaAgent) doCPU(sinkBatcher *sink.Batcher) error {

    procTimes, err := a.process.Times()
    if err != nil { return err }

    now := time.Now()
    total := procTimes.Total()

    if a.hasPrevCPU {
        delta := now.Sub(a.prevCPUTime).Seconds() * float64(a.numCPU)
        percent := a.calculateCPUPercent(a.prevCPUTotal, total, delta, a.numCPU)
        if err = sinkBatcher.Put(a.config.Path + ".cpu_percent", percent); err != nil {
            return err
        }
    }

    a.hasPrevCPU = true
    a.prevCPUTime = now
    a.prevCPUTotal = total

    return nil
}

func (a *metaAgent) calculateCPUPercent(t1, t2 float64, delta float64, numcpu int) float64 {
    if delta == 0 { return 0 }
    return (((t2 - t1) / delta) * 100) * float64(numcpu)
}

func (a *metaAgent) doMem(sinkBatcher *sink.Batcher) error {
    memInfo, err := a.process.MemoryInfo()
    if err != nil { return err }
    return sinkBatcher.Put(a.config.Path + ".rss", float64(memInfo.RSS))
}
