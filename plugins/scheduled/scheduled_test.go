package scheduled
import "testing"
func TestParseCrontab(t *testing.T){x:=ParseCrontab("# hi\n0 2 * * * /bin/backup --quiet\n","user crontab");if len(x)!=1||x[0].Schedule!="0 2 * * *"||x[0].Command!="/bin/backup --quiet"{t.Fatalf("%#v",x)}}
