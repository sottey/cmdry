// Package numbertools provides local, dependency-free numerical utilities.
package numbertools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Summary describes a non-empty set of finite numbers.
type Summary struct {
	Count   int
	Sum     float64
	Average float64
	Minimum float64
	Maximum float64
}

// ParseNumbers accepts values separated by commas, semicolons, or whitespace.
func ParseNumbers(input string) ([]float64, error) {
	tokens := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || r == '\t' || r == ' '
	})
	if len(tokens) == 0 {
		return nil, fmt.Errorf("enter at least one number")
	}
	numbers := make([]float64, 0, len(tokens))
	for index, token := range tokens {
		value, err := strconv.ParseFloat(token, 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
			return nil, fmt.Errorf("value %d (%q) is not a finite number", index+1, token)
		}
		numbers = append(numbers, value)
	}
	return numbers, nil
}

// Summarize returns common aggregate values for input.
func Summarize(input []float64) (Summary, error) {
	if len(input) == 0 {
		return Summary{}, fmt.Errorf("enter at least one number")
	}
	result := Summary{Count: len(input), Minimum: input[0], Maximum: input[0]}
	for _, value := range input {
		result.Sum += value
		if value < result.Minimum {
			result.Minimum = value
		}
		if value > result.Maximum {
			result.Maximum = value
		}
	}
	result.Average = result.Sum / float64(result.Count)
	return result, nil
}
