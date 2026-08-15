package main

import (
	"fmt"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/cronexpr"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.cron-explainer", Name: "Cron Expression Explainer", Version: "0.1.0", Description: "Validate and explain standard cron schedules locally.", Category: "time", Icon: "clock",
		SearchTerms: []string{"cron", "crontab", "schedule", "timer", "expression", "guru"}, Pages: []cmdry.Page{{ID: "overview", Name: "Explain", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New expression", Method: "read"}, {ID: "explain", Name: "Explain cron", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "explain": explain}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Cron Expression Explainer", Components: []cmdry.Component{{Type: "form", Title: "Explain cron schedule", Action: "explain", Submit: "Explain cron", Description: "Supports standard five-field crontab expressions, ranges, lists, steps, month/day names, and @daily-style shortcuts. It does not support Quartz extensions.", Fields: []cmdry.Field{{Name: "input", Label: "Cron expression", Type: "text", Required: true, Placeholder: "*/15 9-17 * * MON-FRI"}}}}}, nil
}

func explain(request cmdry.Request) (cmdry.View, error) {
	fields, description, err := cronexpr.Explain(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("invalid cron expression: %w", err)
	}
	components := []cmdry.Component{{Type: "alert", Level: "success", Title: "Cron expression is valid", Message: description}}
	if len(fields) > 0 {
		rows := make([]map[string]any, 0, len(fields))
		for _, field := range fields {
			rows = append(rows, map[string]any{"field": field.Name, "expression": field.Expression, "meaning": field.Explanation})
		}
		components = append(components, cmdry.Component{Type: "table", ID: "cron-explanation", Columns: []cmdry.Column{{Key: "field", Label: "Field"}, {Key: "expression", Label: "Expression"}, {Key: "meaning", Label: "Meaning"}}, Rows: rows})
	}
	components = append(components, cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Explain another expression"}}})
	return cmdry.View{Title: "Cron explanation", Components: components}, nil
}
