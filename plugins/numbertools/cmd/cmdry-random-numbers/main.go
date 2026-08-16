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
		ProtocolVersion: 1, ID: "com.sottey.random-numbers", Name: "Random Number Generator", Version: "0.1.0", Description: "Generate secure fixed-precision decimal numbers in a local range.", Category: "number", Icon: "random",
		SearchTerms: []string{"random", "numbers", "decimal", "range", "generate", "secure"}, Pages: []cmdry.Page{{ID: "overview", Name: "Generate", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New generation", Method: "read"}, {ID: "generate", Name: "Generate numbers", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "generate": generate}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Random Number Generator", Components: []cmdry.Component{{Type: "form", Title: "Generate random decimal numbers", Action: "generate", Submit: "Generate numbers", Description: "Uses the operating system's cryptographically secure random source. Both bounds are inclusive and inputs may use up to the selected number of decimal places.", Fields: []cmdry.Field{{Name: "min", Label: "Minimum", Type: "text", Value: "0", Required: true}, {Name: "max", Label: "Maximum", Type: "text", Value: "1", Required: true}, {Name: "places", Label: "Decimal places", Type: "number", Value: "2", Min: "0", Max: "9", Required: true}, {Name: "count", Label: "How many", Type: "number", Value: "10", Min: "1", Max: "10000", Required: true}}}}}, nil
}

func generate(request cmdry.Request) (cmdry.View, error) {
	places, err := strconv.Atoi(fmt.Sprint(request.Params["places"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("decimal places must be a whole number")
	}
	count, err := strconv.Atoi(fmt.Sprint(request.Params["count"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("how many must be a whole number")
	}
	values, err := numbertools.GenerateRandomDecimals(fmt.Sprint(request.Params["min"]), fmt.Sprint(request.Params["max"]), places, count)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Random numbers", Components: []cmdry.Component{{Type: "metric", Label: "Numbers generated", Value: fmt.Sprint(len(values))}, {Type: "code", Title: "Random decimal numbers", Text: strings.Join(values, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Generate another set"}}}}}, nil
}
