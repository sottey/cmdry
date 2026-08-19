package plugins

import "encoding/json"

const ProtocolVersion = 1

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
type Response struct {
	OK    bool         `json:"ok"`
	Data  *View        `json:"data,omitempty"`
	Error *PluginError `json:"error,omitempty"`
}
type PluginError struct {
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

func (r Response) JSON() ([]byte, error) { return json.Marshal(r) }
