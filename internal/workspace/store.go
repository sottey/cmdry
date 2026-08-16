// Package workspace persists local visual and navigation preferences.
package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const filename = "workspace-preferences.json"

const (
	ThemeSystem = "system"
	ThemeLight  = "light"
	ThemeDark   = "dark"

	DensityComfortable = "comfortable"
	DensityCompact     = "compact"

	LandingOverview = "overview"
	LandingLastUsed = "last-used"
)

type State struct {
	Theme          string `json:"theme"`
	ReducedMotion  bool   `json:"reduced_motion"`
	SidebarDensity string `json:"sidebar_density"`
	DefaultLanding string `json:"default_landing"`
	LastPluginID   string `json:"last_plugin_id,omitempty"`
}

func DefaultState() State {
	return State{Theme: ThemeSystem, SidebarDensity: DensityComfortable, DefaultLanding: LandingOverview}
}

func (s State) Valid() bool {
	return validTheme(s.Theme) && validDensity(s.SidebarDensity) && validLanding(s.DefaultLanding)
}

type Store struct{ DataDir string }

func (s Store) Load() (State, error) {
	contents, err := os.ReadFile(filepath.Join(s.DataDir, filename))
	if os.IsNotExist(err) {
		return DefaultState(), nil
	}
	if err != nil {
		return State{}, fmt.Errorf("read workspace preferences: %w", err)
	}
	state := DefaultState()
	if err := json.Unmarshal(contents, &state); err != nil {
		return State{}, fmt.Errorf("decode workspace preferences: %w", err)
	}
	if !state.Valid() {
		return State{}, fmt.Errorf("workspace preferences contain invalid values")
	}
	return state, nil
}

func (s Store) Save(state State) error {
	if !state.Valid() {
		return fmt.Errorf("workspace preferences contain invalid values")
	}
	contents, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("encode workspace preferences: %w", err)
	}
	temporary, err := os.CreateTemp(s.DataDir, ".workspace-preferences-*")
	if err != nil {
		return fmt.Errorf("create workspace preferences: %w", err)
	}
	name := temporary.Name()
	defer os.Remove(name)
	if err := temporary.Chmod(0600); err != nil {
		temporary.Close()
		return fmt.Errorf("protect workspace preferences: %w", err)
	}
	if _, err := temporary.Write(contents); err != nil {
		temporary.Close()
		return fmt.Errorf("write workspace preferences: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close workspace preferences: %w", err)
	}
	if err := os.Rename(name, filepath.Join(s.DataDir, filename)); err != nil {
		return fmt.Errorf("replace workspace preferences: %w", err)
	}
	return nil
}

func validTheme(value string) bool {
	return value == ThemeSystem || value == ThemeLight || value == ThemeDark
}
func validDensity(value string) bool { return value == DensityComfortable || value == DensityCompact }
func validLanding(value string) bool { return value == LandingOverview || value == LandingLastUsed }
