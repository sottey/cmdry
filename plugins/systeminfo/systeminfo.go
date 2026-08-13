// Package systeminfo normalizes small, read-only host facts.
package systeminfo

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Info struct {
	OS, Kernel, Architecture, Uptime, CPU, Model string
	Cores                                        int
	MemoryTotal, MemoryAvailable                 int64
}

func (i Info) Rows() []map[string]any {
	rows := []map[string]any{
		{"label": "Operating system", "value": unknown(i.OS)},
		{"label": "Kernel", "value": unknown(i.Kernel)},
		{"label": "Architecture", "value": unknown(i.Architecture)},
		{"label": "Uptime", "value": unknown(i.Uptime)},
		{"label": "CPU", "value": unknown(i.CPU)},
		{"label": "Hardware model", "value": unknown(i.Model)},
	}
	return rows
}

func unknown(value string) string {
	if strings.TrimSpace(value) == "" {
		return "Unavailable"
	}
	return value
}

func HumanBytes(value int64) string {
	if value <= 0 {
		return "Unavailable"
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	divisor := float64(1)
	for _, unit := range units {
		if float64(value) < divisor*1024 || unit == units[len(units)-1] {
			return fmt.Sprintf("%.1f %s", float64(value)/divisor, unit)
		}
		divisor *= 1024
	}
	return strconv.FormatInt(value, 10) + " B"
}

func ParseOSRelease(input string) string {
	values := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		key, value, ok := strings.Cut(scanner.Text(), "=")
		if ok {
			values[key] = strings.Trim(value, `"`)
		}
	}
	if values["PRETTY_NAME"] != "" {
		return values["PRETTY_NAME"]
	}
	return values["NAME"]
}

func ParseMemInfo(input string) (total, available int64) {
	for _, line := range strings.Split(input, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		fields := strings.Fields(value)
		if len(fields) == 0 {
			continue
		}
		amount, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			continue
		}
		switch key {
		case "MemTotal":
			total = amount * 1024
		case "MemAvailable":
			available = amount * 1024
		}
	}
	return total, available
}
