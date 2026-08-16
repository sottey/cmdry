package main

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.shuffle-list", Name: "Shuffle List", Version: "0.1.0", Description: "Randomize pasted list item order locally.", Category: "text", Icon: "list", SearchTerms: []string{"list", "shuffle", "randomize", "order"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New list", Method: "read"}, {ID: "run", Name: "Shuffle list", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Shuffle List", Components: []cmdry.Component{{Type: "form", Title: "Randomize list order", Action: "run", Submit: "Shuffle list", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}}}}}, nil
}
func run(request cmdry.Request) (cmdry.View, error) {
	items, err := listtools.ShuffleLines(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Shuffled list", Components: []cmdry.Component{{Type: "metric", Label: "Items", Value: fmt.Sprint(len(items))}, {Type: "code", Title: "Shuffled list", Text: strings.Join(items, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Shuffle another list"}}}}}, nil
}
