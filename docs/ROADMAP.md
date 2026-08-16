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

## Prioritized product enhancements

### Priority 1

- [x] Add a keyboard-first command palette (`Cmd+K` / `Ctrl+K`) that searches
  plugin names, descriptions, and search terms, with keyboard navigation and
  recently used tools.
- [x] Add persistent Favorites and Recent Tools sections above the grouped sidebar.
  Favorites should be explicit user choices; recent tools should be local-only
  and bounded.
- [x] Add a plugin diagnostics page that lists rejected candidates, refresh time,
  manifest failures, and safely bounded stderr excerpts. It must never execute
  unregistered plugins from the UI.
- [x] Add a plugin detail page with the manifest description, version, declared
  permissions, search terms, binary path, platform notes, and a copyable
  diagnostic command.
- [x] Add user-managed navigation settings: manage Favorites, clear Recent Tools,
  choose the bounded recent-tools limit, and restore default sidebar order and
  group expansion state.
- [x] Add appearance and accessibility settings: system/light/dark theme, reduced
  motion, sidebar density, and an Overview-or-last-used default landing page.
- Add a persisted enabled/disabled state for installed plugins. Disabled tools
  remain visible in Plugins and Diagnostics and can be re-enabled without
  deleting their binary.

### Priority 3

- Add opt-in non-sensitive per-plugin history and named presets, with a clear
  “run again” workflow. Secret inputs such as Bcrypt values must opt out and
  never be retained.
- Design a safe file-input protocol for future plugins: explicit permission,
  schema validation, bounded upload size, temporary-file lifecycle, and no
  unrestricted filesystem access.
- Add privacy and local-data controls that show exactly what Cmdry retains and
  let users clear navigation state, saved history, and presets independently.
  Any future history or presets must remain opt-in and exclude secret inputs.
- Add diagnostics preferences, including safe clearing of retained scan reports
  and a redacted, copyable support snapshot that excludes form data and
  environment secrets.

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

- Extend the bounded text-input model beyond the existing form component only
  with explicit schemas, request validation, size limits, and CSRF protection.
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
| TRUE | Port Inspector | Scans ports |`lsof`, `ss` | macos, linux supported. |
| TRUE | Journal Viewer | Recent system events, severity, and service filters | `journalctl` | Linux-only; pagination is important. |
| FALSE | Systemd Status | Unit status, failures, dependencies, and logs | `systemctl`, `journalctl` | Do not add start/stop/restart controls initially. |
| FALSE | SMART Monitor | Disk identity, health attributes, self-test history | `smartctl` | Hardware access and elevated permissions need careful deployment docs. |
| FALSE | RAID Monitor | Array state, members, rebuild progress | `mdadm`, `/proc/mdstat` | Linux-specific; distinguish degraded from unavailable information. |
| TRUE | Filesystem Inspector | Mounted filesystems, capacity, and available space | `df` | Cross-platform and read-only; mount management remains deferred. |
| FALSE | Backup Status | Recent Restic snapshots, repository health, failures | `restic` | Secrets must remain outside plugin responses and logs. |
| FALSE | Rsync Activity | Parse known job logs and last-result summaries | `rsync`, configured logs | Avoid a general-purpose command runner. |
| FALSE | Rclone Status | Configured remote health and recent job summaries | `rclone` | Treat remote configuration and tokens as sensitive. |
| FALSE | Samba Inspector | Service state and configured shares | `testparm`, `smbstatus` | Display configuration safely; no share editing in an initial version. |
| FALSE | Fail2ban Status | Jails, current bans, and recent events | `fail2ban-client` | Read-only first; unban actions require explicit audit/confirmation design. |
| FALSE | SSH Configuration Inspector | Effective daemon settings and listening status | `sshd -T`, `ss` | Do not expose private keys or enable arbitrary config editing. |
| TRUE | Scheduled Tasks | User cron jobs plus platform-native scheduled tasks | `crontab`, `systemctl list-timers`, launchd plist files | Cross-platform read-only view; loaded status and task editing remain deferred. |
| TRUE | Network Interface Inspector | Interfaces, assigned addresses, and default gateway | `ifconfig`, `route`, `ip` | Cross-platform read-only baseline; DNS and traffic counters remain deferred. |
| TRUE | System Information | OS, uptime, CPU, memory, and hardware facts | `sw_vers`, `sysctl`, `/proc` | Cross-platform read-only baseline; protected fields remain explicitly unavailable. |
| TRUE | Battery and Power Inspector | Battery state and current power source | `pmset`, `/sys/class/power_supply` | Cross-platform read-only baseline; battery health remains deferred. |
| TRUE | Wi-Fi Inspector | Active Wi-Fi network and connection facts | `networksetup`, `nmcli` | Cross-platform read-only baseline; saved networks and credentials remain excluded. |
| TRUE | Process Resource Snapshot | Selected process CPU, memory, and state | `ps` | Keep filtering bounded; it is not a task manager or process killer. |
| TRUE | JSON to CSV | Convert pasted objects or arrays of objects to a browser-local CSV download | Built-in JSON and CSV encoding | Cross-platform and one-shot; no host file access, watch folders, or background service. |
| TRUE | CSV to JSON | Convert pasted header-based CSV records to a browser-local JSON download | Built-in CSV and JSON encoding | Cross-platform and one-shot; values remain strings rather than inferred types. |
| TRUE | JSON Compare | Compare two pasted JSON documents structurally | Built-in JSON decoding | Cross-platform and one-shot; object key order and whitespace are ignored, array order is preserved. |
| TRUE | Hidden Character Detector | Identify invisible Unicode formatting and non-standard whitespace in pasted text | Built-in Unicode tables | Cross-platform and one-shot; normal spaces, tabs, and line breaks are excluded. |

## Candidate selection criteria

Prioritize a candidate when it:

1. Delivers clear value as a narrow read-only view.
2. Has a stable, scriptable host command with parseable output.
3. Can explain missing permissions or unavailable host data accurately.
4. Has a testable parser with fixtures that do not depend on live host state.
5. Does not require Cmdry core changes merely to render its first useful view.

Defer candidates that fundamentally require a broad command runner, arbitrary
file browser, terminal, remote access layer, or a write-capable control panel.
