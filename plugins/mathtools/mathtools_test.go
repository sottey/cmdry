package mathtools

import (
	"math"
	"testing"
	"time"
)

func TestEvaluate(t *testing.T) {
	for input, want := range map[string]float64{"2 + 3 * 4": 14, "2^3^2": 512, "sqrt(81) + abs(-2)": 11, "-2^2": -4} {
		got, err := Evaluate(input)
		if err != nil || math.Abs(got-want) > 1e-9 {
			t.Fatalf("%s = %v, %v", input, got, err)
		}
	}
	if _, err := Evaluate("1 / 0"); err == nil {
		t.Fatal("accepted division by zero")
	}
}

func TestEstimateETA(t *testing.T) {
	eta, err := EstimateETA(25, 100, 30*time.Minute)
	if err != nil || eta.Progress != 25 || eta.RatePerHour != 50 || eta.Remaining != 90*time.Minute {
		t.Fatalf("eta=%#v err=%v", eta, err)
	}
}

func TestPercentages(t *testing.T) {
	if value, _ := PercentageOf(15, 80); value != 12 {
		t.Fatalf("of=%v", value)
	}
	if value, _ := WhatPercent(15, 60); value != 25 {
		t.Fatalf("what=%v", value)
	}
	if value, _ := PercentChange(80, 100); value != 25 {
		t.Fatalf("change=%v", value)
	}
	if value, _ := PercentageDifference(80, 120); value != 40 {
		t.Fatalf("difference=%v", value)
	}
}
