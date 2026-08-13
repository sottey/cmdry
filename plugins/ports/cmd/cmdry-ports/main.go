package main

import (
	"sort"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/ports"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.port.inspector", Name: "Port Inspector", Version: "0.2.0", Description: "Inspect listening network ports and their owning processes.", Category: "system", Icon: "network", Pages: []cmdry.Page{{ID: "overview", Name: "Ports", Default: true, Action: "list"}}, Permissions: []string{"network.read", "process.read"}, Actions: []cmdry.Action{{ID: "list", Name: "List ports", Method: "read"}}}, Actions: map[string]cmdry.Handler{"list": listPorts}})
}
func listPorts(_ cmdry.Request) (cmdry.View, error) {
	items, err := ports.CollectListeningPorts()
	if err != nil {
		return cmdry.View{}, err
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Port == items[j].Port {
			return items[i].Protocol < items[j].Protocol
		}
		return items[i].Port < items[j].Port
	})
	tcp, udp := ports.Summary(items)
	rows := make([]map[string]any, 0, len(items))
	for _, item := range items {
		rows = append(rows, item.Row())
	}
	return cmdry.View{Title: "Listening Ports", Components: []cmdry.Component{{Type: "metric", Label: "Listening ports", Value: strconv.Itoa(len(items))}, {Type: "metric", Label: "TCP", Value: strconv.Itoa(tcp)}, {Type: "metric", Label: "UDP", Value: strconv.Itoa(udp)}, {Type: "table", ID: "ports", Columns: []cmdry.Column{{Key: "port", Label: "Port"}, {Key: "protocol", Label: "Protocol"}, {Key: "address", Label: "Address"}, {Key: "process", Label: "Process"}, {Key: "pid", Label: "PID"}}, Rows: rows}}}, nil
}
