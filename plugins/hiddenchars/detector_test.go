package hiddenchars

import "testing"

func TestDetectFindsInvisibleCharactersWithLocation(t *testing.T) {
	findings := Detect("plain\nzero\u200bwidth\u00a0space")
	if got, want := len(findings), 2; got != want {
		t.Fatalf("findings = %#v", findings)
	}
	if got, want := findings[0], (Finding{Line: 2, Column: 5, CodePoint: "U+200B", Description: "Zero-width space"}); got != want {
		t.Fatalf("first finding = %#v, want %#v", got, want)
	}
}

func TestDetectAllowsNormalWhitespace(t *testing.T) {
	if findings := Detect("one two\tthree\nfour\r\n"); len(findings) != 0 {
		t.Fatalf("findings = %#v", findings)
	}
}
