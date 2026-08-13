//go:build darwin

package processes

import (
	"fmt"
	"os/exec"
)

func Collect() ([]Process, error) {
	out, err := exec.Command("ps", "-axo", "pid=,ppid=,%cpu=,%mem=,state=,comm=").Output()
	if err != nil {
		return nil, fmt.Errorf("unable to execute ps: %w", err)
	}
	return Parse(string(out))
}
