package lint

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestValidFileOutput(t *testing.T) {
	result, err := Check(writeModel(t, validModel()))
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !result.Valid {
		t.Fatalf("Valid = false, findings = %#v", result.Findings)
	}
	if result.Summary != (Summary{}) {
		t.Fatalf("Summary = %#v, want zero counts", result.Summary)
	}
	if len(result.Findings) != 0 {
		t.Fatalf("len(Findings) = %d, want 0", len(result.Findings))
	}
	data, err := result.JSON()
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	if !strings.Contains(string(data), `"findings": []`) {
		t.Fatalf("JSON = %s, want empty findings array", data)
	}
}

func TestJSONDocumentShapeAndLocation(t *testing.T) {
	result, err := Check(writeModel(t, `---
title: Example
ratingScale:
  - level: target
    description: Target.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  "has an assessment":
    assessment: Inspect it.
---
`))
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	data, err := result.JSON()
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	raw := string(data)
	for _, want := range []string{`"schemaVersion": 1`, `"valid": false`, `"ruleId": "missing-criterion"`, `"modelPath": [`} {
		if !strings.Contains(raw, want) {
			t.Fatalf("JSON = %s, want %s", raw, want)
		}
	}
	finding := firstRule(t, result, RuleMissingCriterion)
	if got, want := finding.Location.ModelPath, []PathSegment{"ratingScale", 0, "criterion"}; !slices.EqualFunc(got, want, func(a, b PathSegment) bool { return a == b }) {
		t.Fatalf("modelPath = %#v, want %#v", got, want)
	}
	positionedResult, err := Check(writeModel(t, `---
title: Example
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
  "has an assessment":
    assessment: Inspect it.
    ratings:
      excellent: Exceeds it.
---
`))
	if err != nil {
		t.Fatalf("positioned Check() error = %v", err)
	}
	positioned := firstRule(t, positionedResult, RuleUnknownRatingKey)
	if positioned.Location.Line == 0 || positioned.Location.Column == 0 {
		t.Fatalf("location = %#v, want source position", positioned.Location)
	}
}

func TestDeterministicOrdering(t *testing.T) {
	path := writeModel(t, `---
ratingScale:
  - level: target
  - level: unacceptable
requirements:
  "has an assessment":
    assessment: Inspect it.
---
`)
	first, err := Check(path)
	if err != nil {
		t.Fatalf("first Check() error = %v", err)
	}
	second, err := Check(path)
	if err != nil {
		t.Fatalf("second Check() error = %v", err)
	}
	if len(first.Findings) != len(second.Findings) {
		t.Fatalf("finding counts differ: %d vs %d", len(first.Findings), len(second.Findings))
	}
	for i := range first.Findings {
		if first.Findings[i].RuleID != second.Findings[i].RuleID || first.Findings[i].Location.Label != second.Findings[i].Location.Label {
			t.Fatalf("finding %d differs: %#v vs %#v", i, first.Findings[i], second.Findings[i])
		}
	}
}

func TestMalformedParentBlocksDependentRulesOnly(t *testing.T) {
	result, err := Check(writeModel(t, `---
title: Example
ratingScale:
  target:
    criterion: Meets it.
requirements:
  "missing assessment": {}
---
`))
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !hasRule(result, RuleInvalidFrontmatter) {
		t.Fatalf("findings = %#v, want invalid-frontmatter", result.Findings)
	}
	for _, blocked := range []RuleID{RuleMissingRatingScale, RuleTooFewLevels, RuleMissingLevelName, RuleMissingCriterion} {
		if hasRule(result, blocked) {
			t.Fatalf("findings = %#v, blocked rule %s emitted", result.Findings, blocked)
		}
	}
	if !hasRule(result, RuleInvalidAssessment) {
		t.Fatalf("findings = %#v, want unrelated requirement rule to keep running", result.Findings)
	}
}

func TestLoadRejectsInvalidModel(t *testing.T) {
	_, err := Load(writeModel(t, `---
ratings:
  - level: target
    criterion: Meets it.
---
`))
	if err == nil {
		t.Fatal("Load() error = nil, want lint error")
	}
}

