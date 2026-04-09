package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type InstallEntry struct {
	Name      string `json:"name"`
	Installer string `json:"installer"`
}

type StoredLists struct {
	ActiveList string                    `json:"activeList"`
	Lists      map[string][]InstallEntry `json:"lists"`
}

type JSONStore struct {
	path string
}

func NewJSONStore(path string) *JSONStore {
	return &JSONStore{path: path}
}

func (s *JSONStore) Load() (StoredLists, bool, error) {
	var result StoredLists = StoredLists{
		ActiveList: "None",
		Lists:      make(map[string][]InstallEntry),
	}

	var f *os.File
	var err error
	f, err = os.Open(s.path)
	if err != nil {
		return result, true, nil // first run
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

func (s *JSONStore) Save(data StoredLists) error {
	var b []byte
	var err error
	b, err = json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0644)
}
