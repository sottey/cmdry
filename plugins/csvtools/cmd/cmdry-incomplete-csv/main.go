package main

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/csvtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.incomplete-csv", Name: "Find Incomplete CSV Records", Version: "0.1.0", Description: "Find pasted comma-separated CSV rows with empty values locally.", Category: "data", Icon: "table",
		SearchTerms: []string{"csv", "incomplete", "missing", "empty", "records", "rows"}, Pages: []cmdry.Page{{ID: "overview", Name: "Check", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New CSV check", Method: "read"}, {ID: "check", Name: "Find incomplete records", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "check": check}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Find Incomplete CSV Records", Components: []cmdry.Component{{Type: "form", Title: "Check CSV completeness", Action: "check", Submit: "Find incomplete records", Description: "Checks comma-separated CSV locally. The first row supplies column names; only empty fields are reported as missing.", Fields: []cmdry.Field{{Name: "input", Label: "CSV input", Type: "textarea", Required: true}}}}}, nil
}

func check(request cmdry.Request) (cmdry.View, error) {
	records, rows, err := csvtools.FindIncompleteRecords(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	components := []cmdry.Component{{Type: "metric", Label: "Data rows checked", Value: fmt.Sprint(rows)}, {Type: "metric", Label: "Incomplete rows", Value: fmt.Sprint(len(records))}}
	if len(records) == 0 {
		components = append(components, cmdry.Component{Type: "alert", Level: "success", Title: "No incomplete records", Message: "Every data row has a value for each header column."})
	} else {
		rows := make([]map[string]any, 0, len(records))
		for _, record := range records {
			rows = append(rows, map[string]any{"row": record.Row, "columns": strings.Join(record.Columns, ", ")})
		}
		components = append(components, cmdry.Component{Type: "alert", Level: "info", Title: "Incomplete records found", Message: "Each listed row has one or more empty values."}, cmdry.Component{Type: "table", ID: "incomplete-csv", Columns: []cmdry.Column{{Key: "row", Label: "CSV row"}, {Key: "columns", Label: "Missing columns"}}, Rows: rows})
	}
	components = append(components, cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Check another CSV"}}})
	return cmdry.View{Title: "CSV completeness", Components: components}, nil
}
