// Package buildinfo exposes metadata injected by the release build.
package buildinfo

// Version is replaced by scripts/build.sh and the Docker build. The fallback
// identifies binaries built directly with go build or go run.
var Version = "0.0.0-dev"
