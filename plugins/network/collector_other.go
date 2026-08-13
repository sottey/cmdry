//go:build !darwin && !linux

package network

import "fmt"

func Collect() ([]Interface, string, error) {
	return nil, "", fmt.Errorf("Network Interface Inspector is supported on Linux and macOS")
}
