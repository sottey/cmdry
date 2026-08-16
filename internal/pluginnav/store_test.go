package pluginnav

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSave(t *testing.T) {
	store := Store{DataDir: t.TempDir()}
	if state, err := store.Load(); err != nil || len(state.Favorites) != 0 || len(state.Recents) != 0 || state.RecentLimit != DefaultRecentLimit || !state.ShowFavorites || !state.ShowRecents {
		t.Fatalf("got state=%#v err=%v", state, err)
	}
	want := State{Favorites: []string{"com.sottey.ports"}, Recents: []string{"com.sottey.processes", "com.sottey.ports"}, Hidden: []string{"com.sottey.journal"}, RecentLimit: 4, ShowFavorites: false, ShowRecents: true}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Favorites) != 1 || got.Favorites[0] != want.Favorites[0] || len(got.Recents) != 2 || got.Recents[1] != want.Recents[1] || len(got.Hidden) != 1 || got.Hidden[0] != want.Hidden[0] || got.RecentLimit != 4 || got.ShowFavorites || !got.ShowRecents {
		t.Fatalf("got %#v", got)
	}
}

func TestLoadDefaultsVisibilityForExistingState(t *testing.T) {
	store := Store{DataDir: t.TempDir()}
	if err := os.WriteFile(filepath.Join(store.DataDir, filename), []byte(`{"favorites":[],"recents":[],"recent_limit":4}`), 0600); err != nil {
		t.Fatal(err)
	}
	state, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !state.ShowFavorites || !state.ShowRecents {
		t.Fatalf("existing state should preserve visible defaults: %#v", state)
	}
}
