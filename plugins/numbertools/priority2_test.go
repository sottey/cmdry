package numbertools

import (
	"math"
	"testing"
)

func TestPriority2Calculators(t *testing.T) {
	drop, power, err := VoltageDrop(10, 2.5, 10)
	if err != nil || math.Abs(drop-1.4) > 0.000001 || math.Abs(power-14) > 0.000001 {
		t.Fatalf("drop=%v power=%v err=%v", drop, power, err)
	}
	tension, err := SlacklineTension(75, 0.5, 10)
	if err != nil || math.Abs(tension-3677.49375) > 0.000001 {
		t.Fatalf("tension=%v err=%v", tension, err)
	}
	if _, _, err := VoltageDrop(0, 2.5, 10); err == nil {
		t.Fatal("expected voltage validation error")
	}
	if _, err := SlacklineTension(75, 5, 10); err == nil {
		t.Fatal("expected sag validation error")
	}
}

func TestArithmeticSequence(t *testing.T) {
	values, err := ArithmeticSequence(2, 0.5, 4)
	if err != nil || len(values) != 4 || values[0] != 2 || values[3] != 3.5 {
		t.Fatalf("values=%v err=%v", values, err)
	}
	if _, err := ArithmeticSequence(1, 1, 0); err == nil {
		t.Fatal("expected count validation error")
	}
}
