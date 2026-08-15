//go:build linux

package power

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Collect() (Info, error) {
	entries, err := os.ReadDir("/sys/class/power_supply")
	if err != nil {
		return Info{}, nil
	}
	info := Info{}
	for _, e := range entries {
		d := filepath.Join("/sys/class/power_supply", e.Name())
		typ, _ := os.ReadFile(filepath.Join(d, "type"))
		if strings.TrimSpace(string(typ)) == "Mains" {
			info.Source = "AC Power"
		}
		if strings.TrimSpace(string(typ)) != "Battery" {
			continue
		}
		info.Present = true
		if b, _ := os.ReadFile(filepath.Join(d, "capacity")); b != nil {
			info.Percent, _ = strconv.Atoi(strings.TrimSpace(string(b)))
		}
		if b, _ := os.ReadFile(filepath.Join(d, "status")); b != nil {
			info.State = strings.TrimSpace(string(b))
		}
	}
	return info, nil
}
