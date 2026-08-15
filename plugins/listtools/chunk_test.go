package listtools

import "testing"

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
