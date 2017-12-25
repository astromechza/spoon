# `base_path` configuration

The `base_path` setting in the config provides a global prefix for all metrics emitted by this service.

This should define the namespace for your hosts metrics so its useful to include some identifier for your host.

To make this easier, some special calculated tokens are provided via a `%(token)` substring:

### `%(hostname)`

This uses the Golang `os.Hostname` method to add in the the hostname of the current machine.

For example:

```
"base_path": "example.%(hostname).metrics" -> example.my-host.metrics
```

### `%(hostname-rev)`

To better support environments where hostname returns a dotted hostname like "my-host.local", this token
returns the same data as `%(hostname)` but with parts reversed.

```
"base_path": "example.%(hostname).metrics" -> example.local.my-host.metrics
```

### `%(iface-ipv4-eth0)`

This finds the v4 ip address for the interface named as a suffix (eg: "eth0").

### `%($envvar)`

This simply replaces the string with the environment variable available from `$envvar`. This can be used to implement any extra environmental path information.
