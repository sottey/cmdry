package basicconverters

import "testing"

func TestIntegerValues(t *testing.T) {
	values, err := IntegerValues("FF", 16)
	if err != nil {
		t.Fatal(err)
	}
	if values["binary"] != "11111111" || values["octal"] != "377" || values["decimal"] != "255" || values["hexadecimal"] != "FF" {
		t.Fatalf("values = %#v", values)
	}
	if _, err := IntegerValues("2", 2); err == nil {
		t.Fatal("accepted invalid binary input")
	}
}

func TestRomanConversions(t *testing.T) {
	roman, err := ToRoman(2026)
	if err != nil || roman != "MMXXVI" {
		t.Fatalf("ToRoman = %q, %v", roman, err)
	}
	value, err := FromRoman("mmxxvi")
	if err != nil || value != 2026 {
		t.Fatalf("FromRoman = %d, %v", value, err)
	}
	if _, err := FromRoman("IIV"); err == nil {
		t.Fatal("accepted non-canonical numeral")
	}
}

func TestColorConversions(t *testing.T) {
	color, err := ParseColor("hsl(180, 50%, 50%)")
	if err != nil {
		t.Fatal(err)
	}
	if color.Hex() != "#40BFBF" || color.RGB() != "rgb(64, 191, 191)" || color.HSL() != "hsl(180, 49.8%, 50%)" {
		t.Fatalf("color = %#v, hex=%s rgb=%s hsl=%s", color, color.Hex(), color.RGB(), color.HSL())
	}
	if _, err := ParseColor("rgb(256, 0, 0)"); err == nil {
		t.Fatal("accepted invalid RGB channel")
	}
}
