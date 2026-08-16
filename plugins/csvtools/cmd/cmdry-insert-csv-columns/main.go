package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/csvtools"
	"strconv"
	"strings"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.insert-csv-columns", Name: "Insert CSV Columns", Version: "0.1.0", Description: "Insert blank columns into a CSV table locally.", Category: "data", Icon: "table", SearchTerms: []string{"csv", "insert", "columns", "table"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Insert columns", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": form, "run": run}})
}
func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Insert CSV Columns", Components: []cmdry.Component{{Type: "form", Title: "Add columns", Action: "run", Submit: "Insert columns", Fields: []cmdry.Field{{Name: "input", Label: "CSV input", Type: "textarea", Required: true}, {Name: "position", Label: "Position (1 is first)", Type: "number", Value: "1", Min: "1", Required: true}, {Name: "names", Label: "Column names (comma-separated)", Type: "text", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	p, e := strconv.Atoi(fmt.Sprint(r.Params["position"]))
	if e != nil {
		return cmdry.View{}, fmt.Errorf("position must be a whole number")
	}
	out, e := csvtools.InsertColumns(fmt.Sprint(r.Params["input"]), p, strings.Split(fmt.Sprint(r.Params["names"]), ","))
	if e != nil {
		return cmdry.View{}, e
	}
	return cmdry.View{Title: "CSV ready", Components: []cmdry.Component{{Type: "code", Title: "CSV", Text: out}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Insert more columns"}}}}}, nil
}
