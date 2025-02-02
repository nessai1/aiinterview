package prompt

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Storage struct {
	lang      string
	promptDir fs.FS
	prompts   map[string]string
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

func (s *Storage) LoadPrompt(promptID string, placeholders map[string]string) (string, error) {
	var prompt string
	var err error

	prompt, ok := s.prompts[promptID]
	if !ok {
		prompt, err = s.loadPromptFromDisk(promptID)
		if err != nil {
			return "", fmt.Errorf("cannot load prompt '%s' from disk: %w", prompt, err)
		}
	}

	s.prompts[promptID] = prompt
	result := prompt

	for placeholder, value := range placeholders {
		result = strings.ReplaceAll(result, "{{"+placeholder+"}}", value)
	}

	return result, nil
}

func (s *Storage) loadPromptFromDisk(promptID string) (string, error) {
	filename := promptID + ".txt"
	fd, err := s.promptDir.Open(filename)
	if err != nil {
		return "", fmt.Errorf("cannot open file '%s': %w", filename, err)
	}

	bf := bytes.Buffer{}
	_, err = bf.ReadFrom(fd)
	if err != nil {
		return "", fmt.Errorf("cannot read file '%s': %w", filename, err)
	}

	return bf.String(), nil
}
