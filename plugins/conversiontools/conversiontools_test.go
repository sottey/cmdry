package conversiontools

import (
	"strings"
	"testing"
)

func TestTOMLConversions(t *testing.T) {
	tomlOutput, err := JSONToTOML(`{"name":"Cmdry","count":2}`)
	if err != nil || !strings.Contains(string(tomlOutput), "name = \"Cmdry\"") {
		t.Fatalf("JSONToTOML = %q, %v", tomlOutput, err)
	}
	jsonOutput, err := TOMLToJSON("name = \"Cmdry\"\ncount = 2\n")
	if err != nil || !strings.Contains(string(jsonOutput), `"count": 2`) {
		t.Fatalf("TOMLToJSON = %q, %v", jsonOutput, err)
	}
	if _, err := FormatTOML("invalid = ["); err == nil {
		t.Fatal("FormatTOML accepted invalid TOML")
	}
}

func TestSpreadsheetConversions(t *testing.T) {
	xlsx, rows, err := CSVToXLSX([]byte("name,value\nAda,=1+1\n"), ',')
	if err != nil || rows != 2 {
		t.Fatalf("CSVToXLSX rows=%d err=%v", rows, err)
	}
	csvOutput, rows, err := XLSXToDelimited(xlsx, ',')
	if err != nil || rows != 2 || !strings.Contains(string(csvOutput), "'=1+1") {
		t.Fatalf("XLSXToDelimited rows=%d output=%q err=%v", rows, csvOutput, err)
	}
}

func TestMarkdownAndTypeScriptConversions(t *testing.T) {
	html, err := MarkdownToHTML("# Hello")
	if err != nil || string(html) != "<h1>Hello</h1>\n" {
		t.Fatalf("MarkdownToHTML = %q, %v", html, err)
	}
	types, err := JSONToTypeScript(`{"name":"Ada","active":true,"scores":[1,2]}`, "Person")
	if err != nil || !strings.Contains(types, "export type Person") || !strings.Contains(types, "scores: (number)[]") {
		t.Fatalf("JSONToTypeScript = %q, %v", types, err)
	}
}
