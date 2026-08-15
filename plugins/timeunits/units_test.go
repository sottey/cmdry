package timeunits

import "testing"

func TestDayAndHourConversions(t *testing.T) {
	hours, err := DaysToHours(1.5)
	if err != nil || hours != 36 {
		t.Fatalf("hours = %v, err = %v", hours, err)
	}
	days, err := HoursToDays(-12)
	if err != nil || days != -0.5 {
		t.Fatalf("days = %v, err = %v", days, err)
	}
}
