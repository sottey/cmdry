package pluginorder

import "testing"

func TestLoadSave(t *testing.T) {
	store := Store{DataDir: t.TempDir()}
	if ids, err := store.Load(); err != nil || ids != nil {
		t.Fatalf("got ids=%#v err=%v", ids, err)
	}
	want := []string{"com.sottey.processes", "com.sottey.ports"}
	if err := store.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %#v", got)
	}
}
