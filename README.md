# Cmdry

Cmdry is a small, self-hosted web workbench for Linux command-line tools. The core provides navigation, a consistent server-rendered UI, plugin discovery, and strict executable-plugin boundaries. Plugins are separate processes that return versioned JSON—not arbitrary HTML, JavaScript, or shell commands.

The first included plugin is **Port Inspector**, a read-only view of listening TCP and UDP ports from `ss`.

> Screenshot placeholder: Cmdry’s Port Inspector page shows concise listening-port counts above a sortable-ready data table.

## What Cmdry is (and is not)

Cmdry is an extensible way to expose a focused Linux CLI utility through a local web UI. It is not a server-control panel, SSH terminal, Docker manager, file browser, plugin marketplace, or fleet-management system. V1 has no process-killing, service-control, write, or arbitrary-command actions.

## Run on a Linux host

Build the core and bundled plugins:

```bash
./scripts/build.sh
sudo install -Dm755 dist/cmdry /usr/local/bin/cmdry
sudo install -Dm755 dist/plugins/cmdry-ports /opt/cmdry/plugins/cmdry-ports
sudo install -Dm755 dist/plugins/cmdry-journal /opt/cmdry/plugins/cmdry-journal
sudo mkdir -p /opt/cmdry/data
sudo chown "$USER" /opt/cmdry/data
CMDRY_ADDR=127.0.0.1:8080 cmdry serve
```

Open `http://127.0.0.1:8080`. Cmdry binds to localhost by default. If you expose it on a network, put it behind access controls you manage; each installed plugin can execute a host command.

Port Inspector requires `ss`, normally supplied by the `iproute2` package. It returns port data even when `ss` cannot reveal an owning process. Journal Viewer requires `journalctl` and is available only on Linux.

## Run natively on macOS

Cmdry and Port Inspector run natively on macOS. The plugin uses the built-in
`lsof` command: TCP results are listening sockets, while UDP results are
unconnected sockets with an assigned local port. macOS may hide processes that
the current user cannot inspect; Cmdry preserves missing process and PID values.

```bash
./scripts/build.sh
./scripts/run.sh
```

Open `http://127.0.0.1:8080`. This must run directly on the Mac: Docker Desktop
uses a Linux VM, so its host-network mode cannot inspect the macOS host network.

## Docker

The provided Compose file is for a Linux Docker host. It uses host networking so
Port Inspector sees that host’s network namespace rather than the container’s.
This is intentional. It listens on host port `8087`:

```bash
sudo mkdir -p /home/sottey/docker-data/cmdry/{data,plugins}
./scripts/build.sh
sudo install -Dm755 dist/plugins/cmdry-ports /home/sottey/docker-data/cmdry/plugins/cmdry-ports
sudo install -Dm755 dist/plugins/cmdry-journal /home/sottey/docker-data/cmdry/plugins/cmdry-journal
docker compose up --build -d
```

Then use `http://host:8087`. Run the build script on that Linux host so the
plugins are Linux executables. Do not use a conventional isolated container and interpret its
container-local port list as the host’s port list. Docker Desktop on macOS does
not provide macOS host port inspection.

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

See the [plugin development guide](docs/plugin-development.md) for the complete
v1 manifest, response contract, SDK example, local development workflow, and
troubleshooting steps.

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

## Included plugins

- **Port Inspector**: read-only listening TCP and UDP ports, using `ss` on
  Linux and `lsof` on macOS.
- **Journal Viewer**: the 100 newest local journal entries on Linux, using
  `journalctl`. It reports a clear unsupported-platform error on macOS.
- **Process Resource Snapshot**: visible local processes with CPU, memory,
  state, and parent PID, using `ps` on Linux and macOS.
- **Filesystem Inspector**: mounted filesystem capacity and available space,
  using portable `df` output on Linux and macOS.
- **Network Interface Inspector**: local interfaces, assigned addresses, and
  the default gateway, using native network tools on Linux and macOS.
- **System Information**: local OS, kernel, uptime, CPU, memory, and hardware
  facts, using native system sources on Linux and macOS.

## Build bundled binaries

Build the Cmdry core into `dist/` and stage every bundled plugin in
`dist/plugins/`, which is ready to use as `CMDRY_PLUGIN_DIR`:

```bash
./scripts/build.sh
```

Set `CMDRY_BUILD_DIR` to use a different output directory. To stage plugins
directly into an existing Cmdry installation, set `CMDRY_PLUGIN_DIR` to that
installation's configured plugin directory; the script needs write permission
to it. Restart Cmdry after rebuilding.

Run the locally built core with staged plugins using `./scripts/run.sh`. It
defaults to `127.0.0.1:8080`, `dist/plugins/`, and `.cmdry-data/`; override
those locations with `CMDRY_ADDR`, `CMDRY_PLUGIN_DIR`, `CMDRY_DATA_DIR`, or
`CMDRY_BUILD_DIR`. Extra arguments are passed to `cmdry serve`, for example:

```bash
./scripts/run.sh --addr 127.0.0.1:8090
```

Use **Refresh plugins** on the Plugins page after adding, replacing, or removing
a plugin binary. Cmdry scans and atomically replaces its in-memory registry;
if the plugin directory cannot be scanned, it keeps the existing registry.

Drag installed-tool links in the sidebar to reorder them. Cmdry saves the order
to `plugin-order.json` in `CMDRY_DATA_DIR` and restores it on future starts and
plugin refreshes. Newly discovered plugins appear after the saved entries.

## Future ideas

Systemd Manager, SMART Monitor, RAID Manager, Rsync Manager, Rclone Manager,
Samba Manager, Mount Manager, Fail2ban Manager, SSH Config Manager, Restic
Manager, and Cron Manager are possible future plugins. None are currently
implemented.

See [ROADMAP.md](docs/ROADMAP.md) for proposed core improvements and candidate plugin
scopes. It is a planning document, not a release commitment.
