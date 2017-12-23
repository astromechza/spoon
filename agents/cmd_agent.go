package agents

import (
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/constants"
	"github.com/AstromechZA/spoon/sink"
)

type cmdAgent struct {
	config     conf.SpoonConfigAgent
	cmd        []string
	lineRegexp regexp.Regexp
}

func NewCMDAgent(config *conf.SpoonConfigAgent) (Agent, error) {

	cmditem, ok := config.Settings["cmd"]
	if ok == false {
		return nil, errors.New("cmdAgent requires a 'cmd' array in the settings dictionary")
	}

	cmditems, ok := cmditem.([]interface{})
	if ok == false {
		return nil, errors.New("cmdAgent requires the 'cmd' setting to be an array")
	}

	if len(cmditems) < 1 {
		return nil, errors.New("cmdAgent 'cmd' setting must have at least one item")
	}

	cmdStringItems := make([]string, len(cmditems))
	for i, a := range cmditems {
		sa, ok := a.(string)
		if ok == false {
			return nil, errors.New("cmdAgent 'cmd' setting should only contain strings")
		}
		cmdStringItems[i] = sa
	}

	return &cmdAgent{
		config:     *config,
		cmd:        cmdStringItems,
		lineRegexp: *regexp.MustCompile("^(" + constants.ValidAgentPathRegex + ")*\\s+([\\-0-9\\.]+)\\s*$"),
	}, nil
}

func (a *cmdAgent) GetConfig() conf.SpoonConfigAgent {
	return a.config
}

func (a *cmdAgent) Tick(sinkBatcher *sink.Batcher) error {
	sinkBatcher.Clear()

	out, err := exec.Command(a.cmd[0], a.cmd[1:]...).Output()
	if err != nil {
		log.Printf("%v command failed %v: %s", a.cmd[0], err, err.(*exec.ExitError).Stderr)
		return err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if a.lineRegexp.MatchString(line) {
			groups := a.lineRegexp.FindStringSubmatch(line)
			subpath := groups[1]
			value, err := strconv.ParseFloat(groups[2], 64)
			if err != nil {
				log.Printf("Path %v had value %v which was not a valid 64bit float", subpath, groups[2])
			}
			if subpath[0] == '.' {
				subpath = a.config.Path + subpath
			}
			err = sinkBatcher.Put(subpath, value)
			if err != nil {
				return err
			}
		}
	}

	return sinkBatcher.Flush()
}
