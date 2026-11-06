package priority1tools

import (
	"encoding/base64"
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"strings"
)

func RunCSS() {
	cmdry.Run(p("css-formatter-minifier", "CSS Formatter and Minifier", "Format or minify pasted CSS locally.", cssForm, cssRun))
}
func RunHTMLJSX() {
	cmdry.Run(p("html-to-jsx", "HTML to JSX Converter", "Convert pasted HTML to JSX locally.", jsxForm, jsxRun))
}
func RunQRDecoder() {
	cmdry.Run(p("qr-code-decoder", "QR Code Decoder", "Decode one uploaded QR-code image locally.", qrForm, qrRun))
}
func RunDataURL() {
	cmdry.Run(p("base64-image-data-url", "Base64 Image and Data URL Converter", "Convert one bounded upload or base64 data URL locally.", dataForm, dataRun))
}
func RunSlash() {
	cmdry.Run(p("backslash-escape-unescape", "Backslash Escape Unescape", "Escape or unescape generic backslash sequences locally.", slashForm, slashRun))
}
func RunJSONCode() {
	cmdry.Run(p("json-to-code", "JSON to Code Generator", "Generate bounded code declarations from a JSON sample locally.", jsonForm, jsonRun))
}
func RunStringInspector() {
	cmdry.Run(p("string-inspector", "String Inspector", "Inspect local text length line word Unicode and encoding facts.", inspectForm, inspectRun))
}
func p(id, name, desc string, o, a cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: desc, Category: category(id), Icon: "code", SearchTerms: []string{id, "convert", "local"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": o, "run": a}}
}
func category(id string) string {
	if strings.Contains(id, "qr") {
		return "image"
	}
	if strings.Contains(id, "data-url") || strings.Contains(id, "html-to") {
		return "converter"
	}
	if strings.Contains(id, "json-to") {
		return "developer"
	}
	if strings.Contains(id, "css") {
		return "web"
	}
	return "text"
}
func form(t, d, s string, f []cmdry.Field) (cmdry.View, error) {
	return cmdry.View{Title: t, Components: []cmdry.Component{{Type: "form", Title: t, Description: d + " Your data never leaves this device.", Action: "run", Submit: s, Fields: f}}}, nil
}
func again(n string) cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: n}}}
}
func code(t, m, v, a string) (cmdry.View, error) {
	return cmdry.View{Title: t, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: m}, {Type: "code", Title: "Output", Text: v}, again(a)}}, nil
}
func table(t, m string, r []map[string]any, a string) (cmdry.View, error) {
	return cmdry.View{Title: t, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: m}, {Type: "table", ID: "result", Columns: []cmdry.Column{{Key: "field", Label: "Field"}, {Key: "value", Label: "Value"}}, Rows: r}, again(a)}}, nil
}
func cssForm(cmdry.Request) (cmdry.View, error) {
	return form("CSS Formatter and Minifier", "Format or minify bounded CSS without sending it anywhere.", "Process CSS", []cmdry.Field{{Name: "operation", Label: "Operation", Type: "select", Value: "format", Options: []cmdry.Option{{Value: "format", Label: "Format"}, {Value: "minify", Label: "Minify"}}}, {Name: "input", Label: "CSS", Type: "textarea", Required: true}})
}
func cssRun(r cmdry.Request) (cmdry.View, error) {
	v, e := FormatCSS(fmt.Sprint(r.Params["input"]), fmt.Sprint(r.Params["operation"]) == "minify")
	if e != nil {
		return cmdry.View{}, e
	}
	return code("CSS ready", "CSS was processed locally.", v, "Process another stylesheet")
}
func jsxForm(cmdry.Request) (cmdry.View, error) {
	return form("HTML to JSX Converter", "Converts HTML attributes such as class and for to JSX names. Inline style strings remain strings.", "Convert HTML", []cmdry.Field{{Name: "input", Label: "HTML", Type: "textarea", Required: true}})
}
func jsxRun(r cmdry.Request) (cmdry.View, error) {
	v, e := HTMLToJSX(fmt.Sprint(r.Params["input"]))
	if e != nil {
		return cmdry.View{}, e
	}
	return code("JSX ready", "HTML was converted locally.", v, "Convert another document")
}
func qrForm(cmdry.Request) (cmdry.View, error) {
	return form("QR Code Decoder", "Upload one PNG, JPEG, or GIF containing a QR code. It is held only for this request.", "Decode QR code", []cmdry.Field{{Name: "image", Label: "QR-code image", Type: "file", Accept: "image/png,image/jpeg,image/gif", Required: true}})
}
func qrRun(r cmdry.Request) (cmdry.View, error) {
	_, b, e := r.File("image")
	if e != nil {
		return cmdry.View{}, e
	}
	v, e := DecodeQR(b)
	if e != nil {
		return cmdry.View{}, e
	}
	return code("QR code decoded", "The image was decoded locally.", v, "Decode another QR code")
}
func dataForm(cmdry.Request) (cmdry.View, error) {
	return form("Base64 Image and Data URL Converter", "Choose upload to create a data URL, or data URL to create a download. Provide only the field required by the chosen direction.", "Convert", []cmdry.Field{{Name: "operation", Label: "Direction", Type: "select", Value: "to-data-url", Options: []cmdry.Option{{Value: "to-data-url", Label: "Upload to data URL"}, {Value: "from-data-url", Label: "Data URL to download"}}}, {Name: "file", Label: "Image or file", Type: "file", Accept: "image/*"}, {Name: "input", Label: "Base64 data URL", Type: "textarea"}})
}
func dataRun(r cmdry.Request) (cmdry.View, error) {
	if fmt.Sprint(r.Params["operation"]) == "to-data-url" {
		u, b, e := r.File("file")
		if e != nil {
			return cmdry.View{}, e
		}
		v, e := ToDataURL(b, u.MIMEType)
		if e != nil {
			return cmdry.View{}, e
		}
		return code("Data URL ready", "The upload was encoded locally.", v, "Convert another file")
	}
	b, m, e := FromDataURL(fmt.Sprint(r.Params["input"]))
	if e != nil {
		return cmdry.View{}, e
	}
	return cmdry.View{Title: "File ready", Components: []cmdry.Component{{Type: "download", Filename: FilenameForMIME(m), MIMEType: m, Content: base64.StdEncoding.EncodeToString(b)}, again("Convert another data URL")}}, nil
}
func slashForm(cmdry.Request) (cmdry.View, error) {
	return form("Backslash Escape Unescape", "Supports standard newline, carriage-return, tab, quote, and slash escapes.", "Convert", []cmdry.Field{{Name: "operation", Label: "Operation", Type: "select", Value: "escape", Options: []cmdry.Option{{Value: "escape", Label: "Escape"}, {Value: "unescape", Label: "Unescape"}}}, {Name: "input", Label: "Text", Type: "textarea", Required: true}})
}
func slashRun(r cmdry.Request) (cmdry.View, error) {
	v, e := SlashEscape(fmt.Sprint(r.Params["input"]), fmt.Sprint(r.Params["operation"]) == "unescape")
	if e != nil {
		return cmdry.View{}, e
	}
	return code("Converted text", "Text was converted locally.", v, "Convert another value")
}
func jsonForm(cmdry.Request) (cmdry.View, error) {
	return form("JSON to Code Generator", "TypeScript is structural; Go and Python output conservative top-level declarations.", "Generate code", []cmdry.Field{{Name: "target", Label: "Target", Type: "select", Value: "typescript", Options: []cmdry.Option{{Value: "typescript", Label: "TypeScript"}, {Value: "go", Label: "Go"}, {Value: "python", Label: "Python"}}}, {Name: "name", Label: "Root name", Type: "text", Value: "Root", Required: true}, {Name: "input", Label: "JSON sample", Type: "textarea", Required: true}})
}
func jsonRun(r cmdry.Request) (cmdry.View, error) {
	v, e := JSONCode(fmt.Sprint(r.Params["input"]), fmt.Sprint(r.Params["name"]), fmt.Sprint(r.Params["target"]))
	if e != nil {
		return cmdry.View{}, e
	}
	return code("Code ready", "Generated locally from the JSON sample.", v, "Generate another declaration")
}
func inspectForm(cmdry.Request) (cmdry.View, error) {
	return form("String Inspector", "Inspect bounded text without retaining it.", "Inspect text", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}})
}
func inspectRun(r cmdry.Request) (cmdry.View, error) {
	f := InspectString(fmt.Sprint(r.Params["input"]))
	return table("String facts", "Facts were calculated locally.", []map[string]any{{"field": "Bytes", "value": fmt.Sprint(f.Bytes)}, {"field": "Unicode code points", "value": fmt.Sprint(f.Runes)}, {"field": "Words", "value": fmt.Sprint(f.Words)}, {"field": "Lines", "value": fmt.Sprint(f.Lines)}, {"field": "ASCII only", "value": fmt.Sprint(f.ASCII)}}, "Inspect another string")
}
