package jsonutils

import "testing"

func TestSortKeysAndArrays(t *testing.T) {
	got, err := Sort(`{"z":{"b":1,"a":2},"a":["z","a"]}`, true)
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"a\": [\n    \"a\",\n    \"z\"\n  ],\n  \"z\": {\n    \"a\": 2,\n    \"b\": 1\n  }\n}"
	if string(got) != want {
		t.Fatalf("Sort() = %s, want %s", got, want)
	}
}
