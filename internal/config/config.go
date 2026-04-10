// Package config handles loading and persisting the application configuration.
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/leandrojesus-enterprise/put-in-list/internal/i18n"
)

// Config holds user-level settings that are persisted between sessions.
type Config struct {
	CurrentInstaller string    `json:"currentInstaller"` // The package manager selected by the user.
	Language         i18n.Lang `json:"language"`         // The display language (e.g. "en", "pt", "es").
}

// JSONService reads and writes Config values to a JSON file at a given path.
type JSONService struct {
	path string
	tr   *i18n.Translator
}

// NewJSONService creates a JSONService that operates on the file at path.
func NewJSONService(path string) *JSONService {
	return &JSONService{path: path}
}

// SetTranslator injects the translator used to localise error messages.
func (s *JSONService) SetTranslator(tr *i18n.Translator) {
	s.tr = tr
}

// Load reads the config file and deserialises it into a Config.
// The second return value is true when the file does not exist (first run),
// in which case a default Config and nil error are returned.
func (s *JSONService) Load() (Config, bool, error) {
	var cfg Config = Config{CurrentInstaller: "None"}

	var f *os.File
	var err error
	f, err = os.Open(s.path)
	if err != nil {
		return cfg, true, nil // file absent → treat as first run
	}
	defer f.Close()

	var data []byte
	data, err = io.ReadAll(f)
	if err != nil {
		return cfg, false, fmt.Errorf("%s %w", s.tr.Trans("config_read_error"), err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, false, err
	}

	return cfg, false, nil
}

// Save serialises cfg as indented JSON and writes it to the config file.
func (s *JSONService) Save(cfg Config) error {
	var data []byte
	var err error
	data, err = json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("%s %w", s.tr.Trans("save_config_error_marshal"), err)
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("%s %w", s.tr.Trans("save_config_error_write"), err)
	}
	return nil
}
