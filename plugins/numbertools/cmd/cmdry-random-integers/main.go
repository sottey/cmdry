package main

import (
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/numbertools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.random-integers", Name: "Generate Random Numbers", Version: "0.1.0", Description: "Generate cryptographically secure random integers in a local range.", Category: "number", Icon: "random",
		SearchTerms: []string{"random", "numbers", "generate", "integer", "range", "secure"}, Pages: []cmdry.Page{{ID: "overview", Name: "Generate", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New generation", Method: "read"}, {ID: "generate", Name: "Generate numbers", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "generate": generate}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Generate Random Numbers", Components: []cmdry.Component{{Type: "form", Title: "Generate random integers", Action: "generate", Submit: "Generate numbers", Description: "Uses the operating system's cryptographically secure random source. Bounds are inclusive and may be larger than standard integer fields.", Fields: []cmdry.Field{{Name: "min", Label: "Minimum", Type: "text", Value: "1", Required: true}, {Name: "max", Label: "Maximum", Type: "text", Value: "100", Required: true}, {Name: "count", Label: "How many", Type: "number", Value: "10", Min: "1", Max: "10000", Required: true}}}}}, nil
}

func generate(request cmdry.Request) (cmdry.View, error) {
	count, err := strconv.Atoi(fmt.Sprint(request.Params["count"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("how many must be a whole number")
	}
	values, err := numbertools.GenerateRandomIntegers(fmt.Sprint(request.Params["min"]), fmt.Sprint(request.Params["max"]), count)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Random numbers", Components: []cmdry.Component{{Type: "metric", Label: "Numbers generated", Value: fmt.Sprint(len(values))}, {Type: "code", Title: "Random integers", Text: strings.Join(values, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Generate another set"}}}}}, nil
}
