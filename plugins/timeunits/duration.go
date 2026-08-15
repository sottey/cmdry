package timeunits

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseClockDuration accepts signed HH:MM:SS or MM:SS input and returns the
// total number of whole seconds.
func ParseClockDuration(input string) (int64, error) {
	value := strings.TrimSpace(input)
	if value == "" {
		return 0, fmt.Errorf("time is required")
	}
	sign := int64(1)
	if value[0] == '-' || value[0] == '+' {
		if value[0] == '-' {
			sign = -1
		}
		value = value[1:]
	}
	parts := strings.Split(value, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return 0, fmt.Errorf("enter HH:MM:SS or MM:SS")
	}
	values := make([]int64, len(parts))
	for index, part := range parts {
		parsed, err := parseClockPart(part)
		if err != nil {
			return 0, err
		}
		values[index] = parsed
	}
	var hours, minutes, seconds int64
	if len(values) == 2 {
		minutes, seconds = values[0], values[1]
	} else {
		hours, minutes, seconds = values[0], values[1], values[2]
	}
	if minutes > 59 || seconds > 59 {
		return 0, fmt.Errorf("minutes and seconds must be between 0 and 59")
	}
	if hours > (1<<63-1)/3600 {
		return 0, fmt.Errorf("time is too large")
	}
	total := hours*3600 + minutes*60 + seconds
	if sign < 0 {
		total = -total
	}
	return total, nil
}

func parseClockPart(input string) (int64, error) {
	if input == "" {
		return 0, fmt.Errorf("time contains an empty field")
	}
	value, err := strconv.ParseInt(input, 10, 64)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("time fields must be non-negative whole numbers")
	}
	return value, nil
}

// FormatSeconds returns a signed HH:MM:SS representation of whole seconds.
func FormatSeconds(seconds int64) string {
	negative := seconds < 0
	value := uint64(seconds)
	if negative {
		value = uint64(-(seconds + 1)) + 1
	}
	hours := value / 3600
	minutes := (value % 3600) / 60
	remaining := value % 60
	prefix := ""
	if negative {
		prefix = "-"
	}
	return fmt.Sprintf("%s%02d:%02d:%02d", prefix, hours, minutes, remaining)
}

// DecimalHours converts whole seconds into decimal hours.
func DecimalHours(seconds int64) float64 { return float64(seconds) / 3600 }
