//go:build linux
package scheduled
import("os/exec";"strings")
func Collect()([]Task,string,error){var tasks []Task;notice:="";if out,err:=exec.Command("crontab","-l").Output();err==nil {tasks=append(tasks,ParseCrontab(string(out),"user crontab")...)} else {notice="User crontab is unavailable to this process."};if out,err:=exec.Command("systemctl","list-timers","--all","--no-pager","--no-legend").Output();err==nil {for _,l:=range strings.Split(string(out),"\n"){f:=strings.Fields(l);if len(f)>=2{tasks=append(tasks,Task{"systemd timer",f[0],f[len(f)-1]})}}};return tasks,notice,nil}
