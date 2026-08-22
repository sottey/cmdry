package server

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sottey/cmdry/internal/buildinfo"
	"github.com/sottey/cmdry/internal/config"
	"github.com/sottey/cmdry/internal/groupstate"
	"github.com/sottey/cmdry/internal/pluginnav"
	"github.com/sottey/cmdry/internal/pluginorder"
	"github.com/sottey/cmdry/internal/plugins"
	"github.com/sottey/cmdry/internal/pluginstate"
	"github.com/sottey/cmdry/internal/workspace"
)

//go:embed templates/*.html static/*
var assets embed.FS

type App struct {
	cfg            config.Config
	registry       *plugins.Registry
	discovery      plugins.Discoverer
	order          pluginorder.Store
	orderIDs       []string
	orderMu        sync.RWMutex
	groups         groupstate.Store
	groupState     map[string]bool
	groupMu        sync.RWMutex
	navigation     pluginnav.Store
	navState       pluginnav.State
	navMu          sync.RWMutex
	pluginState    pluginstate.Store
	disabledIDs    []string
	pluginStateMu  sync.RWMutex
	workspace      workspace.Store
	workspaceState workspace.State
	workspaceMu    sync.RWMutex
	diagnostics    plugins.ScanReport
	diagnosticsMu  sync.RWMutex
	runner         plugins.Runner
	logger         *slog.Logger
	templates      *template.Template
}
type pageData struct {
	Title, Section     string
	AppVersion         string
	Plugins            []plugins.Registered
	EnabledPlugins     []plugins.Registered
	PluginGroups       []PluginGroup
	Favorites          []plugins.Registered
	FavoritesCollapsed bool
	Recents            []plugins.Registered
	Hidden             []string
	AllVisible         bool
	SomeVisible        bool
	RecentLimit        int
	ShowFavorites      bool
	ShowRecents        bool
	Theme              string
	ReducedMotion      bool
	SidebarDensity     string
	DefaultLanding     string
	DemoMode           bool
	Diagnostics        []plugins.Diagnostic
	DisabledPlugins    []plugins.Registered
	LastScan           time.Time
	ScanFailure        string
	Current            *plugins.Registered
	View               *plugins.View
	Error              string
	Message            string
	PluginDetail       *PluginDetail
}

type PluginDetail struct {
	Entry             plugins.Registered
	PlatformNotes     []string
	DiagnosticCommand string
}
type PluginGroup struct {
	ID, Name  string
	Collapsed bool
	Plugins   []plugins.Registered
}

