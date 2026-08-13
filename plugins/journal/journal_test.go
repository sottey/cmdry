package journal

import "testing"

func TestParse(t *testing.T) {
	input := []byte(`{"__REALTIME_TIMESTAMP":"1723593845123456","PRIORITY":"3","MESSAGE":"disk failed\nreplace it","_SYSTEMD_UNIT":"smartd.service"}
{"__REALTIME_TIMESTAMP":"1723593846123456","PRIORITY":"6","MESSAGE":"service started"}
`)
	entries, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("got %d entries", len(entries))
	}
	if entries[0].Priority != 3 || entries[0].Unit != "smartd.service" || entries[0].Message != "disk failed\nreplace it" || entries[0].Timestamp.IsZero() {
		t.Fatalf("unexpected first entry: %#v", entries[0])
	}
	row := entries[0].Row()
	if row["priority"] != "error" || row["unit"] != "smartd.service" || row["time"] == "" {
		t.Fatalf("unexpected row: %#v", row)
	}
	if got := ShortMessage(entries[0].Message); got != "disk failed replace it" {
		t.Fatalf("got %q", got)
	}
	errors, warnings := Summary(entries)
	if errors != 1 || warnings != 0 {
		t.Fatalf("got errors=%d warnings=%d", errors, warnings)
	}
}

func TestParseRejectsMalformedRecord(t *testing.T) {
	if _, err := Parse([]byte("not json\n")); err == nil {
		t.Fatal("accepted malformed JSON")
	}
	if _, err := Parse([]byte(`{"PRIORITY":"loud"}` + "\n")); err == nil {
		t.Fatal("accepted an invalid priority")
	}
}

func TestPriorityName(t *testing.T) {
	if got := PriorityName(4); got != "warning" {
		t.Fatalf("got %q", got)
	}
	if got := PriorityName(99); got != "unknown" {
		t.Fatalf("got %q", got)
	}
}
