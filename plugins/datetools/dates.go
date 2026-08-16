package datetools

import (
	"fmt"
	"time"
)

func YearsWithWeekday(month time.Month, day, start, end int, weekday time.Weekday) ([]int, error) {
	if start > end || start < 1 || end > 9999 {
		return nil, fmt.Errorf("use a year range from 1 through 9999")
	}
	if month < time.January || month > time.December || day < 1 || day > 31 {
		return nil, fmt.Errorf("enter a valid month and day")
	}
	years := []int{}
	for year := start; year <= end; year++ {
		value := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		if value.Month() != month || value.Day() != day {
			continue
		}
		if value.Weekday() == weekday {
			years = append(years, year)
		}
	}
	return years, nil
}

func DiscordTimestamp(value time.Time, style string) (string, error) {
	styles := map[string]string{"short-time": "t", "long-time": "T", "short-date": "d", "long-date": "D", "short-date-time": "f", "long-date-time": "F", "relative": "R"}
	format, ok := styles[style]
	if !ok {
		return "", fmt.Errorf("unsupported Discord timestamp style")
	}
	return fmt.Sprintf("<t:%d:%s>", value.Unix(), format), nil
}
