package main

import (
	"fmt"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/csvtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.csv-tsv", Name: "CSV to TSV", Version: "0.1.0", Description: "Convert pasted comma-separated CSV data to tab-separated values locally.", Category: "data", Icon: "convert",
		SearchTerms: []string{"csv", "tsv", "tab", "convert", "delimiter", "spreadsheet"}, Pages: []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert to TSV", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "convert": convert}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "CSV to TSV", Components: []cmdry.Component{{Type: "form", Title: "Convert CSV to TSV", Action: "convert", Submit: "Convert to TSV", Description: "Parses quoted CSV fields correctly and converts locally. Output uses tab separators and standard newline endings.", Fields: []cmdry.Field{{Name: "input", Label: "CSV input", Type: "textarea", Required: true}}}}}, nil
}

func convert(request cmdry.Request) (cmdry.View, error) {
	output, rows, err := csvtools.ChangeSeparator(fmt.Sprint(request.Params["input"]), "comma", "tab")
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Converted TSV", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "CSV converted to TSV", Message: fmt.Sprint(rows) + " row(s) converted locally."}, {Type: "code", Title: "TSV output", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another CSV"}}}}}, nil
}
