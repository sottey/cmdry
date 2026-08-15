package timeunits

import (
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

// RunSecondsToTime serves the seconds-to-clock-time plugin.
func RunSecondsToTime() {
	runDuration("seconds-to-time", "Convert Seconds to Time", "Convert whole seconds to signed HH:MM:SS locally.", []string{"seconds", "time", "duration", "clock", "hh:mm:ss"}, "Seconds", "3661", "Enter a whole number of seconds. Negative offsets are supported.", secondsToTime)
}

// RunTimeToDecimal serves the clock-time-to-decimal-hours plugin.
func RunTimeToDecimal() {
	runDuration("time-to-decimal", "Convert Time to Decimal", "Convert signed clock time to decimal hours locally.", []string{"time", "decimal", "hours", "duration", "hh:mm:ss"}, "Time", "01:30:00", "Use HH:MM:SS or MM:SS. Negative offsets are supported.", timeToDecimal)
}

// RunTimeToSeconds serves the clock-time-to-seconds plugin.
func RunTimeToSeconds() {
	runDuration("time-to-seconds", "Convert Time to Seconds", "Convert signed clock time to total seconds locally.", []string{"time", "seconds", "duration", "clock", "hh:mm:ss"}, "Time", "01:30:00", "Use HH:MM:SS or MM:SS. Negative offsets are supported.", timeToSeconds)
}

func runDuration(id, name, description string, terms []string, label, placeholder, hint string, action cmdry.Handler) {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "time", Icon: "clock",
		SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": durationForm(name, label, placeholder, hint), "convert": action}})
}

func durationForm(name, label, placeholder, hint string) cmdry.Handler {
	return func(cmdry.Request) (cmdry.View, error) {
		return cmdry.View{Title: name, Components: []cmdry.Component{{Type: "form", Title: name, Action: "convert", Submit: "Convert", Description: hint, Fields: []cmdry.Field{{Name: "input", Label: label, Type: "text", Placeholder: placeholder, Required: true}}}}}, nil
	}
}

func secondsToTime(request cmdry.Request) (cmdry.View, error) {
	seconds, err := strconv.ParseInt(strings.TrimSpace(fmt.Sprint(request.Params["input"])), 10, 64)
	if err != nil {
		return cmdry.View{}, fmt.Errorf("seconds must be a whole number")
	}
	return durationOutput("Time", FormatSeconds(seconds)), nil
}

func timeToDecimal(request cmdry.Request) (cmdry.View, error) {
	seconds, err := ParseClockDuration(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return durationOutput("Decimal hours", strconv.FormatFloat(DecimalHours(seconds), 'g', 12, 64)), nil
}

func timeToSeconds(request cmdry.Request) (cmdry.View, error) {
	seconds, err := ParseClockDuration(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return durationOutput("Seconds", strconv.FormatInt(seconds, 10)), nil
}

func durationOutput(title, value string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "code", Title: title, Text: value}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another value"}}}}}
}
