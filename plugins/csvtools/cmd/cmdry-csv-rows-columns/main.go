package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/csvtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: manifest("csv-rows-columns", "CSV Rows to Columns", "Convert a CSV table's rows into columns locally."), Actions: map[string]cmdry.Handler{"overview": form, "convert": convert}})
}
func manifest(id, name, description string) cmdry.Manifest {
	return cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "data", Icon: "convert", SearchTerms: []string{"csv", "transpose", "rows", "columns", "convert"}, Pages: []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert CSV", Method: "write"}}}
}
func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "CSV Rows to Columns", Components: []cmdry.Component{{Type: "form", Title: "Transpose CSV", Action: "convert", Submit: "Convert CSV", Description: "The header becomes the first value of each output row.", Fields: []cmdry.Field{{Name: "input", Label: "CSV input", Type: "textarea", Required: true, Placeholder: "name,age\nAda,37\nLin,30"}}}}}, nil
}
func convert(r cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(r.Params["input"])
	if strings.TrimSpace(input) == "" {
		return cmdry.View{}, fmt.Errorf("CSV input is required")
	}
	output, columns, rows, err := csvtools.RowsToColumns(input)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Transposed CSV", Components: []cmdry.Component{{Type: "metric", Label: "Input rows", Value: fmt.Sprint(rows)}, {Type: "metric", Label: "Input columns", Value: fmt.Sprint(columns)}, {Type: "download", Filename: "cmdry-transposed.csv", MIMEType: "text/csv;charset=utf-8", Content: base64.StdEncoding.EncodeToString([]byte(output))}, {Type: "code", Title: "CSV preview", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another"}}}}}, nil
}
