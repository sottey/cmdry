// Package timediff calculates the elapsed time between two user-supplied dates.
package timediff

import (
	"fmt"
	"strings"
	"time"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func Run() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.time-between-dates", Name: "Time Between Dates", Version: "0.1.0", Description: "Calculate the elapsed time between two dates locally.", Category: "time", Icon: "clock",
		SearchTerms: []string{"time", "between", "dates", "duration", "difference", "days"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "calculate", Name: "Calculate time", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": form, "calculate": calculate}})
}

func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Time Between Dates", Components: []cmdry.Component{{Type: "form", Title: "Calculate elapsed time", Action: "calculate", Submit: "Calculate time", Description: "Use an ISO date (2026-08-16), a local date/time (2026-08-16 14:30:00), or an RFC 3339 timestamp with timezone.", Fields: []cmdry.Field{{Name: "start", Label: "Start date/time", Type: "text", Required: true, Placeholder: "2026-08-16 09:00:00"}, {Name: "end", Label: "End date/time", Type: "text", Required: true, Placeholder: "2026-08-18 17:30:00"}}}}}, nil
}

func calculate(request cmdry.Request) (cmdry.View, error) {
	difference, err := Between(fmt.Sprint(request.Params["start"]), fmt.Sprint(request.Params["end"]))
	if err != nil {
		return cmdry.View{}, err
	}
	direction := "End is after start"
	if difference.EndBeforeStart {
		direction = "End is before start"
	}
	return cmdry.View{Title: "Time between dates", Components: []cmdry.Component{{Type: "alert", Level: "info", Title: direction, Message: "The displayed duration is absolute."}, {Type: "metric", Label: "Total days", Value: fmt.Sprint(difference.Days)}, {Type: "metric", Label: "Total hours", Value: fmt.Sprint(difference.Hours)}, {Type: "metric", Label: "Total minutes", Value: fmt.Sprint(difference.Minutes)}, {Type: "metric", Label: "Total seconds", Value: fmt.Sprint(difference.Seconds)}, {Type: "code", Title: "Exact duration", Text: difference.Duration.String()}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Compare other dates"}}}}}, nil
}

// Difference holds an absolute duration and its convenient whole units.
type Difference struct {
	Duration       time.Duration
	Days           int64
	Hours          int64
	Minutes        int64
	Seconds        int64
	EndBeforeStart bool
}

// Between parses two supported date formats and calculates their elapsed time.
func Between(start, end string) (Difference, error) {
	startTime, err := parseDate(start)
	if err != nil {
		return Difference{}, fmt.Errorf("start date: %w", err)
	}
	endTime, err := parseDate(end)
	if err != nil {
		return Difference{}, fmt.Errorf("end date: %w", err)
	}
	duration := endTime.Sub(startTime)
	before := duration < 0
	if before {
		duration = -duration
	}
	return Difference{Duration: duration, Days: int64(duration / (24 * time.Hour)), Hours: int64(duration / time.Hour), Minutes: int64(duration / time.Minute), Seconds: int64(duration / time.Second), EndBeforeStart: before}, nil
}

func parseDate(input string) (time.Time, error) {
	value := strings.TrimSpace(input)
	if value == "" {
		return time.Time{}, fmt.Errorf("is required")
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"} {
		if layout == time.RFC3339 {
			if parsed, err := time.Parse(layout, value); err == nil {
				return parsed, nil
			}
			continue
		}
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("use YYYY-MM-DD, YYYY-MM-DD HH:MM:SS, or an RFC 3339 timestamp")
}
