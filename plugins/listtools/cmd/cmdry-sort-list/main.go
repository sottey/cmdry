package main

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.sort-list", Name: "Sort List", Version: "0.1.0", Description: "Sort newline-delimited list items locally.", Category: "text", Icon: "list",
		SearchTerms: []string{"sort", "list", "alphabetical", "ascending", "descending", "lines"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "sort", Name: "Sort list", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": form, "sort": sortList}})
}

func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Sort List", Components: []cmdry.Component{{Type: "form", Title: "Sort list items", Action: "sort", Submit: "Sort list", Description: "Sorts non-blank newline-delimited items locally.", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}, {Name: "order", Label: "Order", Type: "select", Value: "ascending", Options: []cmdry.Option{{Value: "ascending", Label: "A to Z"}, {Value: "descending", Label: "Z to A"}}}}}}}, nil
}

func sortList(request cmdry.Request) (cmdry.View, error) {
	items, err := listtools.SortLines(fmt.Sprint(request.Params["input"]), fmt.Sprint(request.Params["order"]) == "descending")
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Sorted list", Components: []cmdry.Component{{Type: "metric", Label: "Items", Value: fmt.Sprint(len(items))}, {Type: "code", Title: "Sorted list", Text: strings.Join(items, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Sort another list"}}}}}, nil
}
