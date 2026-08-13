# Cmdry

Cmdry is a small, self-hosted web workbench for Linux command-line tools. The core provides navigation, a consistent server-rendered UI, plugin discovery, and strict executable-plugin boundaries. Plugins are separate processes that return versioned JSON—not arbitrary HTML, JavaScript, or shell commands.

The first included plugin is **Port Inspector**, a read-only view of listening TCP and UDP ports from `ss`.

> Screenshot placeholder: Cmdry’s Port Inspector page shows concise listening-port counts above a sortable-ready data table.

## What Cmdry is (and is not)

Cmdry is an extensible way to expose a focused Linux CLI utility through a local web UI. It is not a server-control panel, SSH terminal, Docker manager, file browser, plugin marketplace, or fleet-management system. V1 has no process-killing, service-control, write, or arbitrary-command actions.

## Run on a Linux host

Build both executables:

```bash
go build -o cmdry .
go build -o cmdry-ports ./plugins/ports/cmd/cmdry-ports
sudo install -Dm755 cmdry /usr/local/bin/cmdry
sudo install -Dm755 cmdry-ports /opt/cmdry/plugins/cmdry-ports
sudo mkdir -p /opt/cmdry/data
sudo chown "$USER" /opt/cmdry/data
CMDRY_ADDR=127.0.0.1:8080 cmdry serve
```

Open `http://127.0.0.1:8080`. Cmdry binds to localhost by default. If you expose it on a network, put it behind access controls you manage; each installed plugin can execute a host command.

Port Inspector requires `ss`, normally supplied by the `iproute2` package. It returns port data even when `ss` cannot reveal an owning process.

## Docker

The provided Compose file uses Linux host networking so Port Inspector sees the host’s network namespace rather than the container’s. This is intentional. It listens on host port `8087`:

```bash
sudo mkdir -p /home/sottey/docker-data/cmdry/{data,plugins}
go build -o /home/sottey/docker-data/cmdry/plugins/cmdry-ports ./plugins/ports/cmd/cmdry-ports
docker compose up --build -d
```

Then use `http://host:8087`. Host networking is designed for a Linux Docker host. Do not use a conventional isolated container and interpret its container-local port list as the host’s port list.

## Configuration

Configuration is read once at startup:

| Variable | Default | Meaning |
| --- | --- | --- |
| `CMDRY_ADDR` | `127.0.0.1:8080` | HTTP listen address |
| `CMDRY_PLUGIN_DIR` | `/opt/cmdry/plugins` | executable plugin directory |
| `CMDRY_DATA_DIR` | `/opt/cmdry/data` | Cmdry state directory |
| `CMDRY_LOG_LEVEL` | `info` | `debug`, `info`, `warn`, or `error` |

`cmdry serve --addr` and `--plugin-dir` override the matching environment variables.

## Plugins and protocol

At startup Cmdry scans the plugin directory for executable files named `cmdry-*`. It invokes each candidate as:

```text
cmdry-ports manifest
```

The program must print only a valid version-1 JSON manifest to stdout. Cmdry validates the protocol version, safe IDs, version, and unique page/action IDs. Bad binaries are logged and skipped; they never stop Cmdry from starting.

For a page request, Cmdry runs an explicitly registered action with no shell:

```text
cmdry-ports execute list
```

The request arrives as JSON on stdin. The response must be protocol JSON on stdout; diagnostics belong on stderr. Core invocation has a timeout, validates action IDs and UI component types, and treats malformed output, a non-zero exit, and missing binaries as a contained page error.

The supported UI component types are `metric`, `text`, `alert`, `table`, and `actions`. The core renders all values with `html/template`; plugins cannot inject markup, scripts, styles, URLs, or commands.

The included Go SDK is in [`plugin-sdk/go`](plugin-sdk/go). A plugin has this shape:

```go
cmdry.Run(cmdry.Plugin{
    Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "example", Name: "Example", Version: "0.1.0"},
    Actions: map[string]cmdry.Handler{"list": listItems},
})
```

Build a future plugin separately, place the executable in `CMDRY_PLUGIN_DIR`, and restart Cmdry. No Cmdry core modification is needed.

## Security model

Plugins are local executables but are not implicitly trusted. Cmdry uses their absolute discovered path, never uses `sh -c`, passes explicit arguments, validates inputs/outputs, limits request execution time, and HTML-escapes displayed data. Manifest permissions are declarative in V1 and visibly listed on the Plugins page; they are not a sandbox. Install only plugins you trust.

Authentication is deliberately not implemented in V1. Keep Cmdry local by default, or supply your own network and authentication boundary before exposing it.

## Development

```bash
go test ./...
go vet ./...
go build ./...
```

Parser tests use representative `ss` fixtures and never depend on your host’s live ports. On Linux, manually verify discovery, the Ports page, a removed plugin, malformed plugin JSON, and a missing `ss` binary before a release.

## Future ideas

Journal Viewer, Systemd Manager, SMART Monitor, RAID Manager, Rsync Manager, Rclone Manager, Samba Manager, Mount Manager, Fail2ban Manager, SSH Config Manager, Restic Manager, and Cron Manager are possible future plugins. None are currently implemented.
