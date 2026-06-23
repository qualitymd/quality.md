package cli

import (
	"bytes"
	"encoding/json"
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

func TestInitMinimalWritesMinimalSkeleton(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"init", "--minimal", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(got) != string(scaffold.MinimalBytes()) {
		t.Fatal("QUALITY.md did not contain the minimal skeleton")
	}
	if string(got) == string(scaffold.Bytes()) {
		t.Fatal("--minimal wrote the guided scaffold")
	}
	if !strings.Contains(stderr.String(), "Created "+path) {
		t.Fatalf("stderr = %q, want created path", stderr.String())
	}
}

func TestInitMinimalStdoutPassthrough(t *testing.T) {
	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"init", "--minimal", "-"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if out.String() != string(scaffold.MinimalBytes()) {
		t.Fatal("stdout did not contain the minimal skeleton")
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
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

func TestInitJSONReceipt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom.md")
	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"init", "--json", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var receipt InitReceipt
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; stdout = %s", err, out.String())
	}
	if receipt.SchemaVersion != 1 {
		t.Fatalf("schemaVersion = %d, want 1", receipt.SchemaVersion)
	}
	if receipt.Path != path {
		t.Fatalf("path = %q, want %q", receipt.Path, path)
	}
	if !receipt.Created {
		t.Fatal("created = false, want true")
	}
	if len(receipt.NextActions) != 1 || receipt.NextActions[0].Command != "qualitymd lint "+path {
		t.Fatalf("nextActions = %#v, want lint action", receipt.NextActions)
	}
}

func TestInitJSONForceReportsNotCreated(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte("keep me"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"init", "--json", "--force", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var receipt InitReceipt
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if receipt.Created {
		t.Fatal("created = true, want false for overwrite")
	}
}

func TestInitJSONRejectsStdoutPassthrough(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"init", "--json", "-"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if got := codeFor(err); got != ExitUsage {
		t.Fatalf("codeFor(error) = %d, want %d", got, ExitUsage)
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
	if got := codeFor(err); got != ExitInternal {
		t.Fatalf("codeFor(error) = %d, want %d", got, ExitInternal)
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

func TestInitJSONOverwriteRefusalWritesErrorObject(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte("keep me"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"init", "--json", path})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want existing-file error")
	}
	if got := codeFor(err); got != ExitInternal {
		t.Fatalf("codeFor(error) = %d, want %d", got, ExitInternal)
	}
	if out.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", out.String())
	}
	var body struct {
		SchemaVersion int    `json:"schemaVersion"`
		Path          string `json:"path"`
		Reason        string `json:"reason"`
	}
	if err := json.Unmarshal(stderr.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; stderr = %s", err, stderr.String())
	}
	if body.SchemaVersion != 1 || body.Path != path || !strings.Contains(body.Reason, "already exists") {
		t.Fatalf("error body = %#v, want path and reason", body)
	}
}
