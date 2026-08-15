package timeunits

import "testing"

func TestParseAndFormatClockDuration(t *testing.T) {
	seconds, err := ParseClockDuration("01:30:15")
	if err != nil || seconds != 5415 {
		t.Fatalf("seconds = %d, err = %v", seconds, err)
	}
	if got, want := FormatSeconds(-5415), "-01:30:15"; got != want {
		t.Fatalf("formatted = %q, want %q", got, want)
	}
}

func TestParseClockDurationSupportsMinutesAndSeconds(t *testing.T) {
	seconds, err := ParseClockDuration("03:05")
	if err != nil || seconds != 185 {
		t.Fatalf("seconds = %d, err = %v", seconds, err)
	}
}

func TestParseClockDurationRejectsInvalidFields(t *testing.T) {
	if _, err := ParseClockDuration("1:60:00"); err == nil {
		t.Fatal("expected invalid minutes error")
	}
}
