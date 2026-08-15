//go:build darwin

package power

import (
	"fmt"
	"os/exec"
)

func Collect() (Info, error) {
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return Info{}, fmt.Errorf("unable to execute pmset: %w", err)
	}
	return ParsePMSet(string(out)), nil
}
