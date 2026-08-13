package systeminfo

import "testing"

func TestParseOSReleaseAndMemInfo(t *testing.T) {
	if got := ParseOSRelease("NAME=Ubuntu\nPRETTY_NAME=\"Ubuntu 24.04 LTS\"\n"); got != "Ubuntu 24.04 LTS" {
		t.Fatalf("got %q", got)
	}
	total, available := ParseMemInfo("MemTotal: 16384000 kB\nMemAvailable: 8192000 kB\n")
	if total != 16777216000 || available != 8388608000 {
		t.Fatalf("got %d %d", total, available)
	}
}

func TestRowsAndHumanBytes(t *testing.T) {
	info := Info{OS: "macOS 27", Architecture: "arm64"}
	rows := info.Rows()
	if rows[0]["value"] != "macOS 27" || rows[1]["value"] != "Unavailable" {
		t.Fatalf("unexpected rows: %#v", rows)
	}
	if got := HumanBytes(1610612736); got != "1.5 GiB" {
		t.Fatalf("got %q", got)
	}
}
