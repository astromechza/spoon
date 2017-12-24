# Sinks

A sink, in the Spoon world, is a destination for metrics. Two sinks are available
at the time of writing: 'log', and 'statsd'.

The `"log"` sink will just print the gathered metrics in batches to the standard
logging output.

The `"statsd"` sink will send the metrics to a Statsd endpoint.

To configure a sink, add the top-level `"sink"` section to your config file.

## Log sink

Print metrics to the log:

```
"sink": {
    "type": "log",
}
```

## Statsd sink

Send metrics to a Statsd instance:

```
"sink": {
    "type": "statsd",
    "settings": {
        "address": "statsd.domain.com:8125"
    }
}
```
