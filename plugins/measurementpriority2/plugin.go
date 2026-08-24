package measurementpriority2

import (
	"encoding/base64"
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"strconv"
	"strings"
)

func RunChronometer() {
	cmdry.Run(plugin("chronometer", "Chronometer", "Run a browser-local elapsed-time stopwatch.", []string{"stopwatch", "timer", "lap", "chronometer"}, chronoView, chronoView))
}
func RunTemperatureConverter() {
	cmdry.Run(plugin("temperature-converter", "Temperature Converter", "Convert temperatures between Celsius, Fahrenheit, and Kelvin locally.", []string{"temperature", "celsius", "fahrenheit", "kelvin", "convert"}, tempForm, tempRun))
}
func RunBenchmarkBuilder() {
	cmdry.Run(plugin("benchmark-builder", "Benchmark Builder", "Create a bounded, repeatable local benchmark plan without executing commands.", []string{"benchmark", "performance", "test", "plan", "protocol"}, benchmarkForm, benchmarkRun))
}
func RunKitchenMeasurementConverter() {
	cmdry.Run(plugin("kitchen-measurement-converter", "Kitchen Measurement Converter", "Convert common mass or US-volume kitchen units locally.", []string{"kitchen", "convert", "grams", "ounces", "cups", "liters"}, kitchenForm, kitchenRun))
}
func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "measurement", Icon: "ruler", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": action}}
}
func chronoView(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Chronometer", Components: []cmdry.Component{{Type: "chronometer", ID: "local-stopwatch"}}}, nil
}
func tempForm(cmdry.Request) (cmdry.View, error) {
	return form("Temperature Converter", "Conversions run locally and reject values below absolute zero.", "Convert", []cmdry.Field{{Name: "value", Label: "Temperature", Type: "number", Value: "72", Required: true}, {Name: "from", Label: "From", Type: "select", Value: "f", Options: units()}, {Name: "to", Label: "To", Type: "select", Value: "c", Options: units()}})
}
func units() []cmdry.Option {
	return []cmdry.Option{{Value: "c", Label: "Celsius (°C)"}, {Value: "f", Label: "Fahrenheit (°F)"}, {Value: "k", Label: "Kelvin (K)"}}
}
func tempRun(r cmdry.Request) (cmdry.View, error) {
	value, err := ParseTemperature(fmt.Sprint(r.Params["value"]))
	if err != nil {
		return cmdry.View{}, err
	}
	from, to := fmt.Sprint(r.Params["from"]), fmt.Sprint(r.Params["to"])
	output, err := ConvertTemperature(value, from, to)
	if err != nil {
		return cmdry.View{}, err
	}
	return table("Temperature result", "The conversion was calculated locally.", "temperature", []map[string]any{{"field": "Input", "value": fmt.Sprintf("%g %s", value, from)}, {"field": "Result", "value": fmt.Sprintf("%g %s", output, to)}}, "Convert another temperature")
}
func benchmarkForm(cmdry.Request) (cmdry.View, error) {
	return form("Benchmark Builder", "Creates a repeatable benchmark plan only. Cmdry will not execute a command or collect system data.", "Create plan", []cmdry.Field{{Name: "name", Label: "Benchmark name", Type: "text", Value: "JSON parse comparison", Required: true}, {Name: "runs", Label: "Measured runs", Type: "number", Value: "10", Min: "1", Max: "100", Required: true}, {Name: "warmups", Label: "Warm-up runs", Type: "number", Value: "2", Min: "0", Max: "20", Required: true}, {Name: "duration", Label: "Max duration per run (seconds)", Type: "number", Value: "60", Min: "1", Max: "3600", Required: true}, {Name: "notes", Label: "Notes / workload", Type: "textarea", Value: "Use the same 10 MiB JSON fixture for each run."}})
}
func benchmarkRun(r cmdry.Request) (cmdry.View, error) {
	parse := func(key string) (int, error) { return strconv.Atoi(strings.TrimSpace(fmt.Sprint(r.Params[key]))) }
	runs, err := parse("runs")
	if err != nil {
		return cmdry.View{}, fmt.Errorf("runs must be an integer")
	}
	warmups, err := parse("warmups")
	if err != nil {
		return cmdry.View{}, fmt.Errorf("warm-up runs must be an integer")
	}
	duration, err := parse("duration")
	if err != nil {
		return cmdry.View{}, fmt.Errorf("duration must be an integer")
	}
	plan, err := BuildBenchmarkPlan(fmt.Sprint(r.Params["name"]), runs, warmups, duration, fmt.Sprint(r.Params["notes"]))
	if err != nil {
		return cmdry.View{}, err
	}
	output := plan.Markdown()
	return cmdry.View{Title: "Benchmark plan", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: "A repeatable benchmark plan was created locally. Nothing has been run."}, {Type: "download", Filename: "benchmark-plan.md", MIMEType: "text/markdown;charset=utf-8", Content: base64.StdEncoding.EncodeToString([]byte(output))}, {Type: "code", Title: "Plan", Text: output}, restart("Create another plan")}}, nil
}
func kitchenForm(cmdry.Request) (cmdry.View, error) {
	return form("Kitchen Measurement Converter", "Converts within mass or volume only. tsp, tbsp, cup, pt, qt, and gal use US customary measures.", "Convert", []cmdry.Field{{Name: "value", Label: "Value", Type: "number", Value: "1", Required: true}, {Name: "from", Label: "From", Type: "select", Value: "cup", Options: kitchenOptions()}, {Name: "to", Label: "To", Type: "select", Value: "ml", Options: kitchenOptions()}})
}
func kitchenOptions() []cmdry.Option {
	return []cmdry.Option{{Value: "g", Label: "grams (g)"}, {Value: "kg", Label: "kilograms (kg)"}, {Value: "lb", Label: "pounds (lb)"}, {Value: "oz", Label: "ounces (oz)"}, {Value: "ml", Label: "milliliters (mL)"}, {Value: "cl", Label: "centiliters (cL)"}, {Value: "dl", Label: "deciliters (dL)"}, {Value: "l", Label: "liters (L)"}, {Value: "tsp", Label: "teaspoons (tsp)"}, {Value: "tbsp", Label: "tablespoons (tbsp)"}, {Value: "cup", Label: "US cups"}, {Value: "pt", Label: "US pints"}, {Value: "qt", Label: "US quarts"}, {Value: "gal", Label: "US gallons"}}
}
func kitchenRun(r cmdry.Request) (cmdry.View, error) {
	value, err := ParseTemperature(fmt.Sprint(r.Params["value"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("value must be a finite number")
	}
	from, to := fmt.Sprint(r.Params["from"]), fmt.Sprint(r.Params["to"])
	output, err := ConvertKitchen(value, from, to)
	if err != nil {
		return cmdry.View{}, err
	}
	return table("Kitchen measurement", "Calculated locally; mass-to-volume conversions require an ingredient density.", "kitchen-measurement", []map[string]any{{"field": "Input", "value": fmt.Sprintf("%g %s", value, from)}, {"field": "Result", "value": fmt.Sprintf("%g %s", output, to)}}, "Convert another measurement")
}
func form(title, description, submit string, fields []cmdry.Field) (cmdry.View, error) {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Description: description + " Your data never leaves this device.", Action: "run", Submit: submit, Fields: fields}}}, nil
}
func table(title, message, id string, rows []map[string]any, again string) (cmdry.View, error) {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: message}, {Type: "table", ID: id, Columns: []cmdry.Column{{Key: "field", Label: "Field"}, {Key: "value", Label: "Value"}}, Rows: rows}, restart(again)}}, nil
}
func restart(name string) cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: name}}}
}
