package qrcode

import (
	"bytes"
	"testing"
)

func TestGeneratePNG(t *testing.T) {
	png, err := Generate("https://example.test", "medium", 128)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(png, []byte{0x89, 'P', 'N', 'G'}) {
		t.Fatal("expected PNG data")
	}
}
