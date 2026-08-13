//go:build !darwin && !linux

package systeminfo

import "fmt"

func Collect() (Info, error) {
	return Info{}, fmt.Errorf("System Information is supported on Linux and macOS")
}
