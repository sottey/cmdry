package main

import (
	"sort"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/network"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1,
		ID:              "com.sottey.network",
		Name:            "Network Interface Inspector",
		Version:         "0.1.0",
		Description:     "Inspect local interfaces, assigned addresses, and the default gateway.",
		Category:        "server",
		Icon:            "network",
		Pages:           []cmdry.Page{{ID: "overview", Name: "Network", Default: true, Action: "list"}},
		Permissions:     []string{"network.read"},
		Actions:         []cmdry.Action{{ID: "list", Name: "Refresh", Method: "read"}},
	}, Actions: map[string]cmdry.Handler{"list": list}})
}

func list(_ cmdry.Request) (cmdry.View, error) {
	interfaces, gateway, err := network.Collect()
	if err != nil {
		return cmdry.View{}, err
	}
	sort.Slice(interfaces, func(i, j int) bool { return interfaces[i].Name < interfaces[j].Name })
	rows := make([]map[string]any, 0, len(interfaces))
	addressCount := 0
	for _, item := range interfaces {
		rows = append(rows, item.Row())
		addressCount += len(item.Addresses)
	}
	if gateway == "" {
		gateway = "Unavailable"
	}
	return cmdry.View{Title: "Network Interfaces", Components: []cmdry.Component{
		{Type: "metric", Label: "Interfaces", Value: strconv.Itoa(len(interfaces))},
		{Type: "metric", Label: "Assigned addresses", Value: strconv.Itoa(addressCount)},
		{Type: "metric", Label: "Default gateway", Value: gateway},
		{Type: "table", ID: "network", Columns: []cmdry.Column{{Key: "interface", Label: "Interface"}, {Key: "status", Label: "Status"}, {Key: "addresses", Label: "Addresses"}}, Rows: rows},
	}}, nil
}
