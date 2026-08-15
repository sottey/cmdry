#!/bin/sh

# Keep the executable presented to Cmdry's plugin scanner as a regular file,
# while running the real helper from its macOS application bundle so TCC can
# read the bundle's Location Services usage description.
set -eu

plugin_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
exec "$plugin_dir/cmdry-wifi.app/Contents/MacOS/cmdry-wifi" "$@"
