package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/workspace"
)

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
	if err := writeFile(filepath.Join(exampleDir, "README.md"), exampleReadme); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(exampleDir, "QUALITY.md"), exampleModel); err != nil {
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
	fmt.Printf("Generated %s\n", filepath.ToSlash(filepath.Join(galleryRel, created.Path, "report.md")))
	return nil
}

func pinRunIdentity(runPath, createdAt, runID string) error {
	path := filepath.Join(runPath, "data", "run-manifest.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading run manifest: %w", err)
	}
	var manifest evaluation.RunManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return fmt.Errorf("decoding run manifest: %w", err)
	}
	manifest.ID = runID
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

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", filepath.ToSlash(filepath.Dir(path)), err)
	}
	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", filepath.ToSlash(path), err)
	}
	return nil
}

const exampleReadme = `
# Software service report gallery

This gallery is a generated, illustrative QUALITY.md example for a fictional
software service named LedgerLite. It exists to make Evaluation report design
easy to inspect without cutting a release or running a fresh /quality evaluate.

Open the generated run report:

- [Evaluation report](.quality/evaluations/0001-full-eval/report.md)
- [Findings](.quality/evaluations/0001-full-eval/findings.md)
- [Recommendations](.quality/evaluations/0001-full-eval/recommendations.md)

The generated Evaluation uses synthetic routine outputs. Findings, ratings,
roll-ups, recommendations, and synthetic-source:* evidence references are
fictional and demonstrate report structure only; they are not an assessment of a
real system. The concrete source system and source-code tree are intentionally
omitted for now.

Regenerate this gallery from the repository root with ` + "`mise run report-gallery`" + `.

Do not edit files under .quality/evaluations/ by hand.
`

const exampleModel = `
---
title: LedgerLite Service
description: Illustrative quality model for a fictional ledger and payments API.
ratingScale:
  - level: outstanding
    title: 🟢 Outstanding
    description: The service clearly exceeds the shared quality bar.
    criterion: Consistently exceeds the requirement with clear operational margin.
  - level: target
    title: 🔵 Target
    description: The service meets the shared quality bar.
    criterion: Meets the expected quality bar with evidence a maintainer can verify.
  - level: minimum
    title: 🟡 Minimum
    description: The service is usable, but quality gaps need attention.
    criterion: Meets the lowest acceptable bar, with visible gaps or limited evidence.
  - level: unacceptable
    title: 🔴 Unacceptable
    description: The service is below the shared quality bar.
    criterion: Falls below the lowest acceptable bar.
areas:
  api:
    title: Public API
    source: synthetic-source:api
    factors:
      correctness:
        title: Correctness
        description: Requests have clear semantics and preserve ledger intent.
      operability:
        title: Operability
        description: API behavior is understandable to callers and operators.
    requirements:
      idempotent-mutations:
        title: mutation endpoints are idempotent under retry
        factors: [correctness]
        assessment: >
          Inspect the mutation contract and retry tests for idempotency keys,
          replay behavior, and duplicate-write prevention.
      predictable-error-contracts:
        title: error responses are predictable for callers
        factors: [operability]
        assessment: >
          Compare documented error responses with handler behavior for common
          validation, authorization, and conflict cases.
  persistence:
    title: Ledger Persistence
    source: synthetic-source:persistence
    factors:
      integrity:
        title: Integrity
        description: Stored ledger state preserves accounting invariants.
      recoverability:
        title: Recoverability
        description: Data changes can be reversed or recovered when releases fail.
    requirements:
      balance-invariants:
        title: ledger mutations preserve balance invariants
        factors: [integrity]
        assessment: >
          Inspect mutation tests and reconciliation checks for conservation of
          balance across debits, credits, and failed writes.
      migration-rollback:
        title: migrations have rehearsed rollback paths
        factors: [recoverability]
        assessment: >
          Inspect migration runbooks and release notes for rollback instructions,
          rehearsal evidence, and known irreversible changes.
  operations:
    title: Operations
    source: synthetic-source:operations
    factors:
      observability:
        title: Observability
        description: Operators can understand health and customer impact quickly.
      recoverability:
        title: Recoverability
        description: Incidents have clear ownership and practiced recovery paths.
    requirements:
      customer-impact-telemetry:
        title: health signals explain customer impact
        factors: [observability]
        assessment: >
          Inspect dashboards and alerts for signals that connect service health to
          failed customer actions.
      recovery-drill-ownership:
        title: recovery drills have current owners
        factors: [recoverability]
        assessment: >
          Inspect the recovery calendar and incident playbooks for named owners,
          recency, and unresolved drill follow-up.
  agent-harness:
    title: Agent Harness
    source: synthetic-source:agent-harness
    factors:
      agent-accessibility:
        title: Agent Accessibility
        description: Agent-facing instructions expose context, checks, and limits.
    requirements:
      evaluation-entrypoint:
        title: agent guidance routes quality evaluation work
        factors: [agent-accessibility]
        assessment: >
          Inspect agent guidance for a stable entry point to the quality model,
          evaluation checks, generated reports, and permitted mutation boundaries.
---

# Quality model: LedgerLite Service

This illustrative model describes a fictional service that accepts ledger
mutation requests, persists accounting state, operates under production support,
and exposes an agent harness for quality-management work.

The example is intentionally software-specific so report design can exercise a
familiar, non-trivial service shape. It is not a default QUALITY.md template and
does not imply these Factors are universal defaults.

The generated Evaluation report uses synthetic routine outputs and synthetic
evidence references. The concrete source system is omitted.
`

