package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Store owns atomic reads and writes of a run's evaluation.json. A single
// store serializes merges, so concurrent execution strategies can never
// interleave renames.
type Store struct {
	runAbs string
	mu     sync.Mutex
}

// NewStore returns the evaluation.json store for a run folder.
func NewStore(runAbs string) *Store {
	return &Store{runAbs: runAbs}
}

func (s *Store) path() string {
	return filepath.Join(s.runAbs, ArtifactFile)
}

// Exists reports whether the run folder has an evaluation.json artifact.
func (s *Store) Exists() bool {
	_, err := os.Stat(s.path())
	return err == nil
}

// Load reads and decodes the run artifact.
func (s *Store) Load() (*Artifact, error) {
	raw, err := os.ReadFile(s.path())
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", ArtifactFile, err)
	}
	var artifact Artifact
	if err := json.Unmarshal(raw, &artifact); err != nil {
		return nil, fmt.Errorf("decoding %s: %w", ArtifactFile, err)
	}
	return &artifact, nil
}

// Save writes the run artifact atomically through a temp-file plus rename
// sequence, so a crash never leaves a torn evaluation.json.
func (s *Store) Save(artifact *Artifact) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	raw, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding %s: %w", ArtifactFile, err)
	}
	raw = append(raw, '\n')
	tmp, err := os.CreateTemp(s.runAbs, ArtifactFile+".tmp-*")
	if err != nil {
		return fmt.Errorf("staging %s: %w", ArtifactFile, err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(raw); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("staging %s: %w", ArtifactFile, err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("staging %s: %w", ArtifactFile, err)
	}
	if err := os.Rename(tmpName, s.path()); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("writing %s: %w", ArtifactFile, err)
	}
	return nil
}
