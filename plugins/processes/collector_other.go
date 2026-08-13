//go:build !darwin && !linux

package processes

import "fmt"

func Collect() ([]Process, error) {
	return nil, fmt.Errorf("Process Resource Snapshot is supported on Linux and macOS")
}
