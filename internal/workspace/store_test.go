package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSave(t *testing.T) {
	store := Store{DataDir: t.TempDir()}
	if got, err := store.Load(); err != nil || got != DefaultState() {
		t.Fatalf("default state = %#v, %v", got, err)
	}
	want := State{Theme: ThemeDark, ReducedMotion: true, SidebarDensity: DensityCompact, DefaultLanding: LandingLastUsed, LastPluginID: "com.sottey.guid"}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil || got != want {
		t.Fatalf("loaded state = %#v, %v", got, err)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	store := Store{DataDir: t.TempDir()}
	if err := os.WriteFile(filepath.Join(store.DataDir, filename), []byte(`{"theme":"purple","sidebar_density":"compact","default_landing":"overview"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Load(); err == nil {
		t.Fatal("Load should reject invalid preferences")
	}
}
