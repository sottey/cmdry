# Cmdry plugin development guide

Cmdry plugins are standalone executables. A plugin describes itself with a JSON
manifest and returns JSON views for actions. It can be written in Go or any
language that can read standard input, write standard output, and run as an
executable on the same operating system and CPU architecture as Cmdry.

This guide describes the current version-1 protocol. Cmdry loads plugins at
startup and can rescan while running from the **Refresh plugins** button on the
Plugins page. Use that button after installing, replacing, or removing a
plugin; a failed scan preserves the currently registered plugin set.

Administrators can drag installed plugins in the sidebar to set a persistent
display order. Plugins omitted from a saved order, including newly installed
ones, follow those ordered entries alphabetically.

## Before you begin

Plugins are executable code running with the Cmdry process's operating-system
permissions. Only install plugins you trust. Manifest permissions are displayed
in the UI, but they are descriptive in v1; they do not sandbox a plugin.

The core has no command shell. Page loads and ordinary action buttons receive
an empty parameter object. A plugin may return a `form` component to collect
bounded text input; the submitted field values are passed to its declared
action in `params`. Validate every value in the plugin—manifest permissions are
not a sandbox.

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

Cmdry also records the latest discovery scan on its **Diagnostics** page. It
lists rejected executable candidates, the scan time, and a stderr excerpt
capped at 4 KiB. A rejected candidate is never registered and the diagnostics
page has no action that can execute it; correct the binary, then use **Refresh
plugins** on the Plugins page to scan again.

When a user opens a plugin page or submits an action button, Cmdry executes:

```text
/absolute/path/to/cmdry-example execute list
```

It sends this JSON on standard input (a page load or ordinary action has an
empty `params` object):

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
| `id` | Unique among installed plugins; lowercase letters/digits followed by lowercase letters, digits, `_`, `-`, or `.`. Reverse-domain IDs such as `com.sottey.example` are supported. |
| `name` | Non-empty display name. |
| `version` | Semantic-style version such as `0.1.0` or `v1.2.3-beta.1`. |
| `pages` | At least one page; each needs a unique safe `id` and non-empty `name`. One page should be `default: true`. |
| `actions` | At least one action; each needs a unique safe `id` and non-empty `name`. |
| `search_terms` | Optional list of up to 32 short aliases used by Overview plugin search. |

`description`, `category`, `icon`, and `permissions` are optional to the
validator but should be supplied for a useful plugin list. Cmdry groups the
sidebar by `category`, so use a stable, short category such as `text`, `data`,
or `developer` rather than creating a near-duplicate category per plugin. Set a page's
`action` to one of the manifest action IDs. If omitted, Cmdry uses the page ID
as the action ID; if no default page is selected, Cmdry falls back to the first
declared action.

Overview search always matches a plugin's name, ID, category, and description.
Use `search_terms` for aliases people are likely to type, such as `uuid` for a
GUID generator or `checksum` for a hash generator.

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
| `code` | `title`, `text`, rendered in a copyable preformatted block |
| `alert` | `title`, `message`, optional `level` |
| `table` | `id`, `columns` (`key`, `label`), and `rows` |
| `actions` | `actions`, each with `id`, `name`, `method` |
| `form` | `title`, `action`, `submit`, and one or more `fields` |
| `download` | `filename`, `mime_type`, and base64-encoded `content` |

For a table, every column key should have a matching value in each row. Values
are rendered as text. Cmdry uses Go's `html/template`, so plugin values are
HTML-escaped. There is no plugin-provided HTML, CSS, JavaScript, URL, or shell
command capability.

An `actions` component renders POST buttons. The action ID must appear in the
manifest or Cmdry refuses to execute it.

A `form` component renders a POST form. Each field needs a safe `name`,
non-empty `label`, and `type` of `text`, `password`, `textarea`, `number`, `checkbox`, or
`select`; a select field supplies `options` with `value` and `label`. A `file`
field may set `accept` to guide browser selection, and may set `multiple: true`
when it intentionally accepts up to four files. Cmdry accepts one uploaded file
per ordinary file field, or up to four for an explicit multiple field; it holds
them only in request memory, removes multipart temporary data before returning,
and never provides a host path to the plugin. Each upload is limited to 4 MiB;
the complete multipart action body is limited to 18 MiB. Use
Use
`value`, `min`, `max`, and `required: true` where appropriate. Its `action`
must be declared in the manifest. Cmdry limits
the full submitted form body to 6 MiB. Ordinary fields arrive as strings; a
file arrives as `{name, mime_type, content}` where content is standard base64;
a multiple field arrives as an array of those objects.
The Go SDK's `request.File("image")` helper validates and decodes it:

```json
{"action":"resize","params":{"image":{"name":"photo.png","mime_type":"image/png","content":"..."}}}
```

A `download` component provides a browser-local download. Its content must be
standard base64 and is limited to 8 MiB after encoding. Use it for generated,
ephemeral output such as CSV—not for paths on the host filesystem.

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
