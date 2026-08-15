package texttools

import (
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunCompare() {
	run("text-compare", "Text Compare", "Compare two texts locally and view a line-level diff.", []string{"diff", "compare", "text"}, compareForm, compare)
}
func RunBase64() {
	run("text-base64", "Base64 Encoder/Decoder", "Encode or decode pasted text locally.", []string{"base64", "encode", "decode"}, base64Form, base64Action)
}
func RunSlug() {
	run("text-slug", "Slug Generator", "Create a URL-friendly slug from pasted text.", []string{"slug", "url", "permalink"}, slugForm, slug)
}

func run(id, name, description string, terms []string, overview cmdry.Handler, action cmdry.Handler) {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "text", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": action}})
}
func compareForm(cmdry.Request) (cmdry.View, error) {
	return form("Text Compare", "Compare text", "Compare", []cmdry.Field{{Name: "original", Label: "Original", Type: "textarea", Required: true}, {Name: "revised", Label: "Revised", Type: "textarea", Required: true}}), nil
}
func base64Form(cmdry.Request) (cmdry.View, error) {
	return form("Base64 Encoder/Decoder", "Encode or decode", "Run", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}, {Name: "mode", Label: "Operation", Type: "select", Value: "encode", Options: []cmdry.Option{{Value: "encode", Label: "Encode"}, {Value: "decode", Label: "Decode"}}}}), nil
}
func slugForm(cmdry.Request) (cmdry.View, error) {
	return form("Slug Generator", "Create a slug", "Create slug", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}
func form(title, heading, submit string, fields []cmdry.Field) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: heading, Action: "run", Submit: submit, Fields: fields}}}
}
func compare(r cmdry.Request) (cmdry.View, error) {
	a := strings.Split(fmt.Sprint(r.Params["original"]), "\n")
	b := strings.Split(fmt.Sprint(r.Params["revised"]), "\n")
	if len(a) > 1200 || len(b) > 1200 {
		return cmdry.View{}, fmt.Errorf("compare at most 1,200 lines per input")
	}
	dp := make([][]int, len(a)+1)
	for i := range dp {
		dp[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}
	out := []string{"--- Original", "+++ Revised"}
	for i, j := 0, 0; i < len(a) || j < len(b); {
		if i < len(a) && j < len(b) && a[i] == b[j] {
			out = append(out, " "+a[i])
			i++
			j++
		} else if j < len(b) && (i == len(a) || dp[i][j+1] >= dp[i+1][j]) {
			out = append(out, "+"+b[j])
			j++
		} else {
			out = append(out, "-"+a[i])
			i++
		}
	}
	return result("Text diff", strings.Join(out, "\n")), nil
}
func base64Action(r cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(r.Params["input"])
	if r.Params["mode"] == "decode" {
		raw, e := base64.StdEncoding.DecodeString(input)
		if e != nil {
			return cmdry.View{}, fmt.Errorf("decode Base64: %w", e)
		}
		return result("Decoded text", string(raw)), nil
	}
	return result("Base64", base64.StdEncoding.EncodeToString([]byte(input))), nil
}
func slug(r cmdry.Request) (cmdry.View, error) {
	var out []rune
	dash := false
	for _, ch := range strings.ToLower(fmt.Sprint(r.Params["input"])) {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			out = append(out, ch)
			dash = false
		} else if len(out) > 0 && !dash {
			out = append(out, '-')
			dash = true
		}
	}
	value := strings.Trim(string(out), "-")
	if value == "" {
		return cmdry.View{}, fmt.Errorf("text contains no letters or numbers")
	}
	return result("Slug", value), nil
}
func result(title, text string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "code", Title: title, Text: text}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Start over"}}}}}
}
