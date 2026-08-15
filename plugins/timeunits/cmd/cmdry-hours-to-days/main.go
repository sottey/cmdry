package main

import "github.com/sottey/cmdry/plugins/timeunits"

func main() {
	timeunits.Run("hours-to-days", "Convert Hours to Days", "hours", "days", timeunits.HoursToDays)
}
