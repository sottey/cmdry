//go:build linux

package wifi

import (
	"os/exec"
	"strings"
)

func Collect() (Info, error) {
	out, err := exec.Command("nmcli", "-t", "-f", "ACTIVE,SSID,BSSID,SIGNAL,SECURITY", "dev", "wifi").Output()
	if err != nil {
		return Info{}, nil
	}
	info := ParseNMCLI(string(out))
	if !info.Available {
		return info, nil
	}
	if ip, err := exec.Command("hostname", "-I").Output(); err == nil {
		info.Address = strings.Fields(string(ip))[0]
	}
	return info, nil
}
