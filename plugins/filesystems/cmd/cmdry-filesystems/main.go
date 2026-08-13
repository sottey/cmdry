package main

import (
	"sort"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/filesystems"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1,
		ID:              "com.sottey.filesystems",
		Name:            "Filesystem Inspector",
		Version:         "0.1.0",
		Description:     "Inspect mounted filesystem capacity and available space.",
		Category:        "storage",
		Icon:            "storage",
		Pages:           []cmdry.Page{{ID: "overview", Name: "Filesystems", Default: true, Action: "list"}},
		Permissions:     []string{"storage.read"},
		Actions:         []cmdry.Action{{ID: "list", Name: "Refresh", Method: "read"}},
	}, Actions: map[string]cmdry.Handler{"list": list}})
}

func list(_ cmdry.Request) (cmdry.View, error) {
	items, err := filesystems.Collect()
	if err != nil {
		return cmdry.View{}, err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Capacity > items[j].Capacity })
	highestUse := "0%"
	if len(items) > 0 {
		highestUse = strconv.Itoa(items[0].Capacity) + "%"
	}
	rows := make([]map[string]any, 0, len(items))
	for _, item := range items {
		rows = append(rows, item.Row())
	}
	return cmdry.View{Title: "Mounted Filesystems", Components: []cmdry.Component{
		{Type: "metric", Label: "Mounted filesystems", Value: strconv.Itoa(len(items))},
		{Type: "metric", Label: "Highest utilization", Value: highestUse, Description: "Across listed mounts"},
		{Type: "table", ID: "filesystems", Columns: []cmdry.Column{
			{Key: "filesystem", Label: "Filesystem"}, {Key: "total", Label: "Total"}, {Key: "used", Label: "Used"},
			{Key: "available", Label: "Available"}, {Key: "capacity", Label: "Use"}, {Key: "mount", Label: "Mounted on"},
		}, Rows: rows},
	}}, nil
}
