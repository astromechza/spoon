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

    // some cpu vars to track cpu change
    numCPU int
    hasPrevCPU bool
    prevCPUTime time.Time
    prevCPUTotal float64
}

func NewMetaAgent(config *conf.SpoonConfigAgent) (interface{}, error) {
    return &metaAgent{
        pid: int32(os.Getpid()),
        config: (*config),
        numCPU: runtime.NumCPU(),
        hasPrevCPU: false,
    }, nil
}

func (self *metaAgent) GetConfig() conf.SpoonConfigAgent {
    return self.config
}

func (self *metaAgent) Tick(sink sink.Sink) error {
    return self.doCPU(sink)
}

func (self *metaAgent) doCPU(sink sink.Sink) error {
    procInfo, err := process.NewProcess(self.pid)
    if err != nil { return err }
    procTimes, err := procInfo.Times()
    if err != nil { return err }

    now := time.Now()
    total := procTimes.Total()

    if self.hasPrevCPU {
        delta := now.Sub(self.prevCPUTime).Seconds() * float64(self.numCPU)
        percent := calculateCPUPercent(self.prevCPUTotal, total, delta, self.numCPU)
        if err = sink.Put(self.config.Path + ".cpu_percent", percent); err != nil {
            return err
        }
    }

    self.hasPrevCPU = true
    self.prevCPUTime = now
    self.prevCPUTotal = total

    return nil
}

func calculateCPUPercent(t1, t2 float64, delta float64, numcpu int) float64 {
    if delta == 0 { return 0 }
    return (((t2 - t1) / delta) * 100) * float64(numcpu)
}

func (self *metaAgent) doMem(sink sink.Sink) error {
    return nil
}
