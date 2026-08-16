package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const maxDiagnosticStderr = 4096

type Diagnostic struct {
	Path            string
	Reason          string
	Stderr          string
	StderrTruncated bool
}

type ScanReport struct {
	CompletedAt time.Time
	Diagnostics []Diagnostic
	Failure     string
}

type Discoverer struct {
	Directory string
	Timeout   time.Duration
	Logger    *slog.Logger
}

// Discover scans the plugin directory and atomically replaces the registry only
// after a complete scan succeeds.
func (d Discoverer) Discover(ctx context.Context, registry *Registry) error {
	_, err := d.DiscoverWithDiagnostics(ctx, registry)
	return err
}

// DiscoverWithDiagnostics refreshes the registry and returns rejected
// candidates from that one completed scan. Diagnostics never become registry
// entries and cannot be invoked through the core.
func (d Discoverer) DiscoverWithDiagnostics(ctx context.Context, registry *Registry) (ScanReport, error) {
	fresh, report, err := d.scan(ctx)
	if err != nil {
		return report, err
	}
	registry.Replace(fresh.All())
	return report, nil
}

// Scan discovers valid plugins into an isolated registry.
func (d Discoverer) Scan(ctx context.Context) (*Registry, error) {
	registry, _, err := d.scan(ctx)
	return registry, err
}

func (d Discoverer) scan(ctx context.Context) (*Registry, ScanReport, error) {
	report := ScanReport{CompletedAt: time.Now()}
	entries, err := os.ReadDir(d.Directory)
	if err != nil {
		d.Logger.Warn("plugin directory unavailable", "directory", d.Directory, "error", err)
		report.Failure = boundedText(err.Error(), 512)
		return nil, report, err
	}
	registry := NewRegistry()
	for _, entry := range entries {
		info, e := entry.Info()
		if e != nil || info.IsDir() || info.Mode()&0111 == 0 || !strings.HasPrefix(entry.Name(), "cmdry-") {
			continue
		}
		path := filepath.Join(d.Directory, entry.Name())
		m, stderr, stderrTruncated, err := d.readManifest(ctx, path)
		if err != nil {
			d.Logger.Warn("invalid plugin", "path", path, "error", err)
			report.Diagnostics = append(report.Diagnostics, Diagnostic{Path: path, Reason: boundedText(err.Error(), 512), Stderr: stderr, StderrTruncated: stderrTruncated})
			continue
		}
		if !registry.Add(Registered{Manifest: m, Path: path, Status: StatusEnabled}) {
			d.Logger.Warn("duplicate plugin ID", "id", m.ID, "path", path)
			report.Diagnostics = append(report.Diagnostics, Diagnostic{Path: path, Reason: "duplicate plugin ID: " + m.ID})
		} else {
			d.Logger.Info("plugin registered", "id", m.ID, "path", path)
		}
	}
	return registry, report, nil
}
func (d Discoverer) readManifest(parent context.Context, path string) (Manifest, string, bool, error) {
	ctx, cancel := context.WithTimeout(parent, d.Timeout)
	defer cancel()
	stdout := &limitedBuffer{limit: maxDiagnosticStderr}
	stderr := &limitedBuffer{limit: maxDiagnosticStderr}
	command := exec.CommandContext(ctx, path, "manifest")
	command.Stdout = stdout
	command.Stderr = stderr
	err := command.Run()
	if ctx.Err() != nil {
		return Manifest{}, stderr.String(), stderr.truncated, ctx.Err()
	}
	if err != nil {
		return Manifest{}, stderr.String(), stderr.truncated, fmt.Errorf("run manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(stdout.Bytes(), &m); err != nil {
		return Manifest{}, stderr.String(), stderr.truncated, fmt.Errorf("decode manifest: %w", err)
	}
	if err := ValidateManifest(m); err != nil {
		return Manifest{}, stderr.String(), stderr.truncated, fmt.Errorf("validate manifest: %w", err)
	}
	return m, stderr.String(), stderr.truncated, nil
}

type limitedBuffer struct {
	bytes.Buffer
	limit     int
	truncated bool
}

func (b *limitedBuffer) Write(contents []byte) (int, error) {
	remaining := b.limit - b.Len()
	if remaining <= 0 {
		b.truncated = true
		return len(contents), nil
	}
	if len(contents) > remaining {
		_, _ = b.Buffer.Write(contents[:remaining])
		b.truncated = true
		return len(contents), nil
	}
	return b.Buffer.Write(contents)
}

func boundedText(value string, limit int) string {
	value = strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, value)
	if len(value) > limit {
		return value[:limit] + "…"
	}
	return value
}
