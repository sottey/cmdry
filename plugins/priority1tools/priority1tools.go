// Package priority1tools implements bounded local-only priority-one utilities.
package priority1tools

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"regexp"
	"strings"
	"unicode"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/sottey/cmdry/plugins/conversiontools"
	"golang.org/x/net/html"
)

const maxInput = 1 << 20

func FormatCSS(input string, minify bool) (string, error) {
	if len(input) == 0 || len(input) > maxInput {
		return "", fmt.Errorf("CSS must contain between 1 and %d bytes", maxInput)
	}
	var out strings.Builder
	depth := 0
	quote := rune(0)
	comment := false
	space := false
	trimFollowingSpace := false
	writeSpace := func() {
		// Whitespace between ordinary CSS tokens can be meaningful: it may
		// separate selector terms or values such as `0 auto`. Delimiter cases
		// below clear it deliberately, while this path preserves one space.
		if space && out.Len() > 0 {
			out.WriteByte(' ')
		}
		space = false
	}
	runes := []rune(input)
	for i, ch := range runes {
		if comment {
			if ch == '/' && i > 0 && runes[i-1] == '*' {
				comment = false
			}
			continue
		}
		if quote != 0 {
			out.WriteRune(ch)
			if ch == quote && (i == 0 || runes[i-1] != '\\') {
				quote = 0
				trimFollowingSpace = false
			}
			continue
		}
		if ch == '/' && i+1 < len(runes) && runes[i+1] == '*' {
			comment = true
			continue
		}
		if ch == '\'' || ch == '"' {
			writeSpace()
			quote = ch
			out.WriteRune(ch)
			trimFollowingSpace = false
			continue
		}
		if unicode.IsSpace(ch) {
			if minify && trimFollowingSpace {
				continue
			}
			space = true
			continue
		}
		switch ch {
		case '{':
			space = false
			out.WriteString("{")
			trimFollowingSpace = true
			depth++
			if !minify {
				out.WriteByte('\n')
				out.WriteString(strings.Repeat("  ", depth))
			}
		case ';':
			out.WriteByte(';')
			trimFollowingSpace = true
			if !minify {
				out.WriteByte('\n')
				out.WriteString(strings.Repeat("  ", depth))
			}
		case '}':
			depth--
			if depth < 0 {
				return "", fmt.Errorf("CSS has an unmatched closing brace")
			}
			if !minify {
				out.WriteByte('\n')
				out.WriteString(strings.Repeat("  ", depth))
			}
			out.WriteByte('}')
			if !minify {
				out.WriteByte('\n')
			}
		case ':', ',':
			out.WriteRune(ch)
			if ch == ',' && !minify {
				out.WriteByte(' ')
			}
			space = false
			trimFollowingSpace = true
		default:
			writeSpace()
			out.WriteRune(ch)
			trimFollowingSpace = false
		}
	}
	if comment || quote != 0 || depth != 0 {
		return "", fmt.Errorf("CSS has an unfinished comment, string, or block")
	}
	result := strings.TrimSpace(out.String())
	if minify {
		result = strings.ReplaceAll(result, ";}", "}")
	}
	return result + "\n", nil
}

func HTMLToJSX(input string) (string, error) {
	if len(input) == 0 || len(input) > maxInput {
		return "", fmt.Errorf("HTML must contain between 1 and %d bytes", maxInput)
	}
	document, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", fmt.Errorf("parse HTML: %w", err)
	}
	body := findElement(document, "body")
	if body == nil {
		return "", fmt.Errorf("parse HTML: no body element")
	}
	var out strings.Builder
	for node := body.FirstChild; node != nil; node = node.NextSibling {
		renderJSX(&out, node)
	}
	return out.String(), nil
}

func findElement(node *html.Node, name string) *html.Node {
	if node.Type == html.ElementNode && node.Data == name {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findElement(child, name); found != nil {
			return found
		}
	}
	return nil
}

