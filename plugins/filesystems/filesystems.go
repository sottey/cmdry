// Package filesystems parses the portable POSIX df -kP output.
package filesystems

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Filesystem struct {
	Device                               string
	TotalBlocks, UsedBlocks, AvailBlocks int64
	Capacity                             int
	Mount                                string
}

// Parse decodes df -kP output. POSIX output has six fields; the mount point is
// the remainder of the line so paths containing spaces remain intact.
func Parse(input string) ([]Filesystem, error) {
	items := make([]Filesystem, 0)
	scanner := bufio.NewScanner(strings.NewReader(input))
	first := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if first {
			first = false
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		numbersAt := -1
		for index, field := range fields {
			if _, err := strconv.ParseInt(field, 10, 64); err == nil {
				numbersAt = index
				break
			}
		}
		if numbersAt < 1 || len(fields) < numbersAt+5 {
			return nil, fmt.Errorf("parse df entry %q", line)
		}
		device := strings.Join(fields[:numbersAt], " ")
		total, err := strconv.ParseInt(fields[numbersAt], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse total blocks for %q: %w", device, err)
		}
		used, err := strconv.ParseInt(fields[numbersAt+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse used blocks for %q: %w", device, err)
		}
		available, err := strconv.ParseInt(fields[numbersAt+2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse available blocks for %q: %w", device, err)
		}
		capacity, err := strconv.Atoi(strings.TrimSuffix(fields[numbersAt+3], "%"))
		if err != nil {
			return nil, fmt.Errorf("parse capacity for %q: %w", device, err)
		}
		items = append(items, Filesystem{Device: device, TotalBlocks: total, UsedBlocks: used, AvailBlocks: available, Capacity: capacity, Mount: strings.Join(fields[numbersAt+4:], " ")})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read df output: %w", err)
	}
	return items, nil
}

func (f Filesystem) Row() map[string]any {
	return map[string]any{
		"filesystem": f.Device,
		"total":      HumanBytes(f.TotalBlocks * 1024),
		"used":       HumanBytes(f.UsedBlocks * 1024),
		"available":  HumanBytes(f.AvailBlocks * 1024),
		"capacity":   strconv.Itoa(f.Capacity) + "%",
		"mount":      f.Mount,
	}
}

func HumanBytes(value int64) string {
	const unit = int64(1024)
	if value < unit {
		return strconv.FormatInt(value, 10) + " B"
	}
	units := []string{"KiB", "MiB", "GiB", "TiB", "PiB"}
	divisor := unit
	for _, label := range units {
		if value < divisor*unit || label == units[len(units)-1] {
			return fmt.Sprintf("%.1f %s", float64(value)/float64(divisor), label)
		}
		divisor *= unit
	}
	return strconv.FormatInt(value, 10) + " B"
}
