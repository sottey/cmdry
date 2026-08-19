# Cmdry

Cmdry is a small, self-hosted web workbench for focused local tools. The Go core provides navigation, a consistent server-rendered UI, plugin discovery, and strict executable-plugin boundaries. Plugins are separate processes that return versioned JSON—not arbitrary HTML, JavaScript, or shell commands.

It includes host-inspection tools for Linux and macOS alongside local text, list, CSV, JSON, XML, time, and calculator utilities. Cmdry is local-first: transformation plugins work on pasted data and do not read or write host files.

## What Cmdry is (and is not)

Cmdry is an extensible way to expose a focused Linux CLI utility through a local web UI. It is not a server-control panel, SSH terminal, Docker manager, file browser, plugin marketplace, or fleet-management system. V1 has no process-killing, service-control, filesystem-writing, or arbitrary-command actions.

## Run on a Linux host

Build the core and bundled plugins:

```bash
sudo mkdir -p /opt/cmdry/{data,plugins}
sudo CMDRY_PLUGIN_DIR=/opt/cmdry/plugins ./scripts/build.sh
sudo install -Dm755 dist/cmdry /usr/local/bin/cmdry
sudo chown "$USER" /opt/cmdry/data
CMDRY_ADDR=127.0.0.1:8080 cmdry serve
```

Open `http://127.0.0.1:8080`. Cmdry binds to localhost by default. If you expose it on a network, put it behind access controls you manage; each installed plugin can execute a host command.

Port Inspector requires `ss`, normally supplied by the `iproute2` package. It returns port data even when `ss` cannot reveal an owning process. Journal Viewer requires `journalctl` and is available only on Linux.

## Run natively on macOS

Cmdry and its cross-platform plugins run natively on macOS. Port Inspector uses the built-in
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
sudo CMDRY_PLUGIN_DIR=/home/sottey/docker-data/cmdry/plugins ./scripts/build.sh
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

The program must print only a valid version-1 JSON manifest to stdout. Cmdry validates the protocol version, safe IDs, version, and unique page/action IDs. Bad binaries are logged and skipped; they never stop Cmdry from starting. The **Diagnostics** page retains the latest scan's rejected candidates and bounded stderr excerpts. It is read-only: invalid candidates are never registered or executable through the UI.

For a page request, Cmdry runs an explicitly registered action with no shell:

```text
cmdry-ports execute list
```

The request arrives as JSON on stdin. The response must be protocol JSON on stdout; diagnostics belong on stderr. Core invocation has a timeout, validates action IDs and UI component types, and treats malformed output, a non-zero exit, and missing binaries as a contained page error.

The supported UI component types are `metric`, `text`, `code`, `alert`, `table`, `actions`, `form`, and `download`. Forms submit bounded local input to a declared plugin action; downloads are returned as in-memory, browser-local files. The core renders all values with `html/template`; plugins cannot inject markup, scripts, styles, URLs, or commands.

The included Go SDK is in [`plugin-sdk/go`](plugin-sdk/go). A plugin has this shape:

```go
cmdry.Run(cmdry.Plugin{
    Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "example", Name: "Example", Version: "0.1.0"},
    Actions: map[string]cmdry.Handler{"list": listItems},
})
```

Build a future plugin separately, place the executable in `CMDRY_PLUGIN_DIR`, and use **Refresh plugins** from Cmdry’s Plugins page. No Cmdry core modification is needed.

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
- **Scheduled Tasks**: user cron jobs plus Linux systemd timers or macOS launchd
  configuration files. It remains read-only and explicitly reports unavailable
  sources.
- **Battery and Power Inspector**: battery charge state and power source on
  macOS and Linux, with a clear no-battery view for desktop hosts.
- **Wi-Fi Inspector**: active Wi-Fi network facts without exposing saved
  networks or passwords, on macOS and Linux. On its first use on macOS, it
  requests Location Services permission; this is required by macOS before an
  app may read SSID, BSSID, and signal data.
- **JSON to CSV**: converts pasted JSON objects or arrays of objects into a
  browser-local CSV download on macOS and Linux. It does not read, watch, or
  write host files.
- **CSV to JSON**: converts pasted header-based CSV into a browser-local JSON
  download on macOS and Linux. All CSV values remain strings to preserve the
  source data exactly.
- **JSON Compare**: compares two pasted JSON documents structurally on macOS
  and Linux, ignoring object property order and whitespace.
- **Hidden Character Detector**: finds pasted zero-width, bidirectional,
  non-standard-space, and unexpected control characters on macOS and Linux.
- **Byte Converter**: converts local byte values across binary units (KiB,
  MiB, GiB) and decimal units (KB, MB, GB) on macOS and Linux.
- **JSON Validator, Minifier, and Stringifier**: validate JSON, remove its
  unnecessary whitespace, or encode a valid JSON document as a JSON string
  literal—locally on macOS and Linux.
- **Email Extractor, Remove Duplicate Lines, and Text Replacer**: extract
  unique addresses, preserve unique lines, or make exact text replacements
  locally on macOS and Linux.
