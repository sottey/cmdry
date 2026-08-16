package listtools

import "testing"

func TestUniqueLinesPreservesFirstOccurrence(t *testing.T) {
	items, err := UniqueLines("pear\n\napple\npear\nApple\napple")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := FormatChunks([][]string{items}), "pear\napple\nApple"; got != want {
		t.Fatalf("items = %q, want %q", got, want)
	}
}
