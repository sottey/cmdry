package main

import (
	"sort"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/scheduled"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1,
		ID:              "com.sottey.scheduled",
		Name:            "Scheduled Tasks",
		Version:         "0.1.0",
		Description:     "Inspect user cron jobs and platform-native scheduled tasks.",
		Category:        "server",
		Icon:            "clock",
		Pages:           []cmdry.Page{{ID: "overview", Name: "Scheduled", Default: true, Action: "list"}},
		Permissions:     []string{"system.read"},
		Actions:         []cmdry.Action{{ID: "list", Name: "Refresh", Method: "read"}},
	}, Actions: map[string]cmdry.Handler{"list": listTasks}})
}

func listTasks(_ cmdry.Request) (cmdry.View, error) {
	tasks, notice, err := scheduled.Collect()
	if err != nil {
		return cmdry.View{}, err
	}
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Source == tasks[j].Source {
			return tasks[i].Command < tasks[j].Command
		}
		return tasks[i].Source < tasks[j].Source
	})
	rows := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		rows = append(rows, task.Row())
	}
	components := []cmdry.Component{{Type: "metric", Label: "Scheduled tasks", Value: strconv.Itoa(len(tasks))}}
	if notice != "" {
		components = append(components, cmdry.Component{Type: "alert", Level: "warning", Title: "Partial visibility", Message: notice})
	}
	components = append(components, cmdry.Component{Type: "table", ID: "scheduled", Columns: []cmdry.Column{
		{Key: "source", Label: "Source"}, {Key: "schedule", Label: "Schedule"}, {Key: "command", Label: "Task"},
	}, Rows: rows})
	return cmdry.View{Title: "Scheduled Tasks", Components: components}, nil
}
