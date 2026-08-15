//go:build (!darwin || !cgo) && !linux

package wifi

import "fmt"

func Collect() (Info, error) {
	return Info{}, fmt.Errorf("Wi-Fi Inspector is supported on Linux and macOS")
}
