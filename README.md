# Spoon Readme

Spoon is my implementation of a stats daemon
in Go which sends its stats to a Carbon-cache or Carbon-aggregator. It aims to
closely follow how python-diamond works, but quite a bit simpler and easier to
configure.

Because it's written in Go, it will be deployable as a single binary and all of
the agents should support Linux and OSX well. It will probably be able to
compile and run on Windows too, but results may vary.

My aim in writing this software is to learn a bit more about Golang, and also
deploy the daemon to my home servers.

## Usage

```
$ ./spoon --help
Spoon is a simple metric gatherer for Linux systems. Like the popular Diamond
daemon, it runs a configurable number of gathering agents and forwards the
results to a Carbon Aggregator or Cache.

By default, it looks for a config file at /etc/spoon.json but this path can be
specified at the command line using '-config'.

Spoon does not require root permissions to run, but might need them depending on
which agents are configured.

  -config string
        Path to a Spoon config file.
  -debug
        Print debug logging.
  -generate
        Generate a new example config and print it to stdout.
  -validate
        Validate the config passed in via '-config'.
  -version
        Print the version string.
```

The example configuration produced by `-generate` will have all agents enabled
with some sane defaults and will send metrics to the log while logging to stdout.
It should pass the validation provided by `-validate`.

## Agent Types

- `time`: just returns the unix nanoseconds
- `meta`: returns the cpu percent and RSS usage of the Spoon process.
- `cpu`: returns cpu percentage per core
- `disk`: returns disk usage and io counters if available per physical partition and disk
- `mem`: returns system memory and swap usage
- `uptime`: just returns the machines uptime value
- `net`: returns sent/recv info for interfaces
- `cmd`: allows user to run their own command

More detail available on the [agents documentation](doc/agents.md).

## Agents to add

- docker cpu and mem stats? (just cos gopsutil library supports this)

## Building official binaries

Run the `make_official.sh` script to build binaries for:

- Darwin (OSX) x86_64
- Linux x86_64

All binaries are statically compiled and do not have any dependencies on the
host.

## Why 'Spoon'?

It 'feeds' metrics to Graphite?
