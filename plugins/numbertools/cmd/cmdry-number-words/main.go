package main

import (
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/numbertools"
	"strconv"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.number-words", Name: "Number to Words", Version: "0.1.0", Description: "Convert whole numbers to English words locally.", Category: "number", Icon: "text", SearchTerms: []string{"number", "words", "spell", "integer"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Convert number", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": form, "run": run}})
}
func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Number to Words", Components: []cmdry.Component{{Type: "form", Title: "Convert number", Action: "run", Submit: "Convert", Fields: []cmdry.Field{{Name: "input", Label: "Whole number", Type: "number", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	n, e := strconv.ParseInt(fmt.Sprint(r.Params["input"]), 10, 64)
	if e != nil {
		return cmdry.View{}, fmt.Errorf("enter a whole number")
	}
	return cmdry.View{Title: "Number in words", Components: []cmdry.Component{{Type: "code", Title: "Words", Text: numbertools.NumberToWords(n)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another"}}}}}, nil
}
