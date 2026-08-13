//go:build linux

package processes

import (
	"fmt"
	"os/exec"
)

func Collect() ([]Process, error) {
	out, err := exec.Command("ps", "-eo", "pid=,ppid=,pcpu=,pmem=,stat=,comm=").Output()
	if err != nil {
		return nil, fmt.Errorf("unable to execute ps: %w", err)
	}
	return Parse(string(out))
}
