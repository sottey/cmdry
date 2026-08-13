package server

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/sottey/cmdry/internal/config"
	"github.com/sottey/cmdry/internal/plugins"
)

//go:embed templates/*.html static/*
var assets embed.FS

type App struct {
	cfg       config.Config
	registry  *plugins.Registry
	runner    plugins.Runner
	logger    *slog.Logger
	templates *template.Template
}
type pageData struct {
	Title, Section string
	Plugins        []plugins.Registered
	Current        *plugins.Registered
	View           *plugins.View
	Error          string
}

func New(cfg config.Config, registry *plugins.Registry, logger *slog.Logger) (*App, error) {
	t, err := template.New("layout.html").Funcs(template.FuncMap{"row": func(row map[string]any, key string) any { return row[key] }, "category": func(s string) string {
		if s == "" {
			return "other"
		}
		return s
	}}).ParseFS(assets, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &App{cfg: cfg, registry: registry, runner: plugins.Runner{Timeout: cfg.PluginTimeout}, logger: logger, templates: t}, nil
}
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":true}`)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/static/") {
		http.FileServer(http.FS(assets)).ServeHTTP(w, r)
		return
	}
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/":
		a.dashboard(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/plugins":
		a.pluginList(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/settings":
		a.settings(w, r)
	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/actions/"):
		a.pluginAction(w, r)
	case strings.HasPrefix(r.URL.Path, "/plugins/"):
		a.pluginPage(w, r)
	default:
		http.NotFound(w, r)
	}
}
func (a *App) render(w http.ResponseWriter, status int, name string, d pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := a.templates.ExecuteTemplate(w, "layout.html", d); err != nil {
		a.logger.Error("render page", "error", err)
	}
}
func (a *App) base(title, section string) pageData {
	return pageData{Title: title, Section: section, Plugins: a.registry.All()}
}
func (a *App) dashboard(w http.ResponseWriter, r *http.Request) {
	d := a.base("Cmdry", "Overview")
	a.render(w, http.StatusOK, "dashboard.html", d)
}
func (a *App) pluginList(w http.ResponseWriter, r *http.Request) {
	d := a.base("Installed plugins", "Plugins")
	a.render(w, http.StatusOK, "plugins.html", d)
}
func (a *App) settings(w http.ResponseWriter, r *http.Request) {
	d := a.base("Settings", "Settings")
	a.render(w, http.StatusOK, "settings.html", d)
}
func (a *App) pluginPage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(path.Clean(r.URL.Path), "/plugins/"), "/")
	if len(parts) == 0 || !plugins.ValidID(parts[0]) {
		http.NotFound(w, r)
		return
	}
	entry, ok := a.registry.Get(parts[0])
	if !ok {
		http.NotFound(w, r)
		return
	}
	d := a.base(entry.Manifest.Name, entry.Manifest.Name)
	d.Current = &entry
	action := ""
	if len(parts) > 1 {
		for _, p := range entry.Manifest.Pages {
			if p.ID == parts[1] {
				action = p.Action
				if action == "" {
					action = p.ID
				}
			}
		}
	}
	if action == "" {
		for _, p := range entry.Manifest.Pages {
			if p.Default {
				action = p.Action
				if action == "" {
					action = p.ID
				}
				break
			}
		}
	}
	if action == "" {
		action = entry.Manifest.Actions[0].ID
	}
	response, err := a.runner.Run(r.Context(), entry, action, map[string]any{})
	if err != nil {
		d.Error = err.Error()
		a.logger.Warn("plugin invocation failed", "plugin", entry.Manifest.ID, "error", err)
	} else if !response.OK {
		d.Error = response.Error.Message
	} else {
		d.View = response.Data
	}
	a.render(w, http.StatusOK, "plugin.html", d)
}

func (a *App) pluginAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(path.Clean(r.URL.Path), "/plugins/"), "/")
	if len(parts) != 3 || parts[1] != "actions" || !plugins.ValidID(parts[0]) || !plugins.ValidID(parts[2]) {
		http.NotFound(w, r)
		return
	}
	entry, ok := a.registry.Get(parts[0])
	if !ok {
		http.NotFound(w, r)
		return
	}
	d := a.base(entry.Manifest.Name, entry.Manifest.Name)
	d.Current = &entry
	response, err := a.runner.Run(r.Context(), entry, parts[2], map[string]any{})
	if err != nil {
		d.Error = err.Error()
	} else if !response.OK {
		d.Error = response.Error.Message
	} else {
		d.View = response.Data
	}
	a.render(w, http.StatusOK, "plugin.html", d)
}
