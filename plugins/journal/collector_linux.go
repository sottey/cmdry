//go:build linux

package journal

import (
	"fmt"
	"os/exec"
)

// CollectRecent returns the newest 100 local journal entries, newest first.
func CollectRecent() ([]Entry, error) {
	if _, err := exec.LookPath("journalctl"); err != nil {
		return nil, fmt.Errorf("journalctl is unavailable: %w", err)
	}
	out, err := exec.Command("journalctl", "--no-pager", "--output=json", "--lines=100", "--reverse").Output()
	if err != nil {
		return nil, fmt.Errorf("unable to read the system journal: %w", err)
	}
	return Parse(out)
}
