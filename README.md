# Spoon Readme

[![Build Status](https://travis-ci.org/AstromechZA/spoon.svg)](https://travis-ci.org/AstromechZA/spoon)

Spoon is my implementation of a stats daemon in Go which sends its stats to a
Statsd-compatible endpoint. It aims to closely follow how python-diamond
works, but quite a bit simpler and easier to configure.

Because it's written in Go, it will be deployable as a single binary and all of
the agents should support Linux and OSX well. It will probably be able to
compile and run on Windows too, but results may vary.

My aim in writing this software is to learn a bit more about Golang, and also
deploy the daemon to my home servers.

### References

- statsd https://github.com/etsy/statsd
- graphite https://github.com/graphite-project
- grafana https://grafana.com/

## Why should I use this?

- single binary with no dependencies
- simply and straighforward to use with minimal configuration
- pretty reliable
- lightweight
- very easy to deploy

## Why shouldn't I use this?

- hasn't stood the test of time
- it's called Spoon

## Usage

```
$ ./spoon --help
Spoon is a simple metric gatherer for Linux systems. Like the popular Diamond
daemon, it runs a configurable number of gathering agents and forwards the
results to Statsd.

By default, it looks for a config file at /etc/spoon.json but this path can be
specified at the command line using '-config'.

Spoon does not require root permissions to run, but might need them depending on
which agents are configured.

  -config string
    	Path to a Spoon config file.
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

## Installation

Download the binary from the releases page on Github or build it yourself if you're brave.
Add it to `/usr/bin/`, `~/bin` or wherever is appropriate. You can do this all in one step using:

```
$ curl https://raw.githubusercontent.com/AstromechZA/spoon/master/get_spoon.sh | sudo sh
```

Setup your config file in `/etc/spoon.json`:

```
$ spoon -generate > /etc/spoon.json
```

And set up your service definition in `/lib/systemd/system/spoon.service`:

```
[Unit]
Description=spoon metrics daemon
Documentation=https://github.com/AstromechZA/spoon
After=network.target

[Service]
ExecStart=/usr/bin/spoon -config /etc/spoon.json
Restart=always

[Install]
WantedBy=multi-user.target
```

**Note**: if your platform is not supported, please cut a ticket and I'll add it to the build process.

## Configuration

Spoon is configured via a json file passed into it via the `-config` option or
read from `/etc/spoon.json`. This single file configures the agents,
metrics destination, and paths. The example config, [spoon.example.json](spoon.example.json), was
generated via the `-generate` option.

Extra docs:

- [doc/sinks.md](doc/sinks.md) for configuring the metrics sink or destination.
- [doc/agents.md](doc/agents.md) for configuring the active agents.
- [doc/basepath.md](doc/basepath.md) for configuring the base path prefix.

## Running in production

Spoon is not designed to fork itself and make sure it is always running, it
should be controlled by something like `supervisord`, `systemd`, or even just an
rc-local script to launch it at boot time.

Once the agents have been spawned (after config validation) it should not crash
or stop running unless something goes badly wrong (like someone kill -9's it).

It will accept a SIGINT in order to stop gracefully.

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

## How it works

Each agent is spawned as its own goroutine. The goroutine has a while loop
containing a sleep call in order to schedule each metrics call. The sleep time
is adjusted to keep the call rate constant at once every Agent.Interval seconds.

The start time of each agents is randomly delayed by up to half of its interval.
This delay helps to reduce spikey load since most agents will run on a multiple
of 10, 30, or 60 seconds.

Reporting metrics to Carbon are done in batches. One batch of metrics per Agent
call. The connection to Carbon will attempt to reconnect every 10 seconds if the
connection is unsuccessful, and this connection is shared amongst all agents.

## Memory and CPU footprint

These kind of things are subjective depending on the number of agents you have
configured and the frequency of their runs. So far, by the values returned by
the `meta` agent, I've seen cpu usage around 0.03% and memory usage under 8 MB
after running for 24 hours with the example config. So its super lightweight :).

## Why 'Spoon'?

It 'feeds' metrics to Graphite?

## Building it

```
$ dep ensure
$ ./make_official
```
