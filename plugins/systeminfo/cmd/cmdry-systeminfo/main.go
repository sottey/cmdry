package main

import (
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/systeminfo"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1,
		ID:              "com.sottey.systeminfo",
		Name:            "System Information",
		Version:         "0.1.0",
		Description:     "Inspect local operating system, CPU, memory, and hardware facts.",
		Category:        "server",
		Icon:            "computer",
		Pages:           []cmdry.Page{{ID: "overview", Name: "System", Default: true, Action: "list"}},
		Permissions:     []string{"system.read"},
		Actions:         []cmdry.Action{{ID: "list", Name: "Refresh", Method: "read"}},
	}, Actions: map[string]cmdry.Handler{"list": list}})
}

func list(_ cmdry.Request) (cmdry.View, error) {
	info, err := systeminfo.Collect()
	if err != nil {
		return cmdry.View{}, err
	}
	cores := "Unavailable"
	if info.Cores > 0 {
		cores = strconv.Itoa(info.Cores)
	}
	return cmdry.View{Title: "System Information", Components: []cmdry.Component{
		{Type: "metric", Label: "Operating system", Value: fallback(info.OS)},
		{Type: "metric", Label: "CPU cores", Value: cores},
		{Type: "metric", Label: "Total memory", Value: systeminfo.HumanBytes(info.MemoryTotal)},
		{Type: "metric", Label: "Available memory", Value: systeminfo.HumanBytes(info.MemoryAvailable)},
		{Type: "table", ID: "systeminfo", Columns: []cmdry.Column{{Key: "label", Label: "Property"}, {Key: "value", Label: "Value"}}, Rows: info.Rows()},
	}}, nil
}

func fallback(value string) string {
	if value == "" {
		return "Unavailable"
	}
	return value
}