func New(cfg config.Config, registry *plugins.Registry, logger *slog.Logger, initialReport plugins.ScanReport) (*App, error) {
	t, err := template.New("layout.html").Funcs(template.FuncMap{"row": func(row map[string]any, key string) any { return row[key] }, "favorite": func(entries []plugins.Registered, id string) bool {
		for _, entry := range entries {
			if entry.Manifest.ID == id {
				return true
			}
		}
		return false
	}, "hidden": func(ids []string, id string) bool {
		for _, hiddenID := range ids {
			if hiddenID == id {
				return true
			}
		}
		return false
	}, "pluginSearchJSON": func(entries []plugins.Registered) template.JS {
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
	}, "pluginNavigationJSON": func(favorites, recents []plugins.Registered) template.JS {
		ids := func(entries []plugins.Registered) []string {
			result := make([]string, 0, len(entries))
			for _, entry := range entries {
				result = append(result, entry.Manifest.ID)
			}
			return result
		}
		encoded, _ := json.Marshal(struct {
			Favorites []string `json:"favorites"`
			Recents   []string `json:"recents"`
		}{Favorites: ids(favorites), Recents: ids(recents)})
		return template.JS(encoded)
	}, "scanTime": func(value time.Time) string {
		if value.IsZero() {
			return "Not yet scanned"
		}
		return value.Local().Format("Jan 2, 2006 · 3:04:05 PM MST")
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
	navigation := pluginnav.Store{DataDir: cfg.DataDir}
	navState, err := navigation.Load()
	if err != nil {
		return nil, err
	}
	stateStore := pluginstate.Store{DataDir: cfg.DataDir}
	pluginState, err := stateStore.Load()
	if err != nil {
		return nil, err
	}
	workspaceStore := workspace.Store{DataDir: cfg.DataDir}
	workspaceState, err := workspaceStore.Load()
	if err != nil {
		return nil, err
	}
	app := &App{cfg: cfg, registry: registry, discovery: plugins.Discoverer{Directory: cfg.PluginDir, Timeout: cfg.PluginTimeout, Logger: logger, Filter: demoPluginFilter(cfg.DemoMode)}, order: order, orderIDs: orderIDs, groups: groups, groupState: groupState, navigation: navigation, navState: navState, pluginState: stateStore, disabledIDs: pluginState.Disabled, workspace: workspaceStore, workspaceState: workspaceState, diagnostics: initialReport, runner: plugins.Runner{Timeout: cfg.PluginTimeout}, logger: logger, templates: t}
	app.applyDisabledStatuses()
	return app, nil
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
	if a.cfg.DemoMode && r.Method == http.MethodPost && !strings.Contains(r.URL.Path, "/actions/") {
		http.Error(w, "Cmdry demo mode keeps workspace settings read-only", http.StatusForbidden)
		return
	}
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/":
		a.dashboard(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/overview":
		a.dashboard(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/plugins":
		a.pluginList(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/diagnostics":
		a.diagnosticsPage(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/refresh":
		a.refreshPlugins(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/order":
		a.savePluginOrder(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/groups":
		a.saveGroupState(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/favorites":
		a.saveFavorite(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/sidebar-visibility":
		a.savePluginSidebarVisibility(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/sidebar-visibility/all":
		a.saveAllPluginSidebarVisibility(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/plugins/status":
		a.savePluginStatus(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/settings":
		a.settings(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/settings/navigation":
		a.saveNavigationSettings(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/settings/appearance":
		a.saveAppearanceSettings(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/settings/favorites":
		a.removeFavoriteFromSettings(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/settings/recents/clear":
		a.clearRecents(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/settings/navigation/reset":
		a.resetNavigationLayout(w, r)
	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/actions/"):
		a.pluginAction(w, r)
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/details"):
		a.pluginDetails(w, r)
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
	a.navMu.RLock()
	favoriteIDs := append([]string(nil), a.navState.Favorites...)
	recentIDs := append([]string(nil), a.navState.Recents...)
	hiddenIDs := append([]string(nil), a.navState.Hidden...)
	recentLimit := a.navState.RecentLimit
	showFavorites := a.navState.ShowFavorites
	showRecents := a.navState.ShowRecents
	a.navMu.RUnlock()
	a.workspaceMu.RLock()
	workspaceState := a.workspaceState
	a.workspaceMu.RUnlock()
	enabledEntries := entriesWithStatus(entries, plugins.StatusEnabled)
	visibleEntries := entriesWithoutIDs(enabledEntries, hiddenIDs)
	favorites := entriesByID(visibleEntries, favoriteIDs)
	return pageData{Title: title, Section: section, AppVersion: buildinfo.Version, Plugins: entries, EnabledPlugins: enabledEntries, PluginGroups: makeGroups(visibleEntries, state), Favorites: favorites, FavoritesCollapsed: state["favorites"], Recents: entriesByIDExcluding(visibleEntries, recentIDs, favoriteIDs), Hidden: hiddenIDs, AllVisible: len(visibleEntries) == len(enabledEntries), SomeVisible: len(visibleEntries) > 0, RecentLimit: recentLimit, ShowFavorites: showFavorites, ShowRecents: showRecents, Theme: workspaceState.Theme, ReducedMotion: workspaceState.ReducedMotion, SidebarDensity: workspaceState.SidebarDensity, DefaultLanding: workspaceState.DefaultLanding, DemoMode: a.cfg.DemoMode, DisabledPlugins: entriesWithStatus(entries, plugins.StatusDisabled)}
}

func demoPluginFilter(demoMode bool) func(plugins.Manifest) bool {
	if !demoMode {
		return nil
	}
	return func(manifest plugins.Manifest) bool {
		return manifest.Category != "server" && len(manifest.Permissions) == 1 && manifest.Permissions[0] == "data.transform"
	}
}

func entriesWithStatus(entries []plugins.Registered, status plugins.Status) []plugins.Registered {
	result := make([]plugins.Registered, 0, len(entries))
	for _, entry := range entries {
		if entry.Status == status {
			result = append(result, entry)
		}
	}
	return result
}

func entriesWithoutIDs(entries []plugins.Registered, ids []string) []plugins.Registered {
	hidden := make(map[string]bool, len(ids))
	for _, id := range ids {
		hidden[id] = true
	}
	result := make([]plugins.Registered, 0, len(entries))
	for _, entry := range entries {
		if !hidden[entry.Manifest.ID] {
			result = append(result, entry)
		}
	}
	return result
}

func entriesByID(entries []plugins.Registered, ids []string) []plugins.Registered {
	return entriesByIDExcluding(entries, ids, nil)
}

func entriesByIDExcluding(entries []plugins.Registered, ids, excluded []string) []plugins.Registered {
	byID := make(map[string]plugins.Registered, len(entries))
	for _, entry := range entries {
		byID[entry.Manifest.ID] = entry
	}
	skip := make(map[string]bool, len(excluded))
	for _, id := range excluded {
		skip[id] = true
	}
	result := make([]plugins.Registered, 0, len(ids))
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		if entry, ok := byID[id]; ok && !skip[id] && !seen[id] {
			result = append(result, entry)
			seen[id] = true
		}
	}
	return result
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
		if id == "server" {
			continue
		}
		groups = append(groups, *byID[id])
	}
	if server := byID["server"]; server != nil {
		groups = append(groups, *server)
	}
	return groups
}
func (a *App) dashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		a.workspaceMu.RLock()
		defaultLanding, lastPluginID := a.workspaceState.DefaultLanding, a.workspaceState.LastPluginID
		a.workspaceMu.RUnlock()
		if defaultLanding == workspace.LandingLastUsed && lastPluginID != "" {
			if entry, ok := a.registry.Get(lastPluginID); ok && entry.Status == plugins.StatusEnabled {
				http.Redirect(w, r, "/plugins/"+lastPluginID, http.StatusSeeOther)
				return
			}
		}
	}
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
	report, err := a.discovery.DiscoverWithDiagnostics(r.Context(), a.registry)
	a.diagnosticsMu.Lock()
	a.diagnostics = report
	a.diagnosticsMu.Unlock()
	if err != nil {
		a.logger.Warn("plugin refresh failed", "error", err)
		http.Redirect(w, r, "/plugins?refresh=failed", http.StatusSeeOther)
		return
	}
	a.applyDisabledStatuses()
	a.logger.Info("plugins refreshed", "plugins", a.registry.Len())
	http.Redirect(w, r, "/plugins?refreshed=1", http.StatusSeeOther)
}

func (a *App) applyDisabledStatuses() {
	a.pluginStateMu.RLock()
	ids := append([]string(nil), a.disabledIDs...)
	a.pluginStateMu.RUnlock()
	for _, id := range ids {
		a.registry.SetStatus(id, plugins.StatusDisabled)
	}
}

func (a *App) savePluginStatus(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid plugin status request", http.StatusBadRequest)
		return
	}
	id := r.Form.Get("plugin_id")
	status := plugins.Status(r.Form.Get("status"))
	if !plugins.ValidID(id) || (status != plugins.StatusEnabled && status != plugins.StatusDisabled) {
		http.Error(w, "invalid plugin status request", http.StatusBadRequest)
		return
	}
	_, ok := a.registry.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if err := a.setPluginStatus(id, status); err != nil {
		a.logger.Error("save plugin status", "plugin", id, "error", err)
		http.Error(w, "unable to save plugin status", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/plugins", http.StatusSeeOther)
}

func (a *App) setPluginStatus(id string, status plugins.Status) error {
	a.pluginStateMu.Lock()
	state := pluginstate.State{Disabled: append([]string(nil), a.disabledIDs...)}
	if status == plugins.StatusDisabled {
		state.Disabled = appendUnique(state.Disabled, id)
	} else {
		state.Disabled = withoutID(state.Disabled, id)
	}
	if err := a.pluginState.Save(state); err != nil {
		a.pluginStateMu.Unlock()
		return err
	}
	a.disabledIDs = state.Disabled
	a.pluginStateMu.Unlock()
	a.registry.SetStatus(id, status)
	return nil
}

func (a *App) diagnosticsPage(w http.ResponseWriter, r *http.Request) {
	d := a.base("Plugin diagnostics", "Diagnostics")
	a.diagnosticsMu.RLock()
	d.Diagnostics = append([]plugins.Diagnostic(nil), a.diagnostics.Diagnostics...)
	d.LastScan = a.diagnostics.CompletedAt
	d.ScanFailure = a.diagnostics.Failure
	a.diagnosticsMu.RUnlock()
	a.render(w, http.StatusOK, "diagnostics.html", d)
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
	known["favorites"] = true
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

func (a *App) saveFavorite(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid favorite request", http.StatusBadRequest)
		return
	}
	id := r.Form.Get("plugin_id")
	if !plugins.ValidID(id) {
		http.Error(w, "invalid plugin", http.StatusBadRequest)
		return
	}
	if _, ok := a.registry.Get(id); !ok {
		http.NotFound(w, r)
		return
	}
	if err := a.setFavorite(id, r.Form.Get("favorite") == "true"); err != nil {
		a.logger.Error("save plugin favorites", "error", err)
		http.Error(w, "unable to save plugin favorite", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/plugins/"+id, http.StatusSeeOther)
}

func (a *App) savePluginSidebarVisibility(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid sidebar visibility request", http.StatusBadRequest)
		return
	}
	id := r.Form.Get("plugin_id")
	if !plugins.ValidID(id) {
		http.Error(w, "invalid plugin", http.StatusBadRequest)
		return
	}
	entry, ok := a.registry.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if entry.Status != plugins.StatusEnabled {
		http.Error(w, "disabled plugins cannot be shown in the sidebar", http.StatusBadRequest)
		return
	}
	if err := a.setPluginSidebarVisibility(id, r.Form.Get("visible") == "true"); err != nil {
		a.logger.Error("save plugin sidebar visibility", "error", err)
		http.Error(w, "unable to save plugin sidebar visibility", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/plugins", http.StatusSeeOther)
}

func (a *App) saveAllPluginSidebarVisibility(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid sidebar visibility request", http.StatusBadRequest)
		return
	}
	visible := r.Form.Get("visible") == "true"
	ids := make([]string, 0, a.registry.Len())
	if !visible {
		for _, entry := range a.registry.All() {
			if entry.Status == plugins.StatusEnabled {
				ids = append(ids, entry.Manifest.ID)
			}
		}
	}
	if err := a.setAllPluginSidebarVisibility(ids); err != nil {
		a.logger.Error("save all plugin sidebar visibility", "error", err)
		http.Error(w, "unable to save plugin sidebar visibility", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/plugins", http.StatusSeeOther)
}

func (a *App) recordRecent(id string) {
	if a.cfg.DemoMode {
		return
	}
	a.navMu.Lock()
	defer a.navMu.Unlock()
	state := a.navState
	state.Recents = append([]string{id}, withoutID(state.Recents, id)...)
	if len(state.Recents) > state.RecentLimit {
		state.Recents = state.Recents[:state.RecentLimit]
	}
	if slicesEqual(state.Recents, a.navState.Recents) {
		return
	}
	if err := a.navigation.Save(state); err != nil {
		a.logger.Warn("save recently used plugin", "plugin", id, "error", err)
		return
	}
	a.navState = state
}

func (a *App) setFavorite(id string, wantFavorite bool) error {
	a.navMu.Lock()
	defer a.navMu.Unlock()
	state := a.navState
	if wantFavorite {
		state.Favorites = appendUnique(state.Favorites, id)
	} else {
		state.Favorites = withoutID(state.Favorites, id)
	}
	if err := a.navigation.Save(state); err != nil {
		return err
	}
	a.navState = state
	return nil
}

func (a *App) setPluginSidebarVisibility(id string, visible bool) error {
	a.navMu.Lock()
	defer a.navMu.Unlock()
	state := a.navState
	if visible {
		state.Hidden = withoutID(state.Hidden, id)
	} else {
		state.Hidden = appendUnique(state.Hidden, id)
	}
	if err := a.navigation.Save(state); err != nil {
		return err
	}
	a.navState = state
	return nil
}

func (a *App) setAllPluginSidebarVisibility(hidden []string) error {
	a.navMu.Lock()
	defer a.navMu.Unlock()
	state := a.navState
	state.Hidden = hidden
	if err := a.navigation.Save(state); err != nil {
		return err
	}
	a.navState = state
	return nil
}

func appendUnique(ids []string, id string) []string {
	for _, existing := range ids {
		if existing == id {
			return ids
		}
	}
	return append(ids, id)
}

func withoutID(ids []string, id string) []string {
	result := make([]string, 0, len(ids))
	for _, existing := range ids {
		if existing != id {
			result = append(result, existing)
		}
	}
	return result
}

func slicesEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
func (a *App) settings(w http.ResponseWriter, r *http.Request) {
	d := a.base("Settings", "Settings")
	a.render(w, http.StatusOK, "settings.html", d)
}

func (a *App) saveAppearanceSettings(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid appearance settings", http.StatusBadRequest)
		return
	}
	state := workspace.State{
		Theme:          r.Form.Get("theme"),
		ReducedMotion:  r.Form.Get("reduced_motion") == "true",
		SidebarDensity: r.Form.Get("sidebar_density"),
		DefaultLanding: r.Form.Get("default_landing"),
	}
	if !state.Valid() {
		http.Error(w, "invalid appearance settings", http.StatusBadRequest)
		return
	}
	a.workspaceMu.Lock()
	state.LastPluginID = a.workspaceState.LastPluginID
	if err := a.workspace.Save(state); err != nil {
		a.workspaceMu.Unlock()
		a.logger.Error("save appearance settings", "error", err)
		http.Error(w, "unable to save appearance settings", http.StatusInternalServerError)
		return
	}
	a.workspaceState = state
	a.workspaceMu.Unlock()
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (a *App) saveNavigationSettings(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid navigation settings", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(r.Form.Get("recent_limit"))
	if err != nil || !validRecentLimit(limit) {
		http.Error(w, "invalid recent tools limit", http.StatusBadRequest)
		return
	}
	a.navMu.Lock()
	state := a.navState
	state.RecentLimit = limit
	state.ShowFavorites = r.Form.Get("show_favorites") == "true"
	state.ShowRecents = r.Form.Get("show_recents") == "true"
	if len(state.Recents) > limit {
		state.Recents = state.Recents[:limit]
	}
	if err := a.navigation.Save(state); err != nil {
		a.navMu.Unlock()
		a.logger.Error("save navigation settings", "error", err)
		http.Error(w, "unable to save navigation settings", http.StatusInternalServerError)
		return
	}
	a.navState = state
	a.navMu.Unlock()
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (a *App) removeFavoriteFromSettings(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid favorite request", http.StatusBadRequest)
		return
	}
	id := r.Form.Get("plugin_id")
	if !plugins.ValidID(id) {
		http.Error(w, "invalid plugin", http.StatusBadRequest)
		return
	}
	if err := a.setFavorite(id, false); err != nil {
		a.logger.Error("remove plugin favorite", "error", err)
		http.Error(w, "unable to remove plugin favorite", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (a *App) clearRecents(w http.ResponseWriter, r *http.Request) {
	a.navMu.Lock()
	state := a.navState
	state.Recents = nil
	if err := a.navigation.Save(state); err != nil {
		a.navMu.Unlock()
		a.logger.Error("clear recent tools", "error", err)
		http.Error(w, "unable to clear recent tools", http.StatusInternalServerError)
		return
	}
	a.navState = state
	a.navMu.Unlock()
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (a *App) resetNavigationLayout(w http.ResponseWriter, r *http.Request) {
	if err := a.order.Save(nil); err != nil {
		a.logger.Error("reset plugin order", "error", err)
		http.Error(w, "unable to reset sidebar layout", http.StatusInternalServerError)
		return
	}
	if err := a.groups.Save(map[string]bool{}); err != nil {
		a.logger.Error("reset plugin group state", "error", err)
		http.Error(w, "unable to reset sidebar layout", http.StatusInternalServerError)
		return
	}
	a.orderMu.Lock()
	a.orderIDs = nil
	a.orderMu.Unlock()
	a.groupMu.Lock()
	a.groupState = map[string]bool{}
	a.groupMu.Unlock()
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func validRecentLimit(limit int) bool {
	return limit == 4 || limit == 8 || limit == 12
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
	if entry.Status != plugins.StatusEnabled {
		d := a.base(entry.Manifest.Name, entry.Manifest.Name)
		d.Current = &entry
		d.Error = "This plugin is disabled. Re-enable it from the Plugins page before running it."
		a.render(w, http.StatusForbidden, "plugin.html", d)
		return
	}
	a.recordLastPlugin(entry.Manifest.ID)
	a.recordRecent(entry.Manifest.ID)
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

func (a *App) pluginDetails(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(path.Clean(r.URL.Path), "/plugins/"), "/")
	if len(parts) != 2 || parts[1] != "details" || !plugins.ValidID(parts[0]) {
		http.NotFound(w, r)
		return
	}
	entry, ok := a.registry.Get(parts[0])
	if !ok {
		http.NotFound(w, r)
		return
	}
	d := a.base(entry.Manifest.Name+" details", "Plugin details")
	d.Current = &entry
	d.PluginDetail = &PluginDetail{
		Entry:             entry,
		DiagnosticCommand: shellQuote(entry.Path) + " manifest",
		PlatformNotes: []string{
			fmt.Sprintf("Loaded by Cmdry on %s/%s.", runtime.GOOS, runtime.GOARCH),
			"Plugin manifests do not declare platform compatibility; this executable must match the host operating system and architecture.",
		},
	}
	a.render(w, http.StatusOK, "plugin-details.html", d)
}

func (a *App) recordLastPlugin(id string) {
	if a.cfg.DemoMode {
		return
	}
	a.workspaceMu.Lock()
	defer a.workspaceMu.Unlock()
	if a.workspaceState.LastPluginID == id {
		return
	}
	state := a.workspaceState
	state.LastPluginID = id
	if err := a.workspace.Save(state); err != nil {
		a.logger.Warn("save last used plugin", "plugin", id, "error", err)
		return
	}
	a.workspaceState = state
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\\"'\\\"'") + "'"
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
	if entry.Status != plugins.StatusEnabled {
		http.Error(w, "plugin is disabled", http.StatusForbidden)
		return
	}
	params, err := actionParams(w, r)
	if err != nil {
		http.Error(w, "invalid or oversized action input", http.StatusBadRequest)
		return
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

const maxActionBody = 18 * 1024 * 1024
const maxUploadBytes = 4 * 1024 * 1024
const maxUploadFiles = 4

func actionParams(w http.ResponseWriter, r *http.Request) (map[string]any, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxActionBody)
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(maxActionBody); err != nil {
			return nil, err
		}
		defer r.MultipartForm.RemoveAll()
	} else if err := r.ParseForm(); err != nil {
		return nil, err
	}
	fileCount := 0
	if r.MultipartForm != nil {
		fileCount = len(r.MultipartForm.File)
	}
	params := make(map[string]any, len(r.PostForm)+fileCount)
	for key, values := range r.PostForm {
		if !plugins.ValidID(key) || len(values) != 1 {
			return nil, fmt.Errorf("invalid field")
		}
		params[key] = values[0]
	}
	if r.MultipartForm == nil {
		return params, nil
	}
	for key, files := range r.MultipartForm.File {
		if !plugins.ValidID(key) || len(files) < 1 || len(files) > maxUploadFiles {
			return nil, fmt.Errorf("invalid upload")
		}
		uploads := make([]plugins.Upload, 0, len(files))
		for _, header := range files {
			if header.Size < 1 || header.Size > maxUploadBytes {
				return nil, fmt.Errorf("invalid upload")
			}
			file, err := header.Open()
			if err != nil {
				return nil, err
			}
			contents, readErr := io.ReadAll(io.LimitReader(file, maxUploadBytes+1))
			closeErr := file.Close()
			if readErr != nil {
				return nil, readErr
			}
			if closeErr != nil || len(contents) < 1 || len(contents) > maxUploadBytes {
				return nil, fmt.Errorf("invalid upload")
			}
			uploads = append(uploads, plugins.Upload{Name: path.Base(header.Filename), MIMEType: http.DetectContentType(contents), Content: base64.StdEncoding.EncodeToString(contents)})
		}
		if len(uploads) == 1 {
			params[key] = uploads[0]
		} else {
			params[key] = uploads
		}
	}
	return params, nil
}
