package plugins

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Discoverer struct {
	Directory string
	Timeout   time.Duration
	Logger    *slog.Logger
}

func (d Discoverer) Discover(ctx context.Context, registry *Registry) {
	entries, err := os.ReadDir(d.Directory)
	if err != nil {
		d.Logger.Warn("plugin directory unavailable", "directory", d.Directory, "error", err)
		return
	}
	for _, entry := range entries {
		info, e := entry.Info()
		if e != nil || info.IsDir() || info.Mode()&0111 == 0 || !strings.HasPrefix(entry.Name(), "cmdry-") {
			continue
		}
		path := filepath.Join(d.Directory, entry.Name())
		m, err := d.readManifest(ctx, path)
		if err != nil {
			d.Logger.Warn("invalid plugin", "path", path, "error", err)
			continue
		}
		if !registry.Add(Registered{Manifest: m, Path: path, Status: StatusEnabled}) {
			d.Logger.Warn("duplicate plugin ID", "id", m.ID, "path", path)
		} else {
			d.Logger.Info("plugin registered", "id", m.ID, "path", path)
		}
	}
}
func (d Discoverer) readManifest(parent context.Context, path string) (Manifest, error) {
	ctx, cancel := context.WithTimeout(parent, d.Timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, path, "manifest").Output()
	if ctx.Err() != nil {
		return Manifest{}, ctx.Err()
	}
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := json.Unmarshal(out, &m); err != nil {
		return Manifest{}, err
	}
	return m, ValidateManifest(m)
}
