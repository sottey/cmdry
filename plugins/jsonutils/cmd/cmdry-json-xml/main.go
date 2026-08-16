package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"strings"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.json-xml", Name: "JSON to XML", Version: "0.1.0", Description: "Convert JSON to a simple XML document locally.", Category: "data", Icon: "convert", SearchTerms: []string{"json", "xml", "convert"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "run", Name: "Convert to XML", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": form, "run": run}})
}
func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "JSON to XML", Components: []cmdry.Component{{Type: "form", Title: "Paste JSON", Action: "run", Submit: "Convert", Fields: []cmdry.Field{{Name: "input", Label: "JSON input", Type: "textarea", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	var value any
	if e := json.Unmarshal([]byte(fmt.Sprint(r.Params["input"])), &value); e != nil {
		return cmdry.View{}, fmt.Errorf("invalid JSON: %w", e)
	}
	var b strings.Builder
	b.WriteString(xml.Header)
	write(&b, "root", value)
	return cmdry.View{Title: "XML ready", Components: []cmdry.Component{{Type: "code", Title: "XML", Text: b.String()}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another"}}}}}, nil
}
func write(b *strings.Builder, name string, v any) {
	switch x := v.(type) {
	case map[string]any:
		b.WriteString("<" + name + ">")
		for k, v := range x {
			write(b, k, v)
		}
		b.WriteString("</" + name + ">")
	case []any:
		for _, v := range x {
			write(b, "item", v)
		}
	default:
		b.WriteString("<" + name + ">" + xmlEscape(fmt.Sprint(x)) + "</" + name + ">")
	}
}
func xmlEscape(s string) string {
	var b strings.Builder
	xml.EscapeText(&b, []byte(s))
	return b.String()
}
