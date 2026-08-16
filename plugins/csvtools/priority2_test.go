package csvtools

import "testing"

func TestPriority2CSVTools(t *testing.T) {
	swapped, err := SwapColumns("name,age\nAda,37\n", "name", "age")
	if err != nil || swapped != "age,name\n37,Ada\n" {
		t.Fatalf("swap=%q err=%v", swapped, err)
	}
	transposed, columns, rows, err := Transpose("name,age\nAda,37\n")
	if err != nil || columns != 2 || rows != 1 || transposed != "name,Ada\nage,37\n" {
		t.Fatalf("transpose=%q columns=%d rows=%d err=%v", transposed, columns, rows, err)
	}
}
