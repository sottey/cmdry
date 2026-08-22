package plugins

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDiscoverWithDiagnosticsCapturesRejectedManifest(t *testing.T) {
	dir := t.TempDir()
	writeExecutable(t, filepath.Join(dir, "cmdry-good"), "#!/bin/sh\nprintf '%s' '{\"protocol_version\":1,\"id\":\"com.example.good\",\"name\":\"Good\",\"version\":\"1.0.0\",\"actions\":[{\"id\":\"list\",\"name\":\"List\",\"method\":\"read\"}],\"pages\":[{\"id\":\"main\",\"name\":\"Main\",\"default\":true}]}'\n")
	writeExecutable(t, filepath.Join(dir, "cmdry-bad"), "#!/bin/sh\nprintf 'manifest dependency missing' >&2\nprintf '%s' '{not json}'\n")

	registry := NewRegistry()
	discoverer := Discoverer{Directory: dir, Timeout: time.Second, Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	report, err := discoverer.DiscoverWithDiagnostics(context.Background(), registry)
	if err != nil {
		t.Fatal(err)
	}
	if registry.Len() != 1 {
		t.Fatalf("registered plugins = %d, want 1", registry.Len())
	}
	if len(report.Diagnostics) != 1 {
		t.Fatalf("diagnostics = %#v", report.Diagnostics)
	}
	diagnostic := report.Diagnostics[0]
	if !strings.Contains(diagnostic.Reason, "decode manifest") || diagnostic.Stderr != "manifest dependency missing" {
		t.Fatalf("diagnostic = %#v", diagnostic)
	}
	if diagnostic.Path != filepath.Join(dir, "cmdry-bad") || report.CompletedAt.IsZero() {
		t.Fatalf("diagnostic path/time = %#v, %v", diagnostic, report.CompletedAt)
	}
}

func TestLimitedBufferBoundsStderr(t *testing.T) {
	buffer := &limitedBuffer{limit: 4}
	if _, err := buffer.Write([]byte("abcdef")); err != nil {
		t.Fatal(err)
	}
	if got := buffer.String(); got != "abcd" || !buffer.truncated {
		t.Fatalf("got %q truncated=%v", got, buffer.truncated)
	}
}

func TestDiscoverWithFilterSkipsExcludedPlugin(t *testing.T) {
	dir := t.TempDir()
	writeExecutable(t, filepath.Join(dir, "cmdry-transform"), "#!/bin/sh\nprintf '%s' '{\"protocol_version\":1,\"id\":\"com.example.transform\",\"name\":\"Transform\",\"version\":\"1.0.0\",\"category\":\"text\",\"permissions\":[\"data.transform\"],\"actions\":[{\"id\":\"list\",\"name\":\"List\",\"method\":\"read\"}],\"pages\":[{\"id\":\"main\",\"name\":\"Main\",\"default\":true}]}'\n")
	writeExecutable(t, filepath.Join(dir, "cmdry-server"), "#!/bin/sh\nprintf '%s' '{\"protocol_version\":1,\"id\":\"com.example.server\",\"name\":\"Server\",\"version\":\"1.0.0\",\"category\":\"server\",\"permissions\":[\"system.read\"],\"actions\":[{\"id\":\"list\",\"name\":\"List\",\"method\":\"read\"}],\"pages\":[{\"id\":\"main\",\"name\":\"Main\",\"default\":true}]}'\n")

	registry := NewRegistry()
	discoverer := Discoverer{Directory: dir, Timeout: time.Second, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), Filter: func(manifest Manifest) bool {
		return manifest.Category != "server" && len(manifest.Permissions) == 1 && manifest.Permissions[0] == "data.transform"
	}}
	if _, err := discoverer.DiscoverWithDiagnostics(context.Background(), registry); err != nil {
		t.Fatal(err)
	}
	if registry.Len() != 1 {
		t.Fatalf("registered plugins = %d, want 1", registry.Len())
	}
	if _, ok := registry.Get("com.example.transform"); !ok {
		t.Fatal("transform plugin was not registered")
	}
	if _, ok := registry.Get("com.example.server"); ok {
		t.Fatal("server plugin was registered")
	}
}

func writeExecutable(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0700); err != nil {
		t.Fatal(err)
	}
}
