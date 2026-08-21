package conversiontools

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunTOMLTools() {
	cmdry.Run(plugin("toml-tools", "JSON and TOML", "Convert JSON and TOML, or validate and format TOML locally.", []string{"toml", "json", "convert", "format", "validate"}, tomlForm, tomlRun))
}

func RunSpreadsheetConverter() {
	cmdry.Run(plugin("spreadsheet-converter", "CSV, TSV, and XLSX", "Convert one bounded CSV, TSV, or XLSX file locally.", []string{"csv", "tsv", "xlsx", "excel", "spreadsheet", "convert"}, spreadsheetForm, spreadsheetRun))
}

func RunMarkdownHTML() {
	cmdry.Run(plugin("markdown-html", "Markdown to HTML", "Render pasted Markdown as a local HTML download.", []string{"markdown", "md", "html", "render", "convert"}, markdownForm, markdownRun))
}

func RunJSONTypeScript() {
	cmdry.Run(plugin("json-typescript", "JSON to TypeScript", "Generate a TypeScript type from one pasted JSON sample.", []string{"json", "typescript", "ts", "types", "interface", "schema"}, typeScriptForm, typeScriptRun))
}

func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "developer", Icon: "convert", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "convert": action}}
}

func tomlForm(cmdry.Request) (cmdry.View, error) {
	return form("JSON and TOML", "Convert structured text locally. TOML formatting also validates the document.", "Convert", []cmdry.Field{{Name: "operation", Label: "Operation", Type: "select", Value: "json-to-toml", Options: []cmdry.Option{{Value: "json-to-toml", Label: "JSON to TOML"}, {Value: "toml-to-json", Label: "TOML to JSON"}, {Value: "format-toml", Label: "Validate and format TOML"}}}, {Name: "input", Label: "Input", Type: "textarea", Required: true}}), nil
}

func tomlRun(request cmdry.Request) (cmdry.View, error) {
	input := strings.TrimSpace(fmt.Sprint(request.Params["input"]))
	if input == "" {
		return cmdry.View{}, fmt.Errorf("input is required")
	}
	operation := fmt.Sprint(request.Params["operation"])
	var output []byte
	var err error
	filename, mime := "cmdry-output.toml", "application/toml;charset=utf-8"
	switch operation {
	case "json-to-toml":
		output, err = JSONToTOML(input)
	case "toml-to-json":
		output, err = TOMLToJSON(input)
		filename, mime = "cmdry-output.json", "application/json;charset=utf-8"
	case "format-toml":
		output, err = FormatTOML(input)
	default:
		return cmdry.View{}, fmt.Errorf("unsupported operation")
	}
	if err != nil {
		return cmdry.View{}, err
	}
	return downloadView("Conversion ready", filename, mime, output), nil
}

func spreadsheetForm(cmdry.Request) (cmdry.View, error) {
	return form("CSV, TSV, and XLSX", "The first worksheet is used for XLSX. CSV formula-like cells are stored and exported as text.", "Convert", []cmdry.Field{{Name: "source", Label: "Input format", Type: "select", Value: "csv", Options: []cmdry.Option{{Value: "csv", Label: "CSV"}, {Value: "tsv", Label: "TSV"}, {Value: "xlsx", Label: "XLSX (first worksheet)"}}}, {Name: "target", Label: "Output format", Type: "select", Value: "xlsx", Options: []cmdry.Option{{Value: "xlsx", Label: "XLSX"}, {Value: "csv", Label: "CSV"}, {Value: "tsv", Label: "TSV"}}}, {Name: "file", Label: "Spreadsheet file", Type: "file", Accept: ".csv,.tsv,.xlsx,text/csv,text/tab-separated-values,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Required: true}}), nil
}

func spreadsheetRun(request cmdry.Request) (cmdry.View, error) {
	source, target := fmt.Sprint(request.Params["source"]), fmt.Sprint(request.Params["target"])
	if source == target {
		return cmdry.View{}, fmt.Errorf("choose a different output format")
	}
	upload, contents, err := request.File("file")
	if err != nil {
		return cmdry.View{}, err
	}
	_ = upload
	var output []byte
	var rows int
	if source == "xlsx" {
		if target == "xlsx" {
			return cmdry.View{}, fmt.Errorf("choose CSV or TSV output")
		}
		delimiter := ','
		if target == "tsv" {
			delimiter = '\t'
		}
		output, rows, err = XLSXToDelimited(contents, delimiter)
	} else {
		if target != "xlsx" {
			return cmdry.View{}, fmt.Errorf("CSV and TSV can be converted to XLSX")
		}
		delimiter := ','
		if source == "tsv" {
			delimiter = '\t'
		}
		output, rows, err = CSVToXLSX(contents, delimiter)
	}
	if err != nil {
		return cmdry.View{}, err
	}
	filename := strings.TrimSuffix(upload.Name, filepath.Ext(upload.Name)) + "." + target
	mime := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if target == "csv" {
		mime = "text/csv;charset=utf-8"
	}
	if target == "tsv" {
		mime = "text/tab-separated-values;charset=utf-8"
	}
	return cmdry.View{Title: "Spreadsheet ready", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: fmt.Sprintf("%d rows converted", rows), Message: "The converted file is ready for a browser-local download."}, {Type: "download", Filename: filename, MIMEType: mime, Content: base64.StdEncoding.EncodeToString(output)}, restart()}}, nil
}

func markdownForm(cmdry.Request) (cmdry.View, error) {
	return form("Markdown to HTML", "Render pasted Markdown locally. The output is provided as source and a download.", "Render HTML", []cmdry.Field{{Name: "input", Label: "Markdown", Type: "textarea", Required: true}}), nil
}
func markdownRun(request cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(request.Params["input"])
	if strings.TrimSpace(input) == "" {
		return cmdry.View{}, fmt.Errorf("Markdown is required")
	}
	output, err := MarkdownToHTML(input)
	if err != nil {
		return cmdry.View{}, err
	}
	return downloadView("HTML ready", "cmdry-output.html", "text/html;charset=utf-8", output), nil
}

func typeScriptForm(cmdry.Request) (cmdry.View, error) {
	return form("JSON to TypeScript", "Generate a type alias from one JSON sample locally.", "Generate TypeScript", []cmdry.Field{{Name: "name", Label: "Root type name", Type: "text", Value: "Root", Required: true}, {Name: "input", Label: "JSON sample", Type: "textarea", Required: true}}), nil
}
func typeScriptRun(request cmdry.Request) (cmdry.View, error) {
	output, err := JSONToTypeScript(fmt.Sprint(request.Params["input"]), fmt.Sprint(request.Params["name"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return downloadView("TypeScript ready", "cmdry-types.ts", "text/typescript;charset=utf-8", []byte(output)), nil
}

func form(title, description, submit string, fields []cmdry.Field) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Action: "convert", Submit: submit, Description: description + " Your data never leaves this device.", Fields: fields}}}
}
func downloadView(title, filename, mime string, output []byte) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "download", Filename: filename, MIMEType: mime, Content: base64.StdEncoding.EncodeToString(output)}, {Type: "code", Title: "Output preview", Text: string(output)}, restart()}}
}
func restart() cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Start over"}}}
}
