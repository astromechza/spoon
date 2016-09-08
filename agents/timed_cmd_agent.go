package agents

import (
    "errors"
    "os/exec"
    "strconv"
    "fmt"
    "time"

    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/sink"
)

type timedCmdAgent struct {
    config conf.SpoonConfigAgent
    cmd []string
    ignoreError bool
}

func NewTimedCMDAgent(config *conf.SpoonConfigAgent) (interface{}, error) {

    cmditem, ok := config.Settings["cmd"]
    if ok == false { return nil, errors.New("cmdAgent requires a 'cmd' array in the settings dictionary") }

    cmditems, ok := cmditem.([]interface{})
    if ok == false { return nil, errors.New("cmdAgent requires the 'cmd' setting to be an array") }

    if len(cmditems) < 1 {
        return nil, errors.New("cmdAgent 'cmd' setting must have at least one item")
    }

    cmdStringItems := make([]string, len(cmditems))
    for i, a := range cmditems {
        sa, ok := a.(string)
        if ok == false { return nil, errors.New("cmdAgent 'cmd' setting should only contain strings") }
        cmdStringItems[i] = sa
    }

    ignoreErr := false
    var err error
    v, found := config.Settings["ignore_error"]
    if found == true {
        ignoreErr, err = strconv.ParseBool(fmt.Sprintf("%v", v))
        if err != nil {
            return nil, errors.New("Failed to parse ignore_error as boolean")
        }
    }

    return &timedCmdAgent{
        config: *config,
        cmd: cmdStringItems,
        ignoreError: ignoreErr,
    }, nil
}

func (a *timedCmdAgent) GetConfig() conf.SpoonConfigAgent {
    return a.config
}

func (a *timedCmdAgent) Tick(sinkBatcher *sink.Batcher) error {
    start := time.Now().UnixNano()

    _, err := exec.Command(a.cmd[0], a.cmd[1:]...).Output()
    if a.ignoreError == true && err != nil {
        log.Errorf("%v command failed %v: %s", a.cmd[0], err, err.(*exec.ExitError).Stderr)
        return err
    }

    elapsed := time.Now().UnixNano() - start
    return sinkBatcher.PutAndFlush(a.config.Path, float64(elapsed) / float64(time.Second))
}
