package agents

import (
    "fmt"
    "strings"
    "regexp"

    "github.com/shirou/gopsutil/disk"

    "github.com/AstromechZA/spoon/conf"
    sink_ "github.com/AstromechZA/spoon/sink"
)

type diskAgent struct {
    config conf.SpoonConfigAgent
    settings map[string]string
}

func NewDiskAgent(config *conf.SpoonConfigAgent) (interface{}, error) {

    settings := make(map[string]string, 0)
    for k, v := range config.Settings {
        vs, ok := v.(string)
        if ok == false { return nil, fmt.Errorf("Error casting settings value %v to string", v) }
        settings[k] = vs
    }

    return &diskAgent{
        config: (*config),
        settings: settings,
    }, nil
}

func (a *diskAgent) GetConfig() conf.SpoonConfigAgent {
    return a.config
}

func (a *diskAgent) Tick(sink sink_.Sink) error {

    batch := sink_.NewBatch(sink, 10)
    defer batch.Flush()

    devre := a.settings["device_regex"]

    // fetch all the physical disk partitions. the boolean indicates whether
    // non-physical partitions should be returned too.
    parts, err := disk.Partitions(false)
    if err == nil {
        // loop through all the partitions returned
        for _, p := range parts {

            // check against regex if provided
            if devre != "" {
                m, _ := regexp.MatchString(devre, p.Device)
                if m == false {
                    log.Debugf("Skipping usage for %v because it didn't match device_regex", p.Device)
                    continue
                }
            }

            usage, err := disk.Usage(p.Mountpoint)
            if err == nil {
                prefixPath := fmt.Sprintf("%s.%s", a.config.Path, a.formatDeviceName(p.Device))

                err = batch.Put(fmt.Sprintf("%s.total", prefixPath), float64(usage.Total))
                if err != nil { return err }

                err = batch.Put(fmt.Sprintf("%s.free", prefixPath), float64(usage.Free))
                if err != nil { return err }

                err = batch.Put(fmt.Sprintf("%s.used", prefixPath), float64(usage.Used))
                if err != nil { return err }

                err = batch.Put(fmt.Sprintf("%s.used_percent", prefixPath), float64(usage.UsedPercent))
                if err != nil { return err }

                err = batch.Put(fmt.Sprintf("%s.inode_free", prefixPath), float64(usage.InodesFree))
                if err != nil { return err }

                err = batch.Put(fmt.Sprintf("%s.inode_used", prefixPath), float64(usage.InodesUsed))
                if err != nil { return err }

                err = batch.Put(fmt.Sprintf("%s.inode_used_percent", prefixPath), float64(usage.InodesUsedPercent))
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

        for path, iostat := range iocounters {
            deviceName := "/dev/" + path

            // check against regex if provided
            if devre != "" {
                m, _ := regexp.MatchString(devre, deviceName)
                if m == false {
                    log.Debugf("Skipping iocounters for %v because it didn't match device_regex", deviceName)
                    continue
                }
            }

            prefixPath := fmt.Sprintf("%s.%s", a.config.Path, a.formatDeviceName(deviceName))

            err = batch.Put(fmt.Sprintf("%s.read_count", prefixPath), float64(iostat.ReadCount))
            if err != nil { return err }

            err = batch.Put(fmt.Sprintf("%s.write_count", prefixPath), float64(iostat.WriteCount))
            if err != nil { return err }

            err = batch.Put(fmt.Sprintf("%s.read_bytes", prefixPath), float64(iostat.ReadBytes))
            if err != nil { return err }

            err = batch.Put(fmt.Sprintf("%s.write_bytes", prefixPath), float64(iostat.WriteBytes))
            if err != nil { return err }
        }

    } else {
        log.Errorf("Fetching iocounters for system failed: %v", err.Error())
    }
    return nil
}

func (a *diskAgent) formatDeviceName(device string) string {
    // first replace all forward slashes with -
    device = strings.Replace(device, "/", "_", -1)
    // then trim them off
    return strings.Trim(device, "_")
}
