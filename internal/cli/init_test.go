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
	withTempCwd(t)

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
	agents, err := os.ReadFile("AGENTS.md")
	if err != nil {
		t.Fatalf("os.ReadFile(AGENTS.md) error = %v", err)
	}
	if !strings.Contains(string(agents), "See [QUALITY.md](QUALITY.md) for this project's quality model.") {
		t.Fatalf("AGENTS.md = %q, want quality pointer", agents)
	}
	if !strings.Contains(stderr.String(), "Agent instructions: AGENTS.md") {
		t.Fatalf("stderr = %q, want agent instruction update", stderr.String())
	}
}

func TestInitNoAgentInstructionsOptOut(t *testing.T) {
	withTempCwd(t)

	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"init", "--no-agent-instructions"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if _, err := os.Stat("AGENTS.md"); !os.IsNotExist(err) {
		t.Fatalf("os.Stat(AGENTS.md) error = %v, want not exist", err)
	}
	if strings.Contains(stderr.String(), "Agent instructions:") {
		t.Fatalf("stderr = %q, want no agent instruction line", stderr.String())
	}
}

func TestInitMinimalWritesMinimalSkeleton(t *testing.T) {
	withTempCwd(t)
	path := "QUALITY.md"
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
	withTempCwd(t)
	if err := os.Mkdir("docs", 0o755); err != nil {
		t.Fatalf("os.Mkdir() error = %v", err)
	}
	path := filepath.Join("docs", "custom.md")
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
	agents, err := os.ReadFile("AGENTS.md")
	if err != nil {
		t.Fatalf("os.ReadFile(AGENTS.md) error = %v", err)
	}
	if !strings.Contains(string(agents), "See [QUALITY.md](docs/custom.md) for this project's quality model.") {
		t.Fatalf("AGENTS.md = %q, want custom relative pointer", agents)
	}
}

func TestInitJSONReceipt(t *testing.T) {
	withTempCwd(t)
	path := "custom.md"
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
	if len(receipt.AgentInstructionFiles) != 1 {
		t.Fatalf("agentInstructionFiles = %#v, want one update", receipt.AgentInstructionFiles)
	}
	if got := receipt.AgentInstructionFiles[0]; got.Path != "AGENTS.md" || !got.Created || got.Updated {
		t.Fatalf("agentInstructionFiles[0] = %#v, want created AGENTS.md", got)
	}
}

func TestInitJSONForceReportsNotCreated(t *testing.T) {
	withTempCwd(t)
	path := "QUALITY.md"
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
	dir := withTempCwd(t)
	if err := os.WriteFile(filepath.Join(dir, "QUALITY.md"), []byte("keep me"), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
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
	if _, err := os.Stat("AGENTS.md"); !os.IsNotExist(err) {
		t.Fatalf("os.Stat(AGENTS.md) error = %v, want not exist", err)
	}
}

func TestInitRefusesExistingFileUnlessForced(t *testing.T) {
	withTempCwd(t)
	path := "QUALITY.md"
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

func TestInitForceDoesNotDuplicateAgentInstructionPointer(t *testing.T) {
	withTempCwd(t)
	for i := 0; i < 2; i++ {
		cmd := newRootCmd()
		cmd.SetOut(&bytes.Buffer{})
		cmd.SetErr(&bytes.Buffer{})
		cmd.SetArgs([]string{"init", "--force"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() run %d error = %v", i+1, err)
		}
	}
	got, err := os.ReadFile("AGENTS.md")
	if err != nil {
		t.Fatalf("os.ReadFile(AGENTS.md) error = %v", err)
	}
	if count := strings.Count(string(got), "<!-- Added by qualitymd init. -->"); count != 1 {
		t.Fatalf("marker count = %d, want 1 in %q", count, got)
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

func withTempCwd(t *testing.T) string {
	t.Helper()
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
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("os.Chdir(%q) error = %v", dir, err)
	}
	return dir
}
