package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/AstromechZA/spoon/agents"
	"github.com/AstromechZA/spoon/conf"
	"github.com/AstromechZA/spoon/constants"
	"github.com/AstromechZA/spoon/sink"
)

// Load the config information from the file on disk
func Load(path *string) (*conf.SpoonConfig, error) {

	// first read all bytes from file
	data, err := ioutil.ReadFile(*path)
	if err != nil {
		return nil, err
	}

	// now parse config object out
	var cfg conf.SpoonConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	// and return
	return &cfg, nil
}

func GetHostname() (string, error) {
	hn, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return strings.ToLower(hn), err
}

func GetHostnameRev() (string, error) {
	hn, err := GetHostname()
	if err != nil {
		return "", err
	}
	parts := strings.Split(strings.ToLower(hn), ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, "."), nil
}

func GetIfaceIPv4(name string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		if i.Name == name {
			addrs, err := i.Addrs()
			if err != nil {
				return "", err
			}
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if v.IP.To4() != nil {
						return v.IP.To4().String(), nil
					}
				case *net.IPAddr:
					if v.IP.To4() != nil {
						return v.IP.To4().String(), nil
					}
				}
			}
			return "", fmt.Errorf("could not find address for interface '%s'", name)
		}
	}
	return "", fmt.Errorf("unknown iface '%s'", name)
}

func InterpolateBasePath(p string) (output string, err error) {
	r := regexp.MustCompile(`%\(.*?\)`)
	output = string(r.ReplaceAllFunc([]byte(p), func(r []byte) []byte {
		token := string(r)[2 : len(r)-1]
		switch {
		case token == "hostname":
			var hn string
			hn, err = GetHostname()
			if err == nil {
				return []byte(hn)
			}
		case token == "hostname-rev":
			var hn string
			hn, err = GetHostnameRev()
			if err == nil {
				return []byte(hn)
			}
		case strings.HasPrefix(token, "$"):
			return []byte(os.Getenv(string(r[1:])))
		case strings.HasPrefix(token, "iface-ipv4-"):
			var x string
			x, err = GetIfaceIPv4(token[11:])
			if err == nil {
				return []byte(strings.Replace(strings.Replace(x, ":", "_", -1), ".", "_", -1))
			}
		default:
			err = fmt.Errorf("unknown base path interpolation sequence '%s'", token)
		}
		return []byte("?")
	}))
	return
}

func CleanAndValidate(cfg *conf.SpoonConfig) (err error) {
	// check base path
	if cfg.BasePath != "" {

		// interpolate variables into the base path
		cfg.BasePath, err = InterpolateBasePath(cfg.BasePath)
		if err != nil {
			return fmt.Errorf("failed to interpolate base path: %s", err)
		}

		ok, cerr := regexp.MatchString(constants.ValidBasePathRegexStrict, cfg.BasePath)
		if cerr != nil {
			panic(cerr)
		}
		if !ok {
			return fmt.Errorf("Base path %s does not match required format", cfg.BasePath)
		}
	}

	// check Sink config
	_, err = sink.BuildSink(&cfg.Sink)
	if err != nil {
		return fmt.Errorf("failed to build sink: %s", err)
	}

	for _, c := range cfg.Agents {

		// validate agent path
		m, err := regexp.MatchString(constants.ValidAgentPathRegexStrict, c.Path)
		if err != nil {
			panic(err)
		}
		if m == false {
			return fmt.Errorf("%s agent path %s does not match required format", c.Type, c.Path)
		}

		if len(c.Path) > 0 && c.Path[0] == '.' {

			if cfg.BasePath == "" {
				return fmt.Errorf("%s agent path %s is relative, but no base path was specified in config", c.Type, c.Path)
			}

			c.Path = cfg.BasePath + c.Path
		}

		if c.Interval <= 0 {
			return fmt.Errorf("%s agent interval cannot be <= 0", c.Type)
		}

		_, err = agents.BuildAgent(&c)
		if err != nil {
			return err
		}
	}
	return nil
}
