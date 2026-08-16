package listtools

import (
	"strings"
	"testing"
)

func TestChunkLines(t *testing.T) {
	chunks, err := ChunkLines("one\n\ntwo\nthree\nfour", 3)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := FormatChunks(chunks), "one\ntwo\nthree\n\nfour"; got != want {
		t.Fatalf("chunks = %q, want %q", got, want)
	}
}

func TestChunkLinesRejectsEmptyInput(t *testing.T) {
	if _, err := ChunkLines("\n  \n", 2); err == nil {
		t.Fatal("expected empty list error")
	}
}

func TestSortLines(t *testing.T) {
	items, err := SortLines("zebra\nApple\nbanana", false)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.Join(items, ","), "Apple,banana,zebra"; got != want {
		t.Fatalf("sorted = %q, want %q", got, want)
	}
}
