package timediff

import (
	"testing"
	"time"
)

func TestBetween(t *testing.T) {
	difference, err := Between("2026-08-16T00:00:00Z", "2026-08-18T03:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	if difference.Duration != 51*time.Hour+4*time.Minute+5*time.Second || difference.Days != 2 || difference.Hours != 51 {
		t.Fatalf("unexpected difference: %#v", difference)
	}
}

func TestBetweenEarlierEnd(t *testing.T) {
	difference, err := Between("2026-08-18", "2026-08-16")
	if err != nil {
		t.Fatal(err)
	}
	if !difference.EndBeforeStart || difference.Days != 2 {
		t.Fatalf("unexpected difference: %#v", difference)
	}
}
