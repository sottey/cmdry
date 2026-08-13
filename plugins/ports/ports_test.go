package ports

import "testing"

func TestParseSS(t *testing.T) {
	input := `tcp LISTEN 0 4096 0.0.0.0:22 0.0.0.0:* users:(("sshd",pid=881,fd=3))
udp UNCONN 0 0 [::1]:53 [::]:* users:(("resolved",pid=712,fd=2))
tcp LISTEN 0 128 127.0.0.1:8080 0.0.0.0:*`
	got := ParseSS(input)
	if len(got) != 3 {
		t.Fatalf("got %d ports", len(got))
	}
	if got[0].Port != 22 || got[0].Protocol != "TCP" || got[0].Process != "sshd" || got[0].PID == nil || *got[0].PID != 881 {
		t.Fatalf("unexpected TCP port: %#v", got[0])
	}
	if got[1].Address != "::1" || got[1].Port != 53 || got[1].PID == nil {
		t.Fatalf("unexpected IPv6 port: %#v", got[1])
	}
	if got[2].Process != "" || got[2].PID != nil {
		t.Fatalf("missing process should be preserved: %#v", got[2])
	}
}

func TestParseSSIgnoresMalformed(t *testing.T) {
	if got := ParseSS("garbage\ntcp LISTEN"); len(got) != 0 {
		t.Fatalf("got %#v", got)
	}
}
