package wifi

import "testing"

func TestParsers(t *testing.T) {
	if got := ParseAirportNetwork("Current Wi-Fi Network: ExampleNet\n"); got != "ExampleNet" {
		t.Fatal(got)
	}
	if got := ParseHardwarePorts("Hardware Port: Wi-Fi\nDevice: en0\nEthernet Address: aa\n"); got != "en0" {
		t.Fatal(got)
	}
	i := ParseNMCLI("no:x:y:1:--\nyes:Example:aa:78:WPA2\n")
	if !i.Available || i.SSID != "Example" || i.Signal != "78%" {
		t.Fatalf("%#v", i)
	}
}
