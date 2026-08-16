package main

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.reverse-list", Name: "Reverse List", Version: "0.1.0", Description: "Reverse pasted newline-delimited list items locally.", Category: "text", Icon: "list",
		SearchTerms: []string{"list", "reverse", "lines", "order"}, Pages: []cmdry.Page{{ID: "overview", Name: "Reverse", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New list", Method: "read"}, {ID: "reverse", Name: "Reverse list", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "reverse": reverse}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Reverse List", Components: []cmdry.Component{{Type: "form", Title: "Reverse list order", Action: "reverse", Submit: "Reverse list", Description: "Uses one item per line, ignores blank lines, and reverses the remaining item order locally.", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}}}}}, nil
}

func reverse(request cmdry.Request) (cmdry.View, error) {
	items, err := listtools.ReverseLines(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Reversed list", Components: []cmdry.Component{{Type: "metric", Label: "Items reversed", Value: fmt.Sprint(len(items))}, {Type: "code", Title: "Reversed list", Text: strings.Join(items, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Reverse another list"}}}}}, nil
}
