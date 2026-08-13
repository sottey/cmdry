//go:build darwin

package network

import (
	"fmt"
	"os/exec"
)

func Collect() ([]Interface, string, error) {
	out, err := exec.Command("ifconfig").Output()
	if err != nil {
		return nil, "", fmt.Errorf("unable to execute ifconfig: %w", err)
	}
	interfaces := ParseDarwinIfconfig(string(out))
	gateway := ""
	if route, err := exec.Command("route", "-n", "get", "default").Output(); err == nil {
		gateway = ParseGateway(string(route))
	}
	return interfaces, gateway, nil
}