type requirementCase struct {
	Area            string
	AreaTitle       string
	Factor          string
	FactorTitle     string
	Name            string
	Title           string
	Rating          string
	Confidence      string
	FindingID       string
	FindingType     string
	Severity        string
	Statement       string
	Condition       string
	Criterion       string
	BasisStatus     string
	Basis           string
	Effect          string
	RatingEffect    string
	EvidenceRef     string
	Evidence        string
	Assessment      string
	RatingRationale string
}

var requirements = []requirementCase{
	{
		Area: "api", AreaTitle: "Public API", Factor: "correctness", FactorTitle: "Correctness",
		Name: "idempotent-mutations", Title: "mutation endpoints are idempotent under retry",
		Rating: "minimum", Confidence: "medium", FindingID: "gap-001", FindingType: "gap", Severity: "high",
		Statement:   "Mutation retry behavior is not fully specified for duplicate idempotency keys.",
		Condition:   "The synthetic API contract describes idempotency keys, but replayed requests do not have a documented response contract for partial-write recovery.",
		Criterion:   "Mutation endpoints should meet the target correctness bar with retry behavior a maintainer can verify.",
		BasisStatus: "verified", Basis: "The synthetic contract excerpt names idempotency keys but omits partial-write replay behavior.",
		Effect:       "The API reaches the minimum bar but does not meet the target correctness criterion for retry semantics.",
		RatingEffect: "constrains target", EvidenceRef: "synthetic-source:api/idempotency-contract",
		Evidence:        "The synthetic contract covers idempotency-key presence and duplicate detection, but not partial-write replay outcomes.",
		Assessment:      "The idempotency contract is present but incomplete for retry recovery.",
		RatingRationale: "The finding leaves retry semantics below target while preserving a usable minimum contract.",
	},
	{
		Area: "api", AreaTitle: "Public API", Factor: "operability", FactorTitle: "Operability",
		Name: "predictable-error-contracts", Title: "error responses are predictable for callers",
		Rating: "target", Confidence: "high", FindingID: "strength-001", FindingType: "strength", Severity: "low",
		Statement:   "Common caller error cases share a documented response shape.",
		Condition:   "Validation, authorization, and conflict responses use the same synthetic error envelope.",
		Criterion:   "Error contracts should meet the target bar with evidence a maintainer can verify.",
		BasisStatus: "verified", Basis: "The synthetic handler matrix and API reference agree on the error envelope fields.",
		Effect:       "The finding supports the target operability rating for caller-facing errors.",
		RatingEffect: "supports target", EvidenceRef: "synthetic-source:api/error-contracts",
		Evidence:        "The synthetic error matrix maps common failure modes to a stable code, message, and retryable flag.",
		Assessment:      "The error contract is consistent across the sampled API cases.",
		RatingRationale: "The sampled error behavior meets the target operability criterion.",
	},
	{
		Area: "persistence", AreaTitle: "Ledger Persistence", Factor: "integrity", FactorTitle: "Integrity",
		Name: "balance-invariants", Title: "ledger mutations preserve balance invariants",
		Rating: "target", Confidence: "high", FindingID: "strength-002", FindingType: "strength", Severity: "low",
		Statement:   "Ledger mutation checks preserve balance invariants in the sampled paths.",
		Condition:   "Synthetic mutation tests cover debit, credit, failed write, and reconciliation paths.",
		Criterion:   "Ledger mutation evidence should meet the target integrity bar.",
		BasisStatus: "verified", Basis: "The synthetic test matrix includes both success and failure paths for balance preservation.",
		Effect:       "The finding supports the target integrity rating for ledger persistence.",
		RatingEffect: "supports target", EvidenceRef: "synthetic-source:persistence/balance-tests",
		Evidence:        "The synthetic tests assert balanced entries after successful writes and no balance movement after failed writes.",
		Assessment:      "The sampled persistence evidence supports balance preservation.",
		RatingRationale: "The invariant tests satisfy the target integrity criterion for the sampled paths.",
	},
	{
		Area: "persistence", AreaTitle: "Ledger Persistence", Factor: "recoverability", FactorTitle: "Recoverability",
		Name: "migration-rollback", Title: "migrations have rehearsed rollback paths",
		Rating: "minimum", Confidence: "medium", FindingID: "risk-001", FindingType: "risk", Severity: "medium",
		Statement:   "Rollback guidance exists, but rehearsal evidence is stale.",
		Condition:   "The synthetic migration runbook names rollback steps, but the last recorded rehearsal predates two schema changes.",
		Criterion:   "Migration rollback paths should meet the target recoverability bar with current rehearsal evidence.",
		BasisStatus: "plausible", Basis: "Stale rehearsal evidence plausibly misses drift in newer migration behavior.",
		Effect:       "The finding constrains recoverability to minimum until rollback rehearsal is refreshed.",
		RatingEffect: "constrains target", EvidenceRef: "synthetic-source:persistence/migration-runbook",
		Evidence:        "The synthetic runbook contains rollback steps and a rehearsal note older than the latest two migrations.",
		Assessment:      "Rollback instructions are present, but the rehearsal signal is stale.",
		RatingRationale: "The runbook keeps the requirement above unacceptable, but stale rehearsal evidence misses the target bar.",
	},
	{
		Area: "operations", AreaTitle: "Operations", Factor: "observability", FactorTitle: "Observability",
		Name: "customer-impact-telemetry", Title: "health signals explain customer impact",
		Rating: "target", Confidence: "medium", FindingID: "strength-003", FindingType: "strength", Severity: "low",
		Statement:   "Health dashboards connect service errors to failed customer actions.",
		Condition:   "Synthetic dashboards include failed ledger mutations, retry exhaustion, and queue delay panels.",
		Criterion:   "Operational signals should meet the target observability bar by explaining customer impact.",
		BasisStatus: "verified", Basis: "The synthetic dashboard inventory maps technical symptoms to customer-visible failed actions.",
		Effect:       "The finding supports the target observability rating.",
		RatingEffect: "supports target", EvidenceRef: "synthetic-source:operations/customer-impact-dashboard",
		Evidence:        "The synthetic dashboard inventory includes panels for failed customer mutations and retry exhaustion.",
		Assessment:      "The sampled telemetry explains customer impact for core failure modes.",
		RatingRationale: "The dashboard evidence meets the target observability criterion with medium confidence.",
	},
	{
		Area: "operations", AreaTitle: "Operations", Factor: "recoverability", FactorTitle: "Recoverability",
		Name: "recovery-drill-ownership", Title: "recovery drills have current owners",
		Rating: "minimum", Confidence: "low", FindingID: "unknown-001", FindingType: "unknown", Severity: "medium",
		Statement:   "The current owner for ledger recovery drills is ambiguous.",
		Condition:   "The synthetic recovery calendar names a team, while the incident playbook names a former individual owner.",
		Criterion:   "Recovery drills should meet the target recoverability bar with current named ownership.",
		BasisStatus: "not_assessed", Basis: "The available synthetic records conflict, so current ownership cannot be verified.",
		Effect:       "The finding limits confidence and constrains the requirement to minimum.",
		RatingEffect: "constrains target", EvidenceRef: "synthetic-source:operations/recovery-calendar",
		Evidence:        "The synthetic calendar and playbook disagree about who owns the next recovery drill.",
		Assessment:      "Recovery ownership is visible but ambiguous across the sampled records.",
		RatingRationale: "The requirement remains usable because recovery practice exists, but ownership ambiguity misses target.",
	},
	{
		Area: "agent-harness", AreaTitle: "Agent Harness", Factor: "agent-accessibility", FactorTitle: "Agent Accessibility",
		Name: "evaluation-entrypoint", Title: "agent guidance routes quality evaluation work",
		Rating: "target", Confidence: "high", FindingID: "strength-004", FindingType: "strength", Severity: "low",
		Statement:   "Agent guidance points to the quality model and generated report path.",
		Condition:   "The synthetic agent guidance names QUALITY.md, the evaluation command, and the report artifact to inspect.",
		Criterion:   "Agent guidance should meet the target accessibility bar with a stable evaluation entry point.",
		BasisStatus: "verified", Basis: "The synthetic guidance gives an agent a direct path from model to evaluation report.",
		Effect:       "The finding supports the target agent-accessibility rating.",
		RatingEffect: "supports target", EvidenceRef: "synthetic-source:agent-harness/guidance",
		Evidence:        "The synthetic guidance links the model, evaluation workflow, and report output.",
		Assessment:      "The agent-facing entry point is clear enough for quality evaluation work.",
		RatingRationale: "The guidance meets the target agent-accessibility criterion.",
	},
}

