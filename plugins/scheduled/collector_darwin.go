//go:build darwin

package scheduled

import (
	"os"
	"os/exec"
)

func Collect() ([]Task, string, error) {
	tasks := LaunchdTasks([]string{
		"/Library/LaunchAgents",
		os.Getenv("HOME") + "/Library/LaunchAgents",
		"/Library/LaunchDaemons",
	})
	notice := ""
	if output, err := exec.Command("crontab", "-l").Output(); err == nil {
		tasks = append(tasks, ParseCrontab(string(output), "user crontab")...)
	} else {
		notice = "User crontab is unavailable to this process. Launchd configuration files are still listed when readable."
	}
	return tasks, notice, nil
}
