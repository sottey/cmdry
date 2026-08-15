package main

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/jsoncompare"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.json-compare", Name: "JSON Compare", Version: "0.1.0",
		Description: "Compare two JSON documents structurally, ignoring object key order.", Category: "developer", Icon: "compare",
		SearchTerms: []string{"json", "diff", "compare", "structural"},
		Pages:       []cmdry.Page{{ID: "overview", Name: "Compare", Default: true, Action: "overview"}},
		Permissions: []string{"data.transform"},
		Actions:     []cmdry.Action{{ID: "overview", Name: "New comparison", Method: "read"}, {ID: "compare", Name: "Compare JSON", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "compare": compare}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "JSON Compare", Components: []cmdry.Component{{
		Type: "form", Title: "Compare JSON documents", Action: "compare", Submit: "Compare JSON",
		Description: "Object property order and whitespace do not count as changes. Arrays are compared in their existing order.",
		Fields:      []cmdry.Field{{Name: "left", Label: "Original JSON", Type: "textarea", Required: true}, {Name: "right", Label: "Revised JSON", Type: "textarea", Required: true}},
	}}}, nil
}

func compare(request cmdry.Request) (cmdry.View, error) {
	left, _ := request.Params["left"].(string)
	right, _ := request.Params["right"].(string)
	if strings.TrimSpace(left) == "" || strings.TrimSpace(right) == "" {
		return cmdry.View{}, fmt.Errorf("both JSON inputs are required")
	}
	result, err := jsoncompare.Compare(left, right)
	if err != nil {
		return cmdry.View{}, err
	}
	if len(result.Differences) == 0 {
		return cmdry.View{Title: "JSON Compare", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "No structural differences", Message: "The documents contain the same JSON values."}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Compare another pair"}}}}}, nil
	}
	lines := make([]string, 0, len(result.Differences)*2)
	for _, difference := range result.Differences {
		switch difference.Kind {
		case "added":
			lines = append(lines, "+ "+difference.Path+" = "+difference.Right)
		case "removed":
			lines = append(lines, "- "+difference.Path+" = "+difference.Left)
		default:
			lines = append(lines, "~ "+difference.Path, "  - "+difference.Left, "  + "+difference.Right)
		}
	}
	message := fmt.Sprintf("%d difference(s) found.", len(result.Differences))
	if result.Truncated {
		message += " Output stopped after 500 differences."
	}
	return cmdry.View{Title: "JSON differences", Components: []cmdry.Component{{Type: "alert", Level: "info", Title: "Comparison complete", Message: message}, {Type: "code", Title: "Structural diff", Text: strings.Join(lines, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Compare another pair"}}}}}, nil
}