func galleryPayloads() []map[string]any {
	var payloads []map[string]any
	payloads = append(payloads, map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(evaluation.DataKindEvaluationFrame),
		"subject":       map[string]any{"modelLocator": "QUALITY.md"},
		"inputs":        map[string]any{"ratingLevelIds": []any{"rating:outstanding", "rating:target", "rating:minimum", "rating:unacceptable"}},
		"derivedContext": map[string]any{
			"rigor":              "synthetic-gallery",
			"evaluationPolicies": []any{"Synthetic routine outputs are used for report-design demonstration only."},
			"expectedEvaluationLimits": []any{map[string]any{
				"id":          "synthetic-source-omitted",
				"description": "The concrete source system is intentionally omitted from the gallery.",
				"impact":      "Evidence references demonstrate report shape rather than checkable source claims.",
			}},
		},
	})
	payloads = append(payloads, areaFrames()...)
	payloads = append(payloads, requirementPayloads()...)
	payloads = append(payloads, factorPayloads()...)
	payloads = append(payloads, areaAnalyses()...)
	payloads = append(payloads, advicePayloads()...)
	return payloads
}

func areaFrames() []map[string]any {
	areas := []struct {
		id       string
		source   string
		reqs     []string
		factors  []string
		children []string
	}{
		{id: "root", source: "synthetic-source:ledgerlite-service", children: []string{"api", "persistence", "operations", "agent-harness"}},
		{id: "api", source: "synthetic-source:api", reqs: []string{"idempotent-mutations", "predictable-error-contracts"}, factors: []string{"correctness", "operability"}},
		{id: "persistence", source: "synthetic-source:persistence", reqs: []string{"balance-invariants", "migration-rollback"}, factors: []string{"integrity", "recoverability"}},
		{id: "operations", source: "synthetic-source:operations", reqs: []string{"customer-impact-telemetry", "recovery-drill-ownership"}, factors: []string{"observability", "recoverability"}},
		{id: "agent-harness", source: "synthetic-source:agent-harness", reqs: []string{"evaluation-entrypoint"}, factors: []string{"agent-accessibility"}},
	}
	var payloads []map[string]any
	for _, area := range areas {
		areaID := areaRef(area.id)
		inputs := map[string]any{"sourceRefs": []any{area.source}}
		for _, req := range area.reqs {
			inputsArrayAppend(inputs, "localRequirementIds", reqRef(area.id, req))
		}
		for _, factor := range area.factors {
			inputsArrayAppend(inputs, "rootFactorIds", factorRef(area.id, factor))
		}
		for _, child := range area.children {
			inputsArrayAppend(inputs, "childAreaIds", areaRef(child))
		}
		payloads = append(payloads, map[string]any{
			"schemaVersion":  evaluation.SchemaVersion,
			"kind":           string(evaluation.DataKindAreaEvaluationFrame),
			"subject":        map[string]any{"areaId": areaID},
			"inputs":         inputs,
			"derivedContext": map[string]any{"scope": "Synthetic gallery area frame."},
		})
	}
	return payloads
}

