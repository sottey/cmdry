package timeunits

import "testing"

func TestTruncateClockTime(t *testing.T) {
	result, err := TruncateClockTime(5025, "minute")
	if err != nil || result != 4980 {
		t.Fatalf("result=%d err=%v", result, err)
	}
	result, err = TruncateClockTime(-5025, "hour")
	if err != nil || result != -3600 {
		t.Fatalf("negative result=%d err=%v", result, err)
	}
	if _, err := TruncateClockTime(1, "day"); err == nil {
		t.Fatal("expected invalid precision")
	}
}
