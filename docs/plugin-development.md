# Cmdry plugin development guide

Cmdry plugins are standalone executables. A plugin describes itself with a JSON
manifest and returns JSON views for actions. It can be written in Go or any
language that can read standard input, write standard output, and run as an
executable on the same operating system and CPU architecture as Cmdry.

This guide describes the current version-1 protocol. Cmdry loads plugins only
at startup, so restart Cmdry after installing, replacing, or removing one.

## Before you begin

Plugins are executable code running with the Cmdry process's operating-system
permissions. Only install plugins you trust. Manifest permissions are displayed
in the UI, but they are descriptive in v1; they do not sandbox a plugin.

The core has no command shell and does not pass web form data to plugins in v1.
Each page load and button action currently receives an empty parameter object.
Build read-only, fixed-purpose actions unless you control and validate every
input within your own plugin.

## Plugin lifecycle

At startup, Cmdry scans `CMDRY_PLUGIN_DIR` (default
`/opt/cmdry/plugins`) for files that:

1. Are not directories.
2. Have at least one executable permission bit.
3. Have names beginning with `cmdry-`.

For every candidate, Cmdry executes:

```text
/absolute/path/to/cmdry-example manifest
```

The program must write exactly one valid manifest JSON document to standard
output. Diagnostics belong on standard error. An invalid plugin is logged and
skipped without stopping Cmdry.

When a user opens a plugin page or submits an action button, Cmdry executes:

```text
/absolute/path/to/cmdry-example execute list
```

It sends this JSON on standard input:

```json
{"action":"list","params":{}}
```

The action response must be protocol JSON on standard output. Cmdry applies an
eight-second timeout and renders plugin failures as a contained page error.

## Manifest

A v1 manifest has this shape:

```json
{
  "protocol_version": 1,
  "id": "example",
  "name": "Example Inspector",
  "version": "0.1.0",
  "description": "Shows a small, read-only example.",
  "category": "system",
  "icon": "terminal",
  "pages": [
    {"id": "overview", "name": "Overview", "default": true, "action": "list"}
  ],
  "permissions": ["system.read"],
  "actions": [
    {"id": "list", "name": "Refresh", "method": "read"}
  ]
}
```

Required fields and rules:

| Field | Rule |
| --- | --- |
| `protocol_version` | Must be the number `1`. |
| `id` | Unique among installed plugins; lowercase letters/digits followed by lowercase letters, digits, `_`, or `-`. |
| `name` | Non-empty display name. |
| `version` | Semantic-style version such as `0.1.0` or `v1.2.3-beta.1`. |
| `pages` | At least one page; each needs a unique safe `id` and non-empty `name`. One page should be `default: true`. |
| `actions` | At least one action; each needs a unique safe `id` and non-empty `name`. |

`description`, `category`, `icon`, and `permissions` are optional to the
validator but should be supplied for a useful plugin list. Set a page's
`action` to one of the manifest action IDs. If omitted, Cmdry uses the page ID
as the action ID; if no default page is selected, Cmdry falls back to the first
declared action.

`method` is informational in v1. Use `read` for non-mutating actions and name
any write-capable action honestly. The core does not enforce permissions or
method values.

## Responses and UI components

A successful action response has `ok: true` and a `data` view:

```json
{
  "ok": true,
  "data": {
    "title": "Example status",
    "components": [
      {"type": "metric", "label": "Items", "value": "3"},
      {"type": "text", "title": "About", "text": "Everything looks healthy."},
      {
        "type": "table",
        "id": "items",
        "columns": [
          {"key": "name", "label": "Name"},
          {"key": "state", "label": "State"}
        ],
        "rows": [
          {"name": "alpha", "state": "ready"}
        ]
      }
    ]
  }
}
```

Cmdry supports only these component types:

| Type | Fields Cmdry renders |
| --- | --- |
| `metric` | `label`, `value`, optional `description` |
| `text` | `title`, `text` |
| `alert` | `title`, `message`, optional `level` |
| `table` | `id`, `columns` (`key`, `label`), and `rows` |
| `actions` | `actions`, each with `id`, `name`, `method` |

For a table, every column key should have a matching value in each row. Values
are rendered as text. Cmdry uses Go's `html/template`, so plugin values are
HTML-escaped. There is no plugin-provided HTML, CSS, JavaScript, URL, or shell
command capability.

