package main

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/hiddenchars"
)

const maxResults = 500

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.hidden-characters", Name: "Hidden Character Detector", Version: "0.1.0",
		Description: "Find invisible Unicode characters and confusing whitespace in pasted text.", Category: "text", Icon: "search",
		SearchTerms: []string{"unicode", "hidden", "invisible", "zero width", "bidi", "whitespace"},
		Pages:       []cmdry.Page{{ID: "overview", Name: "Detect", Default: true, Action: "overview"}},
		Permissions: []string{"data.transform"},
		Actions:     []cmdry.Action{{ID: "overview", Name: "New scan", Method: "read"}, {ID: "scan", Name: "Scan text", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "scan": scan}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Hidden Character Detector", Components: []cmdry.Component{{
		Type: "form", Title: "Scan pasted text", Action: "scan", Submit: "Find hidden characters",
		Description: "Finds zero-width and bidirectional formatting characters, non-standard spaces, and control characters. Ordinary spaces, tabs, and line breaks are ignored.",
		Fields:      []cmdry.Field{{Name: "text", Label: "Text", Type: "textarea", Required: true, Description: "Your text is inspected locally by the plugin and is never uploaded."}},
	}}}, nil
}

func scan(request cmdry.Request) (cmdry.View, error) {
	input, _ := request.Params["text"].(string)
	if input == "" {
		return cmdry.View{}, fmt.Errorf("text is required")
	}
	findings := hiddenchars.Detect(input)
	if len(findings) == 0 {
		return cmdry.View{Title: "Hidden Character Detector", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "No hidden characters found", Message: "No Unicode formatting, non-standard spacing, or unexpected control characters were found."}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Scan other text"}}}}}, nil
	}
	shown := findings
	if len(shown) > maxResults {
		shown = shown[:maxResults]
	}
	rows := make([]map[string]any, 0, len(shown))
	for _, finding := range shown {
		rows = append(rows, map[string]any{"line": finding.Line, "column": finding.Column, "code_point": finding.CodePoint, "description": finding.Description})
	}
	message := fmt.Sprintf("Found %d potentially hidden character(s).", len(findings))
	if len(findings) > len(shown) {
		message += " Showing the first " + fmt.Sprint(maxResults) + "."
	}
	return cmdry.View{Title: "Hidden characters found", Components: []cmdry.Component{{Type: "alert", Level: "warning", Title: "Review before sharing or parsing", Message: strings.TrimSpace(message)}, {Type: "table", ID: "hidden-characters", Columns: []cmdry.Column{{Key: "line", Label: "Line"}, {Key: "column", Label: "Column"}, {Key: "code_point", Label: "Code point"}, {Key: "description", Label: "Character"}}, Rows: rows}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Scan other text"}}}}}, nil
}
