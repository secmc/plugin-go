# Example config
example config for a golang plugin:
```yaml
# Plugin Server Configuration
# Dragonfly runs a gRPC server that plugins connect to
# Use Unix socket for best performance
server_port: "unix:///tmp/dragonfly_plugin.sock"
# Or use TCP for remote: "127.0.0.1:50050"

# List of plugin IDs that must connect before server starts
# This ensures custom items are registered before the resource pack is built
required_plugins:
  - example-go

# Maximum time to wait for required plugins to connect (milliseconds)
hello_timeout_ms: 5000

plugins:
  - id: example-go
    name: Example Go Plugin
    command: "go"
    args: ["run", "cmd/main.go"]
    work_dir: "/home/restart/projects/plugin-go"

```
