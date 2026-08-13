// Package processes parses a portable, selected-column ps output.
package processes

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Process struct {
	PID, ParentPID int
	CPU, Memory    float64
	State, Command string
}

// Parse decodes lines with PID, PPID, CPU percentage, memory percentage,
// state, and executable name, in that order.
func Parse(input string) ([]Process, error) {
	result := make([]Process, 0)
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 6 {
			continue
		}
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("parse PID %q: %w", fields[0], err)
		}
		parentPID, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("parse parent PID %q: %w", fields[1], err)
		}
		cpu, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, fmt.Errorf("parse CPU %q: %w", fields[2], err)
		}
		memory, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return nil, fmt.Errorf("parse memory %q: %w", fields[3], err)
		}
		result = append(result, Process{PID: pid, ParentPID: parentPID, CPU: cpu, Memory: memory, State: fields[4], Command: strings.Join(fields[5:], " ")})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read process output: %w", err)
	}
	return result, nil
}

func (p Process) Row() map[string]any {
	return map[string]any{
		"pid":     p.PID,
		"parent":  p.ParentPID,
		"cpu":     fmt.Sprintf("%.1f%%", p.CPU),
		"memory":  fmt.Sprintf("%.1f%%", p.Memory),
		"state":   p.State,
		"command": p.Command,
	}
}

func Summary(items []Process) (running int, totalCPU, totalMemory float64) {
	for _, item := range items {
		if strings.Contains(item.State, "R") {
			running++
		}
		totalCPU += item.CPU
		totalMemory += item.Memory
	}
	return running, totalCPU, totalMemory
}
