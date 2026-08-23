package datapriority2

import (
	"encoding/base64"
	"fmt"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunYAMLFamilyConverter() {
	cmdry.Run(plugin("yaml-family-converter", "YAML Family Converter", "Convert YAML, JSON, and TOML document formats locally.", []string{"yaml", "json", "toml", "convert", "format"}, yamlForm, yamlRun))
}
func RunXMLToJSON() {
	cmdry.Run(plugin("xml-to-json", "XML to JSON", "Convert a simple pasted XML document to JSON locally.", []string{"xml", "json", "convert", "data"}, xmlForm, xmlRun))
}
func RunPhoneParserFormatter() {
	cmdry.Run(plugin("phone-parser-formatter", "Phone Parser and Formatter", "Parse NANP numbers or normalize explicit international E.164 numbers locally.", []string{"phone", "telephone", "e164", "nanp", "format"}, phoneForm, phoneRun))
}
func RunIBANValidatorParser() {
	cmdry.Run(plugin("iban-validator-parser", "IBAN Validator and Parser", "Validate an IBAN checksum and show bundled country and bank fields locally.", []string{"iban", "bank", "account", "checksum", "validate"}, ibanForm, ibanRun))
}

func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "data", Icon: "data", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": action}}
}

func yamlForm(cmdry.Request) (cmdry.View, error) {
	return form("YAML Family Converter", "Convert one bounded document. YAML mapping keys must be strings.", "Convert document", []cmdry.Field{{Name: "operation", Label: "Conversion", Type: "select", Value: "yaml-json", Options: []cmdry.Option{{Value: "yaml-json", Label: "YAML to JSON"}, {Value: "yaml-toml", Label: "YAML to TOML"}, {Value: "json-yaml", Label: "JSON to YAML"}, {Value: "toml-yaml", Label: "TOML to YAML"}}}, {Name: "input", Label: "Document", Type: "textarea", Value: "name: Cmdry\nlocal: true\n", Required: true}}), nil
}
func yamlRun(request cmdry.Request) (cmdry.View, error) {
	operation := fmt.Sprint(request.Params["operation"])
	output, err := ConvertYAMLFamily(fmt.Sprint(request.Params["input"]), operation)
	if err != nil {
		return cmdry.View{}, err
	}
	extension, mime := "yaml", "application/yaml;charset=utf-8"
	if operation == "yaml-json" {
		extension, mime = "json", "application/json;charset=utf-8"
	}
	if operation == "yaml-toml" {
		extension, mime = "toml", "application/toml;charset=utf-8"
	}
	return documentResult("Converted document", "The document was converted locally.", "cmdry-output."+extension, mime, output, "Convert another document"), nil
}

func xmlForm(cmdry.Request) (cmdry.View, error) {
	return form("XML to JSON", "Convert one simple XML document. Attributes use @ names and repeated elements become arrays.", "Convert XML", []cmdry.Field{{Name: "input", Label: "XML document", Type: "textarea", Value: "<person id=\"7\"><name>Ada</name></person>", Required: true}}), nil
}
func xmlRun(request cmdry.Request) (cmdry.View, error) {
	output, err := XMLToJSON(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return documentResult("JSON document", "The XML document was converted locally.", "cmdry-output.json", "application/json;charset=utf-8", output, "Convert another XML document"), nil
}

func phoneForm(cmdry.Request) (cmdry.View, error) {
	return form("Phone Parser and Formatter", "NANP validates the US/Canada 10-digit structure. International mode only normalizes explicit E.164 input; it does not verify assignment.", "Parse phone number", []cmdry.Field{{Name: "region", Label: "Validation mode", Type: "select", Value: "nanp", Options: []cmdry.Option{{Value: "nanp", Label: "US / Canada (NANP)"}, {Value: "international", Label: "International E.164 (+ required)"}}}, {Name: "input", Label: "Phone number", Type: "text", Value: "(415) 555-2671", Required: true}}), nil
}
func phoneRun(request cmdry.Request) (cmdry.View, error) {
	phone, err := ParsePhone(fmt.Sprint(request.Params["input"]), fmt.Sprint(request.Params["region"]))
	if err != nil {
		return cmdry.View{}, err
	}
	rows := []map[string]any{{"field": "Input", "value": phone.Input}, {"field": "Mode", "value": phone.Region}, {"field": "E.164", "value": phone.E164}, {"field": "National", "value": phone.National}, {"field": "International", "value": phone.International}, {"field": "Validation note", "value": phone.Note}}
	return tableResult("Parsed phone number", "Formatting and validation were performed locally.", "phone", rows, "Parse another phone number"), nil
}

func ibanForm(cmdry.Request) (cmdry.View, error) {
	return form("IBAN Validator and Parser", "Checks the ISO 13616 mod-97 checksum. Country-specific bank fields are shown only for bundled layouts.", "Validate IBAN", []cmdry.Field{{Name: "input", Label: "IBAN", Type: "text", Value: "DE89 3704 0044 0532 0130 00", Required: true}}), nil
}
func ibanRun(request cmdry.Request) (cmdry.View, error) {
	iban, err := ParseIBAN(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	rows := []map[string]any{{"field": "IBAN", "value": iban.Normalized}, {"field": "Checksum", "value": "Valid"}, {"field": "Country", "value": iban.Country + " (" + iban.CountryCode + ")"}, {"field": "BBAN", "value": iban.BBAN}}
	for _, field := range iban.Fields {
		rows = append(rows, map[string]any{"field": field.Label, "value": field.Value})
	}
	if len(iban.Fields) == 0 {
		rows = append(rows, map[string]any{"field": "Country-specific fields", "value": "Not bundled for this country code"})
	}
	return tableResult("Validated IBAN", "The checksum was validated locally. This does not confirm the account exists.", "iban", rows, "Validate another IBAN"), nil
}

func form(title, description, submit string, fields []cmdry.Field) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Description: description + " Your data never leaves this device.", Action: "run", Submit: submit, Fields: fields}}}
}
func documentResult(title, message, filename, mime string, output []byte, again string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: message}, {Type: "download", Filename: filename, MIMEType: mime, Content: base64.StdEncoding.EncodeToString(output)}, {Type: "code", Title: "Output preview", Text: string(output)}, restart(again)}}
}
func tableResult(title, message, id string, rows []map[string]any, again string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: message}, {Type: "table", ID: id, Columns: []cmdry.Column{{Key: "field", Label: "Field"}, {Key: "value", Label: "Value"}}, Rows: rows}, restart(again)}}
}
func restart(name string) cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: name}}}
}
