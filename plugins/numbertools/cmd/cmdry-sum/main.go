package main

import (
	"fmt"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/numbertools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.sum", Name: "Sum", Version: "0.1.0",
		Description: "Calculate a local sum and summary statistics from pasted numbers.", Category: "number", Icon: "sum",
		SearchTerms: []string{"sum", "add", "total", "average", "numbers", "statistics"},
		Pages:       []cmdry.Page{{ID: "overview", Name: "Calculate", Default: true, Action: "overview"}},
		Permissions: []string{"data.transform"},
		Actions:     []cmdry.Action{{ID: "overview", Name: "New calculation", Method: "read"}, {ID: "calculate", Name: "Calculate", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "calculate": calculate}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Sum", Components: []cmdry.Component{{
		Type: "form", Title: "Add numbers", Action: "calculate", Submit: "Calculate",
		Description: "Separate finite numbers with commas, spaces, semicolons, or line breaks. Decimal points and scientific notation are supported.",
		Fields:      []cmdry.Field{{Name: "input", Label: "Numbers", Type: "textarea", Required: true, Placeholder: "10\n25.5\n-3"}},
	}}}, nil
}

func calculate(request cmdry.Request) (cmdry.View, error) {
	numbers, err := numbertools.ParseNumbers(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	summary, err := numbertools.Summarize(numbers)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Number summary", Components: []cmdry.Component{
		{Type: "metric", Label: "Sum", Value: format(summary.Sum)},
		{Type: "metric", Label: "Count", Value: strconv.Itoa(summary.Count)},
		{Type: "metric", Label: "Average", Value: format(summary.Average)},
		{Type: "metric", Label: "Minimum", Value: format(summary.Minimum)},
		{Type: "metric", Label: "Maximum", Value: format(summary.Maximum)},
		{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Calculate another sum"}}},
	}}, nil
}

func format(value float64) string { return strconv.FormatFloat(value, 'g', 12, 64) }
