package main

import (
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.truncate-list", Name: "Truncate List", Version: "0.1.0", Description: "Keep the first specified number of list items locally.", Category: "text", Icon: "list",
		SearchTerms: []string{"truncate", "list", "limit", "first", "lines"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "truncate", Name: "Truncate list", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": form, "truncate": truncate}})
}

func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Truncate List", Components: []cmdry.Component{{Type: "form", Title: "Keep the first items", Action: "truncate", Submit: "Truncate list", Description: "Blank lines are ignored; all retained item text stays unchanged.", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}, {Name: "limit", Label: "Items to keep", Type: "number", Value: "10", Min: "1", Required: true}}}}}, nil
}

func truncate(request cmdry.Request) (cmdry.View, error) {
	limit, err := strconv.Atoi(fmt.Sprint(request.Params["limit"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("item limit must be a whole number")
	}
	items, err := listtools.TruncateLines(fmt.Sprint(request.Params["input"]), limit)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Truncated list", Components: []cmdry.Component{{Type: "metric", Label: "Items kept", Value: fmt.Sprint(len(items))}, {Type: "code", Title: "Truncated list", Text: strings.Join(items, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Truncate another list"}}}}}, nil
}
