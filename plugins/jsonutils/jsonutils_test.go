package jsonutils

import "testing"

func TestParseCompactsAndIdentifiesRoot(t *testing.T) {
	value, compact, err := Parse(" { \n \"name\": \"Ada\", \"items\": [ 1, 2 ] } ")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(compact), `{"name":"Ada","items":[1,2]}`; got != want {
		t.Fatalf("compact = %q, want %q", got, want)
	}
	if got, want := RootType(value), "object"; got != want {
		t.Fatalf("root type = %q, want %q", got, want)
	}
}

func TestParseRejectsMultipleValues(t *testing.T) {
	if _, _, err := Parse(`{} []`); err == nil {
		t.Fatal("accepted multiple values")
	}
}

func TestStringify(t *testing.T) {
	encoded, err := Stringify(" { \"value\": true } ")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(encoded), `"{\"value\":true}"`; got != want {
		t.Fatalf("stringify = %q, want %q", got, want)
	}
}

func TestEscapeStringAcceptsRawText(t *testing.T) {
	escaped, err := EscapeString("a quote: \"\nand a tab:\t")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(escaped), `"a quote: \"\nand a tab:\t"`; got != want {
		t.Fatalf("escaped = %q, want %q", got, want)
	}
}
