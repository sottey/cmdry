//go:build !linux

package journal

import "fmt"

func CollectRecent() ([]Entry, error) {
	return nil, fmt.Errorf("Journal Viewer requires Linux and journalctl")
}
