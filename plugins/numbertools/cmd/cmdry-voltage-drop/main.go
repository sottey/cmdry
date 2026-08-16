package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/numbertools"
	"strconv"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.voltage-drop", Name: "Round-Trip Voltage Drop", Version: "0.1.0", Description: "Estimate copper-cable voltage drop and power loss locally.", Category: "math", Icon: "calculator", SearchTerms: []string{"voltage", "drop", "cable", "power", "electrical"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New calculation", Method: "read"}, {ID: "run", Name: "Calculate", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Round-Trip Voltage Drop", Components: []cmdry.Component{{Type: "form", Title: "Estimate copper-cable loss", Action: "run", Submit: "Calculate", Description: "Uses a 0.0175 Ω·mm²/m copper resistivity approximation and round-trip cable length.", Fields: []cmdry.Field{{Name: "length", Label: "One-way length (m)", Type: "number", Value: "10", Min: "0", Required: true}, {Name: "area", Label: "Conductor area (mm²)", Type: "number", Value: "2.5", Min: "0", Required: true}, {Name: "current", Label: "Current (A)", Type: "number", Value: "10", Min: "0", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	length, err := value(r, "length")
	if err != nil {
		return cmdry.View{}, err
	}
	area, err := value(r, "area")
	if err != nil {
		return cmdry.View{}, err
	}
	current, err := value(r, "current")
	if err != nil {
		return cmdry.View{}, err
	}
	drop, power, err := numbertools.VoltageDrop(length, area, current)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Voltage drop estimate", Components: []cmdry.Component{{Type: "metric", Label: "Voltage drop", Value: fmt.Sprintf("%.3f V", drop)}, {Type: "metric", Label: "Power loss", Value: fmt.Sprintf("%.3f W", power)}, {Type: "alert", Level: "info", Title: "Estimate only", Message: "Actual loss varies with conductor material, temperature, connections, AC effects, and installation conditions."}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Calculate another cable"}}}}}, nil
}
func value(r cmdry.Request, key string) (float64, error) {
	value, err := strconv.ParseFloat(fmt.Sprint(r.Params[key]), 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", key)
	}
	return value, nil
}
