// Package urltools provides safe local URL inspection and editing helpers.
package urltools

import (
	"fmt"
	"net/url"
	"strings"
)

// Details contains the editable parts of a complete URL.
type Details struct {
	Scheme   string
	Host     string
	Path     string
	Query    string
	Fragment string
}

// Parse validates a complete URL and returns its editable components.
func Parse(input string) (Details, error) {
	parsed, err := url.Parse(strings.TrimSpace(input))
	if err != nil {
		return Details{}, fmt.Errorf("parse URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return Details{}, fmt.Errorf("enter a complete URL including scheme and host")
	}
	return Details{Scheme: parsed.Scheme, Host: parsed.Host, Path: parsed.EscapedPath(), Query: parsed.RawQuery, Fragment: parsed.Fragment}, nil
}

// Build applies non-empty component overrides to an input URL. Empty component
// fields preserve the original component, so the tool cannot accidentally
// discard a URL part when a user is only changing one other field.
func Build(input string, scheme string, host string, path string, query string, fragment string) (string, Details, error) {
	parsed, err := url.Parse(strings.TrimSpace(input))
	if err != nil {
		return "", Details{}, fmt.Errorf("parse URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", Details{}, fmt.Errorf("enter a complete URL including scheme and host")
	}
	if strings.TrimSpace(scheme) != "" {
		parsed.Scheme = strings.TrimSpace(scheme)
	}
	if strings.TrimSpace(host) != "" {
		parsed.Host = strings.TrimSpace(host)
	}
	if path != "" {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		parsed.Path = path
		parsed.RawPath = ""
	}
	if query != "" {
		if _, err := url.ParseQuery(query); err != nil {
			return "", Details{}, fmt.Errorf("query parameters: %w", err)
		}
		parsed.RawQuery = query
	}
	if fragment != "" {
		parsed.Fragment = fragment
	}
	details := Details{Scheme: parsed.Scheme, Host: parsed.Host, Path: parsed.EscapedPath(), Query: parsed.RawQuery, Fragment: parsed.Fragment}
	return parsed.String(), details, nil
}
