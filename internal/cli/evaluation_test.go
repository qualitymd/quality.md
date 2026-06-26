package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/evaluation"
)

func TestEvaluationDataSetGetAndExample(t *testing.T) {
	repo := testEvaluationRepo(t)
	run, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetIn(strings.NewReader(`{
  "schemaVersion": 2,
  "kind": "RequirementAssessmentResult",
  "requirementId": "requirement:root::has-tests",
  "status": "assessed",
  "statusReason": "Evidence was inspected.",
  "evidenceSummary": "A test exists.",
  "evidenceTargetCoverage": [],
  "findings": [],
  "unknowns": [],
  "evaluationLimits": [],
  "confidence": "medium",
  "confidenceReason": "Evidence is narrow."
}`))
	cmd.SetArgs([]string{"evaluation", "data", "set", "--json", runPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("data set Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"kind": "RequirementAssessmentResult"`) ||
		!strings.Contains(out.String(), "data/areas/root/requirements/has-tests/requirement-assessment-result.json") {
		t.Fatalf("data set stdout = %s, want write receipt", out.String())
	}

	cmd = newRootCmd()
	out.Reset()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "data", "get", "--kind", "RequirementAssessmentResult", "--requirement", "requirement:root::has-tests", runPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("data get Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"kind": "RequirementAssessmentResult"`) ||
		!strings.Contains(out.String(), `"requirementId": "requirement:root::has-tests"`) {
		t.Fatalf("data get stdout = %s, want stored artifact JSON", out.String())
	}

	cmd = newRootCmd()
	out.Reset()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "data", "example", "RequirementRatingResult"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("data example Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"kind": "RequirementRatingResult"`) {
		t.Fatalf("data example stdout = %s, want artifact JSON", out.String())
	}

	cmd = newRootCmd()
	out.Reset()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "data", "schema", "requirement-assessment-result"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("data schema Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"#/$defs/RequirementAssessmentResult"`) {
		t.Fatalf("data schema stdout = %s, want kind schema", out.String())
	}

	cmd = newRootCmd()
	out.Reset()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "data", "verify", "--json", runPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("data verify Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"valid": true`) || !strings.Contains(out.String(), `"checked": 1`) {
		t.Fatalf("data verify stdout = %s, want valid receipt", out.String())
	}
}

func TestEvaluationDataArtifactCommandsRejectJSONWrapper(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "data", "example", "RequirementRatingResult", "--json"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if !strings.Contains(err.Error(), "already emits JSON") {
		t.Fatalf("Execute() error = %v, want artifact JSON diagnostic", err)
	}
}

func TestEvaluationDataKindsDocumentsPayloadKinds(t *testing.T) {
	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "data", "kinds", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{
		`"kind": "RequirementAssessmentResult"`,
		`"kind": "RequirementRatingResult"`,
		`"kind": "EvaluationOutputResult"`,
		`"agentWritable": false`,
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %s, want substring %q", out.String(), want)
		}
	}
}

func TestEvaluationListAndLatestStatusCommands(t *testing.T) {
	repo := testEvaluationRepo(t)
	if _, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Model: "QUALITY.md"}); err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})
	if err := os.Chdir(repo); err != nil {
		t.Fatalf("chdir repo: %v", err)
	}

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "list", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("list Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"runs": [`) || !strings.Contains(out.String(), `"reportable": false`) {
		t.Fatalf("list stdout = %s, want run list", out.String())
	}

	cmd = newRootCmd()
	out.Reset()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "status", "--latest", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("status Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"path": `) || !strings.Contains(out.String(), `"reportable": false`) {
		t.Fatalf("status stdout = %s, want latest status", out.String())
	}
	if !strings.Contains(out.String(), `"path": ".quality/evaluations/0001-full-eval"`) {
		t.Fatalf("status stdout = %s, want repository-relative run path", out.String())
	}
	if !strings.Contains(out.String(), `"kind": "missing-evaluation-data"`) {
		t.Fatalf("status stdout = %s, want documented empty-run gap", out.String())
	}
	if strings.Contains(out.String(), repo) {
		t.Fatalf("status stdout = %s, want no absolute repository path", out.String())
	}
}

func TestEvaluationOldCommandNamesAreRejected(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "create-run"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if !strings.Contains(err.Error(), `unknown command "create-run"`) {
		t.Fatalf("Execute() error = %v, want unknown old command", err)
	}
}

func TestEvaluationOldRecordCommandsAreRejected(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "assessment", "add"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if !strings.Contains(err.Error(), `unknown command "assessment"`) {
		t.Fatalf("Execute() error = %v, want unknown old record command", err)
	}
}

func testEvaluationRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "QUALITY.md"), []byte(`---
title: Test model
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    title: Unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-tests:
    assessment: Inspect tests.
---
`), 0o644); err != nil {
		t.Fatalf("write QUALITY.md: %v", err)
	}
	return repo
}
