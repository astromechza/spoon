package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/constants"
	"github.com/AstromechZA/spoon/sink"
	"golang.org/x/net/context"
)

type cmdAgent struct {
	cmdAgentSettings
	config     conf.SpoonConfigAgent
	lineRegexp regexp.Regexp
}

type cmdAgentSettings struct {
	Command []string `json:"cmd"`
}

func NewCMDAgent(config *conf.SpoonConfigAgent) (Agent, error) {
	s := cmdAgentSettings{}
	if err := json.Unmarshal(config.SettingsRaw, &s); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %s", err)
	}

	if len(s.Command) < 1 {
		return nil, errors.New("cmdAgent 'cmd' setting must have at least one item")
	}

	return &cmdAgent{
		cmdAgentSettings: s,
		config:           *config,
		lineRegexp:       *regexp.MustCompile("^(" + constants.ValidAgentPathRegex + ")*\\s+([\\-0-9\\.]+)\\s*$"),
	}, nil
}

func (a *cmdAgent) GetConfig() conf.SpoonConfigAgent {
	return a.config
}

func (a *cmdAgent) Tick(s sink.Sink) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Interval)*time.Second)
	cmd := exec.CommandContext(ctx, a.Command[0], a.Command[1:]...)
	defer cancel()
	start := time.Now()
	exitcode := 0
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			log.Printf("%v command failed %s: %s", a.Command[0], err, ee.Stderr)
			ws := ee.Sys().(syscall.WaitStatus)
			exitcode = ws.ExitStatus()
		} else {
			log.Printf("%v command failed %s", a.Command[0], err)
		}
	}
	s.Gauge(a.config.Path+".exit_code", exitcode)
	elapsed := time.Now().Sub(start)
	s.Gauge(a.config.Path+".elapsed_seconds", elapsed.Seconds())

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
				subpath = a.config.Path + ".values" + subpath
			}
			s.Gauge(subpath, value)
		}
	}

	return nil
}
