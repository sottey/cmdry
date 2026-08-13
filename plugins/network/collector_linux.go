//go:build linux

package network

import (
	"fmt"
	"os/exec"
)

func Collect() ([]Interface, string, error) {
	out, err := exec.Command("ip", "-o", "addr", "show").Output()
	if err != nil {
		return nil, "", fmt.Errorf("unable to execute ip: %w", err)
	}
	interfaces := ParseLinuxAddresses(string(out))
	gateway := ""
	if route, err := exec.Command("ip", "route", "show", "default").Output(); err == nil {
		gateway = ParseGateway(string(route))
	}
	return interfaces, gateway, nil
}
