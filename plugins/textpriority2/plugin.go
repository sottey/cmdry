package textpriority2

import (
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunNATOAlphabetConverter() {
	cmdry.Run(plugin("nato-alphabet-converter", "NATO Alphabet Converter", "Encode or decode pasted text using the NATO phonetic alphabet locally.", []string{"nato", "phonetic", "alphabet", "encode", "decode"}, natoForm, natoRun))
}

func RunASCIIBinaryConverter() {
	cmdry.Run(plugin("ascii-binary-converter", "ASCII Binary Converter", "Convert pasted ASCII text to or from eight-bit binary bytes locally.", []string{"ascii", "binary", "bytes", "encode", "decode"}, binaryForm, binaryRun))
}

func RunLoremIpsumGenerator() {
	cmdry.Run(plugin("lorem-ipsum-generator", "Lorem Ipsum Generator", "Generate bounded Lorem ipsum placeholder text locally.", []string{"lorem", "ipsum", "placeholder", "text", "generate"}, loremForm, loremRun))
}

func RunEmojiPicker() {
	cmdry.Run(plugin("emoji-picker", "Emoji Picker", "Select and copy a Unicode emoji from Cmdry's bundled local catalog.", []string{"emoji", "unicode", "picker", "symbol", "copy"}, emojiForm, emojiRun))
}

func RunStringObfuscator() {
	cmdry.Run(plugin("string-obfuscator", "String Obfuscator", "Apply a reversible ROT47 text transform locally; it is not encryption.", []string{"obfuscate", "rot47", "text", "reversible", "encode"}, obfuscatorForm, obfuscatorRun))
}

func RunNumeronymGenerator() {
	cmdry.Run(plugin("numeronym-generator", "Numeronym Generator", "Create numeronyms such as i18n from pasted words locally.", []string{"numeronym", "i18n", "abbreviation", "text", "words"}, numeronymForm, numeronymRun))
}

func RunASCIIArtTextGenerator() {
	cmdry.Run(plugin("ascii-art-text-generator", "ASCII Art Text Generator", "Render bounded text as ASCII art using a bundled block font locally.", []string{"ascii", "art", "banner", "text", "font"}, artForm, artRun))
}

func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "text", Icon: "text", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": action}}
}

func natoForm(cmdry.Request) (cmdry.View, error) {
	return form("NATO Alphabet Converter", "Letters and digits only; use / between words when decoding.", "Convert", []cmdry.Field{{Name: "operation", Label: "Operation", Type: "select", Value: "encode", Options: []cmdry.Option{{Value: "encode", Label: "Text to NATO"}, {Value: "decode", Label: "NATO to text"}}}, {Name: "input", Label: "Text or NATO words", Type: "textarea", Value: "Cmdry 2", Required: true}}), nil
}

func natoRun(request cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(request.Params["input"])
	if fmt.Sprint(request.Params["operation"]) == "decode" {
		output, err := NATODecode(input)
		if err != nil {
			return cmdry.View{}, err
		}
		return result("Decoded NATO text", output, "Convert another value"), nil
	}
	output, err := NATOEncode(input)
	if err != nil {
		return cmdry.View{}, err
	}
	return result("NATO alphabet", output, "Convert another value"), nil
}

func binaryForm(cmdry.Request) (cmdry.View, error) {
	return form("ASCII Binary Converter", "ASCII only. Binary input uses whitespace-separated eight-bit bytes.", "Convert", []cmdry.Field{{Name: "operation", Label: "Operation", Type: "select", Value: "encode", Options: []cmdry.Option{{Value: "encode", Label: "Text to binary"}, {Value: "decode", Label: "Binary to text"}}}, {Name: "input", Label: "Text or binary bytes", Type: "textarea", Value: "Cmdry", Required: true}}), nil
}

func binaryRun(request cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(request.Params["input"])
	if fmt.Sprint(request.Params["operation"]) == "decode" {
		output, err := ASCIIBinaryDecode(input)
		if err != nil {
			return cmdry.View{}, err
		}
		return result("Decoded ASCII text", output, "Convert another value"), nil
	}
	output, err := ASCIIBinaryEncode(input)
	if err != nil {
		return cmdry.View{}, err
	}
	return result("ASCII binary", output, "Convert another value"), nil
}

func loremForm(cmdry.Request) (cmdry.View, error) {
	return form("Lorem Ipsum Generator", "Generate one to twenty deterministic placeholder paragraphs.", "Generate placeholder text", []cmdry.Field{{Name: "paragraphs", Label: "Paragraphs", Type: "number", Value: "3", Min: "1", Max: "20", Required: true}}), nil
}

func loremRun(request cmdry.Request) (cmdry.View, error) {
	count, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(request.Params["paragraphs"])))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("paragraph count must be a whole number")
	}
	output, err := LoremIpsum(count)
	if err != nil {
		return cmdry.View{}, err
	}
	return result("Lorem ipsum", output, "Generate more text"), nil
}

func emojiForm(cmdry.Request) (cmdry.View, error) {
	items := EmojiOptions()
	options := make([]cmdry.Option, len(items))
	for index, item := range items {
		options[index] = cmdry.Option{Value: item.Value, Label: item.Label}
	}
	return form("Emoji Picker", "Choose an emoji from Cmdry's bundled local catalog, then copy it from the result.", "Select emoji", []cmdry.Field{{Name: "emoji", Label: "Emoji", Type: "select", Value: "😀", Options: options}}), nil
}

func emojiRun(request cmdry.Request) (cmdry.View, error) {
	selected := fmt.Sprint(request.Params["emoji"])
	for _, item := range EmojiOptions() {
		if selected == item.Value {
			return result("Selected emoji", selected, "Choose another emoji"), nil
		}
	}
	return cmdry.View{}, fmt.Errorf("choose an emoji from the bundled catalog")
}

func obfuscatorForm(cmdry.Request) (cmdry.View, error) {
	return form("String Obfuscator", "ROT47 is reversible obfuscation, not encryption. Do not use it to protect secrets.", "Apply ROT47", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}}), nil
}

func obfuscatorRun(request cmdry.Request) (cmdry.View, error) {
	output, err := ROT47(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "ROT47 output", Components: []cmdry.Component{{Type: "alert", Level: "warning", Title: "Obfuscation is not encryption", Message: "Applying ROT47 again restores the original text."}, {Type: "code", Title: "ROT47 output", Text: output}, restart("Obfuscate another string")}}, nil
}

func numeronymForm(cmdry.Request) (cmdry.View, error) {
	return form("Numeronym Generator", "Words longer than three ASCII letters become forms such as internationalization → i18n.", "Create numeronyms", []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Value: "internationalization accessibility", Required: true}}), nil
}

func numeronymRun(request cmdry.Request) (cmdry.View, error) {
	output, err := Numeronymize(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return result("Numeronyms", output, "Create more numeronyms"), nil
}

func artForm(cmdry.Request) (cmdry.View, error) {
	return form("ASCII Art Text Generator", "Render up to 40 letters, digits, spaces, hyphens, exclamation points, or question marks.", "Render ASCII art", []cmdry.Field{{Name: "input", Label: "Text", Type: "text", Value: "CMDRY", Required: true}}), nil
}

func artRun(request cmdry.Request) (cmdry.View, error) {
	output, err := ASCIIArt(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return result("ASCII art", output, "Render more text"), nil
}

func form(title, description, submit string, fields []cmdry.Field) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Description: description + " Your data never leaves this device.", Action: "run", Submit: submit, Fields: fields}}}
}

func result(title, output, again string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "code", Title: title, Text: output}, restart(again)}}
}

func restart(name string) cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: name}}}
}
