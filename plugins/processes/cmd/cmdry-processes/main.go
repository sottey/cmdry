package main

import (
	"sort"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/processes"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1,
		ID:              "com.sottey.processes",
		Name:            "Process Resource Snapshot",
		Version:         "0.1.0",
		Description:     "Inspect CPU and memory usage of visible local processes.",
		Category:        "system",
		Icon:            "activity",
		Pages:           []cmdry.Page{{ID: "overview", Name: "Processes", Default: true, Action: "list"}},
		Permissions:     []string{"process.read"},
		Actions:         []cmdry.Action{{ID: "list", Name: "Refresh", Method: "read"}},
	}, Actions: map[string]cmdry.Handler{"list": list}})
}

func list(_ cmdry.Request) (cmdry.View, error) {
	items, err := processes.Collect()
	if err != nil {
		return cmdry.View{}, err
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].CPU == items[j].CPU {
			return items[i].Memory > items[j].Memory
		}
		return items[i].CPU > items[j].CPU
	})
	running, totalCPU, totalMemory := processes.Summary(items)
	rows := make([]map[string]any, 0, len(items))
	for _, item := range items {
		rows = append(rows, item.Row())
	}
	return cmdry.View{Title: "Process Resource Snapshot", Components: []cmdry.Component{
		{Type: "metric", Label: "Visible processes", Value: strconv.Itoa(len(items))},
		{Type: "metric", Label: "Running", Value: strconv.Itoa(running)},
		{Type: "metric", Label: "Total CPU", Value: strconv.FormatFloat(totalCPU, 'f', 1, 64) + "%"},
		{Type: "metric", Label: "Total memory", Value: strconv.FormatFloat(totalMemory, 'f', 1, 64) + "%"},
		{Type: "table", ID: "processes", Columns: []cmdry.Column{
			{Key: "pid", Label: "PID"}, {Key: "parent", Label: "Parent PID"}, {Key: "cpu", Label: "CPU"},
			{Key: "memory", Label: "Memory"}, {Key: "state", Label: "State"}, {Key: "command", Label: "Command"},
		}, Rows: rows},
	}}, nil
}
