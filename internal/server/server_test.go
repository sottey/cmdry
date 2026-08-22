package server

import (
	"bytes"
	"encoding/base64"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/sottey/cmdry/internal/config"
	"github.com/sottey/cmdry/internal/pluginnav"
	"github.com/sottey/cmdry/internal/plugins"
	"github.com/sottey/cmdry/internal/pluginstate"
)

func TestActionParamsReadsEphemeralUpload(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("width", "12")
	part, err := writer.CreateFormFile("image", "example.png")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("\x89PNG\r\n\x1a\ncontent"))
	_ = writer.Close()
	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	params, err := actionParams(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	if params["width"] != "12" {
		t.Fatalf("width=%#v", params["width"])
	}
	upload, ok := params["image"].(plugins.Upload)
	if !ok {
		t.Fatalf("upload=%T", params["image"])
	}
	contents, err := base64.StdEncoding.DecodeString(upload.Content)
	if err != nil || string(contents) != "\x89PNG\r\n\x1a\ncontent" {
		t.Fatalf("contents=%q err=%v", contents, err)
	}
}

func TestDemoModeRejectsWorkspaceWrites(t *testing.T) {
	app := &App{cfg: config.Config{DemoMode: true}}
	req := httptest.NewRequest("POST", "/settings/appearance", nil)
	response := httptest.NewRecorder()
	app.ServeHTTP(response, req)
	if response.Code != 403 {
		t.Fatalf("status = %d, want 403", response.Code)
	}
}

func TestMakeGroupsKeepsServerLast(t *testing.T) {
	entries := []plugins.Registered{
		{Manifest: plugins.Manifest{ID: "com.sottey.ports", Name: "Port Inspector", Category: "server"}},
		{Manifest: plugins.Manifest{ID: "com.sottey.base64", Name: "Base64", Category: "text"}},
		{Manifest: plugins.Manifest{ID: "com.sottey.processes", Name: "Processes", Category: "server"}},
		{Manifest: plugins.Manifest{ID: "com.sottey.sum", Name: "Sum", Category: "number"}},
	}
	groups := makeGroups(entries, nil)
	if len(groups) != 3 {
		t.Fatalf("group count = %d, want 3", len(groups))
	}
	if groups[0].ID != "text" || groups[1].ID != "number" || groups[2].ID != "server" {
		t.Fatalf("group order = %q, %q, %q", groups[0].ID, groups[1].ID, groups[2].ID)
	}
	if len(groups[2].Plugins) != 2 {
		t.Fatalf("server plugin count = %d, want 2", len(groups[2].Plugins))
	}
}

func TestActionParamsReadsMultipleEphemeralUploads(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for _, name := range []string{"one.png", "two.png"} {
		part, err := writer.CreateFormFile("images", name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write([]byte("image-" + name)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	params, err := actionParams(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	uploads, ok := params["images"].([]plugins.Upload)
	if !ok || len(uploads) != 2 {
		t.Fatalf("uploads=%#v", params["images"])
	}
	if uploads[0].Name != "one.png" || uploads[1].Name != "two.png" {
		t.Fatalf("upload names=%q, %q", uploads[0].Name, uploads[1].Name)
	}
}

func TestDisabledPluginsStayRegisteredButLeaveNavigation(t *testing.T) {
	registry := plugins.NewRegistry()
	registry.Add(plugins.Registered{Manifest: plugins.Manifest{ID: "com.sottey.text", Name: "Text", Category: "text"}, Status: plugins.StatusEnabled})
	registry.Add(plugins.Registered{Manifest: plugins.Manifest{ID: "com.sottey.ports", Name: "Ports", Category: "server"}, Status: plugins.StatusEnabled})
	app := &App{registry: registry, groupState: map[string]bool{}, pluginState: pluginstate.Store{DataDir: t.TempDir()}}
	if err := app.setPluginStatus("com.sottey.ports", plugins.StatusDisabled); err != nil {
		t.Fatal(err)
	}
	ports, ok := registry.Get("com.sottey.ports")
	if !ok || ports.Status != plugins.StatusDisabled {
		t.Fatalf("ports status = %#v", ports)
	}
	persisted, err := app.pluginState.Load()
	if err != nil || len(persisted.Disabled) != 1 || persisted.Disabled[0] != "com.sottey.ports" {
		t.Fatalf("persisted state = %#v, %v", persisted, err)
	}
	data := app.base("Cmdry", "Overview")
	if len(data.Plugins) != 2 || len(data.EnabledPlugins) != 1 || len(data.DisabledPlugins) != 1 {
		t.Fatalf("plugin sets = all:%d enabled:%d disabled:%d", len(data.Plugins), len(data.EnabledPlugins), len(data.DisabledPlugins))
	}
	if len(data.PluginGroups) != 1 || data.PluginGroups[0].ID != "text" {
		t.Fatalf("visible groups = %#v", data.PluginGroups)
	}
	if err := app.setPluginStatus("com.sottey.ports", plugins.StatusEnabled); err != nil {
		t.Fatal(err)
	}
	persisted, err = app.pluginState.Load()
	if err != nil || len(persisted.Disabled) != 0 {
		t.Fatalf("re-enabled state = %#v, %v", persisted, err)
	}
}

func TestFavoritesUsePersistedSectionState(t *testing.T) {
	registry := plugins.NewRegistry()
	registry.Add(plugins.Registered{Manifest: plugins.Manifest{ID: "com.sottey.text", Name: "Text", Category: "text"}, Status: plugins.StatusEnabled})
	app := &App{
		registry:   registry,
		groupState: map[string]bool{"favorites": true},
		navState:   pluginnav.State{Favorites: []string{"com.sottey.text"}, ShowFavorites: true},
	}
	data := app.base("Cmdry", "Overview")
	if len(data.Favorites) != 1 || !data.FavoritesCollapsed {
		t.Fatalf("favorites = %#v, collapsed = %t", data.Favorites, data.FavoritesCollapsed)
	}
}
