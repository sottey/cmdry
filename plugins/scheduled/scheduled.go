package scheduled

import (
 "bufio"
 "os"
 "path/filepath"
 "strings"
)

type Task struct{ Source, Schedule, Command string }
func (t Task) Row() map[string]any { return map[string]any{"source":t.Source,"schedule":t.Schedule,"command":t.Command} }

func ParseCrontab(input, source string) []Task {
 var out []Task; s:=bufio.NewScanner(strings.NewReader(input)); for s.Scan(){ line:=strings.TrimSpace(s.Text()); if line==""||strings.HasPrefix(line,"#"){continue}; f:=strings.Fields(line); if len(f)<6 {continue}; out=append(out,Task{source,strings.Join(f[:5]," "),strings.Join(f[5:]," ")}) }; return out
}
func LaunchdTasks(dirs []string) []Task { var out []Task; for _,dir:=range dirs { entries,err:=os.ReadDir(dir); if err!=nil {continue}; for _,e:=range entries {if !e.IsDir()&&strings.HasSuffix(e.Name(),".plist"){out=append(out,Task{"launchd", "managed by launchd", filepath.Join(dir,e.Name())})}}}; return out }
