package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/datetools"
	"strconv"
	"strings"
	"time"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.date-weekday", Name: "Date Weekday Calculator", Version: "0.1.0", Description: "Find years when a date falls on a selected weekday locally.", Category: "time", Icon: "calendar", SearchTerms: []string{"date", "weekday", "calendar", "years", "day of week"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New calculation", Method: "read"}, {ID: "run", Name: "Find years", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": form, "run": run}})
}
func form(cmdry.Request) (cmdry.View, error) {
	months := []cmdry.Option{}
	for i := 1; i <= 12; i++ {
		months = append(months, cmdry.Option{Value: fmt.Sprint(i), Label: time.Month(i).String()})
	}
	days := []cmdry.Option{}
	for _, d := range []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"} {
		days = append(days, cmdry.Option{Value: strings.ToLower(d), Label: d})
	}
	return cmdry.View{Title: "Date Weekday Calculator", Components: []cmdry.Component{{Type: "form", Title: "Find matching years", Action: "run", Submit: "Find years", Fields: []cmdry.Field{{Name: "month", Label: "Month", Type: "select", Value: "1", Options: months}, {Name: "day", Label: "Day", Type: "number", Value: "1", Min: "1", Max: "31", Required: true}, {Name: "weekday", Label: "Weekday", Type: "select", Value: "monday", Options: days}, {Name: "start", Label: "Start year", Type: "number", Value: "2026", Required: true}, {Name: "end", Label: "End year", Type: "number", Value: "2036", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	month, _ := strconv.Atoi(fmt.Sprint(r.Params["month"]))
	day, _ := strconv.Atoi(fmt.Sprint(r.Params["day"]))
	start, _ := strconv.Atoi(fmt.Sprint(r.Params["start"]))
	end, _ := strconv.Atoi(fmt.Sprint(r.Params["end"]))
	weekdays := map[string]time.Weekday{"sunday": time.Sunday, "monday": time.Monday, "tuesday": time.Tuesday, "wednesday": time.Wednesday, "thursday": time.Thursday, "friday": time.Friday, "saturday": time.Saturday}
	years, err := datetools.YearsWithWeekday(time.Month(month), day, start, end, weekdays[fmt.Sprint(r.Params["weekday"])])
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Matching years", Components: []cmdry.Component{{Type: "metric", Label: "Years found", Value: fmt.Sprint(len(years))}, {Type: "code", Title: "Years", Text: strings.Trim(fmt.Sprint(years), "[]")}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Try another date"}}}}}, nil
}
