# Codex Build Prompt: Modular Linux CLI Web Toolbox

## Project Goal

Build a self-hosted web application in Go that provides a modular web interface for common Linux CLI tools.

The application is **not** intended to become a giant all-purpose server administration panel. Its purpose is narrower:

> Provide a clean, extensible web platform where Linux CLI utilities can be exposed as independently installable plugins.

The core application should handle the shared web experience, plugin lifecycle, security boundaries, navigation, rendering, and communication protocol. Individual plugins should handle interaction with specific Linux CLI tools.

The first release should include:

1. The **Toolbox Core** application.
2. A plugin SDK/protocol.
3. One first-party plugin: **Port Inspector**.
4. Enough architecture to support future plugins without redesigning the application.

Do not implement additional plugins yet.

---

# Core Product Principles

Follow these principles throughout the implementation.

## 1. Modular First

The core application must not contain Linux-tool-specific logic.

For example:

- The core should not know how `ss` works.
- The core should not parse `lsof`.
- The core should not know anything about `smartctl`, `journalctl`, `systemctl`, `rsync`, etc.

Those responsibilities belong to plugins.

The core only needs to know how to:

- discover plugins
- read plugin manifests
- invoke plugin actions
- validate plugin responses
- expose plugin-provided navigation
- render plugin-provided structured UI data
- manage plugin state

---

## 2. Plugins Are Separate Executables

Do **not** use Go's built-in `plugin` package.

Plugins should be standalone executables.

Example installation directory:

```text
/opt/toolbox/plugins/
```

Example binaries:

```text
toolbox-ports
toolbox-journal
toolbox-smart
toolbox-rsync
```

This architecture is intentional because plugins should eventually be able to be written in languages other than Go.

The plugin protocol must therefore be language-neutral.

Use JSON for communication.

---

## 3. Structured UI, Not Arbitrary Plugin HTML

Plugins must not initially be allowed to inject arbitrary HTML, JavaScript, or CSS into the core application.

Plugins return structured UI definitions and structured data.

The Toolbox Core renders those definitions using shared components.

This ensures:

- consistent styling
- consistent layout
- responsive behavior
- dark mode support
- centralized security
- centralized escaping
- easier plugin development

The first version only needs a small set of UI components.

Suggested initial types:

```text
table
metric
text
alert
actions
```

Design the protocol so additional component types can be added later without breaking existing plugins.

---

# Technology

Use:

- Go
- standard Go HTTP stack unless a very small router is genuinely helpful
- server-rendered HTML
- minimal JavaScript
- SQLite for core application state if persistence is required
- JSON as the plugin protocol
- Linux as the initial supported runtime platform

Avoid unnecessary frontend frameworks.

Do not introduce React, Vue, Angular, Node build tooling, or a separate SPA unless there is an unavoidable reason.

Prefer a simple deployable Go binary.

---

# Initial Repository Structure

Use a structure similar to:

```text
toolbox/
├── cmd/
│   └── toolbox/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── plugins/
│   │   ├── discovery.go
│   │   ├── manifest.go
│   │   ├── protocol.go
│   │   ├── registry.go
│   │   └── runner.go
│   │
│   ├── server/
│   ├── storage/
│   └── ui/
│
├── plugin-sdk/
│   └── go/
│
├── plugins/
│   └── ports/
│       └── cmd/
│           └── toolbox-ports/
│               └── main.go
│
├── web/
│   ├── templates/
│   └── static/
│
├── docker/
│
├── go.mod
├── README.md
└── LICENSE
```

Adjust this only where necessary.

Do not reorganize merely for stylistic preference.

---

# Toolbox Core Responsibilities

The core application should implement the following.

## Web Server

Provide a web UI with:

- dashboard/home page
- sidebar navigation
- installed plugin list
- plugin pages
- plugin error states
- basic settings page
- health endpoint

Suggested routes:

```text
GET /
GET /plugins
GET /plugins/{pluginID}
GET /plugins/{pluginID}/{page}
POST /plugins/{pluginID}/actions/{action}
GET /settings
GET /health
```

Exact routes may differ if there is a cleaner implementation.

---

# Plugin Discovery

At startup, scan a configured plugin directory.

Default:

```text
/opt/toolbox/plugins
```

Also support overriding this through an environment variable.

Suggested:

```text
TOOLBOX_PLUGIN_DIR
```

Only executable files matching:

