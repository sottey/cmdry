package jsoncsv

import "testing"

func TestConvertFlattensObjectsAndPreservesArrays(t *testing.T) {
	result, err := Convert(`[{"name":"Ada","contact":{"email":"ada@example.test"},"roles":["admin","ops"]},{"name":"Lin","contact":{"email":"lin@example.test"}}]`)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(result.CSV), "contact.email,name,roles\nada@example.test,Ada,\"[\"\"admin\"\",\"\"ops\"\"]\"\nlin@example.test,Lin,\n"; got != want {
		t.Fatalf("CSV = %q, want %q", got, want)
	}
}

func TestConvertRejectsNonObjectArrays(t *testing.T) {
	if _, err := Convert(`[1,2,3]`); err == nil {
		t.Fatal("accepted scalar array")
	}
}
