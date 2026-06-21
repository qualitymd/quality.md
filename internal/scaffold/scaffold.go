// Package scaffold provides the starter QUALITY.md content used by `qualitymd init`.
package scaffold

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
)

//go:embed skeleton.md
var skeleton []byte

// Bytes returns a copy of the raw QUALITY.md scaffold. The returned slice is the
// caller's to mutate; it does not alias the shared embedded data.
func Bytes() []byte {
	out := make([]byte, len(skeleton))
	copy(out, skeleton)
	return out
}

// Create writes the scaffold to path, refusing to overwrite unless force is set.
func Create(path string, force bool) error {
	flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	if force {
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}

	file, err := os.OpenFile(path, flags, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%s already exists; pass --force to overwrite", path)
		}
		return fmt.Errorf("creating %s: %w", path, err)
	}

	if _, err := file.Write(skeleton); err != nil {
		_ = file.Close()
		return fmt.Errorf("writing %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("closing %s: %w", path, err)
	}
	return nil
}
