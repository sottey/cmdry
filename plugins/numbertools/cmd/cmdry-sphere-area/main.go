package main

import (
	"fmt"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/numbertools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.sphere-area", Name: "Area of a Sphere", Version: "0.1.0", Description: "Calculate a sphere's surface area from its radius locally.", Category: "number", Icon: "calculate", SearchTerms: []string{"sphere", "area", "surface area", "radius", "geometry"}, Pages: []cmdry.Page{{ID: "overview", Name: "Calculate", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New calculation", Method: "read"}, {ID: "calculate", Name: "Calculate area", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "calculate": calculate}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Area of a Sphere", Components: []cmdry.Component{{Type: "form", Title: "Calculate surface area", Action: "calculate", Submit: "Calculate area", Description: "Uses the formula 4πr².", Fields: []cmdry.Field{{Name: "radius", Label: "Radius", Type: "number", Min: "0", Required: true, Placeholder: "2.5"}}}}}, nil
}

func calculate(request cmdry.Request) (cmdry.View, error) {
	radius, err := strconv.ParseFloat(fmt.Sprint(request.Params["radius"]), 64)
	if err != nil {
		return cmdry.View{}, fmt.Errorf("radius must be a number")
	}
	area, err := numbertools.SphereArea(radius)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Sphere surface area", Components: []cmdry.Component{{Type: "metric", Label: "Surface area", Value: strconv.FormatFloat(area, 'g', 12, 64)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Calculate another sphere"}}}}}, nil
}
