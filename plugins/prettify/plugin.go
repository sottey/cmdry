package prettify

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"strconv"
	"strings"
)

func Run(kind, name string) {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + kind + "-prettifier", Name: name, Version: "0.1.0", Description: "Format pasted " + kind + " locally.", Category: "developer", Pages: []cmdry.Page{{ID: "overview", Name: "Format", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New formatting", Method: "read"}, {ID: "format", Name: "Format", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": func(cmdry.Request) (cmdry.View, error) { return form(kind, name), nil }, "format": func(r cmdry.Request) (cmdry.View, error) { return formatView(kind, name, r) }}})
}
func form(kind, name string) cmdry.View {
	return cmdry.View{Title: name, Components: []cmdry.Component{{Type: "form", Title: "Paste " + strings.ToUpper(kind), Action: "format", Submit: "Format", Description: "Choose indentation, then format locally in this plugin.", Fields: []cmdry.Field{{Name: "input", Label: "Input", Type: "textarea", Required: true}, {Name: "indent", Label: "Indent spaces", Type: "number", Value: "4", Min: "1", Max: "8", Required: true}}}}}
}
func formatView(kind, name string, r cmdry.Request) (cmdry.View, error) {
	input, _ := r.Params["input"].(string)
	spaces, err := strconv.Atoi(fmt.Sprint(r.Params["indent"]))
	if strings.TrimSpace(input) == "" {
		return cmdry.View{}, fmt.Errorf("input is required")
	}
	if err != nil {
		return cmdry.View{}, fmt.Errorf("indent must be a number")
	}
	output, err := Format(kind, input, spaces)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: name + " result", Components: []cmdry.Component{{Type: "code", Title: "Formatted output", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Format another"}}}}}, nil
}
