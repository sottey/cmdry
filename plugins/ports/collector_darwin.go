//go:build darwin

package ports

import (
	"fmt"
	"os/exec"
)

// CollectListeningPorts returns TCP listeners and unconnected UDP sockets
// reported by macOS lsof. macOS has no ss equivalent with the same data shape.
func CollectListeningPorts() ([]Port, error) {
	tcp, err := runLsof("-iTCP", "-sTCP:LISTEN")
	if err != nil {
		return nil, err
	}
	udp, err := runLsof("-iUDP")
	if err != nil {
		return nil, err
	}
	return append(ParseLsof(tcp), ParseLsof(udp)...), nil
}

func runLsof(args ...string) (string, error) {
	base := []string{"-nP", "-FpcPnT"}
	out, err := exec.Command("lsof", append(base, args...)...).Output()
	if err != nil {
		return "", fmt.Errorf("unable to execute lsof: %w", err)
	}
	return string(out), nil
}
