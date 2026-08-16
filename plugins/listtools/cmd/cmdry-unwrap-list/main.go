package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
	"strings"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.unwrap-list", Name: "Unwrap List", Version: "0.1.0", Description: "Remove a matching prefix and suffix from pasted list items locally.", Category: "text", Icon: "list", SearchTerms: []string{"list", "unwrap", "prefix", "suffix", "remove"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New list", Method: "read"}, {ID: "run", Name: "Unwrap list", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Unwrap List", Components: []cmdry.Component{{Type: "form", Title: "Remove matching wrappers", Action: "run", Submit: "Unwrap list", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}, {Name: "prefix", Label: "Prefix to remove", Type: "text"}, {Name: "suffix", Label: "Suffix to remove", Type: "text"}}}}}, nil
}
func run(request cmdry.Request) (cmdry.View, error) {
	items, err := listtools.UnwrapLines(fmt.Sprint(request.Params["input"]), fmt.Sprint(request.Params["prefix"]), fmt.Sprint(request.Params["suffix"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Unwrapped list", Components: []cmdry.Component{{Type: "metric", Label: "Items", Value: fmt.Sprint(len(items))}, {Type: "code", Title: "Unwrapped list", Text: strings.Join(items, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Unwrap another list"}}}}}, nil
}
