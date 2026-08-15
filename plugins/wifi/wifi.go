package wifi

import "strings"

type Info struct {
	Available                                         bool
	Interface, SSID, BSSID, Signal, Security, Address string
}

func (i Info) Rows() []map[string]any {
	return []map[string]any{{"label": "Interface", "value": value(i.Interface)}, {"label": "SSID", "value": value(i.SSID)}, {"label": "BSSID", "value": value(i.BSSID)}, {"label": "Signal", "value": value(i.Signal)}, {"label": "Security", "value": value(i.Security)}, {"label": "IPv4 address", "value": value(i.Address)}}
}
func value(s string) string {
	if strings.TrimSpace(s) == "" {
		return "Unavailable"
	}
	return s
}

func ParseAirportNetwork(output string) string {
	const prefix = "Current Wi-Fi Network: "
	line := strings.TrimSpace(output)
	if strings.HasPrefix(line, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(line, prefix))
	}
	return ""
}
func ParseHardwarePorts(output string) string {
	blocks := strings.Split(output, "\n\n")
	for _, b := range blocks {
		if !strings.Contains(b, "Hardware Port: Wi-Fi") {
			continue
		}
		for _, line := range strings.Split(b, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Device: ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "Device: "))
			}
		}
	}
	return ""
}
func ParseNMCLI(output string) Info {
	for _, line := range strings.Split(output, "\n") {
		f := strings.Split(line, ":")
		if len(f) >= 5 && f[0] == "yes" {
			return Info{Available: true, SSID: f[1], BSSID: f[2], Signal: f[3] + "%", Security: f[4]}
		}
	}
	return Info{}
}
