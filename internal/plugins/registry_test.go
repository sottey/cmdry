package plugins

import "testing"

func TestRegistrySetStatus(t *testing.T) {
	registry := NewRegistry()
	if !registry.Add(Registered{Manifest: validManifest(), Status: StatusEnabled}) {
		t.Fatal("failed to add plugin")
	}
	if !registry.SetStatus("ports", StatusDisabled) {
		t.Fatal("failed to update known plugin")
	}
	entry, ok := registry.Get("ports")
	if !ok || entry.Status != StatusDisabled {
		t.Fatalf("entry = %#v", entry)
	}
	if registry.SetStatus("missing", StatusDisabled) {
		t.Fatal("updated unknown plugin")
	}
}
