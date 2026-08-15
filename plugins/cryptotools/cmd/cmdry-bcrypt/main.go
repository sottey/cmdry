package main

import (
	"fmt"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/cryptotools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.bcrypt", Name: "Bcrypt", Version: "0.1.0", Description: "Create or verify bcrypt password hashes locally.", Category: "security", Icon: "lock",
		SearchTerms: []string{"bcrypt", "hash", "password", "verify", "security"}, Pages: []cmdry.Page{{ID: "overview", Name: "Hash", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run Bcrypt", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "run": run}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Bcrypt", Components: []cmdry.Component{{Type: "form", Title: "Create or verify a bcrypt hash", Action: "run", Submit: "Run Bcrypt", Description: "Runs locally. Do not use real production passwords in a shared screen or terminal.", Fields: []cmdry.Field{
		{Name: "mode", Label: "Operation", Type: "select", Value: "hash", Options: []cmdry.Option{{Value: "hash", Label: "Create hash"}, {Value: "verify", Label: "Verify hash"}}},
		{Name: "input", Label: "Text or password", Type: "password", Required: true},
		{Name: "hash", Label: "Bcrypt hash (verification only)", Type: "password"},
		{Name: "cost", Label: "Cost (hashing only)", Type: "number", Value: "12", Min: "4", Max: "31", Required: true},
	}}}}, nil
}

func run(request cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(request.Params["input"])
	if fmt.Sprint(request.Params["mode"]) == "verify" {
		matched, err := cryptotools.Verify(input, fmt.Sprint(request.Params["hash"]))
		if err != nil {
			return cmdry.View{}, err
		}
		level, title, message := "info", "Bcrypt hash does not match", "The supplied text does not match this bcrypt hash."
		if matched {
			level, title, message = "success", "Bcrypt hash matches", "The supplied text matches this bcrypt hash."
		}
		return cmdry.View{Title: "Bcrypt verification", Components: []cmdry.Component{{Type: "alert", Level: level, Title: title, Message: message}, restart()}}, nil
	}
	cost, err := strconv.Atoi(fmt.Sprint(request.Params["cost"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("cost must be a number")
	}
	hash, err := cryptotools.Hash(input, cost)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Bcrypt hash", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Bcrypt hash created", Message: "The hash was created locally with cost " + strconv.Itoa(cost) + "."}, {Type: "code", Title: "Bcrypt hash", Text: hash}, restart()}}, nil
}

func restart() cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Start over"}}}
}
