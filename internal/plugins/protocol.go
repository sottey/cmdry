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
}
type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

func (r Response) JSON() ([]byte, error) { return json.Marshal(r) }
