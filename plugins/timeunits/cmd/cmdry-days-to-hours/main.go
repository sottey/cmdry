package main

import "github.com/sottey/cmdry/plugins/timeunits"

func main() {
	timeunits.Run("days-to-hours", "Convert Days to Hours", "days", "hours", timeunits.DaysToHours)
}
