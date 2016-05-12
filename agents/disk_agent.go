package agents

import (
    "fmt"

    "github.com/shirou/gopsutil/disk"

    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/sink"
)

type diskAgent struct {
    config conf.SpoonConfigAgent
}

func NewDiskAgent(config *conf.SpoonConfigAgent) (interface{}, error) {
    return &diskAgent{
        config: (*config),
    }, nil
}

func (a *diskAgent) GetConfig() conf.SpoonConfigAgent {
    return a.config
}

func (a *diskAgent) Tick(sink sink.Sink) error {

    // fetch all the physical disk partitions. the boolean indicates whether
    // non-physical partitions should be returned too.
    parts, err := disk.Partitions(false)
    if err == nil {
        // loop through all the partitions returned
        for _, p := range parts {
            usage, err := disk.Usage(p.Mountpoint)
            if err == nil {
                prefixPath := fmt.Sprintf("%s.disk.%s", a.config.Path, p.Device)

                err = sink.Put(fmt.Sprintf("%s.total", prefixPath), float64(usage.Total))
                if err != nil { return err }

                err = sink.Put(fmt.Sprintf("%s.free", prefixPath), float64(usage.Free))
                if err != nil { return err }

                err = sink.Put(fmt.Sprintf("%s.used", prefixPath), float64(usage.Used))
                if err != nil { return err }

                err = sink.Put(fmt.Sprintf("%s.used_percent", prefixPath), float64(usage.UsedPercent))
                if err != nil { return err }

                err = sink.Put(fmt.Sprintf("%s.inode_free", prefixPath), float64(usage.InodesFree))
                if err != nil { return err }

                err = sink.Put(fmt.Sprintf("%s.inode_used", prefixPath), float64(usage.InodesUsed))
                if err != nil { return err }

                err = sink.Put(fmt.Sprintf("%s.inode_used_percent", prefixPath), float64(usage.InodesUsedPercent))
                if err != nil { return err }

            } else {
                log.Errorf("Fetching usage for disk %v failed: %v", p.Mountpoint, err.Error())
            }
        }
    } else {
        // just log this error, we can continue
        log.Errorf("Fetching list of physical disk partitions failed: %v", err.Error())
    }

    iocounters, err := disk.IOCounters()
    if err == nil {

        // TODO test this on linux
        for path, iostat := range iocounters {
            prefixPath := fmt.Sprintf("%s.disk.%s", a.config.Path, path)

            err = sink.Put(fmt.Sprintf("%s.read_count", prefixPath), float64(iostat.ReadCount))
            if err != nil { return err }

            err = sink.Put(fmt.Sprintf("%s.write_count", prefixPath), float64(iostat.WriteCount))
            if err != nil { return err }

            err = sink.Put(fmt.Sprintf("%s.read_bytes", prefixPath), float64(iostat.ReadBytes))
            if err != nil { return err }

            err = sink.Put(fmt.Sprintf("%s.write_bytes", prefixPath), float64(iostat.WriteBytes))
            if err != nil { return err }
        }

    } else {
        log.Errorf("Fetching iocounters for system failed: %v", err.Error())
    }
    return nil
}
