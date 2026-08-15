package main

import (
	"fmt"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/csvtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.csv-separator", Name: "Change CSV Separator", Version: "0.1.0", Description: "Convert a pasted CSV delimiter locally while preserving quoted fields.", Category: "data", Icon: "convert",
		SearchTerms: []string{"csv", "separator", "delimiter", "comma", "semicolon", "tab", "tsv"}, Pages: []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Change separator", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "convert": convert}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	options := []cmdry.Option{{Value: "comma", Label: "Comma (, )"}, {Value: "semicolon", Label: "Semicolon (; )"}, {Value: "tab", Label: "Tab"}, {Value: "pipe", Label: "Pipe (|)"}}
	return cmdry.View{Title: "Change CSV Separator", Components: []cmdry.Component{{Type: "form", Title: "Convert CSV delimiter", Action: "convert", Submit: "Change separator", Description: "Parses quoted CSV fields correctly and converts locally. Output always uses standard newline endings.", Fields: []cmdry.Field{{Name: "input", Label: "CSV input", Type: "textarea", Required: true}, {Name: "source", Label: "Current separator", Type: "select", Value: "comma", Options: options}, {Name: "destination", Label: "New separator", Type: "select", Value: "semicolon", Options: options}}}}}, nil
}

func convert(request cmdry.Request) (cmdry.View, error) {
	output, rows, err := csvtools.ChangeSeparator(fmt.Sprint(request.Params["input"]), fmt.Sprint(request.Params["source"]), fmt.Sprint(request.Params["destination"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Converted CSV", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Separator changed", Message: fmt.Sprint(rows) + " row(s) converted locally."}, {Type: "code", Title: "Converted CSV", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another CSV"}}}}}, nil
}
