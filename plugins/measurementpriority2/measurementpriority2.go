package measurementpriority2

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type KitchenUnit struct {
	Name, Dimension string
	Factor          float64
}

var kitchenUnits = map[string]KitchenUnit{
	"g": {"grams", "mass", 1}, "kg": {"kilograms", "mass", 1000}, "lb": {"pounds", "mass", 453.59237}, "oz": {"ounces", "mass", 28.349523125},
	"ml": {"milliliters", "volume", 1}, "cl": {"centiliters", "volume", 10}, "dl": {"deciliters", "volume", 100}, "l": {"liters", "volume", 1000}, "tsp": {"teaspoons", "volume", 4.92892159375}, "tbsp": {"tablespoons", "volume", 14.78676478125}, "cup": {"US cups", "volume", 236.5882365}, "pt": {"US pints", "volume", 473.176473}, "qt": {"US quarts", "volume", 946.352946}, "gal": {"US gallons", "volume", 3785.411784},
}

func ConvertKitchen(value float64, from, to string) (float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("value must be finite")
	}
	fromUnit, fromOK := kitchenUnits[from]
	toUnit, toOK := kitchenUnits[to]
	if !fromOK || !toOK {
		return 0, fmt.Errorf("choose supported source and target units")
	}
	if fromUnit.Dimension != toUnit.Dimension {
		return 0, fmt.Errorf("mass and volume cannot be converted without an ingredient density")
	}
	return value * fromUnit.Factor / toUnit.Factor, nil
}

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
