package main

import (
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
)

func main() {
	cmdry.Run(plugin("rotate-list", "Rotate List", "Rotate pasted list items by a chosen position locally.", "positions", "1", rotate))
}
func plugin(id, name, description, field, value string, handler cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "text", Icon: "list", SearchTerms: []string{"list", "rotate", "shift", "order"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": func(cmdry.Request) (cmdry.View, error) {
		return cmdry.View{Title: name, Components: []cmdry.Component{{Type: "form", Title: name, Action: "run", Submit: "Rotate list", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}, {Name: field, Label: "Positions to rotate left", Type: "number", Value: value, Required: true}}}}}, nil
	}, "run": handler}}
}
func rotate(request cmdry.Request) (cmdry.View, error) {
	positions, err := strconv.Atoi(fmt.Sprint(request.Params["positions"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("positions must be a whole number")
	}
	items, err := listtools.RotateLines(fmt.Sprint(request.Params["input"]), positions)
	if err != nil {
		return cmdry.View{}, err
	}
	return output("Rotated list", items), nil
}
func output(title string, items []string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "metric", Label: "Items", Value: fmt.Sprint(len(items))}, {Type: "code", Title: title, Text: strings.Join(items, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Start over"}}}}}
}
