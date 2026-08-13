//go:build darwin

package systeminfo

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func Collect() (Info, error) {
	info := Info{Architecture: runtime.GOARCH}
	if output, err := exec.Command("sw_vers", "-productName").Output(); err == nil {
		name := strings.TrimSpace(string(output))
		if version, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
			info.OS = name + " " + strings.TrimSpace(string(version))
		} else {
			info.OS = name
		}
	}
	if output, err := exec.Command("uname", "-sr").Output(); err == nil {
		info.Kernel = strings.TrimSpace(string(output))
	}
	if output, err := exec.Command("uptime").Output(); err == nil {
		info.Uptime = strings.TrimSpace(string(output))
	}
	if output, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
		info.CPU = strings.TrimSpace(string(output))
	}
	if output, err := exec.Command("sysctl", "-n", "hw.ncpu").Output(); err == nil {
		info.Cores, _ = strconv.Atoi(strings.TrimSpace(string(output)))
	}
	if output, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
		info.MemoryTotal, _ = strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
	}
	if output, err := exec.Command("sysctl", "-n", "hw.model").Output(); err == nil {
		info.Model = strings.TrimSpace(string(output))
	}
	return info, nil
}