```text
toolbox-*
```

should be considered plugins.

For each candidate binary:

1. execute the binary with:

```bash
toolbox-plugin manifest
```

2. parse JSON from stdout
3. validate the manifest
4. register the plugin if valid
5. log a clear warning and continue if invalid

A broken plugin must **not** prevent the Toolbox Core from starting.

---

# Plugin Manifest

Define a stable manifest format.

Initial example:

```json
{
  "protocol_version": 1,
  "id": "ports",
  "name": "Port Inspector",
  "version": "0.1.0",
  "description": "Inspect listening network ports and their owning processes.",
  "category": "system",
  "icon": "network",
  "pages": [
    {
      "id": "overview",
      "name": "Ports",
      "default": true
    }
  ],
  "permissions": [
    "network.read",
    "process.read"
  ],
  "actions": [
    {
      "id": "list",
      "name": "List Ports",
      "method": "read"
    }
  ]
}
```

Validate:

- protocol version
- plugin ID
- plugin name
- semantic-ish version string
- unique page IDs
- unique action IDs

Plugin IDs should use a safe restricted format such as:

```text
[a-z0-9][a-z0-9-_]*
```

---

# Plugin Invocation Protocol

The initial implementation may invoke a plugin executable separately for each request.

Example:

```bash
toolbox-ports execute list
```

Request parameters should be provided using JSON over stdin.

Example stdin:

```json
{
  "action": "list",
  "params": {}
}
```

Example response:

```json
{
  "ok": true,
  "data": {
    "components": []
  }
}
```

Error example:

```json
{
  "ok": false,
  "error": {
    "code": "COMMAND_FAILED",
    "message": "Unable to execute ss"
  }
}
```

Requirements:

- stdout must contain protocol JSON only
- plugin diagnostic output should go to stderr
- the core should impose an execution timeout
- the core must handle malformed JSON
- the core must handle non-zero exit status
- the core must handle missing plugin binaries
- plugin failure must not crash the server

Design the invocation abstraction so a persistent-process or socket-based plugin protocol could be added later without rewriting the rest of the application.

Do **not** implement persistent plugin processes yet.

---

# UI Response Format

Plugin actions should return structured UI components.

Example:

```json
{
  "ok": true,
  "data": {
    "title": "Listening Ports",
    "components": [
      {
        "type": "metric",
        "label": "Listening Ports",
        "value": "17"
      },
      {
        "type": "table",
        "id": "ports",
        "columns": [
          {
            "key": "port",
            "label": "Port"
          },
          {
            "key": "protocol",
            "label": "Protocol"
          },
          {
            "key": "address",
            "label": "Address"
          },
          {
            "key": "process",
            "label": "Process"
          },
          {
            "key": "pid",
            "label": "PID"
          }
        ],
        "rows": [
          {
            "port": 22,
            "protocol": "tcp",
            "address": "0.0.0.0",
            "process": "sshd",
            "pid": 881
          }
        ]
      }
    ]
  }
}
```

The core UI should render this.

For the first version, implement these component types:

### metric

Fields:

```text
label
value
description optional
```

### text

Fields:

```text
title optional
text
```

### alert

Fields:

```text
level
title optional
message
```

Levels:

```text
info
success
warning
error
```

### table

Fields:

```text
id
columns
rows
```

The table should support client-side sorting if practical without introducing a framework.

### actions

A simple list of buttons linked to registered plugin actions.

Do not allow arbitrary shell commands or URLs to be declared through the UI schema.

---

# Core Navigation

The sidebar should be generated from installed plugin manifests.

Example:

```text
Toolbox

SYSTEM
  Ports
  Services
  Logs

STORAGE
  Disks
  SMART
  RAID

AUTOMATION
  Rsync
  Cron
```

For the first release only `Ports` will exist.

The category value from the plugin manifest determines grouping.

Supported categories can initially include:

```text
system
storage
network
automation
security
other
```

Unknown categories should safely fall back to `other`.

---

# Plugin Registry

Create an in-memory registry containing:

- plugin ID
- manifest
- executable path
- enabled state
- validation status
- last error if applicable

Persist enable/disable state if practical.

The plugin list page should show:

```text
Name
Version
Category
Status
Permissions
Path
```

Possible statuses:

```text
enabled
disabled
invalid
error
```

---

# Security Requirements

This application will execute system commands through plugins.

Treat that as a major security boundary.

