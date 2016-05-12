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

## Agent Types

- 'time': just returns the unix nanoseconds
- 'meta': returns the cpu percent and RSS usage of the Spoon process.
- 'cpu': returns cpu percentage per core
- 'disk': returns disk usage and io counters if available per physical partition and disk
- 'mem': returns system memory and swap usage
- ...

## TODO

- **Add actual carbon packet sending**
- config validation
- improve agent error handling:
    rather than throwing errors in multivalued agents, just log 'warn' for them.
- add 'enabled' option to configs
- add whitelist option to disk agent (whitelist specific disks to show)
- add hostname to path or have it as an optional interpolated field in the main path name
- add a -validate-config flag
- add a -generate-config /path flag that generates an example config file
- version number?

## Agents to add

- net_agent (network usage across one or more interfaces)
- cmd_agent (get float64 output of a command)
    - cmd_agent is what users will use to add their own metrics
    - they can call out to any executeable binary
    - must be aware that this may be slower than other agents due to calling
        a process
    - casting to a float64 must be done carefully. perhaps we should have a bit
        of regex to identify the first possible number and then we use that
        value. must lookout for negative values too.
- load average?
- docker cpu and mem stats? (just cos gopsutil library supports this)

## Why 'Spoon'?

It 'feeds' metrics to Graphite?
