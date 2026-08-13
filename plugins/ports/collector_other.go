//go:build !linux && !darwin

package ports

import "fmt"

func CollectListeningPorts() ([]Port, error) {
	return nil, fmt.Errorf("Port Inspector is supported on Linux and macOS")
}
