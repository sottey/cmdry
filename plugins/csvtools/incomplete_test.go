package csvtools

import "testing"

func TestFindIncompleteRecords(t *testing.T) {
	records, rows, err := FindIncompleteRecords("name,age,city\nAda,,Paris\nLin,30\n")
	if err != nil {
		t.Fatal(err)
	}
	if rows != 2 || len(records) != 2 || records[0].Row != 2 || records[0].Columns[0] != "age" || records[1].Columns[0] != "city" {
		t.Fatalf("records = %#v, rows = %d", records, rows)
	}
}