func requirementPayloads() []map[string]any {
	var payloads []map[string]any
	for _, req := range requirements {
		payloads = append(payloads,
			requirementFrame(req),
			requirementAssessment(req),
			requirementRating(req),
		)
	}
	return payloads
}

func requirementFrame(req requirementCase) map[string]any {
	return map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(evaluation.DataKindRequirementEvaluationFrame),
		"subject": map[string]any{
			"requirementId": reqRef(req.Area, req.Name),
			"factorIds":     []any{factorRef(req.Area, req.Factor)},
		},
		"inputs": map[string]any{
			"ratingLevelIds":             []any{"rating:outstanding", "rating:target", "rating:minimum", "rating:unacceptable"},
			"requirementAssessmentBasis": "Synthetic gallery evidence.",
		},
		"derivedContext": map[string]any{
			"evidenceTargets": []any{map[string]any{
				"id":         "synthetic-evidence",
				"question":   "Does the omitted synthetic source satisfy the modeled requirement?",
				"sourceRefs": []any{req.EvidenceRef},
				"required":   true,
			}},
			"appliedRatingCriteria": []any{
				criteria("rating:target", req.Criterion),
				criteria("rating:minimum", "The requirement remains usable with visible gaps or limited evidence."),
			},
		},
	}
}

func requirementAssessment(req requirementCase) map[string]any {
	return map[string]any{
		"schemaVersion":   evaluation.SchemaVersion,
		"kind":            string(evaluation.DataKindRequirementAssessment),
		"requirementId":   reqRef(req.Area, req.Name),
		"status":          "assessed",
		"evidenceSummary": req.Assessment,
		"summary":         req.Assessment,
		"factors":         []any{factorRef(req.Area, req.Factor)},
		"findings": []any{map[string]any{
			"id":         req.FindingID,
			"type":       req.FindingType,
			"severity":   req.Severity,
			"confidence": req.Confidence,
			"statement":  req.Statement,
			"condition":  req.Condition,
			"criteria": []any{map[string]any{
				"requirementId": reqRef(req.Area, req.Name),
				"ratingLevelId": "rating:target",
				"criterion":     req.Criterion,
				"rationale":     "The gallery records one finding per requirement so report tables stay easy to inspect.",
			}},
			"basis": map[string]any{
				"status":    req.BasisStatus,
				"statement": req.Basis,
			},
			"effect": map[string]any{
				"statement":    req.Effect,
				"ratingEffect": req.RatingEffect,
			},
			"evidence": []any{map[string]any{
				"sourceRef": req.EvidenceRef,
				"statement": req.Evidence,
				"rationale": "Synthetic source reference retained to demonstrate evidence rendering.",
			}},
		}},
		"confidence":       req.Confidence,
		"confidenceReason": "Synthetic evidence is intentionally bounded to report-gallery demonstration data.",
		"evaluationLimits": []any{map[string]any{
			"id":          "synthetic-source",
			"description": "The source artifact behind this finding is omitted from the checked-in gallery.",
			"impact":      "The finding demonstrates report structure, not a verifiable source claim.",
		}},
	}
}

