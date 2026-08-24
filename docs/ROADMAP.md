# Cmdry roadmap

This is a planning document, not a promise or release schedule. Cmdry remains a
small, local-first web workbench for focused command-line tools. It is not
intended to become a general server-control panel, terminal, Docker manager, or
fleet-management product.

## Current baseline

Cmdry v0.12.0 provides:

- A Go, server-rendered core with a local-by-default listen address.
- Startup discovery of executable `cmdry-*` plugins.
- A version-1 JSON manifest and action protocol.
- Structured, escaped UI components rather than plugin-supplied HTML or script.
- The Port Inspector plugin, using `ss` on Linux and `lsof` on macOS.
- 109 focused bundled plugins spanning host inspection and local data, text,
  list, time, and calculator transformations.
- Persistent local navigation preferences, including favorites, recents,
  ordering, group expansion, sidebar visibility, and appearance settings.

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
- [x] Add a persisted enabled/disabled state for installed plugins. Disabled tools
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
- [x] Persist enabled/disabled state in the data directory, while keeping invalid
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

### New plugin-idea intake (2026-08-24)

The following requests are captured for later selection. Unless specifically
noted, each should be a bounded, local-only, one-shot transformation or
inspection; it must not retain uploaded files or entered key material.

| Candidate | Initial scope and constraints |
| --- | --- |
| Certificate parser | Parse a pasted PEM or DER X.509 certificate; show subject, issuer, validity, SANs, serial number, fingerprints, public-key details, extensions, and parse errors. Do not make trust or revocation claims. |
| Public key from private key | Derive and display/export the public key from a locally supplied private key. Treat the input as secret: do not persist, log, or echo it unnecessarily. |
| Public key from certificate | Extract the Subject Public Key Info from a pasted or uploaded certificate, preserving the key algorithm and encoding. |
| CSR parser | Parse a PEM or DER PKCS#10 CSR and show subject, requested SANs/extensions, public-key details, signature algorithm, and validation result. |
| SSH host-key parser | Parse an OpenSSH known_hosts entry or public host-key line, including marker, hosts, key type, fingerprint, and comment; never accept a private key for this tool. |
| Colorblindness simulator | Render bounded image previews for the selected color-vision deficiency model, clearly labeling the model as an approximation. |
| Markdown preview | Render pasted Markdown in a sandboxed, escaped local preview. This complements the existing Markdown-to-HTML download tool. |
| XPath tester | Evaluate a bounded XPath expression against pasted XML and return matching nodes/values and clear syntax errors. Block external entity resolution and network access. |
| HTTP request and response viewer | Send an explicitly configured HTTP request and display a bounded, redacted response. This needs a dedicated outbound-network, SSRF, redirect, header-secret, timeout, and response-size design before implementation. |
| Remove all whitespace | Remove Unicode whitespace from pasted text, with an explicit option for ASCII-only behavior. |
| Add line numbers / Remove line numbers | Add configurable line-number prefixes or remove recognized prefixes from pasted text without modifying the content after the prefix. |
| Kitchen measurement converter | Convert among oz, g, lb, kg, mL, cL, dL, L, tsp, tbsp, cup, pt, qt, and gal. Separate mass from volume and require an explicit ingredient/density when crossing between them. |
| IP-address extractor | Extract and deduplicate IPv4 and IPv6 literals from text, distinguishing valid addresses from loose matches. |
| Email-address extractor | Existing **Email Extractor** covers this request; consider optional deduplication, normalization, and validation modes only if they remain local. |
| MAC-address extractor | Extract and normalize common colon-, hyphen-, and dotted-MAC representations from text. |
| Domain extractor | Extract candidate domain names from text and URLs, with clear handling for public suffixes, IDNs, and non-registrable hostnames. |
| File-path extractor | Extract platform-specific POSIX, Windows-drive, and UNC paths from pasted text; label heuristic matches rather than asserting they exist. |
| Date extractor | Extract candidate dates/times from text and display the assumed locale, timezone, and ambiguous-date interpretations. |
| EXIF metadata extractor | Read bounded image uploads and show available EXIF/XMP/IPTC metadata, including the absence of metadata; no geolocation lookup. |
| Audio metadata extractor | Read bounded audio uploads and show container, codec, duration, bitrate, tags, and embedded artwork presence where supported. |
| Blur image | Apply a configurable local blur to a bounded image and return a browser-local download. |
| Magic-byte file-type detector | Inspect bounded uploaded bytes and report the best supported signature match, MIME type, confidence/ambiguity, and offset. It must not execute, unpack, or otherwise open the file. Support is intended to cover the supplied signature inventory (from `123d` through `zlib`), using a bundled, versioned signature table and fixtures; the source reference is [Wikipedia's file-signature list](https://en.wikipedia.org/wiki/List_of_file_signatures). |
| Regex tester | Existing **Regex Tester** covers this request; retain explicit execution and input limits. |
| Lorem generator | Existing **Lorem Ipsum Generator** covers this request; keep output bounded and local. |

Before selecting file, image, audio, or certificate candidates, confirm the
file-field protocol supports the required MIME types, upload limits, and
browser-local download behavior. Do not expand the protocol implicitly.

### Comparative utility backlog

The following distinct candidates were identified in a review of a similar
utility workbench. Existing Cmdry tools already cover the overlapping items
(for example Bcrypt, UUIDs, Base64 text, URL tooling, JSON/XML/YAML formatting,
QR generation, and text/JSON comparison). The remaining candidates are tracked
at priority 99 in `docs/comp_tools.csv`; they are ideas, not implementation
commitments.

- **Crypto:** tokens, ULIDs, encryption/decryption, BIP39 phrases, HMACs, RSA
  key pairs, password-strength analysis, and PDF signature checks.
- **Converters:** integer bases, Roman numerals, colors, NATO alphabet, ASCII
  binary, YAML-family conversions, list delimiters, and XML-to-JSON.
- **Web/developer:** HTML entities, Basic Auth, Open Graph metadata, OTP,
  MIME lookup, JWT inspection, keycodes, user agents, HTTP status reference,
  Safe Links decoding, SQL formatting, chmod math, Docker Run conversion,
  email normalization, and regular expressions.
- **Network/math/data/text:** IP and MAC utilities, IPv6 ULAs, math/ETA/
  percentage calculators, temperature conversion, phone and IBAN parsing,
  Lorem ipsum, emoji selection, obfuscation, numeronyms, and ASCII art.
- **Browser-capability candidates:** device information, a chronometer,
  benchmark builder, camera recorder, and a WYSIWYG editor. These need a
  separate privacy and protocol design before implementation.

## Candidate selection criteria

Prioritize a candidate when it:

1. Delivers clear value as a narrow read-only view.
2. Has a stable, scriptable host command with parseable output.
3. Can explain missing permissions or unavailable host data accurately.
4. Has a testable parser with fixtures that do not depend on live host state.
5. Does not require Cmdry core changes merely to render its first useful view.

Defer candidates that fundamentally require a broad command runner, arbitrary
file browser, terminal, remote access layer, or a write-capable control panel.
