package plugins

import "sort"

type Status string

const (
	StatusEnabled  Status = "enabled"
	StatusDisabled Status = "disabled"
	StatusInvalid  Status = "invalid"
	StatusError    Status = "error"
)

type Registered struct {
	Manifest  Manifest
	Path      string
	Status    Status
	LastError string
}
type Registry struct{ entries map[string]Registered }

func NewRegistry() *Registry { return &Registry{entries: map[string]Registered{}} }
func (r *Registry) Add(entry Registered) bool {
	if _, ok := r.entries[entry.Manifest.ID]; ok {
		return false
	}
	r.entries[entry.Manifest.ID] = entry
	return true
}
func (r *Registry) Get(id string) (Registered, bool) { x, ok := r.entries[id]; return x, ok }
func (r *Registry) All() []Registered {
	out := make([]Registered, 0, len(r.entries))
	for _, x := range r.entries {
		out = append(out, x)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Manifest.Name < out[j].Manifest.Name })
	return out
}
func (r *Registry) Len() int { return len(r.entries) }
