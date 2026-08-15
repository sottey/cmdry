package csvtools

import "testing"

func TestChangeSeparatorPreservesQuotedFields(t *testing.T) {
	output, rows, err := ChangeSeparator("name;note\nAda;\"one; two\"\n", "semicolon", "comma")
	if err != nil {
		t.Fatal(err)
	}
	if rows != 2 || output != "name,note\nAda,one; two\n" {
		t.Fatalf("rows = %d, output = %q", rows, output)
	}
}

func TestChangeSeparatorRejectsMalformedCSV(t *testing.T) {
	if _, _, err := ChangeSeparator("name;\"unterminated", "semicolon", "comma"); err == nil {
		t.Fatal("expected malformed CSV error")
	}
}
