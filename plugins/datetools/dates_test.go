package datetools

import (
	"testing"
	"time"
)

func TestDates(t *testing.T) {
	years, err := YearsWithWeekday(time.January, 1, 2024, 2025, time.Monday)
	if err != nil || len(years) != 1 || years[0] != 2024 {
		t.Fatalf("%v %v", years, err)
	}
	value, _ := DiscordTimestamp(time.Unix(0, 0), "relative")
	if value != "<t:0:R>" {
		t.Fatal(value)
	}
}
