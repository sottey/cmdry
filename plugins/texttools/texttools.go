package texttools

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
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
func RunEmailExtractor() {
	run("email-extractor", "Email Extractor", "Extract unique email addresses from pasted text locally.", []string{"email", "extract", "addresses", "contacts"}, emailForm, extractEmails)
}
func RunDuplicateLineRemover() {
	run("remove-duplicate-lines", "Remove Duplicate Lines", "Keep the first occurrence of each line in pasted text.", []string{"duplicate", "unique", "lines", "deduplicate"}, duplicateLineForm, removeDuplicateLines)
}
func RunTextReplacer() {
	run("text-replacer", "Text Replacer", "Replace every exact text match in pasted content locally.", []string{"replace", "find", "text", "substitute"}, replacerForm, replaceText)
}
func RunUppercase() {
	run("uppercase", "Convert to Uppercase", "Convert pasted text to uppercase locally.", []string{"uppercase", "case", "text", "capitalize"}, uppercaseForm, uppercase)
}
func RunLowercase() {
	run("lowercase", "Convert to Lowercase", "Convert pasted text to lowercase locally.", []string{"lowercase", "case", "text"}, lowercaseForm, lowercase)
}
func RunURLEncoder() {
	run("url-encoder", "URL Encoder", "Encode pasted text for use as a URL query value locally.", []string{"url", "encode", "percent", "query"}, urlEncoderForm, encodeURL)
}
func RunURLDecoder() {
	run("url-decoder", "URL Decoder", "Decode a URL query value locally.", []string{"url", "decode", "percent", "query"}, urlDecoderForm, decodeURL)
}
func RunTextStatistics() {
	run("text-statistics", "Text Statistics", "Count characters, words, lines, and bytes in pasted text locally.", []string{"text", "statistics", "word count", "character count", "reading time"}, textStatisticsForm, textStatistics)
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
func emailForm(cmdry.Request) (cmdry.View, error) {
	return form("Email Extractor", "Extract email addresses", "Extract emails", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}
func duplicateLineForm(cmdry.Request) (cmdry.View, error) {
	return form("Remove Duplicate Lines", "Remove repeated lines", "Remove duplicates", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}
func replacerForm(cmdry.Request) (cmdry.View, error) {
	return form("Text Replacer", "Find and replace", "Replace text", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}, {Name: "find", Label: "Find", Type: "text", Required: true}, {Name: "replace", Label: "Replace with", Type: "text"}}), nil
}
func uppercaseForm(cmdry.Request) (cmdry.View, error) {
	return form("Convert to Uppercase", "Convert text", "Convert to uppercase", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}
func lowercaseForm(cmdry.Request) (cmdry.View, error) {
	return form("Convert to Lowercase", "Convert text", "Convert to lowercase", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}
func urlEncoderForm(cmdry.Request) (cmdry.View, error) {
	return form("URL Encoder", "Encode a query value", "Encode URL value", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}
func urlDecoderForm(cmdry.Request) (cmdry.View, error) {
	return form("URL Decoder", "Decode a query value", "Decode URL value", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}
func textStatisticsForm(cmdry.Request) (cmdry.View, error) {
	return form("Text Statistics", "Analyze text", "Analyze text", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
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
func extractEmails(r cmdry.Request) (cmdry.View, error) {
	emails := ExtractEmails(fmt.Sprint(r.Params["input"]))
	if len(emails) == 0 {
		return cmdry.View{Title: "Email Extractor", Components: []cmdry.Component{{Type: "alert", Level: "info", Title: "No email addresses found", Message: "No email addresses with a domain were detected in this text."}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Try another text"}}}}}, nil
	}
	return cmdry.View{Title: "Email addresses", Components: []cmdry.Component{{Type: "metric", Label: "Unique addresses", Value: fmt.Sprint(len(emails))}, {Type: "code", Title: "Email addresses", Text: strings.Join(emails, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Extract another list"}}}}}, nil
}
func removeDuplicateLines(r cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(r.Params["input"])
	output, removed := DeduplicateLines(input)
	return cmdry.View{Title: "Unique lines", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Duplicates removed", Message: fmt.Sprint(removed) + " repeated line(s) removed."}, {Type: "code", Title: "Unique lines", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Clean another list"}}}}}, nil
}
func replaceText(r cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(r.Params["input"])
	find := fmt.Sprint(r.Params["find"])
	replacement, _ := r.Params["replace"].(string)
	if find == "" {
		return cmdry.View{}, fmt.Errorf("find text is required")
	}
	output, count := ReplaceAll(input, find, replacement)
	return cmdry.View{Title: "Text replaced", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Replacement complete", Message: fmt.Sprint(count) + " match(es) replaced."}, {Type: "code", Title: "Replaced text", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Replace more text"}}}}}, nil
}
func uppercase(r cmdry.Request) (cmdry.View, error) {
	return result("Uppercase text", strings.ToUpper(fmt.Sprint(r.Params["input"]))), nil
}
func lowercase(r cmdry.Request) (cmdry.View, error) {
	return result("Lowercase text", strings.ToLower(fmt.Sprint(r.Params["input"]))), nil
}
func encodeURL(r cmdry.Request) (cmdry.View, error) {
	return result("URL-encoded text", url.QueryEscape(fmt.Sprint(r.Params["input"]))), nil
}
func decodeURL(r cmdry.Request) (cmdry.View, error) {
	decoded, err := url.QueryUnescape(fmt.Sprint(r.Params["input"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("decode URL value: %w", err)
	}
	return result("URL-decoded text", decoded), nil
}
func textStatistics(r cmdry.Request) (cmdry.View, error) {
	statistics := TextStatistics(fmt.Sprint(r.Params["input"]))
	components := []cmdry.Component{
		{Type: "metric", Label: "Characters", Value: fmt.Sprint(statistics.Characters)},
		{Type: "metric", Label: "Words", Value: fmt.Sprint(statistics.Words)},
		{Type: "metric", Label: "Lines", Value: fmt.Sprint(statistics.Lines)},
		{Type: "metric", Label: "Bytes (UTF-8)", Value: fmt.Sprint(statistics.Bytes)},
		{Type: "metric", Label: "Reading time", Value: statistics.ReadingTime},
		{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Analyze more text"}}},
	}
	return cmdry.View{Title: "Text Statistics", Components: components}, nil
}

var emailPattern = regexp.MustCompile(`(?i)[a-z0-9.!#$%&'*+/=?^_{|}~-]+@[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+`)

// ExtractEmails finds standard domain-qualified email addresses and preserves
// the first spelling of each address.
func ExtractEmails(input string) []string {
	matches := emailPattern.FindAllString(input, -1)
	seen := make(map[string]struct{}, len(matches))
	emails := make([]string, 0, len(matches))
	for _, email := range matches {
		key := strings.ToLower(email)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		emails = append(emails, email)
	}
	return emails
}

// DeduplicateLines preserves the first exact occurrence of every line.
func DeduplicateLines(input string) (string, int) {
	lines := strings.Split(input, "\n")
	seen := make(map[string]struct{}, len(lines))
	unique := make([]string, 0, len(lines))
	removed := 0
	for _, line := range lines {
		if _, exists := seen[line]; exists {
			removed++
			continue
		}
		seen[line] = struct{}{}
		unique = append(unique, line)
	}
	return strings.Join(unique, "\n"), removed
}

// ReplaceAll returns the transformed text and number of exact replacements.
func ReplaceAll(input, find, replacement string) (string, int) {
	if find == "" {
		return input, 0
	}
	return strings.ReplaceAll(input, find, replacement), strings.Count(input, find)
}

// Statistics is a set of useful local text measurements. Characters are Unicode
// code points, while bytes are the UTF-8 encoded byte count.
type Statistics struct {
	Characters  int
	Words       int
	Lines       int
	Bytes       int
	ReadingTime string
}

// TextStatistics counts a pasted text's words, lines, Unicode characters, and
// UTF-8 bytes. Reading time assumes a 200 words-per-minute adult reading pace.
func TextStatistics(input string) Statistics {
	words := len(strings.Fields(input))
	lines := 0
	if input != "" {
		lines = strings.Count(input, "\n") + 1
	}
	minutes := (words + 199) / 200
	readingTime := "Less than 1 minute"
	if minutes > 1 {
		readingTime = fmt.Sprintf("About %d minutes", minutes)
	}
	return Statistics{Characters: len([]rune(input)), Words: words, Lines: lines, Bytes: len([]byte(input)), ReadingTime: readingTime}
}
func result(title, text string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "code", Title: title, Text: text}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Start over"}}}}}
}
