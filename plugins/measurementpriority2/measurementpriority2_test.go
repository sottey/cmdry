package measurementpriority2

import (
	"strings"
	"testing"
)

func TestTemperature(t *testing.T) {
	value, err := ConvertTemperature(212, "f", "c")
	if err != nil || value != 100 {
		t.Fatalf("value=%v err=%v", value, err)
	}
	if _, err := ConvertTemperature(-1, "k", "c"); err == nil {
		t.Fatal("accepted below absolute zero")
	}
}
func TestBenchmarkPlan(t *testing.T) {
	plan, err := BuildBenchmarkPlan("Parser", 10, 2, 60, "fixture")
	if err != nil || !strings.Contains(plan.Markdown(), "10 measured run") {
		t.Fatalf("plan=%#v err=%v", plan, err)
	}
}
