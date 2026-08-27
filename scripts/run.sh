#!/bin/sh

# Run the locally built Cmdry core with its locally staged bundled plugins.
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
output_dir=${CMDRY_BUILD_DIR:-"$repo_root/dist"}
core_binary="$output_dir/cmdry"

printf 'Repo root is %s\n', "$repo_root"

if [ ! -x "$core_binary" ]; then
	printf 'Cmdry has not been built. Run ./scripts/build.sh first.\n' >&2
	exit 1
fi

if [[ "$1" == "demo" ]]; then
    export CMDRY_DEMO_MODE=true
    shift
fi

export CMDRY_PLUGIN_DIR=${CMDRY_PLUGIN_DIR:-"$output_dir/plugins"}
export CMDRY_DATA_DIR=${CMDRY_DATA_DIR:-"$repo_root/.cmdry-data"}
export CMDRY_ADDR=${CMDRY_ADDR:-"127.0.0.1:8087"}

nohup "$core_binary" serve --plugin-dir="$CMDRY_PLUGIN_DIR" --addr="$CMDRY_ADDR" &
