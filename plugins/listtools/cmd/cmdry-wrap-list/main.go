package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
	"strings"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.wrap-list", Name: "Wrap List", Version: "0.1.0", Description: "Add a prefix and suffix around each pasted list item locally.", Category: "text", Icon: "list", SearchTerms: []string{"list", "wrap", "prefix", "suffix"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New list", Method: "read"}, {ID: "run", Name: "Wrap list", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Wrap List", Components: []cmdry.Component{{Type: "form", Title: "Wrap each list item", Action: "run", Submit: "Wrap list", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}, {Name: "prefix", Label: "Prefix", Type: "text"}, {Name: "suffix", Label: "Suffix", Type: "text"}}}}}, nil
}
func run(request cmdry.Request) (cmdry.View, error) {
	items, err := listtools.WrapLines(fmt.Sprint(request.Params["input"]), fmt.Sprint(request.Params["prefix"]), fmt.Sprint(request.Params["suffix"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Wrapped list", Components: []cmdry.Component{{Type: "metric", Label: "Items", Value: fmt.Sprint(len(items))}, {Type: "code", Title: "Wrapped list", Text: strings.Join(items, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Wrap another list"}}}}}, nil
}
