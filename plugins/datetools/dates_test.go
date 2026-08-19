package datetools

import "testing"

func TestLeapYears(t *testing.T) {
	if IsLeapYear(1900) || !IsLeapYear(2000) || !IsLeapYear(2024) || IsLeapYear(2025) {
		t.Fatal("unexpected leap-year results")
	}
	years, err := LeapYears(1999, 2004)
	if err != nil || len(years) != 2 || years[0] != 2000 || years[1] != 2004 {
		t.Fatalf("years=%v err=%v", years, err)
	}
}
