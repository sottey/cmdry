package plugins

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Runner struct{ Timeout time.Duration }

func (r Runner) Run(parent context.Context, plugin Registered, action string, params map[string]any) (Response, error) {
	if !ValidID(action) {
		return Response{}, fmt.Errorf("invalid action ID")
	}
	known := false
	for _, a := range plugin.Manifest.Actions {
		if a.ID == action {
			known = true
		}
	}
	if !known {
		return Response{}, fmt.Errorf("unknown action %q", action)
	}
	body, err := json.Marshal(Request{Action: action, Params: params})
	if err != nil {
		return Response{}, err
	}
	ctx, cancel := context.WithTimeout(parent, r.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, plugin.Path, "execute", action)
	cmd.Stdin = bytes.NewReader(body)
	out, err := cmd.Output()
	if ctx.Err() != nil {
		return Response{}, fmt.Errorf("plugin timeout")
	}
	var response Response
	if decodeErr := json.Unmarshal(out, &response); decodeErr != nil {
		if err != nil {
			return Response{}, fmt.Errorf("plugin command failed: %w", err)
		}
		return Response{}, fmt.Errorf("invalid plugin response: %w", decodeErr)
	}
	if err := ValidateResponse(response); err != nil {
		return Response{}, err
	}
	// A plugin can intentionally return a structured protocol error with a
	// non-zero status. Preserve that useful administrator-facing message.
	if err != nil && response.OK {
		return Response{}, fmt.Errorf("plugin command failed: %w", err)
	}
	return response, nil
}
func ValidateResponse(r Response) error {
	if !r.OK {
		if r.Error == nil || r.Error.Code == "" || r.Error.Message == "" {
			return fmt.Errorf("invalid plugin error response")
		}
		return nil
	}
	if r.Data == nil {
		return fmt.Errorf("successful response has no data")
	}
	for _, c := range r.Data.Components {
		switch c.Type {
		case "metric", "text", "alert", "table", "actions", "code":
		case "form":
			if !ValidID(c.Action) || strings.TrimSpace(c.Submit) == "" || len(c.Fields) == 0 {
				return fmt.Errorf("invalid form component")
			}
			for _, field := range c.Fields {
				if !ValidID(field.Name) || strings.TrimSpace(field.Label) == "" || (field.Type != "text" && field.Type != "password" && field.Type != "textarea" && field.Type != "number" && field.Type != "checkbox" && field.Type != "select") {
					return fmt.Errorf("invalid form field")
				}
				if field.Type == "select" && len(field.Options) == 0 {
					return fmt.Errorf("select field has no options")
				}
			}
		case "download":
			if strings.TrimSpace(c.Filename) == "" || strings.TrimSpace(c.MIMEType) == "" || len(c.Content) == 0 || len(c.Content) > 8*1024*1024 {
				return fmt.Errorf("invalid download component")
			}
			if _, err := base64.StdEncoding.DecodeString(c.Content); err != nil {
				return fmt.Errorf("invalid download content: %w", err)
			}
		default:
			return fmt.Errorf("unsupported UI component %q", c.Type)
		}
	}
	return nil
}
