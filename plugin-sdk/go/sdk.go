// Package cmdry provides helpers for independently compiled Cmdry plugins.
package cmdry

import (
	"encoding/base64"
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
	SearchTerms     []string `json:"search_terms,omitempty"`
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
type Upload struct {
	Name     string `json:"name"`
	MIMEType string `json:"mime_type"`
	Content  string `json:"content"`
}

// File returns an uploaded file carried in a file form field. Content is held
// only in this request as standard base64; plugins must never expect a path.
func (r Request) File(name string) (Upload, []byte, error) {
	value, ok := r.Params[name]
	if !ok {
		return Upload{}, nil, fmt.Errorf("%s is required", name)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return Upload{}, nil, fmt.Errorf("read %s: %w", name, err)
	}
	var upload Upload
	if err := json.Unmarshal(encoded, &upload); err != nil {
		return Upload{}, nil, fmt.Errorf("read %s: %w", name, err)
	}
	if upload.Name == "" || upload.MIMEType == "" || upload.Content == "" {
		return Upload{}, nil, fmt.Errorf("%s is not a valid upload", name)
	}
	contents, err := base64.StdEncoding.DecodeString(upload.Content)
	if err != nil {
		return Upload{}, nil, fmt.Errorf("decode %s: %w", name, err)
	}
	return upload, contents, nil
}

// Files returns uploads from a multiple file form field. Content is held only
// in this request as standard base64; plugins must never expect host paths.
func (r Request) Files(name string) ([]Upload, [][]byte, error) {
	value, ok := r.Params[name]
	if !ok {
		return nil, nil, fmt.Errorf("%s is required", name)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, nil, fmt.Errorf("read %s: %w", name, err)
	}
	var uploads []Upload
	if err := json.Unmarshal(encoded, &uploads); err != nil {
		return nil, nil, fmt.Errorf("read %s: %w", name, err)
	}
	if len(uploads) == 0 {
		return nil, nil, fmt.Errorf("%s is required", name)
	}
	contents := make([][]byte, len(uploads))
	for index, upload := range uploads {
		if upload.Name == "" || upload.MIMEType == "" || upload.Content == "" {
			return nil, nil, fmt.Errorf("%s is not a valid upload", name)
		}
		contents[index], err = base64.StdEncoding.DecodeString(upload.Content)
		if err != nil {
			return nil, nil, fmt.Errorf("decode %s: %w", name, err)
		}
	}
	return uploads, contents, nil
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
	Action      string           `json:"action,omitempty"`
	Submit      string           `json:"submit,omitempty"`
	Fields      []Field          `json:"fields,omitempty"`
	Filename    string           `json:"filename,omitempty"`
	MIMEType    string           `json:"mime_type,omitempty"`
	Content     string           `json:"content,omitempty"`
}
type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}
type Field struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Value       string   `json:"value,omitempty"`
	Min         string   `json:"min,omitempty"`
	Max         string   `json:"max,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
	Description string   `json:"description,omitempty"`
	Accept      string   `json:"accept,omitempty"`
	Multiple    bool     `json:"multiple,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Options     []Option `json:"options,omitempty"`
}
type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
type Handler func(Request) (View, error)
type Plugin struct {
	Manifest Manifest
	Actions  map[string]Handler
}

// BuildVersion is set for bundled plugins by the release build. It is empty
// for independently compiled plugins, which retain their declared version.
var BuildVersion string

func Run(p Plugin) {
	p.Manifest = builtManifest(p.Manifest)
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

func builtManifest(manifest Manifest) Manifest {
	if BuildVersion != "" {
		manifest.Version = BuildVersion
	}
	return manifest
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
