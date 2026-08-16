package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/datetools"
	"time"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.discord-timestamp", Name: "Discord Timestamp Generator", Version: "0.1.0", Description: "Generate Discord timestamp markup from a date locally.", Category: "time", Icon: "clock", SearchTerms: []string{"discord", "timestamp", "date", "relative time", "markup"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New timestamp", Method: "read"}, {ID: "run", Name: "Generate timestamp", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": form, "run": run}})
}
func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Discord Timestamp Generator", Components: []cmdry.Component{{Type: "form", Title: "Generate timestamp", Action: "run", Submit: "Generate", Description: "Enter an RFC 3339 timestamp, such as 2026-08-16T20:00:00Z.", Fields: []cmdry.Field{{Name: "input", Label: "Date/time", Type: "text", Required: true, Placeholder: "2026-08-16T20:00:00Z"}, {Name: "style", Label: "Display style", Type: "select", Value: "relative", Options: []cmdry.Option{{Value: "relative", Label: "Relative"}, {Value: "long-date-time", Label: "Long date/time"}, {Value: "short-date-time", Label: "Short date/time"}, {Value: "long-date", Label: "Long date"}, {Value: "short-date", Label: "Short date"}, {Value: "long-time", Label: "Long time"}, {Value: "short-time", Label: "Short time"}}}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	value, err := time.Parse(time.RFC3339, fmt.Sprint(r.Params["input"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("use an RFC 3339 timestamp: %w", err)
	}
	output, err := datetools.DiscordTimestamp(value, fmt.Sprint(r.Params["style"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Discord timestamp", Components: []cmdry.Component{{Type: "code", Title: "Paste into Discord", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Generate another"}}}}}, nil
}
