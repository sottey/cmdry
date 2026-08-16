package numbertools

import "testing"

func TestGenerateRandomIntegersUsesInclusiveBounds(t *testing.T) {
	values, err := GenerateRandomIntegers("7", "7", 3)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range values {
		if value != "7" {
			t.Fatalf("value = %q", value)
		}
	}
}

func TestGenerateRandomIntegersRejectsInvalidRange(t *testing.T) {
	if _, err := GenerateRandomIntegers("5", "4", 1); err == nil {
		t.Fatal("expected invalid range error")
	}
}
