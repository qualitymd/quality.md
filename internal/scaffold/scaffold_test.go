package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/lint"
)

func TestSkeletonConformsToCurrentSpec(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, Bytes(), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	model, err := lint.Load(path)
	if err != nil {
		t.Fatalf("lint.Load() error = %v", err)
	}
	if model.Title == "" {
		t.Fatal("skeleton title is empty")
	}
	if len(model.RatingScale) != 4 {
		t.Fatalf("len(RatingScale) = %d, want 4", len(model.RatingScale))
	}
	if got, want := model.RatingScale[0].Level, "outstanding"; got != want {
		t.Fatalf("RatingScale[0].Level = %q, want %q", got, want)
	}
	if got, want := model.RatingScale[3].Level, "unacceptable"; got != want {
		t.Fatalf("RatingScale[3].Level = %q, want %q", got, want)
	}
	factor := model.Factors["<name a quality that matters>"]
	if factor.Description == "" {
		t.Fatal("skeleton factor description is empty")
	}
	requirement := factor.Requirements["<state one expectation you can assess>"]
	if requirement.Assessment == "" {
		t.Fatal("skeleton requirement assessment is empty")
	}

	raw := string(Bytes())
	for _, heading := range []string{"# <the system, component, or artifact this model is about>", "## Overview", "## Scope", "## Needs", "## Risks", "## Known gaps"} {
		if !strings.Contains(raw, heading) {
			t.Fatalf("skeleton is missing %q", heading)
		}
	}
}

func TestCreateRefusesExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte("keep me"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	err := Create(path, false)
	if err == nil {
		t.Fatal("Create() error = nil, want existing-file error")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Create() error = %v, want already exists", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != "keep me" {
		t.Fatalf("existing file changed to %q", got)
	}
}

func TestCreateForceOverwritesExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte("replace me"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	if err := Create(path, true); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != string(Bytes()) {
		t.Fatal("forced create did not write the scaffold")
	}
}
