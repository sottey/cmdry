package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
	"strconv"
	"strings"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.duplicate-list", Name: "Duplicate List", Version: "0.1.0", Description: "Duplicate newline-delimited list items locally.", Category: "text", Icon: "list", SearchTerms: []string{"duplicate", "list", "repeat", "lines"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Duplicate items", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": form, "run": run}})
}
func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Duplicate List", Components: []cmdry.Component{{Type: "form", Title: "Repeat list items", Action: "run", Submit: "Duplicate items", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}, {Name: "copies", Label: "Copies per item", Type: "number", Value: "2", Min: "1", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	copies, err := strconv.Atoi(fmt.Sprint(r.Params["copies"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("copies must be a whole number")
	}
	items, err := listtools.DuplicateLines(fmt.Sprint(r.Params["input"]), copies)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Duplicated list", Components: []cmdry.Component{{Type: "metric", Label: "Output items", Value: fmt.Sprint(len(items))}, {Type: "code", Title: "Duplicated list", Text: strings.Join(items, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Duplicate another list"}}}}}, nil
}
