package textpriority2

import "testing"

func TestNATO(t *testing.T) {
	encoded, err := NATOEncode("Go 2")
	if err != nil || encoded != "Golf Oscar / Two" {
		t.Fatalf("encode=%q err=%v", encoded, err)
	}
	decoded, err := NATODecode(encoded)
	if err != nil || decoded != "GO 2" {
		t.Fatalf("decode=%q err=%v", decoded, err)
	}
}

func TestASCIIBinary(t *testing.T) {
	encoded, err := ASCIIBinaryEncode("Hi")
	if err != nil || encoded != "01001000 01101001" {
		t.Fatalf("encode=%q err=%v", encoded, err)
	}
	decoded, err := ASCIIBinaryDecode(encoded)
	if err != nil || decoded != "Hi" {
		t.Fatalf("decode=%q err=%v", decoded, err)
	}
}

func TestTextUtilities(t *testing.T) {
	if value, _ := ROT47("Hello!"); value != "w6==@P" {
		t.Fatalf("ROT47=%q", value)
	}
	if value, _ := Numeronymize("internationalization and accessibility"); value != "i18n and a11y" {
		t.Fatalf("numeronym=%q", value)
	}
	if value, _ := ASCIIArt("A1"); value == "" || len(value) < 10 {
		t.Fatalf("art=%q", value)
	}
	if value, err := LoremIpsum(2); err != nil || len(value) == 0 {
		t.Fatalf("lorem=%q err=%v", value, err)
	}
}
