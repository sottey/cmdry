package texttools

import (
	"net/url"
	"strings"
	"testing"
)

func TestExtractEmailsPreservesFirstUniqueMatch(t *testing.T) {
	emails := ExtractEmails("Contact Ada@Example.test, lin@example.test, then ada@example.test. Ignore not-an-email@localhost.")
	if got, want := len(emails), 2; got != want {
		t.Fatalf("emails = %#v", emails)
	}
	if got, want := emails[0], "Ada@Example.test"; got != want {
		t.Fatalf("first email = %q, want %q", got, want)
	}
}

func TestDeduplicateLines(t *testing.T) {
	output, removed := DeduplicateLines("one\ntwo\none\ntwo\nthree")
	if got, want := output, "one\ntwo\nthree"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
	if got, want := removed, 2; got != want {
		t.Fatalf("removed = %d, want %d", got, want)
	}
}

func TestReplaceAll(t *testing.T) {
	output, count := ReplaceAll("red red red", "red", "blue")
	if got, want := output, "blue blue blue"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
	if got, want := count, 3; got != want {
		t.Fatalf("count = %d, want %d", got, want)
	}
}

func TestCaseAndURLTransforms(t *testing.T) {
	if got, want := strings.ToUpper("Seán café"), "SEÁN CAFÉ"; got != want {
		t.Fatalf("uppercase = %q, want %q", got, want)
	}
	encoded := url.QueryEscape("hello world/a?b")
	if got, want := encoded, "hello+world%2Fa%3Fb"; got != want {
		t.Fatalf("encoded = %q, want %q", got, want)
	}
	decoded, err := url.QueryUnescape(encoded)
	if err != nil || decoded != "hello world/a?b" {
		t.Fatalf("decoded = %q, error = %v", decoded, err)
	}
}

func TestTextStatistics(t *testing.T) {
	statistics := TextStatistics("One café\nTwo")
	if statistics.Characters != 12 || statistics.Words != 3 || statistics.Lines != 2 || statistics.Bytes != 13 {
		t.Fatalf("unexpected statistics: %#v", statistics)
	}
	if statistics.ReadingTime != "Less than 1 minute" {
		t.Fatalf("reading time = %q", statistics.ReadingTime)
	}
}

func TestExtractSubstringUsesCharacterPositions(t *testing.T) {
	output, err := ExtractSubstring("A café", 3, 5)
	if err != nil || output != "caf" {
		t.Fatalf("output = %q, err = %v", output, err)
	}
	if _, err := ExtractSubstring("abc", 2, 4); err == nil {
		t.Fatal("expected out-of-range error")
	}
}

func TestJoinLines(t *testing.T) {
	output, count, err := JoinLines("one\n\ntwo\nthree", "comma", "")
	if err != nil || count != 3 || output != "one, two, three" {
		t.Fatalf("output = %q, count = %d, err = %v", output, count, err)
	}
}

func TestReverseTextUsesUnicodeCharacters(t *testing.T) {
	if got, want := ReverseText("A café"), "éfac A"; got != want {
		t.Fatalf("reversed = %q, want %q", got, want)
	}
}
