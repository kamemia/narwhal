package upload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Manager struct {
	BaseDir string
}

func NewManager(baseDir string) *Manager {
	return &Manager{BaseDir: baseDir}
}

func (m *Manager) Save(filename string, src io.Reader) (string, error) {
	rawDir := m.BaseDir
	if err := os.MkdirAll(rawDir, 0755); err != nil {
		return "", fmt.Errorf("create raw dir: %w", err)
	}

	destPath := filepath.Join(rawDir, filename)
	f, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, src); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return destPath, nil
}
