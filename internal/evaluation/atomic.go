package evaluation

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func writeNumbered(dir, suffix string, data []byte) (string, error) {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		n, err := nextRecordNumber(dir)
		if err != nil {
			return "", err
		}
		path := filepath.Join(dir, fmt.Sprintf("%03d-%s", n, suffix))
		err = writeCreate(path, data)
		if err == nil {
			return path, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return "", err
		}
		lastErr = err
	}
	return "", fmt.Errorf("numbering collision in %s: %w", dir, lastErr)
}

func writeCreate(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		return err
	}
	return os.Link(tmpName, path)
}

func writeReplace(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
