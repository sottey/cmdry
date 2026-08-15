// Package groupstate persists collapsed sidebar categories.
package groupstate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const filename = "plugin-group-state.json"

type Store struct{ DataDir string }

func (s Store) Load() (map[string]bool, error) {
	contents, err := os.ReadFile(filepath.Join(s.DataDir, filename))
	if os.IsNotExist(err) {
		return map[string]bool{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read plugin group state: %w", err)
	}
	state := map[string]bool{}
	if err := json.Unmarshal(contents, &state); err != nil {
		return nil, fmt.Errorf("decode plugin group state: %w", err)
	}
	return state, nil
}

func (s Store) Save(state map[string]bool) error {
	contents, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("encode plugin group state: %w", err)
	}
	temporary, err := os.CreateTemp(s.DataDir, ".plugin-group-state-*")
	if err != nil {
		return fmt.Errorf("create plugin group state: %w", err)
	}
	name := temporary.Name()
	defer os.Remove(name)
	if err := temporary.Chmod(0600); err != nil {
		temporary.Close()
		return fmt.Errorf("protect plugin group state: %w", err)
	}
	if _, err := temporary.Write(contents); err != nil {
		temporary.Close()
		return fmt.Errorf("write plugin group state: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close plugin group state: %w", err)
	}
	if err := os.Rename(name, filepath.Join(s.DataDir, filename)); err != nil {
		return fmt.Errorf("replace plugin group state: %w", err)
	}
	return nil
}
