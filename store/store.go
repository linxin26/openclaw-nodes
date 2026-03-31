package store

import (
	"os"
	"path/filepath"
)

type Store struct {
	dataDir string
}

func New(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}
	return &Store{dataDir: dataDir}, nil
}

func (s *Store) Path(name string) string {
	return filepath.Join(s.dataDir, name)
}

// DefaultDataDir returns the platform-specific default data directory
func DefaultDataDir() (string, error) {
	// Use %APPDATA%\OpenClaw on Windows
	if dir := os.Getenv("APPDATA"); dir != "" {
		return filepath.Join(dir, "OpenClaw"), nil
	}
	// Fallback to current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "openclaw-data"), nil
}
