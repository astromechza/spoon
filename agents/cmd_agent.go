package agents

import (
    "errors"
    "os/exec"
    "strings"
    "regexp"
    "strconv"
    "fmt"

    "github.com/AstromechZA/spoon/conf"
    "github.com/AstromechZA/spoon/sink"
)

type cmdAgent struct {
    config conf.SpoonConfigAgent
    cmd []string
    lineRegexp regexp.Regexp
}

func NewCMDAgent(config *conf.SpoonConfigAgent) (interface{}, error) {

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

    return &cmdAgent{
        config: *config,
        cmd: cmdStringItems,
        lineRegexp: *regexp.MustCompile("^([a-zA-Z0-9\\-\\_]+?(?:\\.[a-zA-Z0-9\\-\\_]+)*)\\s+([\\-0-9\\.]+)\\s*$"),
    }, nil
}

func (a *cmdAgent) GetConfig() conf.SpoonConfigAgent {
    return a.config
}

func (a *cmdAgent) Tick(sink sink.Sink) error {

    out, err := exec.Command(a.cmd[0], a.cmd[1:]...).Output()
    if err != nil {
        log.Errorf("%v command failed %v: %s", a.cmd[0], err, err.(*exec.ExitError).Stderr)
        return err
    }

    var putError error

    lines := strings.Split(string(out), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if a.lineRegexp.MatchString(line) {
            groups := a.lineRegexp.FindStringSubmatch(line)
            subpath := groups[1]
            value, err := strconv.ParseFloat(groups[2], 64)
            if err != nil {
                log.Errorf("Path %v had value %v which was not a valid 64bit float", subpath, groups[2])
            }
            subpath = fmt.Sprintf("%s.%v", a.config.Path, subpath)
            err = sink.Put(subpath, value)
            if err != nil && putError != nil {
                log.Errorf("Error while putting value for %v: %v", subpath, err.Error())
                putError = err
            }
        }
    }

    return putError
}
