package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeEvaluatorConfig(t *testing.T, repo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(repo, ".quality"), 0o755); err != nil {
		t.Fatalf("mkdir .quality: %v", err)
	}
	config := "evaluators:\n  team:\n    kind: anthropic\n    apiKeyEnv: QUALITYMD_TEST_KEY\n"
	if err := os.WriteFile(filepath.Join(repo, ".quality", "config.yaml"), []byte(config), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func TestEvaluationRunDryRunJSON(t *testing.T) {
	repo := testEvaluationRepo(t)
	writeEvaluatorConfig(t, repo)
	t.Setenv("QUALITYMD_TEST_KEY", "test-key")

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "run", "--dry-run", "--json",
		"--model", filepath.Join(repo, "QUALITY.md"), "--evaluator", "team"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("evaluation run --dry-run Execute() error = %v", err)
	}
	var preview struct {
		SchemaVersion     int    `json:"schemaVersion"`
		Evaluator         string `json:"evaluator"`
		EvaluatorKind     string `json:"evaluatorKind"`
		ExecutionStrategy string `json:"executionStrategy"`
		ExpectedRunPath   string `json:"expectedRunPath"`
		WorkUnits         struct {
			Total          int `json:"total"`
			EvaluatorUnits int `json:"evaluatorUnits"`
		} `json:"workUnits"`
		PlannedScope struct {
			AreaID string `json:"areaId"`
		} `json:"plannedScope"`
	}
	if err := json.Unmarshal(out.Bytes(), &preview); err != nil {
		t.Fatalf("preview is not JSON: %v\n%s", err, out.String())
	}
	if preview.Evaluator != "team" || preview.EvaluatorKind != "anthropic" {
		t.Errorf("evaluator = %s/%s, want team/anthropic", preview.Evaluator, preview.EvaluatorKind)
	}
	if preview.ExecutionStrategy != "sequential" {
		t.Errorf("strategy = %q, want sequential", preview.ExecutionStrategy)
	}
	if preview.ExpectedRunPath != ".quality/evaluations/0001-full-eval" {
		t.Errorf("expectedRunPath = %q", preview.ExpectedRunPath)
	}
	if preview.WorkUnits.Total == 0 || preview.WorkUnits.EvaluatorUnits == 0 {
		t.Errorf("workUnits = %+v", preview.WorkUnits)
	}
	if preview.PlannedScope.AreaID != "area:root" {
		t.Errorf("plannedScope.areaId = %q", preview.PlannedScope.AreaID)
	}
	// The dry run must not create the run folder.
	if _, err := os.Stat(filepath.Join(repo, ".quality", "evaluations")); !os.IsNotExist(err) {
		t.Errorf("dry run created evaluations dir: %v", err)
	}
}

func TestEvaluationRunMissingAPIKeyFailure(t *testing.T) {
	repo := testEvaluationRepo(t)
	writeEvaluatorConfig(t, repo)

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "run", "--dry-run", "--json",
		"--model", filepath.Join(repo, "QUALITY.md"), "--evaluator", "team"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() = nil, want missing_api_key failure")
	}
	var receipt struct {
		Status  string `json:"status"`
		Failure struct {
			Category string `json:"category"`
		} `json:"failure"`
	}
	if jsonErr := json.Unmarshal(out.Bytes(), &receipt); jsonErr != nil {
		t.Fatalf("failure receipt is not JSON: %v\n%s", jsonErr, out.String())
	}
	if receipt.Status != "failed" || receipt.Failure.Category != "missing_api_key" {
		t.Errorf("receipt = %+v, want failed/missing_api_key", receipt)
	}
}

func TestEvaluationRunHarnessCheckpointExitsSuccessfully(t *testing.T) {
	repo := testEvaluationRepo(t)

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "run", "--evaluator", "harness", "--json",
		"--model", filepath.Join(repo, "QUALITY.md")})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v, want success: an awaiting checkpoint is expected progress", err)
	}
	var receipt struct {
		Status           string `json:"status"`
		Evaluator        string `json:"evaluator"`
		EvaluatorRequest *struct {
			RequestID  string          `json:"requestId"`
			WorkUnitID string          `json:"workUnitId"`
			InputHash  string          `json:"inputHash"`
			Schema     json.RawMessage `json:"expectedSchema"`
		} `json:"evaluatorRequest"`
	}
	if err := json.Unmarshal(out.Bytes(), &receipt); err != nil {
		t.Fatalf("receipt is not JSON: %v\n%s", err, out.String())
	}
	if receipt.Status != "awaiting_evaluator" || receipt.Evaluator != "harness" {
		t.Errorf("receipt = %+v, want awaiting_evaluator/harness", receipt)
	}
	if receipt.EvaluatorRequest == nil || receipt.EvaluatorRequest.RequestID == "" ||
		receipt.EvaluatorRequest.InputHash == "" || len(receipt.EvaluatorRequest.Schema) == 0 {
		t.Errorf("evaluatorRequest = %+v, want the complete bounded request", receipt.EvaluatorRequest)
	}
}

func TestEvaluationRunRejectsProfileShadowingReservedName(t *testing.T) {
	repo := testEvaluationRepo(t)
	if err := os.MkdirAll(filepath.Join(repo, ".quality"), 0o755); err != nil {
		t.Fatalf("mkdir .quality: %v", err)
	}
	config := "evaluators:\n  claude:\n    kind: anthropic\n"
	if err := os.WriteFile(filepath.Join(repo, ".quality", "config.yaml"), []byte(config), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "run", "--dry-run",
		"--model", filepath.Join(repo, "QUALITY.md")})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() = nil, want reserved-name usage error")
	}
}
