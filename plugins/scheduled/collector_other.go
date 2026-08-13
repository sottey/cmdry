//go:build !darwin && !linux
package scheduled
import "fmt"
func Collect()([]Task,string,error){return nil,"",fmt.Errorf("Scheduled Tasks is supported on Linux and macOS")}
