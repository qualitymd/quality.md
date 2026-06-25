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
	if got, want := model.RatingScale[0].Title, "🟢 Outstanding"; got != want {
		t.Fatalf("RatingScale[0].Title = %q, want %q", got, want)
	}
	if got, want := model.RatingScale[3].Level, "unacceptable"; got != want {
		t.Fatalf("RatingScale[3].Level = %q, want %q", got, want)
	}
	if got, want := model.RatingScale[3].Title, "🔴 Unacceptable"; got != want {
		t.Fatalf("RatingScale[3].Title = %q, want %q", got, want)
	}
	factor := model.Factors["quality-name"]
	if factor.Description == "" {
		t.Fatal("skeleton factor description is empty")
	}
	requirement := factor.Requirements["requirement-name"]
	if requirement.Title == "" {
		t.Fatal("skeleton requirement title is empty")
	}
	if requirement.Assessment == "" {
		t.Fatal("skeleton requirement assessment is empty")
	}

	raw := string(Bytes())
	for _, heading := range []string{"# <the system, component, or artifact this model is about>", "## Overview", "## Scope", "## Needs", "## Risks"} {
		if !strings.Contains(raw, heading) {
			t.Fatalf("skeleton is missing %q", heading)
		}
	}
	if strings.Contains(raw, "## Known gaps") {
		t.Fatal("skeleton still has a standalone \"## Known gaps\" section")
	}
	for _, marker := range []string{"agent-accessible", "*Unknowns*", "*Open questions*", "*Reviewed —"} {
		if !strings.Contains(raw, marker) {
			t.Fatalf("skeleton is missing the %q convention marker", marker)
		}
	}
}

func TestCreateRefusesExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte("keep me"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	err := Create(path, false, false)
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

	if err := Create(path, true, false); err != nil {
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

func TestMinimalSkeletonConformsToCurrentSpec(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, MinimalBytes(), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	model, err := lint.Load(path)
	if err != nil {
		t.Fatalf("lint.Load() error = %v", err)
	}
	if model.Title == "" {
		t.Fatal("minimal skeleton title is empty")
	}
	if len(model.RatingScale) != 4 {
		t.Fatalf("len(RatingScale) = %d, want 4", len(model.RatingScale))
	}
	if got, want := model.RatingScale[1].Title, "🔵 Target"; got != want {
		t.Fatalf("RatingScale[1].Title = %q, want %q", got, want)
	}
	if got, want := model.RatingScale[2].Title, "🟡 Minimum"; got != want {
		t.Fatalf("RatingScale[2].Title = %q, want %q", got, want)
	}
	factor := model.Factors["quality-name"]
	if factor.Description == "" {
		t.Fatal("minimal skeleton factor description is empty")
	}
	requirement := factor.Requirements["requirement-name"]
	if requirement.Title == "" {
		t.Fatal("minimal skeleton requirement title is empty")
	}
	if requirement.Assessment == "" {
		t.Fatal("minimal skeleton requirement assessment is empty")
	}
}

func TestMinimalSkeletonOmitsGuidedBodySections(t *testing.T) {
	raw := string(MinimalBytes())
	if !strings.Contains(raw, "# <the system, component, or artifact this model is about>") {
		t.Fatal("minimal skeleton is missing the root heading")
	}
	for _, section := range []string{"## Overview", "## Scope", "## Needs", "## Risks"} {
		if strings.Contains(raw, section) {
			t.Fatalf("minimal skeleton should omit the guided %q section", section)
		}
	}
}

func TestCreateMinimalWritesMinimalSkeleton(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := Create(path, false, true); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != string(MinimalBytes()) {
		t.Fatal("Create(minimal) did not write the minimal skeleton")
	}
}
