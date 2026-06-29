// Package agentinstructions updates project agent instruction files with a
// concise pointer to the local QUALITY.md model.
package agentinstructions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	marker       = "<!-- Added by qualitymd init. -->"
	pointerTrail = " for this project's quality model."
)

var candidateFiles = []candidate{
	{Path: "AGENTS.md", Create: true},
	{Path: "CLAUDE.md"},
	{Path: "GEMINI.md"},
}

type candidate struct {
	Path   string
	Create bool
}

// FileResult describes an agent instruction file touched by the updater.
type FileResult struct {
	Path    string `json:"path"`
	Created bool   `json:"created"`
	Updated bool   `json:"updated"`
}

// UpdateOptions configures an agent instruction pointer update.
type UpdateOptions struct {
	// Dir is the directory containing the agent instruction files. Empty means
	// the current working directory.
	Dir string
	// ModelPath is the initialized QUALITY.md path. Empty means QUALITY.md in Dir.
	ModelPath string
}

type resolvedOptions struct {
	dir       string
	modelPath string
}

// Update adds the quality model pointer to eligible agent instruction files and
// returns the files that were created or changed.
func Update(opts UpdateOptions) ([]FileResult, error) {
	resolved, err := resolveOptions(opts)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	var results []FileResult
	for _, candidate := range candidateFiles {
		result, changed, err := updateCandidate(resolved, candidate, seen)
		if err != nil {
			return nil, err
		}
		if changed {
			results = append(results, result)
		}
	}
	return results, nil
}

func resolveOptions(opts UpdateOptions) (resolvedOptions, error) {
	dir := opts.Dir
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return resolvedOptions{}, fmt.Errorf("resolving working directory: %w", err)
		}
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return resolvedOptions{}, fmt.Errorf("resolving %s: %w", dir, err)
	}
	modelPath := opts.ModelPath
	if modelPath == "" {
		modelPath = "QUALITY.md"
	}
	absModelPath := modelPath
	if !filepath.IsAbs(absModelPath) {
		absModelPath = filepath.Join(absDir, modelPath)
	}
	absModelPath, err = filepath.Abs(absModelPath)
	if err != nil {
		return resolvedOptions{}, fmt.Errorf("resolving %s: %w", modelPath, err)
	}
	return resolvedOptions{dir: absDir, modelPath: absModelPath}, nil
}

func updateCandidate(opts resolvedOptions, candidate candidate, seen map[string]bool) (FileResult, bool, error) {
	target := filepath.Join(opts.dir, candidate.Path)
	stat, err := os.Stat(target)
	if err != nil {
		if !os.IsNotExist(err) {
			return FileResult{}, false, fmt.Errorf("checking %s: %w", candidate.Path, err)
		}
		if !candidate.Create {
			return FileResult{}, false, nil
		}
	}

	key := targetKey(target, stat)
	if seen[key] {
		return FileResult{}, false, nil
	}
	seen[key] = true

	result, changed, err := updateFile(target, opts.modelPath, stat == nil)
	if err != nil || !changed {
		return result, changed, err
	}
	result.Path = candidate.Path
	return result, true, nil
}

func targetKey(path string, stat os.FileInfo) string {
	if stat != nil {
		if resolved, err := filepath.EvalSymlinks(path); err == nil {
			if abs, err := filepath.Abs(resolved); err == nil {
				return abs
			}
			return resolved
		}
	}
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}
	return filepath.Clean(path)
}

func updateFile(path, modelPath string, create bool) (FileResult, bool, error) {
	var existing []byte
	if !create {
		var err error
		existing, err = os.ReadFile(path)
		if err != nil {
			return FileResult{}, false, fmt.Errorf("reading %s: %w", path, err)
		}
		if hasPointer(string(existing)) {
			return FileResult{}, false, nil
		}
	}

	block, err := pointerBlock(path, modelPath)
	if err != nil {
		return FileResult{}, false, err
	}
	content := appendBlock(string(existing), block)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return FileResult{}, false, fmt.Errorf("writing %s: %w", path, err)
	}
	return FileResult{Created: create, Updated: !create}, true, nil
}

func hasPointer(content string) bool {
	return strings.Contains(content, marker) ||
		(strings.Contains(content, "See [QUALITY.md](") && strings.Contains(content, pointerTrail))
}

func pointerBlock(instructionPath, modelPath string) (string, error) {
	rel, err := filepath.Rel(filepath.Dir(instructionPath), modelPath)
	if err != nil {
		return "", fmt.Errorf("rendering QUALITY.md link for %s: %w", instructionPath, err)
	}
	rel = filepath.ToSlash(rel)
	return marker + "\nSee [QUALITY.md](" + rel + ")" + pointerTrail, nil
}

func appendBlock(existing, block string) string {
	if existing == "" {
		return block + "\n"
	}
	if strings.HasSuffix(existing, "\n\n") {
		return existing + block + "\n"
	}
	if strings.HasSuffix(existing, "\n") {
		return existing + "\n" + block + "\n"
	}
	return existing + "\n\n" + block + "\n"
}
