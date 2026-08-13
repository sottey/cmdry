package processes

import "testing"

func TestParseAndSummary(t *testing.T) {
	items, err := Parse("  42 1 2.5 1.2 R+ /usr/bin/example\n  99 42 0.0 0.3 S /bin/helper\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].PID != 42 || items[0].Command != "/usr/bin/example" {
		t.Fatalf("unexpected processes: %#v", items)
	}
	running, cpu, memory := Summary(items)
	if running != 1 || cpu != 2.5 || memory != 1.5 {
		t.Fatalf("unexpected summary: running=%d cpu=%v memory=%v", running, cpu, memory)
	}
	row := items[0].Row()
	if row["cpu"] != "2.5%" || row["memory"] != "1.2%" {
		t.Fatalf("unexpected row: %#v", row)
	}
}

func TestParseRejectsInvalidValues(t *testing.T) {
	if _, err := Parse("pid 1 0.0 0.0 S command\n"); err == nil {
		t.Fatal("accepted invalid PID")
	}
}
