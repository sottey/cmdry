// Package numbertools provides local, dependency-free numerical utilities.
package numbertools

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
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

// GenerateRandomIntegers returns count cryptographically secure random integers
// within the inclusive min and max bounds. Bounds may exceed int64.
func GenerateRandomIntegers(minInput, maxInput string, count int) ([]string, error) {
	if count < 1 || count > 10000 {
		return nil, fmt.Errorf("count must be between 1 and 10,000")
	}
	min, ok := new(big.Int).SetString(strings.TrimSpace(minInput), 10)
	if !ok {
		return nil, fmt.Errorf("minimum must be a whole number")
	}
	max, ok := new(big.Int).SetString(strings.TrimSpace(maxInput), 10)
	if !ok {
		return nil, fmt.Errorf("maximum must be a whole number")
	}
	if min.Cmp(max) > 0 {
		return nil, fmt.Errorf("minimum must not exceed maximum")
	}
	width := new(big.Int).Sub(max, min)
	width.Add(width, big.NewInt(1))
	values := make([]string, 0, count)
	for range count {
		offset, err := rand.Int(rand.Reader, width)
		if err != nil {
			return nil, fmt.Errorf("generate random integer: %w", err)
		}
		values = append(values, new(big.Int).Add(min, offset).String())
	}
	return values, nil
}

// SphereArea calculates the surface area 4πr² for a non-negative radius.
func SphereArea(radius float64) (float64, error) {
	if math.IsNaN(radius) || math.IsInf(radius, 0) || radius < 0 {
		return 0, fmt.Errorf("radius must be a non-negative finite number")
	}
	return 4 * math.Pi * radius * radius, nil
}

func OhmsLaw(voltage, current, resistance float64) (float64, error) {
	count := 0
	if voltage > 0 {
		count++
	}
	if current > 0 {
		count++
	}
	if resistance > 0 {
		count++
	}
	if count != 2 {
		return 0, fmt.Errorf("enter exactly two positive values")
	}
	if voltage == 0 {
		return current * resistance, nil
	}
	if current == 0 {
		return voltage / resistance, nil
	}
	return voltage / current, nil
}
func NumberToWords(value int64) string {
	if value == 0 {
		return "zero"
	}
	if value < 0 {
		return "minus " + NumberToWords(-value)
	}
	ones := []string{"", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens := []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
	var under func(int64) string
	under = func(n int64) string {
		if n < 20 {
			return ones[n]
		}
		if n < 100 {
			if n%10 == 0 {
				return tens[n/10]
			}
			return tens[n/10] + "-" + ones[n%10]
		}
		if n < 1000 {
			if n%100 == 0 {
				return ones[n/100] + " hundred"
			}
			return ones[n/100] + " hundred " + under(n%100)
		}
		return ""
	}
	scales := []struct {
		n int64
		s string
	}{{1000000000000, "trillion"}, {1000000000, "billion"}, {1000000, "million"}, {1000, "thousand"}}
	parts := []string{}
	for _, scale := range scales {
		if value >= scale.n {
			parts = append(parts, under(value/scale.n)+" "+scale.s)
			value %= scale.n
		}
	}
	if value > 0 {
		parts = append(parts, under(value))
	}
	return strings.Join(parts, " ")
}
