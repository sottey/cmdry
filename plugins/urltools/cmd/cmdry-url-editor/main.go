package main

import (
	"fmt"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/urltools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.url-editor", Name: "URL Editor", Version: "0.1.0", Description: "Inspect and edit URL components and query parameters locally.", Category: "text", Icon: "link",
		SearchTerms: []string{"url", "editor", "query", "parameters", "link", "parse"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "build", Name: "Build URL", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": form, "build": build}})
}

func form(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "URL Editor", Components: []cmdry.Component{{Type: "form", Title: "Inspect and edit a URL", Action: "build", Submit: "Build URL", Description: "Paste a complete URL. Any non-empty component field replaces its matching part; blank component fields retain the pasted URL's existing part.", Fields: []cmdry.Field{{Name: "url", Label: "URL", Type: "text", Required: true, Placeholder: "https://example.com/path?view=list#details"}, {Name: "scheme", Label: "Scheme override", Type: "text", Placeholder: "https"}, {Name: "host", Label: "Host override", Type: "text", Placeholder: "example.com"}, {Name: "path", Label: "Path override", Type: "text", Placeholder: "/new-path"}, {Name: "query", Label: "Query override", Type: "text", Placeholder: "page=2&sort=name"}, {Name: "fragment", Label: "Fragment override", Type: "text", Placeholder: "details"}}}}}, nil
}

func build(request cmdry.Request) (cmdry.View, error) {
	input := fmt.Sprint(request.Params["url"])
	if strings.TrimSpace(input) == "" {
		return cmdry.View{}, fmt.Errorf("URL is required")
	}
	output, details, err := urltools.Build(input, value(request, "scheme"), value(request, "host"), value(request, "path"), value(request, "query"), value(request, "fragment"))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Edited URL", Components: []cmdry.Component{{Type: "code", Title: "URL", Text: output}, {Type: "metric", Label: "Scheme", Value: details.Scheme}, {Type: "metric", Label: "Host", Value: details.Host}, {Type: "metric", Label: "Path", Value: details.Path}, {Type: "metric", Label: "Query parameters", Value: queryCount(details.Query)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Edit another URL"}}}}}, nil
}

func queryCount(query string) string {
	if query == "" {
		return "0"
	}
	return fmt.Sprint(len(strings.Split(query, "&")))
}

func value(request cmdry.Request, name string) string {
	value, _ := request.Params[name].(string)
	return value
}
