package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
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
		case "metric", "text", "alert", "table", "actions":
		default:
			return fmt.Errorf("unsupported UI component %q", c.Type)
		}
	}
	return nil
}
