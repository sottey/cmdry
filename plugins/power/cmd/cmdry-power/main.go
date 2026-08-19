package main

import (
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/power"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1,
		ID:              "com.sottey.power",
		Name:            "Battery and Power Inspector",
		Version:         "0.1.0",
		Description:     "Inspect local battery charge state and power source.",
		Category:        "server",
		Icon:            "battery",
		Pages:           []cmdry.Page{{ID: "overview", Name: "Power", Default: true, Action: "list"}},
		Permissions:     []string{"system.read"},
		Actions:         []cmdry.Action{{ID: "list", Name: "Refresh", Method: "read"}},
	}, Actions: map[string]cmdry.Handler{"list": list}})
}

func list(_ cmdry.Request) (cmdry.View, error) {
	info, err := power.Collect()
	if err != nil {
		return cmdry.View{}, err
	}
	if !info.Present {
		return cmdry.View{Title: "Battery and Power", Components: []cmdry.Component{{
			Type: "alert", Level: "warning", Title: "No battery reported",
			Message: "This host does not report a battery through its native power interface.",
		}}}, nil
	}
	return cmdry.View{Title: "Battery and Power", Components: []cmdry.Component{
		{Type: "metric", Label: "Charge", Value: strconv.Itoa(info.Percent) + "%"},
		{Type: "metric", Label: "Power source", Value: value(info.Source)},
		{Type: "metric", Label: "State", Value: value(info.State)},
		{Type: "table", ID: "power", Columns: []cmdry.Column{{Key: "label", Label: "Property"}, {Key: "value", Label: "Value"}}, Rows: info.Rows()},
	}}, nil
}

func value(input string) string {
	if input == "" {
		return "Unavailable"
	}
	return input
}
