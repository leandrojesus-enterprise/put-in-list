// Package storage handles persisting and loading the user's package lists.
package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// InstallEntry records a single package that was installed via put-in-list.
type InstallEntry struct {
	Name      string `json:"name"`      // Package name as passed to the installer.
	Installer string `json:"installer"` // The installer used (e.g. "apt", "winget").
}

// StoredLists is the top-level structure persisted to lists.json.
type StoredLists struct {
	ActiveList string                    `json:"activeList"` // Name of the currently active list.
	Lists      map[string][]InstallEntry `json:"lists"`      // All named lists and their entries.
}

// JSONStore reads and writes StoredLists to a JSON file at a given path.
type JSONStore struct {
	path string
}

// NewJSONStore creates a JSONStore that operates on the file at path.
func NewJSONStore(path string) *JSONStore {
	return &JSONStore{path: path}
}

// Load reads the lists file and deserialises it into StoredLists.
// The second return value is true when the file does not exist (first run),
// in which case a default StoredLists and nil error are returned.
func (s *JSONStore) Load() (StoredLists, bool, error) {
	var result StoredLists = StoredLists{
		ActiveList: "None",
		Lists:      make(map[string][]InstallEntry),
	}

	var f *os.File
	var err error
	f, err = os.Open(s.path)
	if err != nil {
		return result, true, nil // file absent → treat as first run
	}
	defer f.Close()

	var data []byte
	data, err = io.ReadAll(f)
	if err != nil {
		return result, false, fmt.Errorf("error reading %s: %w", s.path, err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, false, err
	}

	return result, false, nil
}

// Save serialises data as JSON and writes it to the lists file.
func (s *JSONStore) Save(data StoredLists) error {
	var b []byte
	var err error
	b, err = json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0644)
}
