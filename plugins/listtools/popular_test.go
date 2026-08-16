package listtools

import "testing"

func TestFindMostPopular(t *testing.T) {
	items, err := FindMostPopular("pear\napple\npear\nbanana\napple\npear")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 || items[0] != (PopularItem{Item: "pear", Count: 3}) || items[1] != (PopularItem{Item: "apple", Count: 2}) {
		t.Fatalf("items = %#v", items)
	}
}
