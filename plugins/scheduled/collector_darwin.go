//go:build darwin
package scheduled
import("os";"os/exec")
func Collect() ([]Task,string,error){ tasks:=LaunchdTasks([]string{"/Library/LaunchAgents",os.Getenv("HOME")+"/Library/LaunchAgents","/Library/LaunchDaemons"}); notice:=""; if out,err:=exec.Command("crontab","-l").Output();err==nil {tasks=append(tasks,ParseCrontab(string(out),"user crontab")...)} else {notice="User crontab is unavailable to this process."}; return tasks,notice,nil }
