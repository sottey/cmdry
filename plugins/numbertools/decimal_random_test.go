package numbertools

import "testing"

func TestGenerateRandomDecimalsUsesInclusiveBounds(t *testing.T) {
	values, err := GenerateRandomDecimals("1.25", "1.25", 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range values {
		if value != "1.25" {
			t.Fatalf("value = %q", value)
		}
	}
}

func TestGenerateRandomDecimalsChecksPrecision(t *testing.T) {
	if _, err := GenerateRandomDecimals("1.234", "2", 2, 1); err == nil {
		t.Fatal("expected precision error")
	}
}