var voidTags = map[string]bool{"area": true, "base": true, "br": true, "col": true, "embed": true, "hr": true, "img": true, "input": true, "link": true, "meta": true, "param": true, "source": true, "track": true, "wbr": true}

func renderJSX(out *strings.Builder, n *html.Node) {
	switch n.Type {
	case html.TextNode:
		out.WriteString(strings.ReplaceAll(strings.ReplaceAll(n.Data, "{", "{'{'}"), "}", "{'}'}"))
	case html.CommentNode:
	case html.ElementNode:
		tag := n.Data
		out.WriteByte('<')
		out.WriteString(tag)
		for _, a := range n.Attr {
			name := a.Key
			if name == "class" {
				name = "className"
			}
			if name == "for" {
				name = "htmlFor"
			}
			out.WriteByte(' ')
			out.WriteString(name)
			out.WriteString("=")
			out.WriteString(fmt.Sprintf("%q", a.Val))
		}
		if voidTags[tag] {
			out.WriteString(" />")
			return
		}
		out.WriteByte('>')
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			renderJSX(out, c)
		}
		out.WriteString("</")
		out.WriteString(tag)
		out.WriteByte('>')
	}
}

func DecodeQR(contents []byte) (string, error) {
	if len(contents) == 0 || len(contents) > 4<<20 {
		return "", fmt.Errorf("image must contain between 1 byte and 4 MiB")
	}
	img, _, err := image.Decode(bytes.NewReader(contents))
	if err != nil {
		return "", fmt.Errorf("decode image: %w", err)
	}
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	result, err := qrcode.NewQRCodeReader().Decode(bitmap, nil)
	if err != nil {
		return "", fmt.Errorf("decode QR code: %w", err)
	}
	return result.GetText(), nil
}

