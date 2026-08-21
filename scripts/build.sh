#!/bin/sh

# Build the Cmdry core and stage every bundled cmdry-* plugin for discovery.
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
output_dir=${CMDRY_BUILD_DIR:-"$repo_root/dist"}
plugin_dir=${CMDRY_PLUGIN_DIR:-"$output_dir/plugins"}
cache_root=${CMDRY_GO_CACHE_DIR:-"${TMPDIR:-/tmp}/cmdry-go-cache"}
version_file="$repo_root/VERSION"

if [ ! -r "$version_file" ]; then
	printf 'Missing readable VERSION file at %s\n' "$version_file" >&2
	exit 1
fi
version=$(tr -d '\r\n' < "$version_file")
if ! printf '%s' "$version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+([+-][0-9A-Za-z.-]+)?$'; then
	printf 'VERSION must be semantic-style (for example 0.21.0), got %s\n' "$version" >&2
	exit 1
fi
linker_flags="-X github.com/sottey/cmdry/internal/buildinfo.Version=$version -X github.com/sottey/cmdry/plugin-sdk/go.BuildVersion=$version"

printf 'Repo root is %s\n', "$repo_root"


mkdir -p "$output_dir" "$plugin_dir" "$cache_root/cache" "$cache_root/tmp" "$cache_root/mod" "$cache_root/path"

# Use a writable, disposable cache by default. This avoids depending on a
# machine-wide Go configuration while still allowing callers to supply their
# own cache directory through CMDRY_GO_CACHE_DIR.
export GOCACHE="$cache_root/cache"
export GOTMPDIR="$cache_root/tmp"
export GOMODCACHE="$cache_root/mod"
export GOPATH="$cache_root/path"

# The production image is Alpine Linux. Build pure-Go Linux binaries without
# CGO so plugins staged on the host do not require glibc's dynamic loader in
# that musl-based container. Keep CGO available on macOS for cmdry-wifi, which
# links Apple's CoreWLAN framework.
if [ "$(uname -s)" = "Linux" ]; then
	export CGO_ENABLED=0
fi

go build -ldflags "$linker_flags" -o "$output_dir/cmdry" "$repo_root"

find "$repo_root/plugins" -type f -path '*/cmd/cmdry-*/main.go' -print | sort |
while IFS= read -r main_file; do
	package_dir=$(dirname "$main_file")
	binary_name=$(basename "$package_dir")
	if [ "$(uname -s)" = "Darwin" ] && [ "$binary_name" = "cmdry-wifi" ]; then
		app_dir="$plugin_dir/$binary_name.app"
		app_binary="$app_dir/Contents/MacOS/$binary_name"
		mkdir -p "$app_dir/Contents/MacOS"
		go build -ldflags "$linker_flags" -o "$app_binary" "$package_dir"
		cp "$repo_root/scripts/macos/cmdry-wifi-Info.plist" "$app_dir/Contents/Info.plist"
		/usr/libexec/PlistBuddy -c "Set :CFBundleShortVersionString $version" "$app_dir/Contents/Info.plist"
		/usr/libexec/PlistBuddy -c "Set :CFBundleVersion $version" "$app_dir/Contents/Info.plist"
		codesign --force --sign - "$app_dir" >/dev/null
		cp "$repo_root/scripts/macos/cmdry-wifi-launcher.sh" "$plugin_dir/$binary_name"
		chmod 0755 "$plugin_dir/$binary_name"
	else
		go build -ldflags "$linker_flags" -o "$plugin_dir/$binary_name" "$package_dir"
	fi
done

printf 'Built Cmdry core in %s and staged bundled plugins in %s\n' "$output_dir" "$plugin_dir"
