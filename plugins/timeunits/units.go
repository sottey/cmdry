// Package timeunits provides local time-unit conversion helpers.
package timeunits

import (
	"fmt"
	"math"
)

const HoursPerDay = 24.0

// DaysToHours converts a finite count of days into hours.
func DaysToHours(days float64) (float64, error) {
	if math.IsNaN(days) || math.IsInf(days, 0) {
		return 0, fmt.Errorf("value must be a finite number")
	}
	return days * HoursPerDay, nil
}

// HoursToDays converts a finite count of hours into days.
func HoursToDays(hours float64) (float64, error) {
	if math.IsNaN(hours) || math.IsInf(hours, 0) {
		return 0, fmt.Errorf("value must be a finite number")
	}
	return hours / HoursPerDay, nil
}
