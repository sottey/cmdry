// Package timeconvert converts Unix timestamps and date-time strings locally.
package timeconvert

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const localLayout = "2006-01-02 15:04:05 MST"

// FromUnixSeconds converts a whole Unix timestamp to its local and UTC forms.
func FromUnixSeconds(input string) (time.Time, error) {
	seconds, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("Unix timestamp must be a whole number of seconds")
	}
	return time.Unix(seconds, 0), nil
}

// ToUnixSeconds parses RFC 3339 input. For a date without an offset, zone
// controls whether its displayed clock time is interpreted as UTC or local.
func ToUnixSeconds(input, zone string) (time.Time, error) {
	input = strings.TrimSpace(input)
	if parsed, err := time.Parse(time.RFC3339, input); err == nil {
		return parsed, nil
	}
	location := time.UTC
	if zone == "local" {
		location = time.Local
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", input, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("enter an RFC 3339 date or YYYY-MM-DD HH:MM:SS: %w", err)
	}
	return parsed, nil
}

// Details returns stable display values for a moment in time.
func Details(moment time.Time) map[string]string {
	return map[string]string{
		"unix":  strconv.FormatInt(moment.Unix(), 10),
		"utc":   moment.UTC().Format(time.RFC3339),
		"local": moment.Local().Format(localLayout),
	}
}
