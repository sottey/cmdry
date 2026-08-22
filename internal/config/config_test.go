package config

import "testing"

func TestBoolValue(t *testing.T) {
	t.Setenv("CMDRY_DEMO_MODE", "true")
	got, err := boolValue("CMDRY_DEMO_MODE", false)
	if err != nil || !got {
		t.Fatalf("boolValue = %t, %v; want true, nil", got, err)
	}
}

func TestBoolValueRejectsInvalidValue(t *testing.T) {
	t.Setenv("CMDRY_DEMO_MODE", "sometimes")
	if _, err := boolValue("CMDRY_DEMO_MODE", false); err == nil {
		t.Fatal("boolValue accepted invalid value")
	}
}
