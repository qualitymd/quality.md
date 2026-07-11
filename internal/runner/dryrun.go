package runner

import (
	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/evaluator"
	"github.com/qualitymd/quality.md/internal/receipt"
)

// Preview is the deterministic machine-readable dry-run receipt: the
// resolved model, scope, evaluator, concurrency, work-unit counts, and expected
// run path, with no evaluator invocation and no evaluation judgment data
// written.
type Preview struct {
	SchemaVersion   int                        `json:"schemaVersion"`
	Model           string                     `json:"model"`
	RequestedScope  evaluation.RunScope        `json:"requestedScope"`
	PlannedScope    evaluation.PlannedRunScope `json:"plannedScope"`
	Evaluator       string                     `json:"evaluator"`
	EvaluatorKind   string                     `json:"evaluatorKind"`
	EvaluatorReason string                     `json:"evaluatorReason"`
	// EvaluatorCandidates is the readiness evidence for each CLI candidate
	// auto discovery considered, present only for auto selection.
	EvaluatorCandidates []evaluator.CLIReadiness `json:"evaluatorCandidates,omitempty"`
	Concurrency         int                      `json:"concurrency"`
	WorkUnits           WorkUnitCounts           `json:"workUnits"`
	// Sources is the per-area source dispatch plan: each in-scope area's
	// effective selector, its detected kind, and the resolver that would
	// serve it — visible before anything runs.
	Sources         []SourcePlan     `json:"sources"`
	ExpectedRunPath string           `json:"expectedRunPath"`
	NextActions     []receipt.Action `json:"nextActions"`
}

// SourcePlan is one area's source dispatch entry in previews and run
// receipts.
type SourcePlan struct {
	Area     string `json:"area"`
	Selector string `json:"selector"`
	Kind     string `json:"kind"`
	Resolver string `json:"resolver"`
}

// DryRun resolves everything a run would use and previews it without
// creating the run folder or invoking an evaluator for judgment work.
func DryRun(opts Options) (*Preview, error) {
	ws, err := resolveRunnerWorkspace(opts)
	if err != nil {
		return nil, err
	}
	selection, err := selectEvaluator(opts, ws, requestedEvaluator(opts, ws))
	if err != nil {
		return nil, err
	}
	plan, err := evaluation.PlanRun(evaluation.Options{
		RepoRoot:   opts.RepoRoot,
		ResolveDir: opts.EvaluationDir,
		Area:       opts.Area,
		Factors:    opts.Factors,
		Model:      opts.Model,
	})
	if err != nil {
		return nil, wrapEvaluationError(err)
	}
	modelPlan, err := BuildPlan(plan.ModelSpec, plan.Manifest.PlannedScope)
	if err != nil {
		return nil, err
	}
	sourceKinds := map[string]SourceKind{}
	sources := make([]SourcePlan, 0, len(modelPlan.Areas))
	for _, area := range modelPlan.Areas {
		kind := detectSourceKind(ws.WorkspaceRoot.Abs, area.Source)
		sourceKinds[area.Ref] = kind
		sources = append(sources, SourcePlan{
			Area:     area.Ref,
			Selector: area.Source,
			Kind:     string(kind),
			Resolver: resolverForKind(kind),
		})
	}
	graph, err := BuildGraph(plan.ModelSpec, plan.Manifest.PlannedScope, sourceKinds)
	if err != nil {
		return nil, err
	}
	concurrency, err := resolveConcurrency(ws.Evaluation.Concurrency, selection.Evaluator.Capabilities())
	if err != nil {
		return nil, err
	}
	runCommand := "qualitymd evaluation run"
	if opts.Model != "" {
		runCommand += " --model " + opts.Model
	}
	if opts.Area != "" {
		runCommand += " --area " + opts.Area
	}
	for _, factor := range opts.Factors {
		runCommand += " --factor " + factor
	}
	if opts.Evaluator != "" {
		runCommand += " --evaluator " + opts.Evaluator
	}
	return &Preview{
		SchemaVersion:       evaluation.SchemaVersion,
		Model:               plan.Manifest.Model,
		RequestedScope:      plan.Manifest.RequestedScope,
		PlannedScope:        plan.Manifest.PlannedScope,
		Evaluator:           selection.Name,
		EvaluatorKind:       selection.Evaluator.Kind(),
		EvaluatorReason:     selection.Reason,
		EvaluatorCandidates: selection.Candidates,
		Concurrency:         concurrency,
		WorkUnits: WorkUnitCounts{
			Total:          len(graph.Units),
			EvaluatorUnits: graph.EvaluatorUnits(),
		},
		Sources:         sources,
		ExpectedRunPath: plan.RunRel,
		NextActions: []receipt.Action{{
			ID:      "evaluation-run",
			Label:   "Execute the evaluation run",
			Command: runCommand,
		}},
	}, nil
}
