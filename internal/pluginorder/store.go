// Package pluginorder persists a user-selected plugin navigation order.
package pluginorder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const filename = "plugin-order.json"

type Store struct{ DataDir string }

func (s Store) Load() ([]string, error) {
	contents, err := os.ReadFile(filepath.Join(s.DataDir, filename))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read plugin order: %w", err)
	}
	var ids []string
	if err := json.Unmarshal(contents, &ids); err != nil {
		return nil, fmt.Errorf("decode plugin order: %w", err)
	}
	return ids, nil
}

func (s Store) Save(ids []string) error {
	contents, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("encode plugin order: %w", err)
	}
	temporary, err := os.CreateTemp(s.DataDir, ".plugin-order-*")
	if err != nil {
		return fmt.Errorf("create plugin order: %w", err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if err := temporary.Chmod(0600); err != nil {
		temporary.Close()
		return fmt.Errorf("protect plugin order: %w", err)
	}
	if _, err := temporary.Write(contents); err != nil {
		temporary.Close()
		return fmt.Errorf("write plugin order: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close plugin order: %w", err)
	}
	if err := os.Rename(temporaryName, filepath.Join(s.DataDir, filename)); err != nil {
		return fmt.Errorf("replace plugin order: %w", err)
	}
	return nil
}
