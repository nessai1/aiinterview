package prompt

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Storage struct {
	lang         string
	promptDir    fs.FS
	placeholders map[string]string
}

func NewStorage(lang string) (*Storage, error) {
	if lang == "" {
		return nil, fmt.Errorf("prompt storage language must be non-zero string")
	}

	dir := filepath.Join("prompts", lang)
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot get stats for prompt dir '%s': %w", dir, err)
	}

	if !dirInfo.IsDir() {
		return nil, fmt.Errorf("prompt dir '%s' is not a directory", dir)
	}

	return &Storage{lang: lang, promptDir: os.DirFS(dir)}, nil
}

func (s *Storage) LoadPrompt(promptID string, placeholders ...string) (string, error) {
	// TODO: check cache/load data, insert placeholders
}
