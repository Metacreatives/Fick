package frontend

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type ArtifactStore struct {
	root    string
	version string
	mu      sync.Mutex
}

func NewArtifactStore(root, version string) *ArtifactStore {
	return &ArtifactStore{
		root:    root,
		version: version,
	}
}

func (s *ArtifactStore) Read(route string) ([]byte, bool, error) {
	data, err := os.ReadFile(s.filePath(route))

	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return data, true, nil
}

func (s *ArtifactStore) Write(route string, data string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := s.filePath(route)

	if _, err := os.Stat(filePath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	directory := filepath.Dir(filePath)

	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(
		directory,
		".artifact-*.tmp",
	)
	if err != nil {
		return err
	}

	tempPath := tempFile.Name()

	defer os.Remove(tempPath)

	if _, err := tempFile.Write([]byte(data)); err != nil {
		_ = tempFile.Close()
		return err
	}

	if err := tempFile.Chmod(0o644); err != nil {
		_ = tempFile.Close()
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	return os.Rename(tempPath, filePath)
}

func (s *ArtifactStore) filePath(route string) string {
	sum := sha256.Sum256([]byte(route))

	return filepath.Join(
		s.root,
		s.version,
		hex.EncodeToString(sum[:])+".html",
	)
}