func requirementRating(req requirementCase) map[string]any {
	return map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(evaluation.DataKindRequirementRating),
		"requirementId": reqRef(req.Area, req.Name),
		"status":        "rated",
		"ratingLevelId": "rating:" + req.Rating,
		"rationale":     req.RatingRationale,
		"ratingDrivers": []any{map[string]any{
			"description":   req.Statement,
			"effect":        req.RatingEffect,
			"ratingLevelId": "rating:" + req.Rating,
			"inputRefs":     []any{routineRef(evaluation.DataKindRequirementAssessment, "requirementId", reqRef(req.Area, req.Name), "findings["+req.FindingID+"]")},
		}},
		"criteriaResults": []any{
			map[string]any{"ratingLevelId": "rating:target", "matched": req.Rating == "target" || req.Rating == "outstanding", "rationale": req.Criterion},
			map[string]any{"ratingLevelId": "rating:minimum", "matched": req.Rating == "minimum" || req.Rating == "target" || req.Rating == "outstanding", "rationale": "The requirement remains usable in the synthetic evaluation."},
		},
		"confidence":       req.Confidence,
		"confidenceReason": "Confidence reflects the synthetic evidence limits recorded in the assessment.",
	}
}

func factorPayloads() []map[string]any {
	var payloads []map[string]any
	for _, req := range requirements {
		factorID := factorRef(req.Area, req.Factor)
		reqID := reqRef(req.Area, req.Name)
		payloads = append(payloads, map[string]any{
			"schemaVersion": evaluation.SchemaVersion,
			"kind":          string(evaluation.DataKindFactorAnalysisFrame),
			"subject":       map[string]any{"factorId": factorID},
			"inputs": map[string]any{
				"directRequirementRatingRefs": []any{routineRef(evaluation.DataKindRequirementRating, "requirementId", reqID, "")},
			},
			"derivedContext": map[string]any{"synthesisGuidanceRef": "synthetic-gallery-factor-rollup"},
		})
		scope := analyzedScope(req.Rating, req.Confidence, req.FactorTitle+" follows its direct requirement signal.", []any{
			driver(req.Statement, req.RatingEffect, routineRef(evaluation.DataKindRequirementRating, "requirementId", reqID, "")),
		})
		payloads = append(payloads, map[string]any{
			"schemaVersion":              evaluation.SchemaVersion,
			"kind":                       string(evaluation.DataKindFactorAnalysis),
			"factorId":                   factorID,
			"localAnalysis":              scope,
			"localAndDescendantAnalysis": scope,
		})
	}
	return payloads
}