func TestFixAppliesRepairsAndReportsPostRepairResult(t *testing.T) {
	path := writeModel(t, `---
title: Example
description: ""
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
  "has an assessment":
    assessment: Inspect it.
---

# Body

Keep me.
`)
	result, err := Fix(path)
	if err != nil {
		t.Fatalf("Fix() error = %v", err)
	}
	if result.Summary.Fixed != 1 {
		t.Fatalf("Fixed = %d, want 1", result.Summary.Fixed)
	}
	if result.Summary.Fixable != 0 {
		t.Fatalf("Fixable = %d, want post-repair 0", result.Summary.Fixable)
	}
	if len(result.Repairs) != 1 || result.Repairs[0].RuleID != RuleEmptyProperty {
		t.Fatalf("Repairs = %#v, want empty-property repair", result.Repairs)
	}
	raw := readFile(t, path)
	if strings.Contains(raw, `description: ""`) {
		t.Fatalf("fixed file still contains empty description:\n%s", raw)
	}
	if !strings.Contains(raw, "---\n\n# Body\n\nKeep me.\n") {
		t.Fatalf("Markdown body was not preserved:\n%s", raw)
	}
	data, err := result.JSON()
	if err != nil {
		t.Fatalf("JSON() error = %v", err)
	}
	if !strings.Contains(string(data), `"fixed": 1`) || !strings.Contains(string(data), `"repairs": [`) {
		t.Fatalf("JSON repair reporting missing: %s", data)
	}
}

func TestFixNoOpDoesNotRewrite(t *testing.T) {
	path := writeModel(t, validModel())
	before := readFile(t, path)
	result, err := Fix(path)
	if err != nil {
		t.Fatalf("Fix() error = %v", err)
	}
	if result.Summary.Fixed != 0 || len(result.Repairs) != 0 {
		t.Fatalf("result = %#v, want no repairs", result)
	}
	if after := readFile(t, path); after != before {
		t.Fatalf("file changed on no-op fix:\n%s", after)
	}
}

func TestFixRefusesSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.md")
	if err := os.WriteFile(target, []byte(validModel()), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	link := filepath.Join(dir, "QUALITY.md")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("os.Symlink() error = %v", err)
	}
	_, err := Fix(link)
	if err == nil {
		t.Fatal("Fix() error = nil, want symlink refusal")
	}
	if !strings.Contains(err.Error(), "symbolic link") {
		t.Fatalf("Fix() error = %v, want symbolic link", err)
	}
}

func TestFixWriteFailureLeavesOriginalFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "QUALITY.md")
	content := `---
title: Example
description: ""
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
  "has an assessment":
    assessment: Inspect it.
---
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatalf("os.Chmod() error = %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(dir, 0o700)
	})
	_, err := Fix(path)
	if err == nil {
		t.Skip("write unexpectedly succeeded; cannot verify permission failure on this filesystem")
	}
	_ = os.Chmod(dir, 0o700)
	if got := readFile(t, path); got != content {
		t.Fatalf("file changed after failed repair:\n%s", got)
	}
}

func hasRule(result Result, ruleID RuleID) bool {
	return slices.ContainsFunc(result.Findings, func(f Finding) bool {
		return f.RuleID == ruleID
	})
}

func firstRule(t *testing.T, result Result, ruleID RuleID) Finding {
	t.Helper()
	for _, finding := range result.Findings {
		if finding.RuleID == ruleID {
			return finding
		}
	}
	t.Fatalf("findings = %#v, want %s", result.Findings, ruleID)
	return Finding{}
}

func validModel() string {
	return `---
title: Example
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    title: Unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
    requirements:
      "has an assessment":
        assessment: Inspect it.
---

# Example
`
}

func validFrontmatter(body string) string {
	return `---
title: Example
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    title: Unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
` + body + `---
`
}

func writeModel(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	return string(raw)
}
