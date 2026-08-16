package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/csvtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.swap-csv-columns", Name: "Swap CSV Columns", Version: "0.1.0", Description: "Swap two named columns in pasted CSV locally.", Category: "data", Icon: "table", SearchTerms: []string{"csv", "swap", "columns", "reorder"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New CSV", Method: "read"}, {ID: "run", Name: "Swap columns", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Swap CSV Columns", Components: []cmdry.Component{{Type: "form", Title: "Swap two CSV columns", Action: "run", Submit: "Swap columns", Fields: []cmdry.Field{{Name: "input", Label: "CSV", Type: "textarea", Required: true}, {Name: "first", Label: "First header", Type: "text", Required: true}, {Name: "second", Label: "Second header", Type: "text", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	output, err := csvtools.SwapColumns(fmt.Sprint(r.Params["input"]), fmt.Sprint(r.Params["first"]), fmt.Sprint(r.Params["second"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Swapped CSV", Components: []cmdry.Component{{Type: "code", Title: "CSV", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Swap another CSV"}}}}}, nil
}