func areaAnalyses() []map[string]any {
	areaRatings := map[string]string{
		"api":           "minimum",
		"persistence":   "minimum",
		"operations":    "minimum",
		"agent-harness": "target",
	}
	areaConfidence := map[string]string{
		"api":           "medium",
		"persistence":   "medium",
		"operations":    "low",
		"agent-harness": "high",
	}
	var payloads []map[string]any
	for _, area := range []string{"api", "persistence", "operations", "agent-harness"} {
		var refs []any
		var drivers []any
		for _, req := range requirements {
			if req.Area != area {
				continue
			}
			ref := routineRef(evaluation.DataKindFactorAnalysis, "factorId", factorRef(req.Area, req.Factor), "localAndDescendantAnalysis")
			refs = append(refs, ref)
			drivers = append(drivers, driver(req.FactorTitle+" is driven by "+req.Title+".", req.RatingEffect, ref))
		}
		payloads = append(payloads, map[string]any{
			"schemaVersion":  evaluation.SchemaVersion,
			"kind":           string(evaluation.DataKindAreaAnalysisFrame),
			"subject":        map[string]any{"areaId": areaRef(area)},
			"inputs":         map[string]any{"factorAnalysisRefs": refs},
			"derivedContext": map[string]any{"synthesisGuidanceRef": "synthetic-gallery-area-rollup"},
		})
		scope := analyzedScope(areaRatings[area], areaConfidence[area], areaSummary(area), drivers)
		payloads = append(payloads, map[string]any{
			"schemaVersion":              evaluation.SchemaVersion,
			"kind":                       string(evaluation.DataKindAreaAnalysis),
			"areaId":                     areaRef(area),
			"localAnalysis":              scope,
			"localAndDescendantAnalysis": scope,
		})
	}

	var childRefs []any
	var rootDrivers []any
	for _, area := range []string{"api", "persistence", "operations", "agent-harness"} {
		ref := routineRef(evaluation.DataKindAreaAnalysis, "areaId", areaRef(area), "localAndDescendantAnalysis")
		childRefs = append(childRefs, ref)
		rootDrivers = append(rootDrivers, driver(area+" area contributes to the full service roll-up.", "contributes to minimum", ref))
	}
	payloads = append(payloads, map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(evaluation.DataKindAreaAnalysisFrame),
		"subject":       map[string]any{"areaId": "area:root"},
		"inputs":        map[string]any{"childAreaAnalysisRefs": childRefs},
		"derivedContext": map[string]any{
			"synthesisGuidanceRef": "synthetic-gallery-root-rollup",
			"expectedEvaluationLimits": []any{map[string]any{
				"id":          "synthetic-source-omitted",
				"description": "The full-service roll-up is based on synthetic routine outputs.",
			}},
		},
	})
	payloads = append(payloads, map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(evaluation.DataKindAreaAnalysis),
		"areaId":        "area:root",
		"localAnalysis": map[string]any{
			"status":       "empty",
			"statusReason": "The root area has no direct local requirements in this illustrative model.",
			"confidence":   "none",
		},
		"localAndDescendantAnalysis": analyzedScope("minimum", "medium", "LedgerLite is usable in the synthetic evaluation, but API idempotency, rollback rehearsal, and recovery ownership keep the overall service below target.", rootDrivers),
	})
	return payloads
}

