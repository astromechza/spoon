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
- ...

## TODO

1. Add more agent types
2. Add carbon packet sending
3. config validation
4. improve agent error handling:
    rather than throwing errors in multivalued agents, just log 'warn' for them.

## Agents to add

- disk_agent (disk usage and IO of a given partition)
- net_agent (network usage across one or more interfaces)
- cmd_agent (get float64 output of a command)
    - cmd_agent is what users will use to add their own metrics
    - they can call out to any executeable binary
    - must be aware that this may be slower than other agents due to calling
        a process
    - casting to a float64 must be done carefully. perhaps we should have a bit
        of regex to identify the first possible number and then we use that
        value. must lookout for negative values too.

## Why 'Spoon'?

It 'feeds' metrics to Graphite?
