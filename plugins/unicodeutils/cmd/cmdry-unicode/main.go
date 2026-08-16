package main

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/unicodeutils"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.unicode", Name: "Unicode Encoder / Decoder", Version: "0.1.0", Description: "Encode text as Unicode escapes or decode Unicode escapes locally.", Category: "text", Icon: "text",
		SearchTerms: []string{"unicode", "encode", "decode", "escape", "utf-16", "code point"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Transform Unicode", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": form, "run": transform}})
}

func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Unicode Encoder / Decoder", Components: []cmdry.Component{{Type: "form", Title: "Transform Unicode", Action: "run", Submit: "Transform text", Description: "Encoding uses \\uXXXX escapes and UTF-16 surrogate pairs when needed. Decoding accepts \\uXXXX and \\UXXXXXXXX escapes.", Fields: []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}, {Name: "mode", Label: "Operation", Type: "select", Value: "encode", Options: []cmdry.Option{{Value: "encode", Label: "Encode as Unicode escapes"}, {Value: "decode", Label: "Decode Unicode escapes"}}}}}}}, nil
}

func transform(request cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(request.Params["input"])
	if strings.TrimSpace(input) == "" {
		return cmdry.View{}, fmt.Errorf("text is required")
	}
	output := ""
	if fmt.Sprint(request.Params["mode"]) == "decode" {
		var err error
		output, err = unicodeutils.Decode(input)
		if err != nil {
			return cmdry.View{}, err
		}
	} else {
		output = unicodeutils.Encode(input)
	}
	return cmdry.View{Title: "Unicode result", Components: []cmdry.Component{{Type: "code", Title: "Unicode result", Text: output}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Transform more text"}}}}}, nil
}
