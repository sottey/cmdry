package main

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.unique-list", Name: "Find Unique", Version: "0.1.0", Description: "Keep the first occurrence of each pasted newline-delimited item locally.", Category: "text", Icon: "list",
		SearchTerms: []string{"list", "unique", "duplicates", "deduplicate", "lines"}, Pages: []cmdry.Page{{ID: "overview", Name: "Find", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New list", Method: "read"}, {ID: "find", Name: "Find unique items", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "find": find}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Find Unique", Components: []cmdry.Component{{Type: "form", Title: "Keep unique items", Action: "find", Submit: "Find unique items", Description: "Uses one item per line, ignores blank lines, and preserves the first exact occurrence locally.", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}}}}}, nil
}

func find(request cmdry.Request) (cmdry.View, error) {
	items, err := listtools.UniqueLines(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Unique list items", Components: []cmdry.Component{{Type: "metric", Label: "Unique items", Value: fmt.Sprint(len(items))}, {Type: "code", Title: "Unique items", Text: strings.Join(items, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Find another list"}}}}}, nil
}
