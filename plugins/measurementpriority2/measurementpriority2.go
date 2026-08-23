package measurementpriority2

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ConvertTemperature(value float64, from, to string) (float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("temperature must be finite")
	}
	toCelsius := map[string]func(float64) float64{"c": func(v float64) float64 { return v }, "f": func(v float64) float64 { return (v - 32) * 5 / 9 }, "k": func(v float64) float64 { return v - 273.15 }}
	fromCelsius := map[string]func(float64) float64{"c": func(v float64) float64 { return v }, "f": func(v float64) float64 { return v*9/5 + 32 }, "k": func(v float64) float64 { return v + 273.15 }}
	convertIn, ok := toCelsius[from]
	if !ok {
		return 0, fmt.Errorf("choose a source unit")
	}
	convertOut, ok := fromCelsius[to]
	if !ok {
		return 0, fmt.Errorf("choose a target unit")
	}
	celsius := convertIn(value)
	if celsius < -273.15-1e-9 {
		return 0, fmt.Errorf("temperature is below absolute zero")
	}
	return convertOut(celsius), nil
}
func ParseTemperature(input string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("temperature must be a finite number")
	}
	return value, nil
}

type BenchmarkPlan struct {
	Name            string
	Runs            int
	Warmups         int
	DurationSeconds int
	Notes           string
}

func BuildBenchmarkPlan(name string, runs, warmups, duration int, notes string) (BenchmarkPlan, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 120 {
		return BenchmarkPlan{}, fmt.Errorf("benchmark name must contain 1 to 120 characters")
	}
	if runs < 1 || runs > 100 {
		return BenchmarkPlan{}, fmt.Errorf("runs must be between 1 and 100")
	}
	if warmups < 0 || warmups > 20 {
		return BenchmarkPlan{}, fmt.Errorf("warm-up runs must be between 0 and 20")
	}
	if duration < 1 || duration > 3600 {
		return BenchmarkPlan{}, fmt.Errorf("duration must be between 1 and 3,600 seconds")
	}
	if len(notes) > 4000 {
		return BenchmarkPlan{}, fmt.Errorf("notes must not exceed 4,000 characters")
	}
	return BenchmarkPlan{Name: name, Runs: runs, Warmups: warmups, DurationSeconds: duration, Notes: strings.TrimSpace(notes)}, nil
}
func (p BenchmarkPlan) Markdown() string {
	return fmt.Sprintf("# %s\n\n## Protocol\n\n1. Record machine, OS, application version, power mode, and test input.\n2. Close unrelated workloads and use the same input for every run.\n3. Perform %d warm-up run(s); do not include them in results.\n4. Perform %d measured run(s), each capped at %d second(s).\n5. Record elapsed time, output correctness, and relevant resource observations for each run.\n6. Report median and range; do not compare results from materially different environments.\n\n## Notes\n\n%s\n", p.Name, p.Warmups, p.Runs, p.DurationSeconds, defaultNotes(p.Notes))
}
func defaultNotes(notes string) string {
	if notes == "" {
		return "Add the exact command or workload, input size, and environment details before running."
	}
	return notes
}
