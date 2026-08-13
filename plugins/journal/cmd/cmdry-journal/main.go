package main

import (
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/journal"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1,
		ID:              "journal",
		Name:            "Journal Viewer",
		Version:         "0.1.0",
		Description:     "Review the 100 most recent local Linux journal entries.",
		Category:        "system",
		Icon:            "journal",
		Pages:           []cmdry.Page{{ID: "overview", Name: "Journal", Default: true, Action: "list"}},
		Permissions:     []string{"system.journal.read"},
		Actions:         []cmdry.Action{{ID: "list", Name: "Refresh", Method: "read"}},
	}, Actions: map[string]cmdry.Handler{"list": listEntries}})
}

func listEntries(_ cmdry.Request) (cmdry.View, error) {
	entries, err := journal.CollectRecent()
	if err != nil {
		return cmdry.View{}, err
	}
	errors, warnings := journal.Summary(entries)
	rows := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		row := entry.Row()
		row["message"] = journal.ShortMessage(entry.Message)
		rows = append(rows, row)
	}
	return cmdry.View{Title: "Recent Journal Entries", Components: []cmdry.Component{
		{Type: "metric", Label: "Entries", Value: strconv.Itoa(len(entries)), Description: "Newest 100 entries"},
		{Type: "metric", Label: "Errors", Value: strconv.Itoa(errors), Description: "Priority error or higher"},
		{Type: "metric", Label: "Warnings", Value: strconv.Itoa(warnings), Description: "Priority warning"},
		{Type: "table", ID: "journal", Columns: []cmdry.Column{
			{Key: "time", Label: "Time"},
			{Key: "priority", Label: "Priority"},
			{Key: "unit", Label: "Unit"},
			{Key: "message", Label: "Message"},
		}, Rows: rows},
	}}, nil
}
