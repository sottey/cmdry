// Package cmdry provides helpers for independently compiled Cmdry plugins.
package cmdry

import (
	"encoding/json"
	"fmt"
	"os"
)

type Manifest struct {
	ProtocolVersion int      `json:"protocol_version"`
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	Icon            string   `json:"icon,omitempty"`
	Pages           []Page   `json:"pages"`
	Permissions     []string `json:"permissions"`
	Actions         []Action `json:"actions"`
}
type Page struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default,omitempty"`
	Action  string `json:"action,omitempty"`
}
type Action struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Method string `json:"method"`
}
type Request struct {
	Action string         `json:"action"`
	Params map[string]any `json:"params"`
}
type Response struct {
	OK    bool   `json:"ok"`
	Data  *View  `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type View struct {
	Title      string      `json:"title"`
	Components []Component `json:"components"`
}
type Component struct {
	Type        string           `json:"type"`
	ID          string           `json:"id,omitempty"`
	Label       string           `json:"label,omitempty"`
	Value       string           `json:"value,omitempty"`
	Description string           `json:"description,omitempty"`
	Title       string           `json:"title,omitempty"`
	Text        string           `json:"text,omitempty"`
	Level       string           `json:"level,omitempty"`
	Message     string           `json:"message,omitempty"`
	Columns     []Column         `json:"columns,omitempty"`
	Rows        []map[string]any `json:"rows,omitempty"`
	Actions     []Action         `json:"actions,omitempty"`
}
type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}
type Handler func(Request) (View, error)
type Plugin struct {
	Manifest Manifest
	Actions  map[string]Handler
}

func Run(p Plugin) {
	if len(os.Args) < 2 {
		fail("USAGE", "expected manifest or execute", 2)
		return
	}
	switch os.Args[1] {
	case "manifest":
		write(p.Manifest)
	case "execute":
		if len(os.Args) < 3 {
			fail("USAGE", "missing action", 2)
			return
		}
		var r Request
		if err := json.NewDecoder(os.Stdin).Decode(&r); err != nil {
			fail("INVALID_REQUEST", err.Error(), 2)
			return
		}
		h, ok := p.Actions[os.Args[2]]
		if !ok {
			fail("UNKNOWN_ACTION", "action is not registered", 2)
			return
		}
		v, err := h(r)
		if err != nil {
			fail("COMMAND_FAILED", err.Error(), 1)
			return
		}
		write(Response{OK: true, Data: &v})
	default:
		fail("USAGE", "unknown command", 2)
	}
}
func fail(code, message string, status int) {
	write(Response{OK: false, Error: &Error{Code: code, Message: message}})
	if status != 0 {
		os.Exit(status)
	}
}
func write(v any) {
	if err := json.NewEncoder(os.Stdout).Encode(v); err != nil {
		fmt.Fprintln(os.Stderr, "encode response:", err)
		os.Exit(1)
	}
}
