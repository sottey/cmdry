package main

import (
	"fmt"
	"time"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/timeconvert"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.unix-date", Name: "Unix Date Converter", Version: "0.1.0",
		Description: "Convert Unix timestamps and date-times locally.", Category: "time", Icon: "clock",
		SearchTerms: []string{"unix", "timestamp", "epoch", "date", "time", "utc"},
		Pages:       []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}},
		Permissions: []string{"data.transform"},
		Actions:     []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "convert": convert}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Unix Date Converter", Components: []cmdry.Component{{
		Type: "form", Title: "Convert date and time", Action: "convert", Submit: "Convert",
		Description: "Unix timestamps are whole seconds. Dates accept RFC 3339 or YYYY-MM-DD HH:MM:SS; choose how dates without an offset should be interpreted.",
		Fields: []cmdry.Field{
			{Name: "mode", Label: "Convert", Type: "select", Value: "unix-to-date", Options: []cmdry.Option{{Value: "unix-to-date", Label: "Unix timestamp to date"}, {Value: "date-to-unix", Label: "Date to Unix timestamp"}}},
			{Name: "input", Label: "Timestamp or date", Type: "text", Required: true, Placeholder: "1704164645 or 2024-01-02T03:04:05Z"},
			{Name: "zone", Label: "Date without offset", Type: "select", Value: "utc", Options: []cmdry.Option{{Value: "utc", Label: "UTC"}, {Value: "local", Label: "Local time"}}},
		},
	}}}, nil
}

func convert(request cmdry.Request) (cmdry.View, error) {
	moment, err := timeconvertMoment(request)
	if err != nil {
		return cmdry.View{}, err
	}
	details := timeconvert.Details(moment)
	rows := []map[string]any{{"format": "Unix timestamp (seconds)", "value": details["unix"]}, {"format": "UTC (RFC 3339)", "value": details["utc"]}, {"format": "Local time", "value": details["local"]}}
	return cmdry.View{Title: "Converted date", Components: []cmdry.Component{{Type: "table", ID: "unix-date", Columns: []cmdry.Column{{Key: "format", Label: "Format"}, {Key: "value", Label: "Value"}}, Rows: rows}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another value"}}}}}, nil
}

func timeconvertMoment(request cmdry.Request) (time.Time, error) {
	if fmt.Sprint(request.Params["mode"]) == "date-to-unix" {
		return timeconvert.ToUnixSeconds(fmt.Sprint(request.Params["input"]), fmt.Sprint(request.Params["zone"]))
	}
	return timeconvert.FromUnixSeconds(fmt.Sprint(request.Params["input"]))
}
