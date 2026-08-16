package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/csvtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.transpose-csv", Name: "Transpose CSV", Version: "0.1.0", Description: "Transpose pasted header-based CSV locally.", Category: "data", Icon: "table", SearchTerms: []string{"csv", "transpose", "rows", "columns"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New CSV", Method: "read"}, {ID: "run", Name: "Transpose CSV", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Transpose CSV", Components: []cmdry.Component{{Type: "form", Title: "Transpose CSV rows and columns", Action: "run", Submit: "Transpose CSV", Fields: []cmdry.Field{{Name: "input", Label: "CSV", Type: "textarea", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	output, columns, rows, err := csvtools.Transpose(fmt.Sprint(r.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Transposed CSV", Components: []cmdry.Component{{Type: "metric", Label: "Original columns", Value: fmt.Sprint(columns)}, {Type: "metric", Label: "Original data rows", Value: fmt.Sprint(rows)}, {Type: "code", Title: "CSV", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Transpose another CSV"}}}}}, nil
}
