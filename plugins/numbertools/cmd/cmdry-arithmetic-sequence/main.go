package main

import (
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/numbertools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.arithmetic-sequence", Name: "Arithmetic Sequence", Version: "0.1.0", Description: "Generate a finite arithmetic sequence locally.", Category: "number", Icon: "calculator", SearchTerms: []string{"arithmetic", "sequence", "numbers", "progression", "series"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New sequence", Method: "read"}, {ID: "run", Name: "Generate sequence", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Arithmetic Sequence", Components: []cmdry.Component{{Type: "form", Title: "Generate a sequence", Action: "run", Submit: "Generate sequence", Fields: []cmdry.Field{{Name: "start", Label: "Start", Type: "number", Value: "1", Required: true}, {Name: "step", Label: "Step", Type: "number", Value: "1", Required: true}, {Name: "count", Label: "Values", Type: "number", Value: "10", Min: "1", Max: "10000", Required: true}}}}}, nil
}
func run(request cmdry.Request) (cmdry.View, error) {
	start, err := strconv.ParseFloat(fmt.Sprint(request.Params["start"]), 64)
	if err != nil {
		return cmdry.View{}, fmt.Errorf("start must be a number")
	}
	step, err := strconv.ParseFloat(fmt.Sprint(request.Params["step"]), 64)
	if err != nil {
		return cmdry.View{}, fmt.Errorf("step must be a number")
	}
	count, err := strconv.Atoi(fmt.Sprint(request.Params["count"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("values must be a whole number")
	}
	values, err := numbertools.ArithmeticSequence(start, step, count)
	if err != nil {
		return cmdry.View{}, err
	}
	output := make([]string, len(values))
	for index, value := range values {
		output[index] = strconv.FormatFloat(value, 'g', 12, 64)
	}
	return cmdry.View{Title: "Arithmetic sequence", Components: []cmdry.Component{{Type: "metric", Label: "Values generated", Value: fmt.Sprint(len(values))}, {Type: "code", Title: "Sequence", Text: strings.Join(output, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Generate another sequence"}}}}}, nil
}
