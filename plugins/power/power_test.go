package power

import "testing"

func TestParsePMSet(t *testing.T) {
	i := ParsePMSet("Now drawing from 'AC Power'\n -InternalBattery-0 (id=1)\t90%; discharging; 4:21 remaining present: true\n")
	if !i.Present || i.Source != "AC Power" || i.Percent != 90 || i.State != "discharging" || i.Remaining != "4:21 remaining" {
		t.Fatalf("%#v", i)
	}
}
