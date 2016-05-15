# Sinks

A sink, in the Spoon world, is a destination for metrics. Two sinks are available
at the time of writing: 'log', and 'carbon'.

The `"log"` sink will just print the gathered metrics in batches to the standard
logging output.

The `"carbon"` sink will send the metrics in batches to a Carbon aggregator/relay/cache
specified by the `"carbon_host"` and `"carbon_port"` settings.

To configure a sink, add the top-level `"sink"` section to your config file.

## Log sink

Print metrics to the log:

```
"sink": {
    "type": "log",
}
```

## Carbon sink

Send metrics to a carbon instance:

```
"sink": {
    "type": "carbon",
    "settings": {
        "carbon_host": "carbon.mydomain.com",
        "carbon_port": 2003
    }
}
```
