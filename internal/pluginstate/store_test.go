package pluginstate

import "testing"

func TestLoadSave(t *testing.T) {
	store := Store{DataDir: t.TempDir()}
	state, err := store.Load()
	if err != nil || len(state.Disabled) != 0 {
		t.Fatalf("default state = %#v, %v", state, err)
	}
	want := State{Disabled: []string{"com.sottey.journal", "com.sottey.ports"}}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil || len(got.Disabled) != 2 || got.Disabled[1] != "com.sottey.ports" {
		t.Fatalf("loaded state = %#v, %v", got, err)
	}
}
