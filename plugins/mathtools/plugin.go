package mathtools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunMathEvaluator() {
	cmdry.Run(plugin("math-evaluator", "Math Evaluator", "Evaluate bounded arithmetic expressions locally without arbitrary code execution.", []string{"math", "calculator", "expression", "arithmetic", "sqrt", "formula"}, evaluatorForm, evaluatorRun))
}

func RunETACalculator() {
	cmdry.Run(plugin("eta-calculator", "ETA Calculator", "Estimate completion time from elapsed work and progress values locally.", []string{"eta", "estimate", "progress", "completion", "time", "rate"}, etaForm, etaRun))
}

func RunPercentageCalculator() {
	cmdry.Run(plugin("percentage-calculator", "Percentage Calculator", "Calculate percentages, changes, and percentage differences locally.", []string{"percentage", "percent", "change", "difference", "calculator"}, percentageForm, percentageRun))
}

func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "math", Icon: "calculator", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": action}}
}

func evaluatorForm(cmdry.Request) (cmdry.View, error) {
	return form("Math Evaluator", "Supports +, -, *, /, %, ^, parentheses, pi, e, sqrt, abs, floor, ceil, and round. It cannot execute code.", "Evaluate", []cmdry.Field{{Name: "expression", Label: "Expression", Type: "text", Value: "sqrt(81) + 2^3", Required: true}}), nil
}

func evaluatorRun(request cmdry.Request) (cmdry.View, error) {
	expression := fmt.Sprint(request.Params["expression"])
	value, err := Evaluate(expression)
	if err != nil {
		return cmdry.View{}, err
	}
	return result("Math result", "Expression", expression, "Result", format(value), "Evaluate another expression"), nil
}

func etaForm(cmdry.Request) (cmdry.View, error) {
	return form("ETA Calculator", "Estimate remaining time from completed work, total work, and elapsed minutes.", "Estimate ETA", []cmdry.Field{{Name: "completed", Label: "Completed work", Type: "number", Value: "25", Min: "0", Required: true}, {Name: "total", Label: "Total work", Type: "number", Value: "100", Min: "0.000001", Required: true}, {Name: "elapsed_minutes", Label: "Elapsed minutes", Type: "number", Value: "30", Min: "0.000001", Required: true}}), nil
}

func etaRun(request cmdry.Request) (cmdry.View, error) {
	completed, err := number(request, "completed")
	if err != nil {
		return cmdry.View{}, err
	}
	total, err := number(request, "total")
	if err != nil {
		return cmdry.View{}, err
	}
	minutes, err := number(request, "elapsed_minutes")
	if err != nil {
		return cmdry.View{}, err
	}
	eta, err := EstimateETA(completed, total, time.Duration(minutes*float64(time.Minute)))
	if err != nil {
		return cmdry.View{}, err
	}
	finish := time.Now().Add(eta.Remaining).Format("Mon, Jan 2 3:04 PM")
	rows := []map[string]any{{"field": "Progress", "value": format(eta.Progress) + "%"}, {"field": "Rate", "value": format(eta.RatePerHour) + " units/hour"}, {"field": "Remaining", "value": formatDuration(eta.Remaining)}, {"field": "Estimated completion", "value": finish}}
	return tableResult("Estimated completion", "This estimate assumes the current rate remains constant.", "eta", rows, "Estimate another ETA"), nil
}

func percentageForm(cmdry.Request) (cmdry.View, error) {
	return form("Percentage Calculator", "Choose a common percentage calculation. All calculations run locally.", "Calculate percentage", []cmdry.Field{{Name: "operation", Label: "Calculation", Type: "select", Value: "of", Options: []cmdry.Option{{Value: "of", Label: "What is A% of B?"}, {Value: "what", Label: "A is what percent of B?"}, {Value: "change", Label: "Percent change from A to B"}, {Value: "difference", Label: "Percentage difference between A and B"}}}, {Name: "a", Label: "A", Type: "number", Value: "15", Required: true}, {Name: "b", Label: "B", Type: "number", Value: "80", Required: true}}), nil
}

func percentageRun(request cmdry.Request) (cmdry.View, error) {
	a, err := number(request, "a")
	if err != nil {
		return cmdry.View{}, err
	}
	b, err := number(request, "b")
	if err != nil {
		return cmdry.View{}, err
	}
	operation := fmt.Sprint(request.Params["operation"])
	var value float64
	var label, explanation string
	switch operation {
	case "of":
		value, err = PercentageOf(a, b)
		label, explanation = "A% of B", fmt.Sprintf("%s%% of %s", format(a), format(b))
	case "what":
		value, err = WhatPercent(a, b)
		label, explanation = "A as a percentage of B", fmt.Sprintf("%s is what percent of %s", format(a), format(b))
	case "change":
		value, err = PercentChange(a, b)
		label, explanation = "Percent change", fmt.Sprintf("from %s to %s", format(a), format(b))
	case "difference":
		value, err = PercentageDifference(a, b)
		label, explanation = "Percentage difference", fmt.Sprintf("between %s and %s", format(a), format(b))
	default:
		return cmdry.View{}, fmt.Errorf("choose a supported percentage calculation")
	}
	if err != nil {
		return cmdry.View{}, err
	}
	return result("Percentage result", "Calculation", explanation, label, format(value)+"%", "Calculate another percentage"), nil
}

func form(title, description, submit string, fields []cmdry.Field) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Description: description + " Your data never leaves this device.", Action: "run", Submit: submit, Fields: fields}}}
}
func result(title, leftLabel, leftValue, rightLabel, rightValue, again string) cmdry.View {
	return tableResult(title, "Calculation complete.", "math-result", []map[string]any{{"field": leftLabel, "value": leftValue}, {"field": rightLabel, "value": rightValue}}, again)
}
func tableResult(title, message, id string, rows []map[string]any, again string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: message}, {Type: "table", ID: id, Columns: []cmdry.Column{{Key: "field", Label: "Field"}, {Key: "value", Label: "Value"}}, Rows: rows}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: again}}}}}
}

func number(request cmdry.Request, key string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(request.Params[key])), 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("%s must be a finite number", strings.ReplaceAll(key, "_", " "))
	}
	return value, nil
}
func format(value float64) string { return strconv.FormatFloat(value, 'f', -1, 64) }
func formatDuration(duration time.Duration) string {
	duration = duration.Round(time.Second)
	hours := duration / time.Hour
	minutes := duration % time.Hour / time.Minute
	seconds := duration % time.Minute / time.Second
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
