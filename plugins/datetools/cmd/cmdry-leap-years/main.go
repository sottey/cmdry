package main

import (
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/datetools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.leap-years", Name: "Check Leap Years", Version: "0.1.0", Description: "Find Gregorian leap years in a selected range locally.", Category: "time", Icon: "calendar", SearchTerms: []string{"leap", "year", "calendar", "february", "date"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New range", Method: "read"}, {ID: "run", Name: "Check leap years", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}
func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Check Leap Years", Components: []cmdry.Component{{Type: "form", Title: "Find leap years", Action: "run", Submit: "Check leap years", Fields: []cmdry.Field{{Name: "start", Label: "Start year", Type: "number", Value: "2000", Min: "1", Required: true}, {Name: "end", Label: "End year", Type: "number", Value: "2030", Min: "1", Required: true}}}}}, nil
}
func run(request cmdry.Request) (cmdry.View, error) {
	start, err := strconv.Atoi(fmt.Sprint(request.Params["start"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("start year must be a whole number")
	}
	end, err := strconv.Atoi(fmt.Sprint(request.Params["end"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("end year must be a whole number")
	}
	years, err := datetools.LeapYears(start, end)
	if err != nil {
		return cmdry.View{}, err
	}
	values := make([]string, len(years))
	for index, year := range years {
		values[index] = fmt.Sprint(year)
	}
	message := "No leap years occur in this range."
	if len(years) > 0 {
		message = "Leap years have 366 days, including February 29."
	}
	return cmdry.View{Title: "Leap years", Components: []cmdry.Component{{Type: "metric", Label: "Leap years found", Value: fmt.Sprint(len(years))}, {Type: "alert", Level: "info", Title: "Gregorian calendar", Message: message}, {Type: "code", Title: "Leap years", Text: strings.Join(values, "\n")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Check another range"}}}}}, nil
}
