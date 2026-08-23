package networktools

import (
	"strings"
	"testing"
)

func TestCalculateSubnet(t *testing.T) {
	result, err := CalculateSubnet("192.168.7.233", 24)
	if err != nil {
		t.Fatal(err)
	}
	if result.Network != "192.168.7.0" || result.Broadcast != "192.168.7.255" || result.FirstHost != "192.168.7.1" || result.LastHost != "192.168.7.254" || result.Usable != 254 {
		t.Fatalf("unexpected subnet: %#v", result)
	}
}

func TestIPv4Representations(t *testing.T) {
	values, err := IPv4Representations("192.168.1.1")
	if err != nil {
		t.Fatal(err)
	}
	if values["decimal"] != "3232235777" || values["hex"] != "0xC0A80101" || values["binary"] != "11000000.10101000.00000001.00000001" {
		t.Fatalf("unexpected representations: %#v", values)
	}
}

func TestExpandIPv4Range(t *testing.T) {
	addresses, err := ExpandIPv4Range("192.168.7.0/30")
	if err != nil {
		t.Fatal(err)
	}
	if len(addresses) != 4 || addresses[0].String() != "192.168.7.0" || addresses[3].String() != "192.168.7.3" {
		t.Fatalf("unexpected addresses: %#v", addresses)
	}
	if _, err := ExpandIPv4Range("10.0.0.0/8"); err == nil {
		t.Fatal("accepted unbounded expansion")
	}
}

func TestMACUtilities(t *testing.T) {
	canonical, vendor, err := LookupMACVendor("00-50-56-aa-bb-cc")
	if err != nil || canonical != "00:50:56:AA:BB:CC" || vendor != "VMware" {
		t.Fatalf("unexpected lookup: %q, %q, %v", canonical, vendor, err)
	}
	generated, err := GenerateMAC()
	if err != nil {
		t.Fatal(err)
	}
	_, bytes, err := NormalizeMAC(generated)
	if err != nil || bytes[0]&0x03 != 0x02 {
		t.Fatalf("generated MAC is not locally administered unicast: %q, %v", generated, err)
	}
}

func TestGenerateULAPrefix(t *testing.T) {
	prefix, globalID, err := GenerateULAPrefix()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(prefix, "fd") || !strings.HasSuffix(prefix, "::/48") || len(globalID) != 10 {
		t.Fatalf("unexpected ULA: %q / %q", prefix, globalID)
	}
}
