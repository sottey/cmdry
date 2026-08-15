package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/sottey/cmdry/internal/config"
	"github.com/sottey/cmdry/internal/groupstate"
	"github.com/sottey/cmdry/internal/pluginorder"
	"github.com/sottey/cmdry/internal/plugins"
)

//go:embed templates/*.html static/*
var assets embed.FS

type App struct {
	cfg        config.Config
	registry   *plugins.Registry
	discovery  plugins.Discoverer
	order      pluginorder.Store
	orderIDs   []string
	orderMu    sync.RWMutex
	groups     groupstate.Store
	groupState map[string]bool
	groupMu    sync.RWMutex
	runner     plugins.Runner
	logger     *slog.Logger
	templates  *template.Template
}
type pageData struct {
	Title, Section string
	Plugins        []plugins.Registered
	PluginGroups   []PluginGroup
	Current        *plugins.Registered
	View           *plugins.View
	Error          string
	Message        string
}
type PluginGroup struct {
	ID, Name  string
	Collapsed bool
	Plugins   []plugins.Registered
}

func New(cfg config.Config, registry *plugins.Registry, logger *slog.Logger) (*App, error) {
	t, err := template.New("layout.html").Funcs(template.FuncMap{"row": func(row map[string]any, key string) any { return row[key] }, "pluginSearchJSON": func(entries []plugins.Registered) template.JS {
		type result struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Category    string   `json:"category"`
			Href        string   `json:"href"`
			Terms       []string `json:"terms"`
		}
		items := make([]result, 0, len(entries))
		for _, entry := range entries {
			items = append(items, result{ID: entry.Manifest.ID, Name: entry.Manifest.Name, Description: entry.Manifest.Description, Category: entry.Manifest.Category, Href: "/plugins/" + entry.Manifest.ID, Terms: append([]string{}, entry.Manifest.SearchTerms...)})
		}
		encoded, _ := json.Marshal(items)
		return template.JS(encoded)
	}, "downloadURL": func(mimeType, content string) template.URL {
		return template.URL("data:" + mimeType + ";base64," + content)
	}, "category": func(s string) string {
		if s == "" {
			return "other"
		}
		return s
	}}).ParseFS(assets, "templates/*.html")
	if err != nil {
		return nil, err
	}
	order := pluginorder.Store{DataDir: cfg.DataDir}
	orderIDs, err := order.Load()
	if err != nil {
		return nil, err
	}
	groups := groupstate.Store{DataDir: cfg.DataDir}
	groupState, err := groups.Load()
	if err != nil {
		return nil, err
	}
	return &App{cfg: cfg, registry: registry, discovery: plugins.Discoverer{Directory: cfg.PluginDir, Timeout: cfg.PluginTimeout, Logger: logger}, order: order, orderIDs: orderIDs, groups: groups, groupState: groupState, runner: plugins.Runner{Timeout: cfg.PluginTimeout}, logger: logger, templates: t}, nil
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
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/refresh":
		a.refreshPlugins(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/order":
		a.savePluginOrder(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/groups":
		a.saveGroupState(w, r)
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
	a.orderMu.RLock()
	ids := append([]string(nil), a.orderIDs...)
	a.orderMu.RUnlock()
	entries := a.registry.Ordered(ids)
	a.groupMu.RLock()
	state := make(map[string]bool, len(a.groupState))
	for key, value := range a.groupState {
		state[key] = value
	}
	a.groupMu.RUnlock()
	return pageData{Title: title, Section: section, Plugins: entries, PluginGroups: makeGroups(entries, state)}
}

func makeGroups(entries []plugins.Registered, state map[string]bool) []PluginGroup {
	byID := map[string]*PluginGroup{}
	order := []string{}
	for _, entry := range entries {
		id := strings.TrimSpace(entry.Manifest.Category)
		if id == "" {
			id = "other"
		}
		group := byID[id]
		if group == nil {
			group = &PluginGroup{ID: id, Name: strings.ToUpper(id[:1]) + id[1:], Collapsed: state[id]}
			byID[id] = group
			order = append(order, id)
		}
		group.Plugins = append(group.Plugins, entry)
	}
	groups := make([]PluginGroup, 0, len(order))
	for _, id := range order {
		groups = append(groups, *byID[id])
	}
	return groups
}
func (a *App) dashboard(w http.ResponseWriter, r *http.Request) {
	d := a.base("Cmdry", "Overview")
	a.render(w, http.StatusOK, "dashboard.html", d)
}
func (a *App) pluginList(w http.ResponseWriter, r *http.Request) {
	d := a.base("Installed plugins", "Plugins")
	if r.URL.Query().Get("refreshed") == "1" {
		d.Message = "Plugin scan completed."
	} else if r.URL.Query().Get("refresh") == "failed" {
		d.Error = "Plugin scan failed. Existing registered plugins were kept. Check the server logs and plugin directory."
	}
	a.render(w, http.StatusOK, "plugins.html", d)
}
func (a *App) refreshPlugins(w http.ResponseWriter, r *http.Request) {
	if err := a.discovery.Discover(r.Context(), a.registry); err != nil {
		a.logger.Warn("plugin refresh failed", "error", err)
		http.Redirect(w, r, "/plugins?refresh=failed", http.StatusSeeOther)
		return
	}
	a.logger.Info("plugins refreshed", "plugins", a.registry.Len())
	http.Redirect(w, r, "/plugins?refreshed=1", http.StatusSeeOther)
}
func (a *App) savePluginOrder(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	defer r.Body.Close()
	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "invalid plugin order", http.StatusBadRequest)
		return
	}
	registered := a.registry.All()
	if len(ids) != len(registered) {
		http.Error(w, "plugin order must include every registered plugin", http.StatusBadRequest)
		return
	}
	known := make(map[string]bool, len(registered))
	for _, entry := range registered {
		known[entry.Manifest.ID] = true
	}
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		if !plugins.ValidID(id) || !known[id] || seen[id] {
			http.Error(w, "plugin order contains an unknown or duplicate plugin", http.StatusBadRequest)
			return
		}
		seen[id] = true
	}
	if err := a.order.Save(ids); err != nil {
		a.logger.Error("save plugin order", "error", err)
		http.Error(w, "unable to save plugin order", http.StatusInternalServerError)
		return
	}
	a.orderMu.Lock()
	a.orderIDs = append([]string(nil), ids...)
	a.orderMu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func (a *App) saveGroupState(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	defer r.Body.Close()
	var state map[string]bool
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, "invalid plugin group state", http.StatusBadRequest)
		return
	}
	known := map[string]bool{}
	for _, entry := range a.registry.All() {
		id := strings.TrimSpace(entry.Manifest.Category)
		if id == "" {
			id = "other"
		}
		known[id] = true
	}
	for id := range state {
		if !plugins.ValidID(id) || !known[id] {
			http.Error(w, "unknown plugin group", http.StatusBadRequest)
			return
		}
	}
	if err := a.groups.Save(state); err != nil {
		a.logger.Error("save plugin group state", "error", err)
		http.Error(w, "unable to save plugin group state", http.StatusInternalServerError)
		return
	}
	a.groupMu.Lock()
	a.groupState = state
	a.groupMu.Unlock()
	w.WriteHeader(http.StatusNoContent)
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
	r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid or oversized action input", http.StatusBadRequest)
		return
	}
	params := make(map[string]any, len(r.PostForm))
	for key, values := range r.PostForm {
		if !plugins.ValidID(key) || len(values) != 1 {
			http.Error(w, "invalid action input", http.StatusBadRequest)
			return
		}
		params[key] = values[0]
	}
	d := a.base(entry.Manifest.Name, entry.Manifest.Name)
	d.Current = &entry
	response, err := a.runner.Run(r.Context(), entry, parts[2], params)
	if err != nil {
		d.Error = err.Error()
	} else if !response.OK {
		d.Error = response.Error.Message
	} else {
		d.View = response.Data
	}
	a.render(w, http.StatusOK, "plugin.html", d)
}
