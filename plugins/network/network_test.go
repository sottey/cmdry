package network

import "testing"

func TestParseDarwinIfconfig(t *testing.T) {
	input := `en0: flags=8863<UP,BROADCAST,RUNNING> mtu 1500
	inet6 fe80::1234%en0 prefixlen 64
	inet 192.168.7.169 netmask 0xfffffc00
	status: active
lo0: flags=8049<UP,LOOPBACK,RUNNING> mtu 16384
	inet 127.0.0.1 netmask 0xff000000
`
	items := ParseDarwinIfconfig(input)
	if len(items) != 2 || items[0].Name != "en0" || items[0].Status != "active" || len(items[0].Addresses) != 2 || items[0].Addresses[0] != "fe80::1234" {
		t.Fatalf("unexpected interfaces: %#v", items)
	}
}

func TestParseLinuxAddressesAndGateway(t *testing.T) {
	items := ParseLinuxAddresses("2: eth0    inet 192.168.1.10/24 brd 192.168.1.255 scope global eth0\n2: eth0    inet6 fe80::1/64 scope link\n")
	if len(items) != 1 || items[0].Name != "eth0" || len(items[0].Addresses) != 2 || items[0].Addresses[1] != "fe80::1" {
		t.Fatalf("unexpected interfaces: %#v", items)
	}
	if got := ParseGateway("default via 192.168.1.1 dev eth0\n"); got != "192.168.1.1" {
		t.Fatalf("got %q", got)
	}
	if got := ParseGateway("   gateway: 192.168.7.1\n"); got != "192.168.7.1" {
		t.Fatalf("got %q", got)
	}
}
