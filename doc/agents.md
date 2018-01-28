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

The `settings` for a cmd agent should specify the cmd string array to execute:

```
"settings": {
    "cmd": [
        "python",
        "-c",
        "import random; print '.test.path', random.randint(-100, 100)"
    ]
}
```

For each line in the stdout that contains a relative or absolute metric path
followed by a numeric value, a metric will be generated.

For example, output like the following:

```
this.line.will -123.125125
.relative.path.example 0
This line won't generate a value
```

Could generate the following:

```
this.line.will -> -123.125125
.values.relative.path.example -> 0
.elapsed_seconds -> 0.123
.exit_code -> 0
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

The `disk` agent provides usage metrics about physical disk partitions and IO
counter metrics about physical disks.

For example you might get output like the following:

```
# usage metrics for physical partition /dev/sda1
.dev_sda1.total -> bytes
.dev_sda1.free -> bytes
.dev_sda1.used -> bytes
.dev_sda1.used_percent -> 0-100%
.dev_sda1.inode_free -> count
.dev_sda1.inode_used -> count
.dev_sda1.inode_used_percent -> 0-100%

# usage metrics for physical partition /dev/sdb1
.dev_sdb1.total -> bytes
.dev_sdb1.free -> bytes
.dev_sdb1.used -> bytes
.dev_sdb1.used_percent -> 0-100%
.dev_sdb1.inode_free -> count
.dev_sdb1.inode_used -> count
.dev_sdb1.inode_used_percent -> 0-100%

# io counters for disk /dev/sda
.dev_sda.read_count -> count
.dev_sda.write_count -> count
.dev_sda.read_bytes -> bytes
.dev_sda.write_bytes -> bytes

# io counters for partition /dev/sda1
.dev_sda1.read_count -> count
.dev_sda1.write_count -> count
.dev_sda1.read_bytes -> bytes
.dev_sda1.write_bytes -> bytes

# io counters for partition /dev/sda2
.dev_sda2.read_count -> count
.dev_sda2.write_count -> count
.dev_sda2.read_bytes -> bytes
.dev_sda2.write_bytes -> bytes

# io counters for partition /dev/sda5
.dev_sda5.read_count -> count
.dev_sda5.write_count -> count
.dev_sda5.read_bytes -> bytes
.dev_sda5.write_bytes -> bytes

# io counters for disk /dev/sdb
.dev_sdb.read_count -> count
.dev_sdb.write_count -> count
.dev_sdb.read_bytes -> bytes
.dev_sdb.write_bytes -> bytes

# io counters for partition /dev/sdb1
.dev_sdb1.read_count -> count
.dev_sdb1.write_count -> count
.dev_sdb1.read_bytes -> bytes
.dev_sdb1.write_bytes -> bytes
```

If you don't want to report on all devices and partitions, then add a regular
expression in the `"settings"` parameter as follows:

```
"settings": {
    "device_regex": "sda[12]"
}
```

## `mem` Agent

The `mem` agent reports memory usage information for the whole system. We report
high level metrics for normal ram and swap.

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

## `meta` Agent

The `meta` agent reports the share of cpu that the Spoon process is using as well
as the RSS (resident set size) of memory allocated to it.

```
.cpu_percent -> 0-100%
.rss -> bytes
```

## `net` Agent

The `net` agent reports bytes and packets send and received by each network
interface.

If you don't want to report on all interfaces, then add a regular expression in
the `"settings"` parameter as follows:

```
"settings": {
    "nic_regex": "eth0|eth1"
}
```

This will generate something like:

```
.eth0.bytes_sent -> bytes
.eth0.bytes_recv -> bytes
.eth0.packets_sent -> bytes
.eth0.packets_recv -> bytes
```

## `random` Agent

The `random` agent generated random floating point numbers between the min
and max that you specify. By default between 0 and 100. This agent is very
useful for testing your metric aggregation systems.

```
"settings": {
    "min": 0,
    "max": 100
}
```


```
. -> float
```

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

## `docker` Agent

The `docker` reports resource usage of running containers on the host.

The `container_filters` setting can be used to restrict which containers are
reported on:

```
"settings": {
    "container_filters": {
        "label": "has_stats",
    }
}
```

Example output:

```
.<container name>.cpus.usage.percent -> float
.<container name>.memory.usage.bytes -> integer
.<container name>.memory.usage.percent -> float
.<container name>.networks.eth0.rx.bytes -> integer
.<container name>.networks.eth0.tx.bytes -> integer
```
