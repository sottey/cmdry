package main

import (
	"fmt"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/xmlutils"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.xml-validator", Name: "XML Validator", Version: "0.1.0", Description: "Validate pasted XML locally and report parser errors.", Category: "developer", Icon: "xml",
		SearchTerms: []string{"xml", "validate", "validator", "syntax", "well formed"}, Pages: []cmdry.Page{{ID: "overview", Name: "Validate", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New validation", Method: "read"}, {ID: "validate", Name: "Validate XML", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "validate": validate}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "XML Validator", Components: []cmdry.Component{{Type: "form", Title: "Validate XML", Action: "validate", Submit: "Validate XML", Description: "Checks well-formed XML locally. Parser errors include a line number when XML provides one.", Fields: []cmdry.Field{{Name: "input", Label: "XML input", Type: "textarea", Required: true}}}}}, nil
}

func validate(request cmdry.Request) (cmdry.View, error) {
	info, err := xmlutils.Validate(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("invalid XML: %w", err)
	}
	return cmdry.View{Title: "Valid XML", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "XML is well formed", Message: "The input contains one well-formed XML root element."}, {Type: "metric", Label: "Root element", Value: info.Root}, {Type: "metric", Label: "Root attributes", Value: fmt.Sprint(info.Attributes)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Validate another document"}}}}}, nil
}
