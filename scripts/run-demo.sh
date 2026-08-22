#!/bin/sh

# Run the locally built Cmdry public-demo profile. It uses the normal staged
# plugin directory but registers only transform-only plugins and writes no
# workspace state.
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
output_dir=${CMDRY_BUILD_DIR:-"$repo_root/dist"}
core_binary="$output_dir/cmdry"

if [ ! -x "$core_binary" ]; then
	printf 'Cmdry has not been built. Run ./scripts/build.sh first.\n' >&2
	exit 1
fi

export CMDRY_DEMO_MODE=true
export CMDRY_PLUGIN_DIR=${CMDRY_PLUGIN_DIR:-"$output_dir/plugins"}
export CMDRY_DATA_DIR=${CMDRY_DATA_DIR:-"$repo_root/.cmdry-demo-data"}
export CMDRY_ADDR=${CMDRY_ADDR:-"127.0.0.1:8088"}

exec "$core_binary" serve "$@"