- **Uppercase, Lowercase, URL Encoder, and URL Decoder**: change Unicode text
  case or encode/decode URL query values locally on macOS and Linux.
- **Unix Date Converter**: converts whole-second Unix timestamps and RFC 3339
  dates locally, displaying UTC and local-time values on macOS and Linux.
- **Text Statistics**: counts characters, words, lines, UTF-8 bytes, and an
  estimated reading time for pasted text locally on macOS and Linux.
- **Sum**: calculates a local sum, count, average, minimum, and maximum from
  pasted numbers on macOS and Linux.
- **Cron Expression Explainer**: validates standard five-field cron schedules,
  including ranges, lists, steps, names, and common `@` shortcuts, locally on
  macOS and Linux.
- **JSON String Escaper**: escapes arbitrary text into a JSON string literal
  locally on macOS and Linux, without requiring the input itself to be JSON.
- **XML Validator**: validates pasted XML locally and reports parser errors,
  including line locations where the XML parser provides them.
- **Bcrypt**: creates and verifies bcrypt password hashes locally, with a
  configurable work factor and masked browser form fields.
- **Change CSV Separator**: converts comma, semicolon, tab, or pipe-separated
  CSV locally while preserving quoted fields.
- **Chunk List**: splits a newline-delimited list into fixed-size groups
  locally, ignoring blank lines.
- **CSV to TSV**: converts quoted comma-separated CSV to tab-separated values
  locally on macOS and Linux.
- **Days to Hours and Hours to Days**: convert finite decimal durations locally
  on macOS and Linux, including negative offsets.
- **Seconds to Time, Time to Decimal, and Time to Seconds**: convert signed
  whole-second and `HH:MM:SS` / `MM:SS` durations locally on macOS and Linux.
- **Extract Substring**: extracts a one-based, inclusive Unicode character
  range from pasted text locally on macOS and Linux.
- **Find Incomplete CSV Records**: reports comma-separated CSV rows and header
  columns with empty values locally on macOS and Linux.
- **Find Most Popular**: counts and ranks exact newline-delimited list items
  locally on macOS and Linux.
- **Find Unique**: keeps the first exact occurrence of each nonblank
  newline-delimited list item locally on macOS and Linux.
- **Generate Random Numbers**: creates cryptographically secure integers in an
  inclusive configured range locally on macOS and Linux.
- **Join Text**: combines nonblank newline-delimited items with a common or
  custom separator locally on macOS and Linux.
- **QR Code Generator**: encodes pasted text or URLs into local PNG downloads
  on macOS and Linux.
- **Random Number Generator**: creates cryptographically secure fixed-precision
  decimal values in inclusive local ranges on macOS and Linux.
- **Reverse List and Reverse Text**: reverse nonblank list item order or
  Unicode text characters locally on macOS and Linux.
- **Repeat Text, Rotate Text, and Text to Morse**: repeat text, rotate Unicode
  characters, or encode supported text as International Morse code locally.
- **Text Censor and Text Quoter**: mask selected whole words or add single,
  double, or curly quotes around each nonblank text line locally.
- **Rotate List and Shuffle List**: rotate list positions or use a
  cryptographically secure shuffle of pasted list order locally.
- **Wrap List and Unwrap List**: add or remove matching prefixes and suffixes
  around newline-delimited list items locally.
- **Swap CSV Columns and Transpose CSV**: swap two named headers or transpose
  a complete header-based CSV table locally.
- **Round-Trip Voltage Drop**: estimate copper-cable DC voltage drop and power
  loss from length, conductor area, and current; it is an estimate, not an
  electrical design result.
- **Slackline Tension**: estimates a simplified centered static load only and
  is explicitly not appropriate for safety, rigging, or equipment decisions.
- **Truncate Clock Time**: drop smaller units from a signed `HH:MM:SS` or
  `MM:SS` duration locally.
- **Arithmetic Sequence and Check Leap Years**: generate finite numeric
  progressions or find Gregorian leap years in a bounded range locally.
- **Resize Image, Crop Image, and Rotate Image**: transform one uploaded PNG,
  JPEG, or GIF entirely in memory and return a browser-local PNG download.
- **Image Opacity, Color Replacement, and Transparent PNG**: adjust alpha,
  replace a matching color, or make a matching color transparent in one
  uploaded PNG, JPEG, or GIF entirely in memory.
- **JPEG, PNG, and WebP Conversion/Compression**: re-encode one uploaded image
  locally as JPEG, lossless PNG, or WebP with the appropriate quality or
  compression control.

- **Watermark Images**: add centered text with a selected color and opacity to
  one uploaded image and download a local PNG result.

- **Images to PDF**: turn up to four uploaded PNG, JPEG, or GIF images into a
  scaled A4 portrait or landscape PDF entirely in memory.

  Uploaded images are limited to 4 MiB and are never exposed to plugins as host
  paths or retained after the request.

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
