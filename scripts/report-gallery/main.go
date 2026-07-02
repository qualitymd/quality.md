// Command report-gallery regenerates the checked-in LedgerLite example:
// the exemplar QUALITY.md, the gallery README, the synthetic quality
// changelog, and a single evaluation run built from synthetic payloads
// through the real evaluation pipeline.
package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/workspace"
)

//go:embed content
var contentFS embed.FS

const galleryRel = "examples/report-gallery/software-service"
const galleryCreatedAt = "2026-06-29T12:00:00Z"
const galleryRunID = "20260629T120000Z-0123456789ab"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "report-gallery: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	repoRoot, err := workspace.FindRepoRoot("")
	if err != nil {
		return err
	}
	exampleDir := filepath.Join(repoRoot, filepath.FromSlash(galleryRel))
	if err := os.MkdirAll(exampleDir, 0o755); err != nil {
		return fmt.Errorf("creating example directory: %w", err)
	}
	if err := writeEmbedded("content/README.md", filepath.Join(exampleDir, "README.md")); err != nil {
		return err
	}
	if err := writeEmbedded("content/QUALITY.md", filepath.Join(exampleDir, "QUALITY.md")); err != nil {
		return err
	}
	if err := writeChangelog(exampleDir); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(exampleDir, ".quality", "evaluations")); err != nil {
		return fmt.Errorf("removing generated evaluations: %w", err)
	}

	modelRel := filepath.ToSlash(filepath.Join(galleryRel, "QUALITY.md"))
	created, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repoRoot, Model: modelRel})
	if err != nil {
		return fmt.Errorf("creating gallery run: %w", err)
	}
	runPath := filepath.Join(exampleDir, filepath.FromSlash(created.Path))
	if err := pinRunIdentity(runPath, galleryCreatedAt, galleryRunID); err != nil {
		return err
	}
	payloads, err := json.MarshalIndent(galleryPayloads(), "", "  ")
	if err != nil {
		return fmt.Errorf("encoding payload batch: %w", err)
	}
	if _, err := evaluation.SetData(runPath, payloads, evaluation.DataSetOptions{DryRun: true}); err != nil {
		return fmt.Errorf("validating synthetic payloads: %w", err)
	}
	if _, err := evaluation.SetData(runPath, payloads, evaluation.DataSetOptions{}); err != nil {
		return fmt.Errorf("writing synthetic payloads: %w", err)
	}
	if _, err := evaluation.BuildReport(runPath); err != nil {
		return fmt.Errorf("building gallery reports: %w", err)
	}
	if err := trimGeneratedMarkdown(runPath); err != nil {
		return err
	}
	fmt.Printf("Generated %s\n", filepath.ToSlash(filepath.Join(galleryRel, created.Path, "report.md")))
	return nil
}

// writeChangelog rewrites the example's quality changelog from the embedded
// entries so the whole example stays generator-owned.
func writeChangelog(exampleDir string) error {
	changelogDir := filepath.Join(exampleDir, ".quality", "changelog")
	if err := os.RemoveAll(changelogDir); err != nil {
		return fmt.Errorf("removing generated changelog: %w", err)
	}
	entries, err := fs.ReadDir(contentFS, "content/changelog")
	if err != nil {
		return fmt.Errorf("reading embedded changelog: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := "content/changelog/" + entry.Name()
		if err := writeEmbedded(src, filepath.Join(changelogDir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func writeEmbedded(src, dst string) error {
	data, err := contentFS.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading embedded %s: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", filepath.ToSlash(filepath.Dir(dst)), err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", filepath.ToSlash(dst), err)
	}
	return nil
}

func trimGeneratedMarkdown(root string) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading generated markdown %s: %w", filepath.ToSlash(path), err)
		}
		trimmed := bytes.TrimRight(data, "\n")
		trimmed = append(trimmed, '\n')
		if bytes.Equal(data, trimmed) {
			return nil
		}
		if err := os.WriteFile(path, trimmed, 0o644); err != nil {
			return fmt.Errorf("writing generated markdown %s: %w", filepath.ToSlash(path), err)
		}
		return nil
	})
}

func pinRunIdentity(runPath, createdAt, runID string) error {
	path := filepath.Join(runPath, "data", "evaluation-manifest.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading run manifest: %w", err)
	}
	var manifest evaluation.EvaluationManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return fmt.Errorf("decoding run manifest: %w", err)
	}
	manifest.EvaluationID = runID
	manifest.CreatedAt = createdAt
	pinned, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding run manifest: %w", err)
	}
	if err := os.WriteFile(path, append(pinned, '\n'), 0o644); err != nil {
		return fmt.Errorf("writing run manifest: %w", err)
	}
	return nil
}
