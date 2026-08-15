#!/bin/sh

# Build the Cmdry core and stage every bundled cmdry-* plugin for discovery.
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
output_dir=${CMDRY_BUILD_DIR:-"$repo_root/dist"}
plugin_dir=${CMDRY_PLUGIN_DIR:-"$output_dir/plugins"}
cache_root=${CMDRY_GO_CACHE_DIR:-"${TMPDIR:-/tmp}/cmdry-go-cache"}

mkdir -p "$output_dir" "$plugin_dir" "$cache_root/cache" "$cache_root/tmp" "$cache_root/mod" "$cache_root/path"

# Use a writable, disposable cache by default. This avoids depending on a
# machine-wide Go configuration while still allowing callers to supply their
# own cache directory through CMDRY_GO_CACHE_DIR.
export GOCACHE="$cache_root/cache"
export GOTMPDIR="$cache_root/tmp"
export GOMODCACHE="$cache_root/mod"
export GOPATH="$cache_root/path"

go build -o "$output_dir/cmdry" "$repo_root"

find "$repo_root/plugins" -type f -path '*/cmd/cmdry-*/main.go' -print | sort |
while IFS= read -r main_file; do
	package_dir=$(dirname "$main_file")
	binary_name=$(basename "$package_dir")
	if [ "$(uname -s)" = "Darwin" ] && [ "$binary_name" = "cmdry-wifi" ]; then
		app_dir="$plugin_dir/$binary_name.app"
		app_binary="$app_dir/Contents/MacOS/$binary_name"
		mkdir -p "$app_dir/Contents/MacOS"
		go build -o "$app_binary" "$package_dir"
		cp "$repo_root/scripts/macos/cmdry-wifi-Info.plist" "$app_dir/Contents/Info.plist"
		codesign --force --sign - "$app_dir" >/dev/null
		cp "$repo_root/scripts/macos/cmdry-wifi-launcher.sh" "$plugin_dir/$binary_name"
		chmod 0755 "$plugin_dir/$binary_name"
	else
		go build -o "$plugin_dir/$binary_name" "$package_dir"
	fi
done

printf 'Built Cmdry core in %s and staged bundled plugins in %s\n' "$output_dir" "$plugin_dir"
