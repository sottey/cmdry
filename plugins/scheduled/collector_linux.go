//go:build linux

package scheduled

import (
	"os/exec"
)

func Collect() ([]Task, string, error) {
	var tasks []Task
	notice := ""
	if output, err := exec.Command("crontab", "-l").Output(); err == nil {
		tasks = append(tasks, ParseCrontab(string(output), "user crontab")...)
	} else {
		notice = "User crontab is unavailable to this process. Systemd timers are still listed when available."
	}
	if output, err := exec.Command("systemctl", "list-timers", "--all", "--no-pager", "--no-legend", "--plain").Output(); err == nil {
		tasks = append(tasks, ParseSystemdTimers(string(output))...)
	}
	return tasks, notice, nil
}
