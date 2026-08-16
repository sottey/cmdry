package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/numbertools"
	"strconv"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.ohms-law", Name: "Ohm's Law", Version: "0.1.0", Description: "Calculate voltage, current, or resistance locally.", Category: "number", Icon: "calculate", SearchTerms: []string{"ohm", "voltage", "current", "resistance", "electricity"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New calculation", Method: "read"}, {ID: "run", Name: "Calculate", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": form, "run": run}})
}
func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Ohm's Law", Components: []cmdry.Component{{Type: "form", Title: "Enter any two values", Action: "run", Submit: "Calculate", Fields: []cmdry.Field{{Name: "voltage", Label: "Voltage (V)", Type: "number"}, {Name: "current", Label: "Current (A)", Type: "number"}, {Name: "resistance", Label: "Resistance (Ω)", Type: "number"}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	parse := func(k string) float64 { v, _ := strconv.ParseFloat(fmt.Sprint(r.Params[k]), 64); return v }
	v, c, o := parse("voltage"), parse("current"), parse("resistance")
	value, e := numbertools.OhmsLaw(v, c, o)
	if e != nil {
		return cmdry.View{}, e
	}
	label := "Resistance (Ω)"
	if v == 0 {
		label = "Voltage (V)"
	} else if c == 0 {
		label = "Current (A)"
	}
	return cmdry.View{Title: "Ohm's Law result", Components: []cmdry.Component{{Type: "metric", Label: label, Value: fmt.Sprint(value)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Calculate again"}}}}}, nil
}