## Command Execution

Never invoke commands using:

```bash
sh -c
```

unless absolutely unavoidable.

Prefer Go's:

```go
exec.CommandContext()
```

with explicit argument arrays.

Never concatenate user input into shell commands.

---

## Plugin Execution

The core should:

- use explicit executable paths
- enforce timeouts
- limit request sizes
- validate plugin IDs
- validate action IDs
- validate JSON responses
- escape all rendered output

Do not trust plugins merely because they are local.

---

## Permissions

For v1, permissions may be declarative only.

Example:

```json
{
  "permissions": [
    "network.read",
    "process.read"
  ]
}
```

The UI should display requested permissions.

Do not implement complex sandboxing yet.

However, keep permissions in the manifest and internal model from the beginning.

Future permissions may include:

```text
network.read
network.write
process.read
process.kill
systemd.read
systemd.control
storage.read
storage.write
filesystem.read
filesystem.write
backup.execute
```

---

# Authentication

Do not overbuild authentication in the first pass.

At minimum, structure the server so authentication middleware can be added cleanly.

If implementing authentication now, use a simple local user/password system with secure password hashing and session cookies.

Do not implement:

- OAuth
- SSO
- LDAP
- multi-tenant permissions
- external identity providers

unless the basic architecture requires an abstraction point.

---

# Configuration

Use environment variables.

Suggested configuration:

```text
TOOLBOX_ADDR=:8080
TOOLBOX_PLUGIN_DIR=/opt/toolbox/plugins
TOOLBOX_DATA_DIR=/opt/toolbox/data
TOOLBOX_LOG_LEVEL=info
```

Provide reasonable defaults.

Configuration should be loaded once at startup into a typed struct.

---

# Logging

Use structured or consistently formatted logging.

The core should log:

- startup
- plugin discovery
- plugin registration
- invalid plugins
- plugin invocation
- invocation duration
- plugin errors
- server errors

Do not log passwords, tokens, or sensitive request bodies.

---

# First-Party Plugin: Port Inspector

Implement the first plugin:

```text
toolbox-ports
```

Purpose:

> Show listening ports on the Linux host and identify the process that owns each port.

The plugin should initially be **read-only**.

Do not implement process killing or service restart actions yet.

---

# Port Inspector Data Sources

Primary source:

```bash
ss
```

Use:

```bash
ss -H -l -n -t -u -p
```

or another appropriate noninteractive invocation.

The exact command should be chosen based on reliable parsing.

The plugin should gather:

```text
protocol
local address
port
process name
PID
```

Where possible.

If process information cannot be retrieved because of permissions, the plugin should still return the port and protocol data rather than fail completely.

---

# Optional Port Enrichment

If practical without making v1 fragile, enrich port data using:

```text
/proc
systemctl
docker
```

However, this enrichment is secondary.

The first requirement is reliable port discovery.

Do not make Docker a required dependency.

If Docker exists, the plugin may attempt to correlate ports with containers.

If Docker is not installed, the plugin should behave normally.

A possible future row might eventually contain:

```json
{
  "port": 8081,
  "protocol": "tcp",
  "address": "0.0.0.0",
  "process": "docker-proxy",
  "pid": 22101,
  "container": "8081-logoura",
  "service": null
}
```

For the first implementation, only include fields that can be obtained reliably.

Do not fake or infer missing relationships.

---

# Port Inspector UI

The default page should show:

## Summary

Examples:

```text
Listening Ports: 17
TCP: 13
UDP: 4
```

## Ports Table

Columns:

```text
Port
Protocol
Address
Process
PID
```

Allow sorting by:

```text
Port
Protocol
Process
PID
```

If practical, include a simple search/filter box.

Example displayed data:

```text
22     TCP    0.0.0.0    sshd          881
53     UDP    127.0.0.53 systemd-resolved 712
8081   TCP    0.0.0.0    docker-proxy  22101
8888   TCP    0.0.0.0    bark2ntfy     3543012
```

---

# Port Inspector Plugin SDK Usage

Create a small Go SDK so the plugin implementation does not manually parse every protocol detail.

Desired developer experience should resemble:

```go
func main() {
    toolbox.Run(toolbox.Plugin{
        Manifest: toolbox.Manifest{
            ProtocolVersion: 1,
            ID:              "ports",
            Name:            "Port Inspector",
            Version:         "0.1.0",
            Description:     "Inspect listening ports and owning processes.",
            Category:        "system",
        },

        Actions: map[string]toolbox.Handler{
            "list": listPorts,
        },
    })
}
```

