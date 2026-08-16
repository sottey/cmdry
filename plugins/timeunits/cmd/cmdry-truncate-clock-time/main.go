package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/timeunits"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.truncate-clock-time", Name: "Truncate Clock Time", Version: "0.1.0", Description: "Truncate a clock duration to hours, minutes, or seconds locally.", Category: "time", Icon: "clock", SearchTerms: []string{"time", "clock", "truncate", "duration", "round"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New time", Method: "read"}, {ID: "run", Name: "Truncate time", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Truncate Clock Time", Components: []cmdry.Component{{Type: "form", Title: "Drop smaller time units", Action: "run", Submit: "Truncate time", Description: "Accepts signed HH:MM:SS or MM:SS. Truncation moves toward zero.", Fields: []cmdry.Field{{Name: "input", Label: "Clock time", Type: "text", Value: "01:23:45", Required: true}, {Name: "precision", Label: "Keep through", Type: "select", Value: "minute", Options: []cmdry.Option{{Value: "hour", Label: "Hours"}, {Value: "minute", Label: "Minutes"}, {Value: "second", Label: "Seconds"}}}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	seconds, err := timeunits.ParseClockDuration(fmt.Sprint(r.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	result, err := timeunits.TruncateClockTime(seconds, fmt.Sprint(r.Params["precision"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Truncated clock time", Components: []cmdry.Component{{Type: "metric", Label: "Result", Value: timeunits.FormatSeconds(result)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Truncate another time"}}}}}, nil
}
