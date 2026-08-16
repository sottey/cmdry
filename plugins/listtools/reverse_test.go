package listtools

import "testing"

func TestReverseLines(t *testing.T) {
	items, err := ReverseLines("one\n\ntwo\nthree")
	if err != nil || FormatChunks([][]string{items}) != "three\ntwo\none" {
		t.Fatalf("items = %#v, err = %v", items, err)
	}
}