The SDK should handle:

- `manifest`
- `execute`
- stdin decoding
- stdout JSON encoding
- common errors
- exit codes

Keep the SDK small.

Do not hide important behavior behind excessive abstraction.

---

# Plugin Developer Experience

A future developer should be able to build a plugin by:

1. creating a standalone executable
2. implementing the manifest command
3. implementing one or more actions
4. copying the binary into the plugin directory

Example:

```bash
go build -o toolbox-example ./cmd/toolbox-example
cp toolbox-example /opt/toolbox/plugins/
```

After restart, Toolbox should discover it automatically.

Eventually hot-reload may be useful, but do not implement it now unless extremely simple.

---

# Docker Support

The application should be easy to run in Docker.

Provide:

```text
Dockerfile
docker-compose.yml
```

The main application container must be able to access plugin binaries.

Example mount:

```yaml
volumes:
  - /home/sottey/docker-data/toolbox/data:/opt/toolbox/data
  - /home/sottey/docker-data/toolbox/plugins:/opt/toolbox/plugins
```

Use:

```yaml
restart: unless-stopped
```

For the initial compose file, use host port:

```text
8087
```

and container port:

```text
8080
```

Container name:

```text
8087-toolbox
```

Preferred host bind mounts:

```text
/home/sottey/docker-data/toolbox/data
/home/sottey/docker-data/toolbox/plugins
```

Be aware that a Port Inspector running *inside a normal isolated Docker container* will inspect the container's network namespace rather than necessarily the host's complete network state.

Therefore:

- document this clearly
- design deployment so host inspection is possible
- do not silently return misleading data

For a first implementation, either:

1. recommend host-network deployment for host port inspection, or
2. support running Toolbox directly on the host

Choose the approach that provides accurate results with the least complexity.

Do not pretend container-local results represent the host.

---

# Error Handling

Errors should be useful to a Linux administrator.

Bad:

```text
Internal Server Error
```

Better:

```text
Port Inspector failed

Unable to execute `ss`.

The command was not found on this system.

Install the `iproute2` package and try again.
```

Whenever possible distinguish:

```text
command missing
permission denied
command failed
plugin timeout
invalid plugin response
plugin unavailable
unsupported protocol version
```

Do not expose stack traces in the browser.

---

# Testing

Add tests for the important boundaries.

At minimum:

## Core

- manifest validation
- plugin discovery
- duplicate plugin IDs
- invalid executable handling
- unsupported protocol version
- plugin timeout
- malformed JSON
- non-zero plugin exit
- valid plugin response
- UI component decoding

## Port Plugin

- parsing representative `ss` output
- TCP entries
- UDP entries
- IPv4 entries
- IPv6 entries
- entries without process information
- malformed/unexpected lines

Do not make tests depend on the host's actual current network state.

Use fixtures.

---

# README

Write a practical README containing:

1. what Toolbox is
2. project goals
3. screenshots placeholder section
4. installation
5. direct host usage
6. Docker usage
7. plugin directory
8. how plugin discovery works
9. Port Inspector requirements
10. plugin SDK example
11. plugin protocol
12. security model
13. future plugin ideas
14. development instructions

Future plugin ideas can include:

```text
Journal Viewer
Systemd Manager
SMART Monitor
RAID Manager
Rsync Manager
Rclone Manager
Samba Manager
Mount Manager
Fail2ban Manager
SSH Config Manager
Restic Manager
Cron Manager
```

Clearly mark these as future ideas, not implemented features.

---

# Visual Design

Aim for a calm, polished self-hosted admin interface.

Requirements:

- clean sidebar
- restrained styling
- good whitespace
- dark mode friendly
- readable tables
- clear status indicators
- desktop-first but responsive

Avoid:

- excessive gradients
- dashboard clutter
- fake analytics
- giant cards everywhere
- overly animated UI
- copying the visual style of Portainer or Cockpit too closely

The product should feel like a focused modern utility.

---

# Explicit Non-Goals for V1

Do **not** build any of these unless needed as a tiny supporting abstraction:

