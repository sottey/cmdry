package priority1tools

import (
	"strings"
	"testing"

	qrcode "github.com/skip2/go-qrcode"
)

func TestTransforms(t *testing.T) {
	css, err := FormatCSS("/* note */ body { color: red; }", true)
	if err != nil || css != "body{color:red}\n" {
		t.Fatalf("css=%q err=%v", css, err)
	}
	css, err = FormatCSS("@media screen and (min-width: 10px) { .card .title { margin: 0 auto; } }", true)
	if err != nil || css != "@media screen and (min-width:10px){.card .title{margin:0 auto}}\n" {
		t.Fatalf("semantic css=%q err=%v", css, err)
	}
	jsx, err := HTMLToJSX(`<label class="name" for="email">Email</label>`)
	if err != nil || !strings.Contains(jsx, `className="name"`) || !strings.Contains(jsx, `htmlFor="email"`) {
		t.Fatalf("jsx=%q err=%v", jsx, err)
	}
	url, err := ToDataURL([]byte("hello"), "text/plain")
	if err != nil {
		t.Fatal(err)
	}
	decoded, kind, err := FromDataURL(url)
	if err != nil || string(decoded) != "hello" || kind != "text/plain" {
		t.Fatalf("data url=%q kind=%q err=%v", decoded, kind, err)
	}
	escaped, err := SlashEscape("a\nb", false)
	if err != nil || escaped != `a\nb` {
		t.Fatalf("escaped=%q err=%v", escaped, err)
	}
	unescaped, err := SlashEscape(escaped, true)
	if err != nil || unescaped != "a\nb" {
		t.Fatalf("unescaped=%q err=%v", unescaped, err)
	}
	goCode, err := JSONCode(`{"name":"Ada","active":true}`, "Person", "go")
	if err != nil || !strings.Contains(goCode, "Name string") || !strings.Contains(goCode, "Active bool") {
		t.Fatalf("go code=%q err=%v", goCode, err)
	}
	pythonCode, err := JSONCode(`{"name":"Ada"}`, "Person", "python")
	if err != nil || !strings.Contains(pythonCode, "class Person(TypedDict)") || !strings.Contains(pythonCode, "name: str") {
		t.Fatalf("python code=%q err=%v", pythonCode, err)
	}
}

func TestQRAndFacts(t *testing.T) {
	png, err := qrcode.Encode("https://cmdry.test", qrcode.Medium, 128)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeQR(png)
	if err != nil || decoded != "https://cmdry.test" {
		t.Fatalf("decoded=%q err=%v", decoded, err)
	}
	facts := InspectString("hi\n🙂")
	if facts.Lines != 2 || facts.Runes != 4 || facts.ASCII {
		t.Fatalf("facts=%#v", facts)
	}
}
