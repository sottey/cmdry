// Package scheduled reads configured scheduled tasks without modifying them.
package scheduled

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Task struct {
	Source, Schedule, Command string
}

func (t Task) Row() map[string]any {
	return map[string]any{"source": t.Source, "schedule": t.Schedule, "command": t.Command}
}

// ParseCrontab parses standard five-field cron entries and @special entries.
// Comments, blank lines, and malformed entries are omitted.
func ParseCrontab(input, source string) []Task {
	var tasks []Task
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if strings.HasPrefix(fields[0], "@") {
			tasks = append(tasks, Task{Source: source, Schedule: fields[0], Command: strings.Join(fields[1:], " ")})
			continue
		}
		if len(fields) < 6 {
			continue
		}
		tasks = append(tasks, Task{Source: source, Schedule: strings.Join(fields[:5], " "), Command: strings.Join(fields[5:], " ")})
	}
	return tasks
}

// LaunchdTasks lists launchd configuration files. Their presence establishes
// only that they are installed at that path; it does not infer loaded status.
func LaunchdTasks(directories []string) []Task {
	var tasks []Task
	for _, directory := range directories {
		entries, err := os.ReadDir(directory)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".plist") {
				continue
			}
			tasks = append(tasks, Task{
				Source:   "launchd configuration",
				Schedule: "configured by launchd",
				Command:  filepath.Join(directory, entry.Name()),
			})
		}
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].Command < tasks[j].Command })
	return tasks
}

// ParseSystemdTimers parses the final UNIT column from systemctl list-timers.
// The preceding timestamp fields vary by locale, so the parser does not try to
// reinterpret them beyond showing the first field as the next run value.
func ParseSystemdTimers(input string) []Task {
	var tasks []Task
	for _, line := range strings.Split(input, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		tasks = append(tasks, Task{Source: "systemd timer", Schedule: fields[0], Command: fields[len(fields)-2]})
	}
	return tasks
}
