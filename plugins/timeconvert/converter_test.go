package timeconvert

import "testing"

func TestFromUnixSeconds(t *testing.T) {
	moment, err := FromUnixSeconds("0")
	if err != nil || moment.UTC().Format("2006-01-02T15:04:05Z") != "1970-01-01T00:00:00Z" {
		t.Fatalf("moment = %v, err = %v", moment, err)
	}
}

func TestToUnixSeconds(t *testing.T) {
	moment, err := ToUnixSeconds("2024-01-02T03:04:05Z", "local")
	if err != nil || moment.Unix() != 1704164645 {
		t.Fatalf("moment = %v, err = %v", moment, err)
	}
}

func TestToUnixSecondsRejectsInvalidDate(t *testing.T) {
	if _, err := ToUnixSeconds("not a date", "utc"); err == nil {
		t.Fatal("expected invalid date error")
	}
}