An `actions` component renders POST buttons. The action ID must appear in the
manifest or Cmdry refuses to execute it. These actions do not receive a form
payload in v1.

To return an expected operational error, use this shape:

```json
{
  "ok": false,
  "error": {
    "code": "COMMAND_MISSING",
    "message": "examplectl is not installed; install the example-tools package and try again."
  }
}
```

Both `code` and `message` are required. Keep messages useful to the local
administrator and never include credentials, tokens, or stack traces.

## Go SDK example

The repository includes a Go helper at `plugin-sdk/go`. The following complete
example may live in a separate Go module, provided it imports a compatible
Cmdry SDK version:

```go
package main

import cmdry "github.com/sottey/cmdry/plugin-sdk/go"

func main() {
	cmdry.Run(cmdry.Plugin{
		Manifest: cmdry.Manifest{
			ProtocolVersion: 1,
			ID:              "example",
			Name:            "Example Inspector",
			Version:         "0.1.0",
			Description:     "A read-only example plugin.",
			Category:        "system",
			Pages: []cmdry.Page{{
				ID: "overview", Name: "Overview", Default: true, Action: "list",
			}},
			Permissions: []string{"system.read"},
			Actions:     []cmdry.Action{{ID: "list", Name: "Refresh", Method: "read"}},
		},
		Actions: map[string]cmdry.Handler{"list": list},
	})
}

func list(_ cmdry.Request) (cmdry.View, error) {
	return cmdry.View{
		Title: "Example status",
		Components: []cmdry.Component{
			{Type: "metric", Label: "Items", Value: "3"},
			{Type: "text", Title: "About", Text: "Everything looks healthy."},
		},
	}, nil
}
```

The SDK handles the `manifest` and `execute` command-line contract, JSON input,
JSON output, and structured error responses. A handler error becomes a
`COMMAND_FAILED` response.

Build and install it on the same platform as Cmdry:

```bash
go build -o cmdry-example .
install -Dm755 cmdry-example /opt/cmdry/plugins/cmdry-example
```

For local development without installing to `/opt`, use a project-local plugin
directory:

```bash
mkdir -p .cmdry-data .cmdry-plugins
go build -o .cmdry-plugins/cmdry-example .
CMDRY_PLUGIN_DIR="$PWD/.cmdry-plugins" CMDRY_DATA_DIR="$PWD/.cmdry-data" cmdry serve
```

On macOS, build the plugin natively and run Cmdry natively. A Linux Docker
container cannot execute a macOS binary or inspect the Mac host's ports.

## Language-neutral implementation

You do not have to use Go. Your executable needs only this behavior:

1. With `manifest`, write the manifest JSON to standard output and exit zero.
2. With `execute <action-id>`, decode the request JSON from standard input.
3. Confirm that `<action-id>` is an action your plugin supports.
4. Write either a successful response or a structured error response to
   standard output.
5. Write logs and diagnostics only to standard error.

Do not print banners, debug messages, or progress lines to standard output;
they make the JSON invalid. Do not rely on Cmdry's current working directory,
and use explicit paths or trusted configuration for any files your plugin
needs.

## Test and troubleshoot

Test the executable before installing it:

```bash
./cmdry-example manifest
printf '{"action":"list","params":{}}' | ./cmdry-example execute list
```

Then restart Cmdry and inspect its stderr logs. Common causes of a plugin not
appearing are:

| Symptom | Check |
| --- | --- |
| Plugin does not appear | Filename starts with `cmdry-`, file is executable, and it is in `CMDRY_PLUGIN_DIR`. |
| `invalid plugin` log | Run `./cmdry-example manifest` directly; it must emit only valid v1 manifest JSON. |
| Plugin page shows an error | Run the `execute` command above; verify dependent commands and permissions. |
| Plugin times out | Make the action complete within eight seconds or return a faster, bounded result. |
| Works on one machine only | Rebuild for the target OS and CPU architecture; do not copy a macOS executable to Linux or vice versa. |

Test both success and failure paths. In particular, verify malformed command
output, missing dependencies, permission-limited results, and unknown action
handling. Keep all command invocation inside the plugin explicit—never build a
shell command from untrusted input.
