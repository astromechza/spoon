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
path like `".sub.agent"` which will be combined with the base path in the configuration
(if no basepath is available, the configuration will fail validation).

The `interval` parameter is a numeric number of seconds. It can have decimals. It
may not be <= 0.

The `enabled` parameter is pretty straightforward. If true, the agent will collect
and report metrics.

## ?? Agent
