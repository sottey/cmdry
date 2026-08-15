//go:build !darwin && !linux

package power

import "fmt"

func Collect() (Info, error) {
	return Info{}, fmt.Errorf("Battery and Power Inspector is supported on Linux and macOS")
}
