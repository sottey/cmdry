package numbertools

import (
	"fmt"
	"math"
)

// VoltageDrop calculates round-trip DC voltage drop and power loss for a
// copper conductor using its cross-sectional area in mm².
func VoltageDrop(lengthMetres, areaMM2, currentAmps float64) (float64, float64, error) {
	if !finitePositive(lengthMetres) || !finitePositive(areaMM2) || !finitePositive(currentAmps) {
		return 0, 0, fmt.Errorf("length, conductor area, and current must be positive finite numbers")
	}
	resistance := 0.0175 * (2 * lengthMetres) / areaMM2
	drop := currentAmps * resistance
	return drop, currentAmps * currentAmps * resistance, nil
}

// SlacklineTension calculates an approximate static line tension. It assumes a
// centered load, a level line, and a small sag angle; it is not safety advice.
func SlacklineTension(loadKG, sagMetres, spanMetres float64) (float64, error) {
	if !finitePositive(loadKG) || !finitePositive(sagMetres) || !finitePositive(spanMetres) {
		return 0, fmt.Errorf("load, sag, and span must be positive finite numbers")
	}
	if sagMetres >= spanMetres/2 {
		return 0, fmt.Errorf("sag must be less than half the span")
	}
	return (loadKG * 9.80665 * spanMetres) / (4 * sagMetres), nil
}

func finitePositive(value float64) bool {
	return value > 0 && !math.IsNaN(value) && !math.IsInf(value, 0)
}

// ArithmeticSequence returns count values beginning at start with step added
// for each successive value.
func ArithmeticSequence(start, step float64, count int) ([]float64, error) {
	if math.IsNaN(start) || math.IsInf(start, 0) || math.IsNaN(step) || math.IsInf(step, 0) {
		return nil, fmt.Errorf("start and step must be finite numbers")
	}
	if count < 1 || count > 10000 {
		return nil, fmt.Errorf("count must be between 1 and 10,000")
	}
	values := make([]float64, count)
	for index := range values {
		values[index] = start + float64(index)*step
		if math.IsInf(values[index], 0) {
			return nil, fmt.Errorf("sequence exceeds the supported number range")
		}
	}
	return values, nil
}
