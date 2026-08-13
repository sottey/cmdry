package filesystems

import (
	"fmt"
	"os/exec"
)

// Collect reports mounted filesystems using the portable POSIX df format.
func Collect() ([]Filesystem, error) {
	out, err := exec.Command("df", "-kP").Output()
	if err != nil {
		return nil, fmt.Errorf("unable to execute df: %w", err)
	}
	return Parse(string(out))
}
