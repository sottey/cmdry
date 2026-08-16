package urltools

import "testing"

func TestBuild(t *testing.T) {
	output, details, err := Build("https://example.test/old?one=1#top", "", "", "/new path", "page=2", "section")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := output, "https://example.test/new%20path?page=2#section"; got != want {
		t.Fatalf("URL = %q, want %q", got, want)
	}
	if details.Query != "page=2" || details.Fragment != "section" {
		t.Fatalf("details = %#v", details)
	}
}

func TestParseRequiresCompleteURL(t *testing.T) {
	if _, err := Parse("example.test/path"); err == nil {
		t.Fatal("expected complete URL error")
	}
}
