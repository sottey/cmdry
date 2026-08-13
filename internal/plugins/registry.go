package plugins

import (
	"sort"
	"sync"
)

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
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Registered
}

func NewRegistry() *Registry { return &Registry{entries: map[string]Registered{}} }
func (r *Registry) Add(entry Registered) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[entry.Manifest.ID]; ok {
		return false
	}
	r.entries[entry.Manifest.ID] = entry
	return true
}
func (r *Registry) Get(id string) (Registered, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	x, ok := r.entries[id]
	return x, ok
}
func (r *Registry) All() []Registered {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Registered, 0, len(r.entries))
	for _, x := range r.entries {
		out = append(out, x)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Manifest.Name < out[j].Manifest.Name })
	return out
}
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

// Replace atomically swaps the registry contents after a complete plugin scan.
func (r *Registry) Replace(entries []Registered) {
	next := make(map[string]Registered, len(entries))
	for _, entry := range entries {
		next[entry.Manifest.ID] = entry
	}
	r.mu.Lock()
	r.entries = next
	r.mu.Unlock()
}
