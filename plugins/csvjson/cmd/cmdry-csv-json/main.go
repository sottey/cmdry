package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/csvjson"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.csv-json", Name: "CSV to JSON", Version: "0.1.0",
		Description: "Convert pasted header-based CSV data to a local JSON download.", Category: "data", Icon: "convert",
		SearchTerms: []string{"csv", "json", "convert", "spreadsheet"},
		Pages:       []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}},
		Permissions: []string{"data.transform"},
		Actions:     []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert to JSON", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "convert": convert}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "CSV to JSON", Components: []cmdry.Component{{
		Type: "form", Title: "Paste CSV", Action: "convert", Submit: "Convert to JSON",
		Description: "The first row supplies JSON property names. Values stay as strings so source data is preserved exactly.",
		Fields:      []cmdry.Field{{Name: "csv", Label: "CSV input", Type: "textarea", Placeholder: "name,active\nAda,true", Description: "Your data is converted locally by the plugin and is never uploaded.", Required: true}},
	}}}, nil
}

func convert(request cmdry.Request) (cmdry.View, error) {
	input, _ := request.Params["csv"].(string)
	if strings.TrimSpace(input) == "" {
		return cmdry.View{}, fmt.Errorf("CSV input is required")
	}
	result, err := csvjson.Convert(input)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "JSON ready", Components: []cmdry.Component{
		{Type: "alert", Level: "success", Title: strconv.Itoa(len(result.Rows)) + " rows converted", Message: "JSON is ready to inspect or download."},
		{Type: "metric", Label: "Columns", Value: strconv.Itoa(len(result.Headers))},
		{Type: "download", Filename: "cmdry-export.json", MIMEType: "application/json;charset=utf-8", Content: base64.StdEncoding.EncodeToString(result.JSON)},
		{Type: "code", Title: "JSON preview", Text: string(result.JSON)},
		{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another"}}},
	}}, nil
}
