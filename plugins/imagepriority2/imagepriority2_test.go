package imagepriority2

import (
	"strings"
	"testing"
)

func TestWiFiPayload(t *testing.T) {
	payload, err := WiFiPayload("Cafe;Net", "pa:ss", "WPA", true)
	if err != nil || payload != "WIFI:T:WPA;S:Cafe\\;Net;P:pa\\:ss;H:true;;" {
		t.Fatalf("payload=%q err=%v", payload, err)
	}
}
func TestSVGPlaceholder(t *testing.T) {
	svg, err := SVGPlaceholder(640, 360, "Hello <Cmdry>", "#112233", "#ffffff")
	if err != nil || !strings.Contains(string(svg), "Hello &lt;Cmdry&gt;") {
		t.Fatalf("svg=%q err=%v", svg, err)
	}
}
