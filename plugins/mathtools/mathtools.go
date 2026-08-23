// Package mathtools implements bounded, local-only mathematical utilities.
package mathtools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Evaluate parses a bounded arithmetic expression without evaluating code.
// It supports +, -, *, /, %, ^, parentheses, pi, e, and common unary functions.
func Evaluate(input string) (float64, error) {
	if len([]rune(input)) == 0 || len([]rune(input)) > 512 {
		return 0, fmt.Errorf("expression must contain between 1 and 512 characters")
	}
	parser := expressionParser{input: []rune(input)}
	value, err := parser.expression()
	if err != nil {
		return 0, err
	}
	parser.space()
	if parser.position != len(parser.input) {
		return 0, fmt.Errorf("unexpected character %q", parser.input[parser.position])
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("expression result must be finite")
	}
	return value, nil
}

type expressionParser struct {
	input    []rune
	position int
	depth    int
}

func (p *expressionParser) expression() (float64, error) {
	value, err := p.term()
	if err != nil {
		return 0, err
	}
	for {
		p.space()
		if !p.take('+') && !p.take('-') {
			return value, nil
		}
		operator := p.input[p.position-1]
		right, err := p.term()
		if err != nil {
			return 0, err
		}
		if operator == '+' {
			value += right
		} else {
			value -= right
		}
	}
}

func (p *expressionParser) term() (float64, error) {
	value, err := p.unary()
	if err != nil {
		return 0, err
	}
	for {
		p.space()
		if !p.take('*') && !p.take('/') && !p.take('%') {
			return value, nil
		}
		operator := p.input[p.position-1]
		right, err := p.unary()
		if err != nil {
			return 0, err
		}
		if (operator == '/' || operator == '%') && right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		switch operator {
		case '*':
			value *= right
		case '/':
			value /= right
		case '%':
			value = math.Mod(value, right)
		}
	}
}

func (p *expressionParser) unary() (float64, error) {
	p.space()
	if p.take('+') {
		return p.unary()
	}
	if p.take('-') {
		value, err := p.unary()
		return -value, err
	}
	return p.power()
}

func (p *expressionParser) power() (float64, error) {
	value, err := p.primary()
	if err != nil {
		return 0, err
	}
	p.space()
	if p.take('^') {
		right, err := p.unary()
		if err != nil {
			return 0, err
		}
		return math.Pow(value, right), nil
	}
	return value, nil
}

func (p *expressionParser) primary() (float64, error) {
	p.space()
	if p.take('(') {
		p.depth++
		if p.depth > 32 {
			return 0, fmt.Errorf("expression nesting is limited to 32 levels")
		}
		value, err := p.expression()
		p.depth--
		if err != nil {
			return 0, err
		}
		p.space()
		if !p.take(')') {
			return 0, fmt.Errorf("missing closing parenthesis")
		}
		return value, nil
	}
	if p.position < len(p.input) && (unicode.IsLetter(p.input[p.position]) || p.input[p.position] == 'π') {
		name := p.identifier()
		switch strings.ToLower(name) {
		case "pi", "π":
			return math.Pi, nil
		case "e":
			return math.E, nil
		}
		p.space()
		if !p.take('(') {
			return 0, fmt.Errorf("unknown constant or function %q", name)
		}
		argument, err := p.expression()
		if err != nil {
			return 0, err
		}
		p.space()
		if !p.take(')') {
			return 0, fmt.Errorf("missing closing parenthesis after %s", name)
		}
		switch strings.ToLower(name) {
		case "sqrt":
			if argument < 0 {
				return 0, fmt.Errorf("sqrt requires a non-negative value")
			}
			return math.Sqrt(argument), nil
		case "abs":
			return math.Abs(argument), nil
		case "floor":
			return math.Floor(argument), nil
		case "ceil":
			return math.Ceil(argument), nil
		case "round":
			return math.Round(argument), nil
		default:
			return 0, fmt.Errorf("unsupported function %q", name)
		}
	}
	start := p.position
	for p.position < len(p.input) && (unicode.IsDigit(p.input[p.position]) || p.input[p.position] == '.') {
		p.position++
	}
	if start == p.position {
		return 0, fmt.Errorf("expected a number or parenthesized expression")
	}
	value, err := strconv.ParseFloat(string(p.input[start:p.position]), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", string(p.input[start:p.position]))
	}
	return value, nil
}

func (p *expressionParser) identifier() string {
	start := p.position
	for p.position < len(p.input) && (unicode.IsLetter(p.input[p.position]) || p.input[p.position] == 'π') {
		p.position++
	}
	return string(p.input[start:p.position])
}
func (p *expressionParser) space() {
	for p.position < len(p.input) && unicode.IsSpace(p.input[p.position]) {
		p.position++
	}
}
func (p *expressionParser) take(character rune) bool {
	if p.position < len(p.input) && p.input[p.position] == character {
		p.position++
		return true
	}
	return false
}

// ETA describes an estimate derived from elapsed work and completed progress.
type ETA struct {
	Progress, RatePerHour float64
	Remaining             time.Duration
}

func EstimateETA(completed, total float64, elapsed time.Duration) (ETA, error) {
	if math.IsNaN(completed) || math.IsNaN(total) || math.IsInf(completed, 0) || math.IsInf(total, 0) || completed < 0 || total <= 0 || completed > total {
		return ETA{}, fmt.Errorf("progress must be finite with 0 <= completed <= total and total > 0")
	}
	if elapsed <= 0 {
		return ETA{}, fmt.Errorf("elapsed time must be greater than zero")
	}
	if completed == 0 {
		return ETA{}, fmt.Errorf("at least some work must be completed to estimate an ETA")
	}
	rate := completed / elapsed.Hours()
	remaining := time.Duration(float64(time.Hour) * (total - completed) / rate)
	return ETA{Progress: completed / total * 100, RatePerHour: rate, Remaining: remaining}, nil
}

func PercentageOf(percent, value float64) (float64, error) {
	return finiteResult(percent / 100 * value)
}
func WhatPercent(part, whole float64) (float64, error) {
	if whole == 0 {
		return 0, fmt.Errorf("whole value must not be zero")
	}
	return finiteResult(part / whole * 100)
}
func PercentChange(original, updated float64) (float64, error) {
	if original == 0 {
		return 0, fmt.Errorf("original value must not be zero")
	}
	return finiteResult((updated - original) / math.Abs(original) * 100)
}
func PercentageDifference(a, b float64) (float64, error) {
	denominator := (math.Abs(a) + math.Abs(b)) / 2
	if denominator == 0 {
		return 0, fmt.Errorf("values must not both be zero")
	}
	return finiteResult(math.Abs(a-b) / denominator * 100)
}
func finiteResult(value float64) (float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("result must be finite")
	}
	return value, nil
}
