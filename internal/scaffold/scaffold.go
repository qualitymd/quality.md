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

//go:embed skeleton-minimal.md
var skeletonMinimal []byte

// Bytes returns a copy of the raw QUALITY.md scaffold. The returned slice is the
// caller's to mutate; it does not alias the shared embedded data.
func Bytes() []byte {
	return clone(skeleton)
}

// MinimalBytes returns a copy of the minimal QUALITY.md scaffold: a valid
// frontmatter skeleton without the guided template prose. The returned slice is
// the caller's to mutate; it does not alias the shared embedded data.
func MinimalBytes() []byte {
	return clone(skeletonMinimal)
}

func clone(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

// Create writes the scaffold to path, refusing to overwrite unless force is set.
// When minimal is set it writes the minimal skeleton instead of the guided one.
func Create(path string, force, minimal bool) error {
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

	content := skeleton
	if minimal {
		content = skeletonMinimal
	}
	if _, err := file.Write(content); err != nil {
		_ = file.Close()
		return fmt.Errorf("writing %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("closing %s: %w", path, err)
	}
	return nil
}
