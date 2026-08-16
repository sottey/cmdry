package main

import (
	"crypto/rand"
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"math/big"
	"strconv"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.random-port", Name: "Random Port Generator", Version: "0.1.0", Description: "Generate random unprivileged TCP/UDP port numbers locally.", Category: "number", Icon: "network", SearchTerms: []string{"random", "port", "tcp", "udp", "network"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Generate ports", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": form, "run": run}})
}
func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Random Port Generator", Components: []cmdry.Component{{Type: "form", Title: "Generate ports", Action: "run", Submit: "Generate", Fields: []cmdry.Field{{Name: "count", Label: "Ports", Type: "number", Value: "5", Min: "1", Max: "100", Required: true}}}}}, nil
}
func run(r cmdry.Request) (cmdry.View, error) {
	count, e := strconv.Atoi(fmt.Sprint(r.Params["count"]))
	if e != nil || count < 1 || count > 100 {
		return cmdry.View{}, fmt.Errorf("count must be between 1 and 100")
	}
	values := []string{}
	seen := map[int]bool{}
	for len(values) < count {
		n, _ := rand.Int(rand.Reader, big.NewInt(65535-1024+1))
		p := int(n.Int64()) + 1024
		if !seen[p] {
			seen[p] = true
			values = append(values, fmt.Sprint(p))
		}
	}
	return cmdry.View{Title: "Random ports", Components: []cmdry.Component{{Type: "code", Title: "Ports", Text: fmt.Sprint(values)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Generate again"}}}}}, nil
}
