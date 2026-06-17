package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/scaffold"
)

func TestInitWritesDefaultQualityFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("os.Chdir(%q) error = %v", cwd, err)
		}
	})
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}

	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"init"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got, err := os.ReadFile("QUALITY.md")
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != string(scaffold.Bytes()) {
		t.Fatal("QUALITY.md did not contain the scaffold")
	}
	if out.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", out.String())
	}
	if !strings.Contains(stderr.String(), "Created QUALITY.md") {
		t.Fatalf("stderr = %q, want created path", stderr.String())
	}
	if !strings.Contains(stderr.String(), "qualitymd lint QUALITY.md") {
		t.Fatalf("stderr = %q, want next action", stderr.String())
	}
}

func TestInitWritesCustomPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom.md")
	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"init", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != string(scaffold.Bytes()) {
		t.Fatal("custom path did not contain the scaffold")
	}
	if out.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", out.String())
	}
	if !strings.Contains(stderr.String(), "Created "+path) {
		t.Fatalf("stderr = %q, want custom created path", stderr.String())
	}
}

func TestInitStdoutDoesNotTouchExistingFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("os.Chdir(%q) error = %v", cwd, err)
		}
	})
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "QUALITY.md"), []byte("keep me"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}

	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"init", "-"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if out.String() != string(scaffold.Bytes()) {
		t.Fatal("stdout did not contain the scaffold")
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	got, err := os.ReadFile("QUALITY.md")
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != "keep me" {
		t.Fatalf("existing file changed to %q", got)
	}
}

func TestInitRefusesExistingFileUnlessForced(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte("keep me"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"init", path})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want existing-file error")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Execute() error = %v, want already exists", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != "keep me" {
		t.Fatalf("existing file changed to %q", got)
	}

	cmd = newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"init", "--force", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("forced Execute() error = %v", err)
	}
	got, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != string(scaffold.Bytes()) {
		t.Fatal("--force did not write the scaffold")
	}
}
