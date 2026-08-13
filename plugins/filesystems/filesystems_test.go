package filesystems

import "testing"

func TestParseAndSummary(t *testing.T) {
	input := `Filesystem 1024-blocks Used Available Capacity Mounted on
/dev/sda1 1048576 524288 524288 50% /
server:/share 2048 1024 1024 50% /Volumes/Shared Files
map auto_home 0 0 0 100% /System/Volumes/Data/home
`
	items, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 || items[1].Mount != "/Volumes/Shared Files" || items[1].Capacity != 50 || items[2].Device != "map auto_home" {
		t.Fatalf("unexpected filesystems: %#v", items)
	}
	row := items[0].Row()
	if row["total"] != "1.0 GiB" || row["capacity"] != "50%" {
		t.Fatalf("unexpected row: %#v", row)
	}
}

func TestParseRejectsInvalidData(t *testing.T) {
	if _, err := Parse("Filesystem 1024-blocks Used Available Capacity Mounted on\n/dev/a 1 1 1 loud /\n"); err == nil {
		t.Fatal("accepted invalid capacity")
	}
}

func TestHumanBytes(t *testing.T) {
	if got := HumanBytes(1536); got != "1.5 KiB" {
		t.Fatalf("got %q", got)
	}
}
