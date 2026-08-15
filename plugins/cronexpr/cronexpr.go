// Package cronexpr validates and explains standard five-field cron expressions.
package cronexpr

import (
	"fmt"
	"strconv"
	"strings"
)

type Field struct {
	Name        string
	Expression  string
	Explanation string
}

var shortcuts = map[string]string{
	"@yearly":   "At midnight on January 1 every year.",
	"@annually": "At midnight on January 1 every year.",
	"@monthly":  "At midnight on the first day of every month.",
	"@weekly":   "At midnight every Sunday.",
	"@daily":    "At midnight every day.",
	"@midnight": "At midnight every day.",
	"@hourly":   "At the start of every hour.",
}

type fieldSpec struct {
	name, singular, plural string
	min, max               int
	names                  map[string]int
}

var specs = []fieldSpec{
	{name: "Minute", singular: "minute", plural: "minutes", min: 0, max: 59},
	{name: "Hour", singular: "hour", plural: "hours", min: 0, max: 23},
	{name: "Day of month", singular: "day of the month", plural: "days of the month", min: 1, max: 31},
	{name: "Month", singular: "month", plural: "months", min: 1, max: 12, names: map[string]int{"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4, "MAY": 5, "JUN": 6, "JUL": 7, "AUG": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12}},
	{name: "Day of week", singular: "day of the week", plural: "days of the week", min: 0, max: 7, names: map[string]int{"SUN": 0, "MON": 1, "TUE": 2, "WED": 3, "THU": 4, "FRI": 5, "SAT": 6}},
}

// Explain validates an expression and returns an explanation for each field.
// It supports standard crontab syntax, not Quartz's seconds/year extensions.
func Explain(input string) ([]Field, string, error) {
	expression := strings.TrimSpace(input)
	if description, ok := shortcuts[strings.ToLower(expression)]; ok {
		return nil, description, nil
	}
	parts := strings.Fields(expression)
	if len(parts) != 5 {
		return nil, "", fmt.Errorf("enter exactly five fields (minute hour day-of-month month day-of-week), or a supported @shortcut")
	}
	fields := make([]Field, 0, len(parts))
	for index, part := range parts {
		explanation, err := explainField(part, specs[index])
		if err != nil {
			return nil, "", fmt.Errorf("%s field: %w", strings.ToLower(specs[index].name), err)
		}
		fields = append(fields, Field{Name: specs[index].name, Expression: part, Explanation: explanation})
	}
	return fields, "Runs when all shown field conditions match. Traditional cron implementations may treat day-of-month and day-of-week with OR semantics; confirm the behavior of the cron service that will run it.", nil
}

func explainField(expression string, spec fieldSpec) (string, error) {
	parts := strings.Split(expression, ",")
	descriptions := make([]string, 0, len(parts))
	for _, part := range parts {
		description, err := explainPart(part, spec)
		if err != nil {
			return "", err
		}
		descriptions = append(descriptions, description)
	}
	if len(descriptions) == 1 {
		return descriptions[0], nil
	}
	return strings.Join(descriptions, "; "), nil
}

func explainPart(part string, spec fieldSpec) (string, error) {
	if part == "" {
		return "", fmt.Errorf("empty list item")
	}
	base, step, hasStep := part, 0, false
	if strings.Contains(part, "/") {
		pieces := strings.Split(part, "/")
		if len(pieces) != 2 || pieces[0] == "" || pieces[1] == "" {
			return "", fmt.Errorf("invalid step %q", part)
		}
		parsed, err := strconv.Atoi(pieces[1])
		if err != nil || parsed < 1 {
			return "", fmt.Errorf("step must be a positive number")
		}
		base, step, hasStep = pieces[0], parsed, true
	}
	if base == "*" {
		if hasStep {
			return fmt.Sprintf("every %d %s", step, spec.plural), nil
		}
		return "every " + spec.singular, nil
	}
	if strings.Contains(base, "-") {
		bounds := strings.Split(base, "-")
		if len(bounds) != 2 {
			return "", fmt.Errorf("invalid range %q", base)
		}
		start, err := parseValue(bounds[0], spec)
		if err != nil {
			return "", err
		}
		end, err := parseValue(bounds[1], spec)
		if err != nil {
			return "", err
		}
		if start > end {
			return "", fmt.Errorf("range start must not be after range end")
		}
		if hasStep {
			return fmt.Sprintf("every %d %s from %s through %s", step, spec.plural, bounds[0], bounds[1]), nil
		}
		return fmt.Sprintf("%s %s through %s", spec.plural, bounds[0], bounds[1]), nil
	}
	if hasStep {
		return "", fmt.Errorf("a step requires * or a range")
	}
	if _, err := parseValue(base, spec); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", spec.singular, base), nil
}

func parseValue(value string, spec fieldSpec) (int, error) {
	if named, ok := spec.names[strings.ToUpper(value)]; ok {
		return named, nil
	}
	number, err := strconv.Atoi(value)
	if err != nil || number < spec.min || number > spec.max {
		return 0, fmt.Errorf("%q must be between %d and %d", value, spec.min, spec.max)
	}
	return number, nil
}
