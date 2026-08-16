package pluginnav

import "testing"

func TestLoadSave(t *testing.T) {
	store := Store{DataDir: t.TempDir()}
	if state, err := store.Load(); err != nil || len(state.Favorites) != 0 || len(state.Recents) != 0 {
		t.Fatalf("got state=%#v err=%v", state, err)
	}
	want := State{Favorites: []string{"com.sottey.ports"}, Recents: []string{"com.sottey.processes", "com.sottey.ports"}}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Favorites) != 1 || got.Favorites[0] != want.Favorites[0] || len(got.Recents) != 2 || got.Recents[1] != want.Recents[1] {
		t.Fatalf("got %#v", got)
	}
}
