package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/leandrojesus-enterprise/put-in-list/internal/i18n"
)

type Config struct {
	CurrentInstaller string   `json:"currentInstaller"`
	Language         i18n.Lang `json:"language"`
}

type JSONService struct {
	path string
	tr   *i18n.Translator
}

func NewJSONService(path string) *JSONService {
	return &JSONService{path: path}
}

func (s *JSONService) SetTranslator(tr *i18n.Translator) {
	s.tr = tr
}

func (s *JSONService) Load() (Config, bool, error) {
	cfg := Config{CurrentInstaller: "None"}

	f, err := os.Open(s.path)
	if err != nil {
		return cfg, true, nil // first run
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return cfg, false, fmt.Errorf("%s %w", s.tr.Trans("config_read_error"), err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, false, err
	}

	return cfg, false, nil
}

func (s *JSONService) Save(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("%s %w", s.tr.Trans("save_config_error_marshal"), err)
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("%s %w", s.tr.Trans("save_config_error_write"), err)
	}
	return nil
}
