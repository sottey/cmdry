package unicodeutils

import "testing"

func TestEncodeDecode(t *testing.T) {
	encoded := Encode("A😀")
	if got, want := encoded, "\\u0041\\uD83D\\uDE00"; got != want {
		t.Fatalf("Encode() = %q, want %q", got, want)
	}
	decoded, err := Decode(encoded)
	if err != nil || decoded != "A😀" {
		t.Fatalf("Decode() = %q, err = %v", decoded, err)
	}
}

func TestDecodePreservesOrdinaryBackslashes(t *testing.T) {
	decoded, err := Decode(`C:\work\file \u00E9`)
	if err != nil || decoded != "C:\\work\\file é" {
		t.Fatalf("Decode() = %q, err = %v", decoded, err)
	}
}
