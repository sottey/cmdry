package main

import (
	"fmt"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.most-popular", Name: "Find Most Popular", Version: "0.1.0", Description: "Count and rank pasted newline-delimited list items locally.", Category: "text", Icon: "list",
		SearchTerms: []string{"list", "popular", "frequency", "count", "duplicates", "rank"}, Pages: []cmdry.Page{{ID: "overview", Name: "Count", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New list", Method: "read"}, {ID: "count", Name: "Rank items", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "count": count}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Find Most Popular", Components: []cmdry.Component{{Type: "form", Title: "Rank list items", Action: "count", Submit: "Rank items", Description: "Uses one item per line, ignores blank lines, and counts exact item matches locally.", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}}}}}, nil
}

func count(request cmdry.Request) (cmdry.View, error) {
	items, err := listtools.FindMostPopular(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	rows := make([]map[string]any, 0, len(items))
	for index, item := range items {
		rows = append(rows, map[string]any{"rank": index + 1, "item": item.Item, "count": item.Count})
	}
	return cmdry.View{Title: "Most popular items", Components: []cmdry.Component{{Type: "metric", Label: "Unique items", Value: fmt.Sprint(len(items))}, {Type: "table", ID: "popular-items", Columns: []cmdry.Column{{Key: "rank", Label: "Rank"}, {Key: "item", Label: "Item"}, {Key: "count", Label: "Occurrences"}}, Rows: rows}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Rank another list"}}}}}, nil
}
