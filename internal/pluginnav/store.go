// Package pluginnav persists local favorites and recently used plugins.
package pluginnav

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const filename = "plugin-navigation.json"

const DefaultRecentLimit = 8

type State struct {
	Favorites     []string `json:"favorites"`
	Recents       []string `json:"recents"`
	Hidden        []string `json:"hidden"`
	RecentLimit   int      `json:"recent_limit"`
	ShowFavorites bool     `json:"show_favorites"`
	ShowRecents   bool     `json:"show_recents"`
}

type Store struct{ DataDir string }

func (s Store) Load() (State, error) {
	contents, err := os.ReadFile(filepath.Join(s.DataDir, filename))
	if os.IsNotExist(err) {
		return State{RecentLimit: DefaultRecentLimit, ShowFavorites: true, ShowRecents: true}, nil
	}
	if err != nil {
		return State{}, fmt.Errorf("read plugin navigation: %w", err)
	}
	var state State
	if err := json.Unmarshal(contents, &state); err != nil {
		return State{}, fmt.Errorf("decode plugin navigation: %w", err)
	}
	if state.RecentLimit <= 0 {
		state.RecentLimit = DefaultRecentLimit
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(contents, &keys); err != nil {
		return State{}, fmt.Errorf("decode plugin navigation fields: %w", err)
	}
	if _, ok := keys["show_favorites"]; !ok {
		state.ShowFavorites = true
	}
	if _, ok := keys["show_recents"]; !ok {
		state.ShowRecents = true
	}
	return state, nil
}

func (s Store) Save(state State) error {
	contents, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("encode plugin navigation: %w", err)
	}
	temporary, err := os.CreateTemp(s.DataDir, ".plugin-navigation-*")
	if err != nil {
		return fmt.Errorf("create plugin navigation: %w", err)
	}
	name := temporary.Name()
	defer os.Remove(name)
	if err := temporary.Chmod(0600); err != nil {
		temporary.Close()
		return fmt.Errorf("protect plugin navigation: %w", err)
	}
	if _, err := temporary.Write(contents); err != nil {
		temporary.Close()
		return fmt.Errorf("write plugin navigation: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close plugin navigation: %w", err)
	}
	if err := os.Rename(name, filepath.Join(s.DataDir, filename)); err != nil {
		return fmt.Errorf("replace plugin navigation: %w", err)
	}
	return nil
}
