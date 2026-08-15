package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/byteconvert"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.byte-converter", Name: "Byte Converter", Version: "0.1.0",
		Description: "Convert byte values between binary and decimal storage units.", Category: "developer", Icon: "convert",
		SearchTerms: []string{"bytes", "storage", "kilobytes", "megabytes", "gigabytes", "kib", "mib", "gib"},
		Pages:       []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}},
		Permissions: []string{"data.transform"},
		Actions:     []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert bytes", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "convert": convert}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	options := make([]cmdry.Option, 0, len(byteconvert.Units))
	for _, unit := range byteconvert.Units {
		options = append(options, cmdry.Option{Value: unit.ID, Label: unit.Label})
	}
	return cmdry.View{Title: "Byte Converter", Components: []cmdry.Component{{
		Type: "form", Title: "Convert a byte value", Action: "convert", Submit: "Convert",
		Description: "Binary units use powers of 1,024 (KiB, MiB, GiB); decimal units use powers of 1,000 (KB, MB, GB).",
		Fields:      []cmdry.Field{{Name: "value", Label: "Value", Type: "number", Value: "1", Min: "0", Required: true}, {Name: "unit", Label: "Starting unit", Type: "select", Value: "mb", Options: options}},
	}}}, nil
}

func convert(request cmdry.Request) (cmdry.View, error) {
	rawValue := strings.TrimSpace(fmt.Sprint(request.Params["value"]))
	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return cmdry.View{}, fmt.Errorf("value must be a finite number greater than or equal to zero")
	}
	unit := fmt.Sprint(request.Params["unit"])
	values, err := byteconvert.Convert(value, unit)
	if err != nil {
		return cmdry.View{}, err
	}
	rows := make([]map[string]any, 0, len(byteconvert.Units))
	for index, target := range byteconvert.Units {
		rows = append(rows, map[string]any{"unit": target.Label, "standard": target.Standard, "value": format(values[index])})
	}
	return cmdry.View{Title: "Byte conversion", Components: []cmdry.Component{
		{Type: "alert", Level: "success", Title: "Conversion complete", Message: "Equivalent values are shown in both binary and decimal units."},
		{Type: "table", ID: "byte-conversion", Columns: []cmdry.Column{{Key: "unit", Label: "Unit"}, {Key: "standard", Label: "Standard"}, {Key: "value", Label: "Equivalent value"}}, Rows: rows},
		{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another value"}}},
	}}, nil
}

func format(value float64) string {
	return strconv.FormatFloat(value, 'g', 10, 64)
}
