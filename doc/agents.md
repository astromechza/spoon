# Agents

Each agent is configured to collect and report a specific set of metrics. Each
agent has a configuration section that defines its type, path, collection
interval, whether it's enabled or not, and various other settings. This document
lists the types of agents available, what metrics they collect, and how to
configure them.

### Common configuration

```
{
    "type": <the agent type string>,
    "path": <relative or absolute '.' separated path>,
    "interval": <the number of seconds between collections>,
    "enabled": <whether this agent is enabled or not>,
    "settings": <settings structure if required>
}
```

The `type` parameter is the name of the agent type such as `"cpu"`, `"mem"`, etc.

The `path` parameter is either an absolute path, such as `"my.cpu.agent"`, or a relative
path like `".sub.agent"` which will be combined with the base path in the main configuration
(if no basepath is available, the configuration will fail validation). All metrics generated
by this agent will have this path as a prefix.

The `interval` parameter is a numeric number of seconds. It can have decimals. It
may not be <= 0.

The `enabled` parameter is pretty straightforward. If true, the agent will collect
and report metrics.

The `settings` parameter is only optionally needed for some agents. For all
others it is ignored and can be left out.

## `cmd` Agent

The `cmd` agent allows the user to write their own agents in the form of other
scripts and binaries that can be called on the machine. This makes up for the
fact that we cant load client collectors like Diamond does.

For each line in the stdout that contains a relative or absolute metric path
followed by a numeric value, a metric will be generated.

For example:

```
this.line.will -123.125125
.relative.path.example 0
This line won't generate a value
```

Will generate

```
this.line.will -> -123.125125
.relative.path.example -> 0
```

The `cmd` agent is the only agent allowed to post metrics with absolute path
names without being influenced by the base path or agent path.

## `cpu` Agent

The `cpu` agent simply returns the cpu usage percentage across all of the
individual cores.

For example on a 4 core machine this will generate:

```
.0.cpu_percentage -> 0-100%
.1.cpu_percentage -> 0-100%
.2.cpu_percentage -> 0-100%
.3.cpu_percentage -> 0-100%
```

## `disk` Agent

The `disk` agent provides usage metrics about partitions and IO counter metrics
about physical disks.

TODO add more here

## `mem` Agent

```
.total -> bytes
.used ->  bytes
.used_percent ->  0-100%
.available ->  bytes

.swap.total -> bytes
.swap.used -> bytes
.swap.used_percent -> 0-100%
.swap.free -> bytes
```

TODO add more here

## `meta` Agent

The `meta` agent reports the share of cpu that the Spoon process is using as well
as the RSS (resident set size) of memory allocated to it.

```
.cpu_percent -> 0-100%
.rss -> bytes
```

## `net` Agent

```
em1.bytes_sent -> bytes
em1.bytes_recv -> bytes
em1.packets_sent -> bytes
em1.packets_recv -> bytes
```

TODO add more here

## `time` Agent

The `time` agent can be used as a testing metric, or when collected between
a large range of machines, can help to reveal NTP issues. It simply returns the
current unix time as number of nanoseconds since January 1, 1970 UTC.

```
. -> nanoseconds
```

## `uptime` Agent

The `uptime` agent reports the uptime of the machine in seconds.

```
. -> seconds
```
