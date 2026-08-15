package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/jsoncsv"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1,
		ID:              "com.sottey.jsoncsv",
		Name:            "JSON to CSV",
		Version:         "0.1.0",
		Description:     "Convert pasted JSON to a local CSV download.",
		Category:        "data",
		Icon:            "convert",
		SearchTerms:     []string{"json", "csv", "convert", "spreadsheet"},
		Pages:           []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}},
		Permissions:     []string{"data.transform"},
		Actions: []cmdry.Action{
			{ID: "overview", Name: "New conversion", Method: "read"},
			{ID: "convert", Name: "Convert to CSV", Method: "write"},
		},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "convert": convert}})
}

func overview(_ cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "JSON to CSV", Components: []cmdry.Component{{
		Type: "form", Title: "Paste JSON", Action: "convert", Submit: "Convert to CSV",
		Description: "Accepts one object or an array of objects. Nested object fields become dot-path columns; arrays stay together as JSON in one cell.",
		Fields:      []cmdry.Field{{Name: "json", Label: "JSON input", Type: "textarea", Placeholder: "[{\"name\": \"Ada\", \"active\": true}]", Description: "Your data is converted locally by the plugin and is never uploaded.", Required: true}},
	}}}, nil
}

func convert(request cmdry.Request) (cmdry.View, error) {
	input, _ := request.Params["json"].(string)
	if strings.TrimSpace(input) == "" {
		return cmdry.View{}, fmt.Errorf("JSON input is required")
	}
	result, err := jsoncsv.Convert(input)
	if err != nil {
		return cmdry.View{}, err
	}
	previewSize := len(result.Rows)
	if previewSize > 25 {
		previewSize = 25
	}
	rows := make([]map[string]any, 0, previewSize)
	for _, row := range result.Rows[:previewSize] {
		item := make(map[string]any, len(row))
		for key, value := range row {
			item[key] = value
		}
		rows = append(rows, item)
	}
	columns := make([]cmdry.Column, 0, len(result.Columns))
	for _, column := range result.Columns {
		columns = append(columns, cmdry.Column{Key: column, Label: column})
	}
	message := "CSV is ready to download."
	if len(result.Rows) > previewSize {
		message += " Showing the first " + strconv.Itoa(previewSize) + " rows below."
	}
	return cmdry.View{Title: "CSV ready", Components: []cmdry.Component{
		{Type: "alert", Level: "success", Title: strconv.Itoa(len(result.Rows)) + " rows converted", Message: message},
		{Type: "metric", Label: "Columns", Value: strconv.Itoa(len(result.Columns))},
		{Type: "metric", Label: "Rows", Value: strconv.Itoa(len(result.Rows))},
		{Type: "download", Filename: "cmdry-export.csv", MIMEType: "text/csv;charset=utf-8", Content: base64.StdEncoding.EncodeToString(result.CSV)},
		{Type: "table", ID: "jsoncsv-preview", Columns: columns, Rows: rows},
		{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another"}}},
	}}, nil
}
