package scheduled

import "testing"

func TestParseCrontab(t *testing.T) {
	tasks := ParseCrontab("# comment\n0 2 * * * /bin/backup --quiet\n@daily /usr/local/bin/report\ninvalid\n", "user crontab")
	if len(tasks) != 2 {
		t.Fatalf("got %#v", tasks)
	}
	if tasks[0].Schedule != "0 2 * * *" || tasks[0].Command != "/bin/backup --quiet" {
		t.Fatalf("unexpected standard task: %#v", tasks[0])
	}
	if tasks[1].Schedule != "@daily" || tasks[1].Command != "/usr/local/bin/report" {
		t.Fatalf("unexpected special task: %#v", tasks[1])
	}
}

func TestParseSystemdTimers(t *testing.T) {
	tasks := ParseSystemdTimers("Wed 2026-08-13 12:00:00 PDT 1h left Wed 2026-08-13 11:00:00 PDT 1h ago fstrim.timer fstrim.service\n")
	if len(tasks) != 1 || tasks[0].Source != "systemd timer" || tasks[0].Command != "fstrim.timer" {
		t.Fatalf("got %#v", tasks)
	}
}
