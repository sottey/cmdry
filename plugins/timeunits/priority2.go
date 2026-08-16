package timeunits

import "fmt"

// TruncateClockTime drops smaller units from a signed whole-second duration.
func TruncateClockTime(seconds int64, precision string) (int64, error) {
	unit := int64(1)
	switch precision {
	case "hour":
		unit = 3600
	case "minute":
		unit = 60
	case "second":
	default:
		return 0, fmt.Errorf("choose hour, minute, or second precision")
	}
	return seconds / unit * unit, nil
}
