package webpriority2

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"strings"
	"time"
)

func RunBasicAuth() {
	cmdry.Run(p("basic-auth-generator", "BASIC Auth Generator", "Generate a Basic Authorization header locally without retaining credentials.", basicForm, basicRun))
}
func RunHTMLEntity() {
	cmdry.Run(p("html-entity", "HTML Entity Encoder Decoder", "Encode or decode HTML entities locally.", entityForm, entityRun))
}
func RunHTTPStatus() {
	cmdry.Run(p("http-status-reference", "HTTP Status Reference", "Look up common HTTP status meanings locally.", statusForm, statusRun))
}
func RunJWT() {
	cmdry.Run(p("jwt-inspector", "JWT Inspector", "Decode JWT header and claims locally without verifying the signature.", jwtForm, jwtRun))
}
func RunMIME() {
	cmdry.Run(p("mime-type-lookup", "MIME Type Lookup", "Look up common file extensions and MIME types locally.", mimeForm, mimeRun))
}
func RunOTP() {
	cmdry.Run(p("otp-code-generator", "OTP Code Generator", "Generate a current TOTP code locally from a supplied secret.", otpForm, otpRun))
}
func RunOG() {
	cmdry.Run(p("open-graph-meta-generator", "Open Graph Meta Generator", "Generate bounded Open Graph meta tags locally.", ogForm, ogRun))
}
func RunSafeLink() {
	cmdry.Run(p("outlook-safe-link-decoder", "Outlook Safe Link Decoder", "Extract the original destination from a pasted Outlook Safe Links URL locally.", safeForm, safeRun))
}
func RunUA() {
	cmdry.Run(p("user-agent-parser", "User Agent Parser", "Parse a pasted user-agent string locally.", uaForm, uaRun))
}
func p(id, name, desc string, o, a cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: desc, Category: "web", Icon: "web", Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": o, "run": a}}
}
func form(t, d, s string, f []cmdry.Field) (cmdry.View, error) {
	return cmdry.View{Title: t, Components: []cmdry.Component{{Type: "form", Title: t, Description: d + " Your data never leaves this device.", Action: "run", Submit: s, Fields: f}}}, nil
}
func code(t, m, v, a string) (cmdry.View, error) {
	return cmdry.View{Title: t, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: m}, {Type: "code", Title: "Output", Text: v}, again(a)}}, nil
}
func table(t, m string, r []map[string]any, a string) (cmdry.View, error) {
	return cmdry.View{Title: t, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: m}, {Type: "table", ID: "result", Columns: []cmdry.Column{{Key: "field", Label: "Field"}, {Key: "value", Label: "Value"}}, Rows: r}, again(a)}}, nil
}
func again(n string) cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: n}}}
}
func basicForm(cmdry.Request) (cmdry.View, error) {
	return form("BASIC Auth Generator", "Creates an Authorization header; it does not send credentials.", "Generate", []cmdry.Field{{Name: "user", Label: "Username", Type: "text", Required: true}, {Name: "password", Label: "Password", Type: "password"}})
}
func basicRun(r cmdry.Request) (cmdry.View, error) {
	v, e := BasicAuth(fmt.Sprint(r.Params["user"]), fmt.Sprint(r.Params["password"]))
	if e != nil {
		return cmdry.View{}, e
	}
	return code("Authorization header", "Generated locally; credentials are not retained.", v, "Generate another header")
}
func entityForm(cmdry.Request) (cmdry.View, error) {
	return form("HTML Entity Encoder Decoder", "Convert pasted text locally.", "Convert", []cmdry.Field{{Name: "op", Label: "Operation", Type: "select", Value: "encode", Options: []cmdry.Option{{Value: "encode", Label: "Encode"}, {Value: "decode", Label: "Decode"}}}, {Name: "input", Label: "Text", Type: "textarea", Required: true}})
}
func entityRun(r cmdry.Request) (cmdry.View, error) {
	v, e := Entity(fmt.Sprint(r.Params["input"]), fmt.Sprint(r.Params["op"]))
	if e != nil {
		return cmdry.View{}, e
	}
	return code("Converted text", "Converted locally.", v, "Convert another value")
}
func statusForm(cmdry.Request) (cmdry.View, error) {
	return form("HTTP Status Reference", "Looks up common status codes from a bundled reference.", "Look up", []cmdry.Field{{Name: "code", Label: "Status code", Type: "number", Value: "404", Required: true}})
}
func statusRun(r cmdry.Request) (cmdry.View, error) {
	c, m := Status(fmt.Sprint(r.Params["code"]))
	if m == "" {
		return cmdry.View{}, fmt.Errorf("status code is not in the bundled reference")
	}
	return table("HTTP status", "Reference result.", []map[string]any{{"field": "Code", "value": c}, {"field": "Meaning", "value": m}}, "Look up another code")
}
func jwtForm(cmdry.Request) (cmdry.View, error) {
	return form("JWT Inspector", "Decodes JSON header and claims only; it does not verify the signature.", "Inspect JWT", []cmdry.Field{{Name: "input", Label: "JWT", Type: "textarea", Required: true}})
}
func jwtRun(r cmdry.Request) (cmdry.View, error) {
	v, e := InspectJWT(fmt.Sprint(r.Params["input"]))
	if e != nil {
		return cmdry.View{}, e
	}
	return code("JWT inspection", "Decoded locally; signature not verified.", "Header:\n"+v.Header+"\n\nClaims:\n"+v.Claims+"\n\nSignature (base64url):\n"+v.Signature, "Inspect another JWT")
}
func mimeForm(cmdry.Request) (cmdry.View, error) {
	return form("MIME Type Lookup", "Search common extensions or MIME types from a bundled table.", "Look up", []cmdry.Field{{Name: "input", Label: "Extension or MIME type", Type: "text", Value: "png", Required: true}})
}
func mimeRun(r cmdry.Request) (cmdry.View, error) {
	ext, m := MIME(fmt.Sprint(r.Params["input"]))
	if m == "" {
		return cmdry.View{}, fmt.Errorf("value is not in the bundled MIME table")
	}
	return table("MIME type", "Reference result.", []map[string]any{{"field": "Extension", "value": "." + ext}, {"field": "MIME type", "value": m}}, "Look up another value")
}
func otpForm(cmdry.Request) (cmdry.View, error) {
	return form("OTP Code Generator", "Generates the current RFC 6238 SHA-1 TOTP code locally. The secret is not retained.", "Generate code", []cmdry.Field{{Name: "secret", Label: "Base32 secret", Type: "password", Required: true}})
}
func otpRun(r cmdry.Request) (cmdry.View, error) {
	v, e := OTP(fmt.Sprint(r.Params["secret"]), time.Now())
	if e != nil {
		return cmdry.View{}, e
	}
	return code("Current TOTP code", "Generated locally. This code changes every 30 seconds.", v, "Generate another code")
}
func ogForm(cmdry.Request) (cmdry.View, error) {
	return form("Open Graph Meta Generator", "Generate Open Graph tags from local form values.", "Generate tags", []cmdry.Field{{Name: "title", Label: "Title", Type: "text", Required: true}, {Name: "description", Label: "Description", Type: "textarea", Required: true}, {Name: "image", Label: "Image URL", Type: "text"}, {Name: "url", Label: "Page URL", Type: "text"}})
}
func ogRun(r cmdry.Request) (cmdry.View, error) {
	return code("Open Graph tags", "Generated locally.", OG(fmt.Sprint(r.Params["title"]), fmt.Sprint(r.Params["description"]), fmt.Sprint(r.Params["image"]), fmt.Sprint(r.Params["url"])), "Generate another set")
}
func safeForm(cmdry.Request) (cmdry.View, error) {
	return form("Outlook Safe Link Decoder", "Extracts a destination URL from a pasted Safe Links URL locally; it never opens the link.", "Decode link", []cmdry.Field{{Name: "input", Label: "Safe Links URL", Type: "textarea", Required: true}})
}
func safeRun(r cmdry.Request) (cmdry.View, error) {
	v, e := SafeLink(fmt.Sprint(r.Params["input"]))
	if e != nil {
		return cmdry.View{}, e
	}
	return code("Original destination", "Extracted locally; not opened.", v, "Decode another link")
}
func uaForm(cmdry.Request) (cmdry.View, error) {
	return form("User Agent Parser", "Classifies a pasted user-agent with local, conservative heuristics.", "Parse user agent", []cmdry.Field{{Name: "input", Label: "User-Agent", Type: "textarea", Required: true}})
}
func uaRun(r cmdry.Request) (cmdry.View, error) {
	v := UA(fmt.Sprint(r.Params["input"]))
	return table("User-Agent summary", "Parsed locally; classifications are heuristic.", []map[string]any{{"field": "Browser", "value": v[0]}, {"field": "Operating system", "value": v[1]}, {"field": "Device", "value": v[2]}}, "Parse another user-agent")
}

var _ = strings.TrimSpace
