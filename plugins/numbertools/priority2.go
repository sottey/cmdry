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
