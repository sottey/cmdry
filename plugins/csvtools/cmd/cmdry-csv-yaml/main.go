package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/csvtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.csv-yaml", Name: "CSV to YAML", Version: "0.1.0", Description: "Convert pasted header-based CSV to YAML locally.", Category: "data", Icon: "convert", SearchTerms: []string{"csv", "yaml", "yml", "convert", "table", "spreadsheet"}, Pages: []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert to YAML", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": form, "convert": convert}})
}
func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "CSV to YAML", Components: []cmdry.Component{{Type: "form", Title: "Paste CSV", Action: "convert", Submit: "Convert to YAML", Description: "Each data row becomes a YAML mapping. Cell values remain strings.", Fields: []cmdry.Field{{Name: "input", Label: "CSV input", Type: "textarea", Required: true}}}}}, nil
}
func convert(r cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(r.Params["input"])
	if strings.TrimSpace(input) == "" {
		return cmdry.View{}, fmt.Errorf("CSV input is required")
	}
	output, rows, err := csvtools.CSVToYAML(input)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "YAML ready", Components: []cmdry.Component{{Type: "metric", Label: "Rows converted", Value: fmt.Sprint(rows)}, {Type: "download", Filename: "cmdry-export.yaml", MIMEType: "application/yaml;charset=utf-8", Content: base64.StdEncoding.EncodeToString([]byte(output))}, {Type: "code", Title: "YAML preview", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another"}}}}}, nil
}
