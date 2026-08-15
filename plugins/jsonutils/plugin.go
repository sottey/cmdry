package jsonutils

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunValidator() {
	run("json-validator", "JSON Validator", "Validate pasted JSON locally and identify its root value.", []string{"json", "validate", "validator", "syntax"}, validate)
}
func RunMinifier() {
	run("json-minifier", "JSON Minifier", "Remove unnecessary whitespace from pasted JSON locally.", []string{"json", "minify", "compact", "compress"}, minify)
}
func RunStringifier() {
	run("json-stringifier", "JSON Stringifier", "Encode a JSON document as a JSON string literal locally.", []string{"json", "stringify", "escape", "string"}, stringify)
}
func RunEscaper() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.json-escaper", Name: "JSON String Escaper", Version: "0.1.0", Description: "Escape arbitrary pasted text as a JSON string literal locally.", Category: "developer", Icon: "json",
		SearchTerms: []string{"json", "escape", "string", "special characters", "quote"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Escape JSON string", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": escaperForm, "run": escape}})
}

func run(id, name, description string, terms []string, action cmdry.Handler) {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "developer", Icon: "json",
		SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: name, Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": func(cmdry.Request) (cmdry.View, error) { return form(name, description), nil }, "run": action}})
}

func form(name, description string) cmdry.View {
	return cmdry.View{Title: name, Components: []cmdry.Component{{Type: "form", Title: name, Action: "run", Submit: name, Description: description + " Your data never leaves this device.", Fields: []cmdry.Field{{Name: "input", Label: "JSON input", Type: "textarea", Required: true}}}}}
}

func input(request cmdry.Request) (string, error) {
	value, _ := request.Params["input"].(string)
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("JSON input is required")
	}
	return value, nil
}

func validate(request cmdry.Request) (cmdry.View, error) {
	value, err := input(request)
	if err != nil {
		return cmdry.View{}, err
	}
	root, _, err := Parse(value)
	if err != nil {
		return cmdry.View{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return cmdry.View{Title: "Valid JSON", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "JSON is valid", Message: "The input contains exactly one valid JSON value."}, {Type: "metric", Label: "Root value", Value: RootType(root)}, restart()}}, nil
}

func minify(request cmdry.Request) (cmdry.View, error) {
	value, err := input(request)
	if err != nil {
		return cmdry.View{}, err
	}
	_, compact, err := Parse(value)
	if err != nil {
		return cmdry.View{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return output("Minified JSON", string(compact)), nil
}

func stringify(request cmdry.Request) (cmdry.View, error) {
	value, err := input(request)
	if err != nil {
		return cmdry.View{}, err
	}
	encoded, err := Stringify(value)
	if err != nil {
		return cmdry.View{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return output("JSON string", string(encoded)), nil
}

func escaperForm(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "JSON String Escaper", Components: []cmdry.Component{{Type: "form", Title: "Escape text", Action: "run", Submit: "Escape JSON string", Description: "Escapes arbitrary text as one JSON string literal locally. The input does not need to be JSON.", Fields: []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}}}}, nil
}

func escape(request cmdry.Request) (cmdry.View, error) {
	value := fmt.Sprint(request.Params["input"])
	escaped, err := EscapeString(value)
	if err != nil {
		return cmdry.View{}, fmt.Errorf("escape JSON string: %w", err)
	}
	return output("Escaped JSON string", string(escaped)), nil
}

func output(title, value string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "code", Title: title, Text: value}, restart()}}
}

func restart() cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Start over"}}}
}
