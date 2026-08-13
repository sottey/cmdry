//go:build linux

package ports

import (
	"fmt"
	"os/exec"
)

// CollectListeningPorts returns the listening TCP and UDP sockets reported by
// ss on Linux.
func CollectListeningPorts() ([]Port, error) {
	out, err := exec.Command("ss", "-H", "-l", "-n", "-t", "-u", "-p").Output()
	if err != nil {
		return nil, fmt.Errorf("unable to execute ss: %w", err)
	}
	return ParseSS(string(out)), nil
}
