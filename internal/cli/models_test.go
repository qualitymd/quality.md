package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/models"
)

func TestModelsListHuman(t *testing.T) {
	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"models", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), "NAME") || !strings.Contains(out.String(), "quality-meta-model") {
		t.Fatalf("stdout = %q, want catalog table", out.String())
	}
	if hasTerminalEscape(out.String()) {
		t.Fatalf("stdout = %q, want plain table for non-terminal writer", out.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRenderModelsListStyledUsesSharedPalette(t *testing.T) {
	var out bytes.Buffer
	if err := renderModelsListStyled(&out, models.Catalog()); err != nil {
		t.Fatalf("renderModelsListStyled() error = %v", err)
	}
	got := out.String()
	for _, want := range []string{"NAME", "TITLE", "quality-meta-model"} {
		if !strings.Contains(got, want) {
			t.Fatalf("styled table = %q, want substring %q", got, want)
		}
	}
	if !hasTerminalEscape(got) {
		t.Fatalf("styled table = %q, want terminal styling", got)
	}
}

func TestModelsListJSON(t *testing.T) {
	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"models", "list", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var got []models.Entry
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; stdout = %s", err, out.String())
	}
	if len(got) != 1 || got[0].Name != "quality-meta-model" {
		t.Fatalf("decoded catalog = %#v, want quality-meta-model", got)
	}
	if hasTerminalEscape(out.String()) {
		t.Fatalf("stdout = %q, want JSON without terminal escapes", out.String())
	}
}

func TestModelsViewMarkdown(t *testing.T) {
	want, err := models.Markdown("quality-meta-model", "")
	if err != nil {
		t.Fatalf("models.Markdown() error = %v", err)
	}

	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"models", "view", "quality-meta-model"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !bytes.Equal(out.Bytes(), want) {
		t.Fatal("stdout did not match bundled model Markdown")
	}
}

func TestModelsViewSourceRewriteJSON(t *testing.T) {
	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"models", "view", "quality-meta-model", "--source", "QUALITY.md", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var got models.View
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; stdout = %s", err, out.String())
	}
	if got.Name != "quality-meta-model" || got.Model.Source != "QUALITY.md" {
		t.Fatalf("decoded view = %#v, want quality-meta-model with rewritten source", got)
	}
}

func TestModelsViewUnknownModelMapsToUsage(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"models", "view", "missing"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want unknown model error")
	}
	if got := codeFor(err); got != ExitUsage {
		t.Fatalf("codeFor(error) = %d, want %d; err = %v", got, ExitUsage, err)
	}
}
