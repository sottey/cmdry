package jsoncompare

import "testing"

func TestCompareIgnoresObjectKeyOrder(t *testing.T) {
	result, err := Compare(`{"name":"Ada","active":true}`, `{"active":true,"name":"Ada"}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Differences) != 0 {
		t.Fatalf("differences = %#v", result.Differences)
	}
}

func TestCompareReportsNestedChanges(t *testing.T) {
	result, err := Compare(`{"user":{"name":"Ada","role":"admin"},"old":true}`, `{"user":{"name":"Lin","email":"lin@example.test"},"new":false}`)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(result.Differences), 5; got != want {
		t.Fatalf("difference count = %d, want %d: %#v", got, want, result.Differences)
	}
	if result.Differences[0].Path != `$["new"]` || result.Differences[4].Path != `$["user"]["role"]` {
		t.Fatalf("unexpected paths: %#v", result.Differences)
	}
}
