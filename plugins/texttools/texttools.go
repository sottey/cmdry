package texttools

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
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
func RunSubstringExtractor() {
	run("extract-substring", "Extract Substring", "Extract a Unicode-safe character range from pasted text locally.", []string{"substring", "extract", "text", "characters", "range"}, substringForm, extractSubstring)
}
func RunTextJoiner() {
	run("join-text", "Join Text", "Join pasted newline-delimited text items locally.", []string{"join", "text", "separator", "combine", "lines"}, joinTextForm, joinText)
}
func RunTextReverser() {
	run("reverse-text", "Reverse Text", "Reverse pasted Unicode text character by character locally.", []string{"reverse", "text", "characters", "unicode"}, reverseTextForm, reverseText)
}
func RunROT13() {
	run("rot13", "ROT13 Encoder/Decoder", "Encode or decode pasted text using the ROT13 cipher locally.", []string{"rot13", "cipher", "encode", "decode", "text"}, rot13Form, rot13)
}
func RunSplitter() {
	run("split-text", "Split Text", "Split pasted text into parts using a separator locally.", []string{"split", "text", "separator", "delimiter", "parts"}, splitForm, splitText)
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
func substringForm(cmdry.Request) (cmdry.View, error) {
	return form("Extract Substring", "Extract character range", "Extract substring", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}, {Name: "start", Label: "Start character", Type: "number", Value: "1", Min: "1", Required: true}, {Name: "end", Label: "End character (inclusive)", Type: "number", Value: "1", Min: "1", Required: true}}), nil
}
func joinTextForm(cmdry.Request) (cmdry.View, error) {
	return form("Join Text", "Join text items", "Join text", []cmdry.Field{{Name: "input", Label: "Text items", Type: "textarea", Required: true}, {Name: "separator", Label: "Separator", Type: "select", Value: "newline", Options: []cmdry.Option{{Value: "newline", Label: "New line"}, {Value: "space", Label: "Space"}, {Value: "comma", Label: "Comma and space"}, {Value: "none", Label: "No separator"}, {Value: "custom", Label: "Custom"}}}, {Name: "custom", Label: "Custom separator", Type: "text"}}), nil
}
func reverseTextForm(cmdry.Request) (cmdry.View, error) {
	return form("Reverse Text", "Reverse text", "Reverse text", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}
func rot13Form(cmdry.Request) (cmdry.View, error) {
	return form("ROT13 Encoder/Decoder", "Transform text", "Apply ROT13", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}
func splitForm(cmdry.Request) (cmdry.View, error) {
	return form("Split Text", "Split text", "Split text", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}, {Name: "separator", Label: "Separator", Type: "text", Required: true, Placeholder: "For example: , or |"}, {Name: "omit_empty", Label: "Omit empty parts", Type: "checkbox", Value: "true"}}), nil
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
func extractSubstring(r cmdry.Request) (cmdry.View, error) {
	start, err := strconv.Atoi(fmt.Sprint(r.Params["start"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("start character must be a whole number")
	}
	end, err := strconv.Atoi(fmt.Sprint(r.Params["end"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("end character must be a whole number")
	}
	output, err := ExtractSubstring(fmt.Sprint(r.Params["input"]), start, end)
	if err != nil {
		return cmdry.View{}, err
	}
	return result("Extracted substring", output), nil
}
func joinText(r cmdry.Request) (cmdry.View, error) {
	separator := fmt.Sprint(r.Params["separator"])
	custom, _ := r.Params["custom"].(string)
	output, count, err := JoinLines(fmt.Sprint(r.Params["input"]), separator, custom)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Joined text", Components: []cmdry.Component{{Type: "metric", Label: "Items joined", Value: fmt.Sprint(count)}, {Type: "code", Title: "Joined text", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Join more text"}}}}}, nil
}
func reverseText(r cmdry.Request) (cmdry.View, error) {
	return result("Reversed text", ReverseText(fmt.Sprint(r.Params["input"]))), nil
}
func rot13(r cmdry.Request) (cmdry.View, error) {
	return result("ROT13 text", ROT13(fmt.Sprint(r.Params["input"]))), nil
}
func splitText(r cmdry.Request) (cmdry.View, error) {
	parts, err := SplitText(fmt.Sprint(r.Params["input"]), fmt.Sprint(r.Params["separator"]), fmt.Sprint(r.Params["omit_empty"]) == "true")
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Split text", Components: []cmdry.Component{{Type: "metric", Label: "Parts", Value: fmt.Sprint(len(parts))}, {Type: "code", Title: "Split parts", Text: strings.Join(parts, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Split more text"}}}}}, nil
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

// ExtractSubstring returns one-based inclusive Unicode character positions.
func ExtractSubstring(input string, start, end int) (string, error) {
	characters := []rune(input)
	if start < 1 || end < 1 {
		return "", fmt.Errorf("character positions must be at least 1")
	}
	if start > end {
		return "", fmt.Errorf("start character must not be after end character")
	}
	if end > len(characters) {
		return "", fmt.Errorf("end character %d exceeds the text length of %d characters", end, len(characters))
	}
	return string(characters[start-1 : end]), nil
}

// JoinLines joins non-blank newline-delimited text values with the chosen
// separator. Exact item spacing is otherwise preserved.
func JoinLines(input, separator, custom string) (string, int, error) {
	items := make([]string, 0)
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) != "" {
			items = append(items, line)
		}
	}
	if len(items) == 0 {
		return "", 0, fmt.Errorf("enter at least one non-blank text item")
	}
	separators := map[string]string{"newline": "\n", "space": " ", "comma": ", ", "none": "", "custom": custom}
	value, ok := separators[separator]
	if !ok {
		return "", 0, fmt.Errorf("unsupported separator")
	}
	return strings.Join(items, value), len(items), nil
}

// ReverseText reverses Unicode code points while preserving their UTF-8 form.
func ReverseText(input string) string {
	characters := []rune(input)
	for left, right := 0, len(characters)-1; left < right; left, right = left+1, right-1 {
		characters[left], characters[right] = characters[right], characters[left]
	}
	return string(characters)
}

// ROT13 shifts ASCII letters by 13 positions. Applying it twice restores the
// original text.
func ROT13(input string) string {
	return strings.Map(func(value rune) rune {
		switch {
		case value >= 'a' && value <= 'z':
			return 'a' + (value-'a'+13)%26
		case value >= 'A' && value <= 'Z':
			return 'A' + (value-'A'+13)%26
		default:
			return value
		}
	}, input)
}

// SplitText splits text on a non-empty exact separator. Empty parts may be
// omitted, while all other whitespace remains unchanged.
func SplitText(input, separator string, omitEmpty bool) ([]string, error) {
	if separator == "" {
		return nil, fmt.Errorf("separator is required")
	}
	parts := strings.Split(input, separator)
	if !omitEmpty {
		return parts, nil
	}
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return filtered, nil
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
