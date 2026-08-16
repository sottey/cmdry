package numbertools

import (
	"math"
	"testing"
)

func TestParseNumbersAndSummarize(t *testing.T) {
	numbers, err := ParseNumbers("1, 2.5; -3\n4e1")
	if err != nil {
		t.Fatal(err)
	}
	summary, err := Summarize(numbers)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Count != 4 || summary.Sum != 40.5 || summary.Average != 10.125 || summary.Minimum != -3 || summary.Maximum != 40 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}

func TestParseNumbersRejectsInvalidValues(t *testing.T) {
	if _, err := ParseNumbers("1 two 3"); err == nil {
		t.Fatal("expected invalid value error")
	}
}

func TestSphereArea(t *testing.T) {
	area, err := SphereArea(2)
	if err != nil || math.Abs(area-16*math.Pi) > 0.000001 {
		t.Fatalf("area = %v, err = %v", area, err)
	}
}
