package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/numbertools"
	"strconv"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.slackline-tension", Name: "Slackline Tension", Version: "0.1.0", Description: "Estimate static centered-load line tension locally; never use it for safety decisions.", Category: "math", Icon: "calculator", SearchTerms: []string{"slackline", "tension", "clothesline", "load", "sag"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New estimate", Method: "read"}, {ID: "run", Name: "Estimate tension", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Slackline Tension", Components: []cmdry.Component{{Type: "alert", Level: "warning", Title: "Not for safety decisions", Message: "This simplified static estimate excludes dynamic loads, anchors, line rating, knots, and safety factors."}, {Type: "form", Title: "Estimate centered-load tension", Action: "run", Submit: "Estimate tension", Fields: []cmdry.Field{{Name: "load", Label: "Centered load (kg)", Type: "number", Value: "75", Min: "0", Required: true}, {Name: "sag", Label: "Sag (m)", Type: "number", Value: "0.5", Min: "0", Required: true}, {Name: "span", Label: "Span (m)", Type: "number", Value: "10", Min: "0", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	load, err := decimal(r, "load")
	if err != nil {
		return cmdry.View{}, err
	}
	sag, err := decimal(r, "sag")
	if err != nil {
		return cmdry.View{}, err
	}
	span, err := decimal(r, "span")
	if err != nil {
		return cmdry.View{}, err
	}
	tension, err := numbertools.SlacklineTension(load, sag, span)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Tension estimate", Components: []cmdry.Component{{Type: "metric", Label: "Approximate static tension", Value: fmt.Sprintf("%.0f N (%.1f kN)", tension, tension/1000)}, {Type: "alert", Level: "warning", Title: "Do not use for rigging", Message: "Consult qualified guidance and rated equipment; real-world tension can be substantially higher."}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Estimate another line"}}}}}, nil
}
func decimal(r cmdry.Request, key string) (float64, error) {
	value, err := strconv.ParseFloat(fmt.Sprint(r.Params[key]), 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", key)
	}
	return value, nil
}