func areaSummary(area string) string {
	switch area {
	case "api":
		return "The API has predictable errors, but idempotency retry semantics need a tighter contract."
	case "persistence":
		return "Ledger integrity is well covered, while rollback rehearsal evidence is stale."
	case "operations":
		return "Customer-impact telemetry is useful, but recovery drill ownership is ambiguous."
	case "agent-harness":
		return "Agent guidance exposes the quality evaluation entry point clearly."
	default:
		return "Synthetic area summary."
	}
}

func advicePayloads() []map[string]any {
	var ordered []any
	for i, req := range requirements {
		tier := "P3"
		if i < 3 {
			tier = []string{"P1", "P1", "P2"}[i]
		}
		if req.FindingType == "gap" || req.FindingType == "risk" || req.FindingType == "unknown" {
			tier = map[string]string{"gap": "P1", "risk": "P1", "unknown": "P2"}[req.FindingType]
		}
		ordered = append(ordered, map[string]any{
			"rank":       i + 1,
			"findingRef": findingRef(req),
			"tier":       tier,
			"rationale":  "Ranked by expected impact on the service quality bar and report-gallery usefulness.",
		})
	}
	recs := []map[string]any{
		recommendation(1, "Tighten the idempotency replay contract", "Specify and test partial-write replay behavior for mutation endpoints.", "The highest-ranked synthetic gap leaves retry semantics below the target API correctness bar.", "Callers and agents can verify retry behavior without inferring undocumented recovery semantics.", "The API contract and retry tests describe duplicate, replayed, and partial-write idempotency outcomes.", "high", "medium", requirements[0]),
		recommendation(2, "Rehearse migration rollback after schema changes", "Run and record a rollback rehearsal for the latest ledger migrations.", "The synthetic runbook has rollback steps, but stale rehearsal evidence leaves recoverability below target.", "Release risk drops because rollback instructions are proven against current migrations.", "A current rollback rehearsal record exists for the latest migration set.", "high", "medium", requirements[3]),
		recommendation(3, "Assign a current recovery drill owner", "Resolve the conflicting recovery-owner records and name the owner in the calendar and playbook.", "Ambiguous ownership limits confidence in recovery practice.", "Incident preparation has a clear owner agents and maintainers can route to.", "The recovery calendar and playbook agree on the current owner and next drill date.", "medium", "low", requirements[5]),
	}
	var rankedRecs []any
	for i, rec := range recs {
		rankedRecs = append(rankedRecs, map[string]any{
			"rank":              i + 1,
			"recommendationRef": rec["number"],
			"impact":            rec["impact"],
			"confidence":        rec["confidence"],
			"rationale":         "Recommendation rank follows the synthetic finding priority and expected quality-management value.",
		})
	}
	var coverage []any
	for _, req := range requirements {
		entry := map[string]any{"findingRef": findingRef(req)}
		switch req.FindingID {
		case "gap-001":
			entry["disposition"] = "addressed_by_recommendation"
			entry["recommendationRefs"] = []any{1}
			entry["rationale"] = "The idempotency recommendation directly addresses this gap."
		case "risk-001":
			entry["disposition"] = "addressed_by_recommendation"
			entry["recommendationRefs"] = []any{2}
			entry["rationale"] = "The rollback rehearsal recommendation directly addresses this risk."
		case "unknown-001":
			entry["disposition"] = "addressed_by_recommendation"
			entry["recommendationRefs"] = []any{3}
			entry["rationale"] = "The ownership recommendation resolves the unknown."
		default:
			entry["disposition"] = "not_advice_driving"
			entry["rationale"] = "The strength supports the rating but does not require follow-up advice."
		}
		coverage = append(coverage, entry)
	}
	payloads := []map[string]any{{
		"schemaVersion":   evaluation.SchemaVersion,
		"kind":            string(evaluation.DataKindFindingRanking),
		"orderedFindings": ordered,
		"rationale":       "Synthetic findings are ranked to exercise report tables and recommendation traceability.",
	}}
	payloads = append(payloads, recs...)
	payloads = append(payloads, map[string]any{
		"schemaVersion":          evaluation.SchemaVersion,
		"kind":                   string(evaluation.DataKindRecommendationRanking),
		"orderedRecommendations": rankedRecs,
		"findingCoverage":        coverage,
		"rationale":              "Synthetic recommendations focus on the highest-value gaps, risks, and unknowns.",
	})
	return payloads
}

