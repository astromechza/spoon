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

## Why 'Spoon'?

It 'feeds' metrics to Graphite?
