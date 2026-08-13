//go:build linux

package systeminfo

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func Collect() (Info, error) {
	info := Info{Architecture: runtime.GOARCH}
	if output, err := os.ReadFile("/etc/os-release"); err == nil {
		info.OS = ParseOSRelease(string(output))
	}
	if output, err := exec.Command("uname", "-sr").Output(); err == nil {
		info.Kernel = strings.TrimSpace(string(output))
	}
	if output, err := os.ReadFile("/proc/uptime"); err == nil {
		fields := strings.Fields(string(output))
		if len(fields) > 0 {
			if seconds, err := strconv.ParseFloat(fields[0], 64); err == nil {
				info.Uptime = formatDuration(int64(seconds))
			}
		}
	}
	if output, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		info.CPU = firstValue(string(output), "model name", "Hardware")
		info.Cores = strings.Count(string(output), "\nprocessor")
	}
	if output, err := os.ReadFile("/proc/meminfo"); err == nil {
		info.MemoryTotal, info.MemoryAvailable = ParseMemInfo(string(output))
	}
	if output, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_name"); err == nil {
		info.Model = strings.TrimSpace(string(output))
	}
	return info, nil
}

func firstValue(input string, keys ...string) string {
	for _, line := range strings.Split(input, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		for _, wanted := range keys {
			if strings.TrimSpace(key) == wanted {
				return strings.TrimSpace(value)
			}
		}
	}
	return ""
}

func formatDuration(seconds int64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	if days > 0 {
		return strconv.FormatInt(days, 10) + "d " + strconv.FormatInt(hours, 10) + "h"
	}
	return strconv.FormatInt(hours, 10) + "h " + strconv.FormatInt(minutes, 10) + "m"
}
