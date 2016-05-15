# Logging

Spoon has a very simple logging system. It can either log to stdout or to a
rotating file.

Logging is configured in the `"logging"` part of the config file.

## Logging to stdout

If the path is `"-"` or `""`, logs will be written to stdout.

```
"logging": {
    "path": "-"
}
```

## Logging to a path

If the path is not empty, it will log to the given file. No logs will be written
to stdout except for the few lines before the config file is read.

```
"logging": {
    "path": "/var/log/spoon.log"
}
```

If you specify a rotate size, the files will be rotated after that many bytes.
There is no maximum number of log rotations.

```
"logging": {
    "path": "/var/log/spoon.log",
    "rotate_size": 1024 * 1024
}
```

The rotation format is `<logpath>.20160515T123713`. The datetime is the local
time when the log file was rotated.