func ToDataURL(contents []byte, mimeType string) (string, error) {
	if len(contents) == 0 || len(contents) > 4<<20 {
		return "", fmt.Errorf("file must contain between 1 byte and 4 MiB")
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(contents), nil
}
func FromDataURL(input string) ([]byte, string, error) {
	if len(input) > 6<<20 {
		return nil, "", fmt.Errorf("data URL exceeds 6 MiB")
	}
	parts := strings.SplitN(input, ",", 2)
	if len(parts) != 2 || !strings.HasPrefix(parts[0], "data:") || !strings.HasSuffix(parts[0], ";base64") {
		return nil, "", fmt.Errorf("enter a base64 data URL")
	}
	contents, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil || len(contents) == 0 || len(contents) > 4<<20 {
		return nil, "", fmt.Errorf("data URL has invalid or oversized base64 content")
	}
	return contents, strings.TrimSuffix(strings.TrimPrefix(parts[0], "data:"), ";base64"), nil
}

func SlashEscape(input string, unescape bool) (string, error) {
	if len(input) > maxInput {
		return "", fmt.Errorf("text exceeds 1 MiB")
	}
	if !unescape {
		r := strings.NewReplacer("\\", "\\\\", "\n", "\\n", "\r", "\\r", "\t", "\\t", "\"", "\\\"")
		return r.Replace(input), nil
	}
	var out strings.Builder
	run := []rune(input)
	for i := 0; i < len(run); i++ {
		if run[i] != '\\' {
			out.WriteRune(run[i])
			continue
		}
		if i+1 == len(run) {
			return "", fmt.Errorf("input ends with an escape")
		}
		i++
		switch run[i] {
		case 'n':
			out.WriteByte('\n')
		case 'r':
			out.WriteByte('\r')
		case 't':
			out.WriteByte('\t')
		case '\\':
			out.WriteByte('\\')
		case '"':
			out.WriteByte('"')
		default:
			return "", fmt.Errorf("unsupported escape \\%c", run[i])
		}
	}
	return out.String(), nil
}

func JSONCode(input, name, target string) (string, error) {
	if len(input) == 0 || len(input) > maxInput {
		return "", fmt.Errorf("JSON must contain between 1 and %d bytes", maxInput)
	}
	if target == "typescript" {
		return conversiontools.JSONToTypeScript(input, name)
	}
	var value any
	dec := json.NewDecoder(strings.NewReader(input))
	dec.UseNumber()
	if err := dec.Decode(&value); err != nil {
		return "", fmt.Errorf("parse JSON: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return "", fmt.Errorf("JSON must contain exactly one value")
		}
		return "", fmt.Errorf("parse JSON: %w", err)
	}
	root := goName(name)
	if target == "go" {
		object, ok := value.(map[string]any)
		if !ok {
			return "type " + root + " " + goType(value) + "\n", nil
		}
		var out strings.Builder
		out.WriteString("type ")
		out.WriteString(root)
		out.WriteString(" struct {\n")
		for key, fieldValue := range object {
			out.WriteString("\t")
			out.WriteString(goName(key))
			out.WriteByte(' ')
			out.WriteString(goType(fieldValue))
			out.WriteString(" `json:\"")
			out.WriteString(key)
			out.WriteString("\"`\n")
		}
		out.WriteString("}\n")
		return out.String(), nil
	}
	if target == "python" {
		object, ok := value.(map[string]any)
		if !ok {
			return root + " = " + pythonType(value) + "\n", nil
		}
		var out strings.Builder
		out.WriteString("from typing import Any, TypedDict\n\nclass ")
		out.WriteString(root)
		out.WriteString("(TypedDict):\n")
		if len(object) == 0 {
			out.WriteString("    pass\n")
		}
		for key, fieldValue := range object {
			if regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(key) {
				out.WriteString("    ")
				out.WriteString(key)
				out.WriteString(": ")
				out.WriteString(pythonType(fieldValue))
				out.WriteByte('\n')
			} else {
				out.WriteString("    # ")
				out.WriteString(fmt.Sprintf("%q", key))
				out.WriteString(" cannot be represented as a TypedDict attribute.\n")
			}
		}
		return out.String(), nil
	}
	return "", fmt.Errorf("choose a supported target")
}

func goType(value any) string {
	switch value.(type) {
	case nil:
		return "any"
	case bool:
		return "bool"
	case string:
		return "string"
	case json.Number:
		return "float64"
	case []any:
		return "[]any"
	case map[string]any:
		return "map[string]any"
	default:
		return "any"
	}
}

func pythonType(value any) string {
	switch value.(type) {
	case nil:
		return "None"
	case bool:
		return "bool"
	case string:
		return "str"
	case json.Number:
		return "float"
	case []any:
		return "list[Any]"
	case map[string]any:
		return "dict[str, Any]"
	default:
		return "Any"
	}
}
func goName(value string) string {
	parts := regexp.MustCompile(`[^A-Za-z0-9]+`).Split(value, -1)
	var out strings.Builder
	for _, part := range parts {
		if part != "" {
			out.WriteString(strings.ToUpper(part[:1]))
			out.WriteString(part[1:])
		}
	}
	if out.Len() == 0 {
		return "Root"
	}
	return out.String()
}

type StringFacts struct {
	Bytes, Runes, Words, Lines int
	ASCII                      bool
}

func InspectString(input string) StringFacts {
	f := StringFacts{Bytes: len(input), Runes: len([]rune(input)), ASCII: true}
	if input != "" {
		f.Lines = strings.Count(input, "\n") + 1
		f.Words = len(strings.Fields(input))
	}
	for _, r := range input {
		if r > 127 {
			f.ASCII = false
			break
		}
	}
	return f
}
func FilenameForMIME(mimeType string) string {
	exts, _ := mime.ExtensionsByType(mimeType)
	if len(exts) > 0 {
		return "cmdry-file" + exts[0]
	}
	return "cmdry-file.bin"
}