func recommendation(number int, title, description, background, expectedValue, doneCriterion, impact, confidence string, req requirementCase) map[string]any {
	return map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(evaluation.DataKindRecommendation),
		"number":        number,
		"title":         title,
		"description":   description,
		"background":    background,
		"expectedValue": expectedValue,
		"doneCriterion": doneCriterion,
		"impact":        impact,
		"confidence":    confidence,
		"traceRefs":     []any{findingRef(req)},
	}
}

func analyzedScope(rating, confidence, rationale string, drivers []any) map[string]any {
	return map[string]any{
		"status":           "analyzed",
		"ratingLevelId":    "rating:" + rating,
		"rationale":        rationale,
		"ratingDrivers":    drivers,
		"confidence":       confidence,
		"confidenceReason": "Synthetic report-gallery judgment.",
		"evaluationLimits": []any{map[string]any{
			"id":          "synthetic-evaluation",
			"description": "This analysis is generated from synthetic routine outputs.",
			"impact":      "Use for report design and example browsing only.",
		}},
	}
}

func driver(description, effect string, ref map[string]any) map[string]any {
	return map[string]any{
		"description": description,
		"effect":      effect,
		"inputRefs":   []any{ref},
	}
}

func criteria(ratingLevelID, criterion string) map[string]any {
	return map[string]any{"ratingLevelId": ratingLevelID, "criterion": criterion}
}

func findingRef(req requirementCase) map[string]any {
	return routineRef(evaluation.DataKindRequirementAssessment, "requirementId", reqRef(req.Area, req.Name), "findings["+req.FindingID+"]")
}

func routineRef(kind evaluation.DataKind, subjectKey, id, selector string) map[string]any {
	ref := map[string]any{
		"kind":    string(kind),
		"subject": map[string]any{subjectKey: id},
	}
	if selector != "" {
		ref["selector"] = selector
	}
	return ref
}

func areaRef(area string) string {
	if area == "root" || area == "" {
		return "area:root"
	}
	return "area:" + area
}

func factorRef(area, factor string) string {
	return "factor:" + area + "::" + factor
}

func reqRef(area, req string) string {
	return "requirement:" + area + "::" + req
}

func inputsArrayAppend(target map[string]any, key string, value string) {
	values, _ := target[key].([]any)
	target[key] = append(values, value)
}
