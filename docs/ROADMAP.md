# Cmdry roadmap

This is a planning document, not a promise or release schedule. Cmdry remains a
small, local-first web workbench for focused command-line tools. It is not
intended to become a general server-control panel, terminal, Docker manager, or
fleet-management product.

## Current baseline

Cmdry v0.2.0 provides:

- A Go, server-rendered core with a local-by-default listen address.
- Startup discovery of executable `cmdry-*` plugins.
- A version-1 JSON manifest and action protocol.
- Structured, escaped UI components rather than plugin-supplied HTML or script.
- The Port Inspector plugin, using `ss` on Linux and `lsof` on macOS.

## Core improvements

### Reliability and test coverage

- Add tests for plugin discovery, duplicate IDs, malformed manifests and action
  responses, non-zero exits, timeouts, removed binaries, and server routes.
- Add end-to-end fixture plugins that exercise normal, malformed, slow, and
  error-returning behavior without relying on host commands.
- Capture plugin stderr in structured logs, with bounded output and no sensitive
  request data.
- Add explicit request-body size limits before future parameterized actions are
  introduced.
- Add platform and architecture checks at discovery so incompatible binaries
  produce a clear administrator-facing log message.

### Plugin developer experience

- Publish the Go SDK as a versioned module and document compatibility policy.
- Add a plugin scaffold command or template with manifest, tests, and a local
  development script.
- Define a protocol JSON Schema and use it for fixture validation.
- Add a plugin diagnostics page that lists rejected candidates and the reason
  each was skipped; do not execute unregistered plugins from the UI.
- Document a stable process for plugin upgrades, rollbacks, and compatibility
  checks.

### Plugin lifecycle and operations

- Add an explicit, administrator-controlled plugin rescan endpoint or command;
  preserve startup-only discovery as the safe default until this is designed.
- Persist enabled/disabled state in the data directory, while keeping invalid
  and failed plugins visible as diagnostics.
- Add bounded per-plugin execution telemetry: invocation count, duration, last
  successful invocation, and last sanitized failure.
- Add a plugin health check convention, only if it remains read-only and has a
  clear timeout/response contract.

### UI and accessibility

- Make the sidebar and tables work well on narrow screens.
- Add accessible labels, keyboard-visible sorting, and sorting direction
  indicators to data tables.
- Show when Port Inspector data was last refreshed and provide a clear manual
  refresh action.
- Improve error pages with a concise next step specific to missing commands,
  permission limits, malformed output, unavailable plugins, and timeouts.

### Security and deployment

- Provide a documented reverse-proxy example for deployments that need network
  access, including an authentication boundary managed outside Cmdry.
- Document a minimal-permission service account model for Linux host installs.
- Define a signed/checksummed plugin distribution process before any plugin
  catalog or remote installation workflow is considered.
- Keep permissions declarative until there is a real enforcement design; do not
  imply that the existing manifest permissions sandbox a plugin.
- Add supported deployment examples for native Linux, native macOS, and Linux
  Docker hosts; keep Docker Desktop on macOS explicitly separate because it
  cannot inspect the macOS host network namespace.

### Protocol evolution (only with a versioned design)

- Specify constrained action parameters for filters, pagination, and confirmed
  write operations. This requires request validation, size limits, CSRF
  protection, and an explicit component/input model first.
- Add pagination and server-side sorting metadata for large result sets.
- Add component-level schemas and validation beyond the current component type
  allowlist.
- Define compatibility rules before introducing protocol version 2; existing
  v1 plugins must fail clearly rather than behaving ambiguously.

## Proposed plugin candidates

All candidates below are proposals only. Start each as a read-only plugin,
preserve uncertainty in its output, and add mutating controls only after a
separate security and interaction review.

| Done? | Candidate | Initial read-only scope | Likely host command(s) | Notes |
| --- | --- | --- | --- | --- |
| TRUE | Journal Viewer | Recent system events, severity, and service filters | `journalctl` | Linux-only; pagination is important. |
| FALSE | Systemd Status | Unit status, failures, dependencies, and logs | `systemctl`, `journalctl` | Do not add start/stop/restart controls initially. |
| FALSE | SMART Monitor | Disk identity, health attributes, self-test history | `smartctl` | Hardware access and elevated permissions need careful deployment docs. |
| FALSE | RAID Monitor | Array state, members, rebuild progress | `mdadm`, `/proc/mdstat` | Linux-specific; distinguish degraded from unavailable information. |
| FALSE | Filesystem and Mount Inspector | Mounted filesystems, capacity, options, and sources | `findmnt`, `df` | Read-only first; never expose arbitrary mount commands. |
| FALSE | Backup Status | Recent Restic snapshots, repository health, failures | `restic` | Secrets must remain outside plugin responses and logs. |
| FALSE | Rsync Activity | Parse known job logs and last-result summaries | `rsync`, configured logs | Avoid a general-purpose command runner. |
| FALSE | Rclone Status | Configured remote health and recent job summaries | `rclone` | Treat remote configuration and tokens as sensitive. |
| FALSE | Samba Inspector | Service state and configured shares | `testparm`, `smbstatus` | Display configuration safely; no share editing in an initial version. |
| FALSE | Fail2ban Status | Jails, current bans, and recent events | `fail2ban-client` | Read-only first; unban actions require explicit audit/confirmation design. |
| FALSE | SSH Configuration Inspector | Effective daemon settings and listening status | `sshd -T`, `ss` | Do not expose private keys or enable arbitrary config editing. |
| FALSE | Cron and Timer Viewer | System/user cron entries and systemd timers | `crontab`, `systemctl list-timers` | Scope visibility carefully by service account permissions. |
| FALSE | Network Interface Inspector | Addresses, routes, interfaces, and DNS state | `ip`, `networksetup`, `scutil` | Needs separate Linux and macOS collectors like Port Inspector. |
| TRUE | Process Resource Snapshot | Selected process CPU, memory, and state | `ps` | Keep filtering bounded; it is not a task manager or process killer. |

## Candidate selection criteria

Prioritize a candidate when it:

1. Delivers clear value as a narrow read-only view.
2. Has a stable, scriptable host command with parseable output.
3. Can explain missing permissions or unavailable host data accurately.
4. Has a testable parser with fixtures that do not depend on live host state.
5. Does not require Cmdry core changes merely to render its first useful view.

Defer candidates that fundamentally require a broad command runner, arbitrary
file browser, terminal, remote access layer, or a write-capable control panel.