- full server management
- file browser
- Docker manager
- Kubernetes support
- SSH terminal
- remote server fleet management
- metrics collection
- Prometheus integration
- Grafana integration
- package management
- user management
- firewall editing
- service restart
- process killing
- automatic plugin marketplace
- plugin downloads
- plugin auto-updating
- plugin sandboxing
- arbitrary plugin HTML
- arbitrary plugin JavaScript
- remote plugin execution
- cluster support
- multi-host agents
- WebSockets unless actually needed
- React
- Vue
- frontend build pipeline
- gRPC

Keep the first release focused.

---

# Development Order

Implement in this order.

## Phase 1: Skeleton

- repository structure
- config
- web server
- templates
- basic navigation
- health route

## Phase 2: Plugin Protocol

- manifest structs
- validation
- discovery
- registry
- runner
- timeout handling
- error handling

Create a fake/test plugin to validate the protocol.

## Phase 3: Structured UI

- component schema
- table renderer
- metric renderer
- text renderer
- alert renderer
- action renderer

## Phase 4: Go Plugin SDK

- manifest support
- action registration
- request decoding
- response encoding
- common error types

## Phase 5: Port Inspector

- invoke `ss`
- parse output
- normalize data
- return structured components
- render summary
- render table

## Phase 6: Plugin Management UI

- installed plugin list
- status
- version
- permissions
- errors
- enable/disable if persistence is implemented

## Phase 7: Packaging

- Dockerfile
- compose
- sample config
- README
- build instructions

## Phase 8: Tests and Cleanup

- unit tests
- fixture tests
- error paths
- lint/build verification

---

# Acceptance Criteria

The project is successful when the following workflow works.

A user starts Toolbox.

Toolbox scans:

```text
/opt/toolbox/plugins
```

and finds:

```text
toolbox-ports
```

Toolbox executes:

```bash
toolbox-ports manifest
```

and registers the plugin.

The browser sidebar automatically shows:

```text
SYSTEM
  Ports
```

The user clicks `Ports`.

Toolbox invokes the plugin.

The plugin executes `ss`, parses the results, and returns structured JSON.

Toolbox renders:

```text
Listening Ports

Port   Protocol   Address      Process        PID
22     TCP        0.0.0.0      sshd           881
8081   TCP        0.0.0.0      docker-proxy   22101
8888   TCP        0.0.0.0      bark2ntfy      3543012
```

If the plugin binary is deleted or broken, Toolbox remains running and displays an appropriate plugin error.

A second independently written executable that implements protocol version 1 should be installable later without modifying Toolbox Core.

That extensibility is the primary architectural requirement.

---

# Coding Expectations

Use idiomatic Go.

Prefer:

- explicit interfaces
- small packages
- clear data structures
- dependency injection where it meaningfully helps testing
- `context.Context`
- `exec.CommandContext`
- `html/template`
- standard error wrapping
- testable parsers

Avoid:

- giant interfaces
- unnecessary frameworks
- reflection-heavy abstractions
- premature generics
- clever metaprogramming
- unnecessary concurrency
- global mutable state
- hidden shell execution
- speculative features

Keep the architecture understandable to another Go developer reading the repository.

---

# Important Implementation Rule

When information is unavailable, preserve that uncertainty.

For example, if `ss` identifies a port but not its owning process, return:

```json
{
  "port": 8080,
  "process": null,
  "pid": null
}
```

Do not infer or invent values.

The same rule should apply to all future plugins.

---

# Deliverables

Produce a working repository containing:

- Toolbox Core
- plugin protocol
- Go plugin SDK
- Port Inspector plugin
- server-rendered web UI
- plugin discovery
- plugin registry
- structured UI renderer
- error handling
- tests
- Dockerfile
- docker-compose.yml
- README

Before declaring the implementation complete:

1. run `go test ./...`
2. run `go vet ./...`
3. build the core binary
4. build the Port Inspector plugin
5. verify plugin discovery manually
6. verify the Ports page manually
7. verify behavior when the plugin is removed
8. verify behavior when the plugin returns invalid JSON
9. verify behavior when `ss` is unavailable

Do not mark the project complete if any of these fail.

---

# Final Guidance

The most important thing about this project is **not Port Inspector**.

Port Inspector is only the first proof that the architecture works.

The real product is:

> A safe, simple, modular web platform for turning Linux CLI tools into consistent self-hosted web interfaces.

Every architectural decision should therefore be evaluated using this question:

> Could another developer write and install a new plugin without modifying the Toolbox Core?

If the answer is no, reconsider the design.
