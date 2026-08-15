package timeunits

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

type Converter func(float64) (float64, error)

// Run serves one independently compiled time-unit converter plugin.
func Run(id, name, sourceUnit, targetUnit string, convert Converter) {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: "Convert pasted " + sourceUnit + " to " + targetUnit + " locally.", Category: "time", Icon: "clock",
		SearchTerms: []string{"time", sourceUnit, targetUnit, "convert", "duration"}, Pages: []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": form(name, sourceUnit, targetUnit), "convert": conversion(name, sourceUnit, targetUnit, convert)}})
}

func form(name, sourceUnit, targetUnit string) cmdry.Handler {
	return func(cmdry.Request) (cmdry.View, error) {
		return cmdry.View{Title: name, Components: []cmdry.Component{{Type: "form", Title: "Convert " + sourceUnit + " to " + targetUnit, Action: "convert", Submit: "Convert", Description: "Accepts finite decimal values locally. Negative values are supported for offsets.", Fields: []cmdry.Field{{Name: "value", Label: strings.Title(sourceUnit), Type: "number", Value: "1", Required: true}}}}}, nil
	}
}

func conversion(name, sourceUnit, targetUnit string, convert Converter) cmdry.Handler {
	return func(request cmdry.Request) (cmdry.View, error) {
		input, err := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(request.Params["value"])), 64)
		if err != nil || math.IsNaN(input) || math.IsInf(input, 0) {
			return cmdry.View{}, fmt.Errorf("%s must be a finite number", sourceUnit)
		}
		output, err := convert(input)
		if err != nil {
			return cmdry.View{}, err
		}
		return cmdry.View{Title: name, Components: []cmdry.Component{{Type: "metric", Label: strings.Title(sourceUnit), Value: format(input)}, {Type: "metric", Label: strings.Title(targetUnit), Value: format(output)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another value"}}}}}, nil
	}
}

func format(value float64) string { return strconv.FormatFloat(value, 'g', 12, 64) }
