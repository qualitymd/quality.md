package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"time"

	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/evaluator"
	"github.com/qualitymd/quality.md/internal/receipt"
	"github.com/qualitymd/quality.md/internal/workspace"
)

// Options configures one `qualitymd evaluation run` invocation.
type Options struct {
	RepoRoot      string
	Model         string
	EvaluationDir string
	Area          string
	Factors       []string
	// Evaluator is the --evaluator flag value ("" means the config default,
	// then auto).
	Evaluator string
	// Resume names an existing run folder to resume.
	Resume string
	// EvaluatorResult is the --evaluator-result path ("-" for stdin) of a
	// harness result envelope submitted for an awaiting run. Valid only with
	// Resume.
	EvaluatorResult string
	// Stdin supplies the envelope body when EvaluatorResult is "-".
	Stdin io.Reader
	// Stderr receives human progress diagnostics.
	Stderr io.Writer
	// SelectEvaluator overrides evaluator selection in tests.
	SelectEvaluator func(opts evaluator.Options) (*evaluator.Selection, error)
}

// UsageError marks a runner error caused by invalid invocation input.
type UsageError struct{ Err error }

func (e *UsageError) Error() string { return e.Err.Error() }
func (e *UsageError) Unwrap() error { return e.Err }

// RunError is a classified run failure surfaced alongside the JSON receipt.
type RunError struct {
	Category evaluator.FailureCategory
	Detail   string
}

func (e *RunError) Error() string {
	return fmt.Sprintf("%s: %s", e.Category, e.Detail)
}

// WorkUnitCounts summarizes work-graph execution for receipts.
type WorkUnitCounts struct {
	Total          int `json:"total"`
	EvaluatorUnits int `json:"evaluatorUnits"`
	Completed      int `json:"completed"`
	Reused         int `json:"reused,omitempty"`
	Failed         int `json:"failed,omitempty"`
}

// Result is the JSON receipt for a completed, awaiting, failed, or cancelled
// run.
type Result struct {
	SchemaVersion int            `json:"schemaVersion"`
	Path          string         `json:"path"`
	Status        string         `json:"status"`
	Evaluator     string         `json:"evaluator"`
	EvaluatorKind string         `json:"evaluatorKind"`
	Concurrency   int            `json:"concurrency"`
	WorkUnits     WorkUnitCounts `json:"workUnits"`
	// Sources is the per-area source dispatch plan pinned at run creation:
	// selector, detected kind, and serving resolver, in plan order.
	Sources []SourcePlan `json:"sources,omitempty"`
	Failure *Failure     `json:"failure,omitempty"`
	// EvaluatorRequest is the pending bounded work request when Status is
	// awaiting_evaluator.
	EvaluatorRequest *EvaluatorRequest        `json:"evaluatorRequest,omitempty"`
	ReportMD         string                   `json:"reportMd,omitempty"`
	RatingResult     *evaluation.RatingResult `json:"ratingResult,omitempty"`
	NextActions      []receipt.Action         `json:"nextActions,omitempty"`
}

// Run executes (or resumes) a complete evaluation run.
func Run(ctx context.Context, opts Options) (*Result, error) {
	if opts.Resume != "" {
		return resumeRun(ctx, opts)
	}
	if opts.EvaluatorResult != "" {
		return nil, &UsageError{Err: fmt.Errorf("--evaluator-result submits a harness judgment for an awaiting run and requires --resume")}
	}
	return startRun(ctx, opts)
}

func startRun(ctx context.Context, opts Options) (*Result, error) {
	ws, err := resolveRunnerWorkspace(opts)
	if err != nil {
		return nil, err
	}
	selection, err := selectEvaluator(opts, ws, requestedEvaluator(opts, ws))
	if err != nil {
		return nil, err
	}
	prepared, err := evaluation.PrepareRun(evaluation.Options{
		RepoRoot:   opts.RepoRoot,
		ResolveDir: opts.EvaluationDir,
		Area:       opts.Area,
		Factors:    opts.Factors,
		Model:      opts.Model,
	})
	if err != nil {
		return nil, wrapEvaluationError(err)
	}
	concurrency, err := resolveConcurrency(ws.Evaluation.Concurrency, selection.Evaluator.Capabilities())
	if err != nil {
		return nil, err
	}
	artifact := &Artifact{
		SchemaVersion: ArtifactSchemaVersion,
		Kind:          ArtifactKind,
		Manifest: Manifest{
			EvaluationID:   prepared.Manifest.EvaluationID,
			CreatedAt:      prepared.Manifest.CreatedAt,
			Model:          prepared.Manifest.Model,
			RequestedScope: prepared.Manifest.RequestedScope,
			PlannedScope:   prepared.Manifest.PlannedScope,
			Run:            prepared.Manifest.Run,
			Evaluator:      selection.Name,
			EvaluatorKind:  selection.Evaluator.Kind(),
			Concurrency:    concurrency,
		},
		State: State{
			Status:    StatusRunning,
			WorkUnits: map[string]*UnitState{},
			StartedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}
	store := NewStore(prepared.RunAbs)
	if err := store.Save(artifact); err != nil {
		return nil, err
	}
	return executeRun(ctx, opts, store, artifact, prepared.RunAbs, prepared.RunRel, selection, ws, nil)
}

func resumeRun(ctx context.Context, opts Options) (*Result, error) {
	runAbs, displayPath, err := resolveResumePath(opts)
	if err != nil {
		return nil, err
	}
	store := NewStore(runAbs)
	if !store.Exists() {
		return nil, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail:   fmt.Sprintf("%s has no %s; only runner-created runs can be resumed — start a new run instead", displayPath, ArtifactFile),
		}
	}
	artifact, err := store.Load()
	if err != nil {
		return nil, &RunError{Category: evaluator.FailureRunStateInvalid, Detail: err.Error() + "; start a new run instead"}
	}
	if artifact.SchemaVersion != ArtifactSchemaVersion {
		return nil, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail: fmt.Sprintf("%s schemaVersion %d is not supported by this qualitymd (want %d); start a new run instead",
				ArtifactFile, artifact.SchemaVersion, ArtifactSchemaVersion),
		}
	}
	ws, err := resolveRunnerWorkspace(opts)
	if err != nil {
		return nil, err
	}
	if artifact.Manifest.Model != ws.Model.Rel {
		return nil, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail: fmt.Sprintf("run manifest model %q does not resolve to the selected model %q; start a new run instead",
				artifact.Manifest.Model, ws.Model.Rel),
		}
	}
	if opts.Evaluator != "" && opts.Evaluator != artifact.Manifest.Evaluator {
		return nil, &UsageError{Err: fmt.Errorf(
			"--resume run was evaluated by %q; re-evaluating with --evaluator %s is a new run, not a resume",
			artifact.Manifest.Evaluator, opts.Evaluator)}
	}
	selection, err := selectEvaluator(opts, ws, artifact.Manifest.Evaluator)
	if err != nil {
		return nil, err
	}
	var harnessResult *evaluator.HarnessResultEnvelope
	if opts.EvaluatorResult != "" {
		if selection.Evaluator.Kind() != "harness" {
			return nil, &UsageError{Err: fmt.Errorf(
				"--evaluator-result submits harness judgment, but this run's evaluator is %q", artifact.Manifest.Evaluator)}
		}
		if artifact.State.PendingEvaluatorCall == nil {
			return nil, &RunError{
				Category: evaluator.FailureRunStateInvalid,
				Detail:   "no work request is awaiting a harness result for this run",
			}
		}
		harnessResult, err = loadHarnessResult(opts.EvaluatorResult, opts.Stdin)
		if err != nil {
			return nil, err
		}
	}
	artifact.State.Status = StatusRunning
	artifact.State.Failure = nil
	artifact.State.Cancelled = false
	artifact.State.CompletedAt = ""
	if err := store.Save(artifact); err != nil {
		return nil, err
	}
	return executeRun(ctx, opts, store, artifact, runAbs, displayPath, selection, ws, harnessResult)
}

func executeRun(ctx context.Context, opts Options, store *Store, artifact *Artifact, runAbs, displayPath string, selection *evaluator.Selection, ws *workspace.Workspace, harnessResult *evaluator.HarnessResultEnvelope) (*Result, error) {
	spec, err := evaluation.LoadRunModel(runAbs)
	if err != nil {
		return nil, err
	}
	plan, err := BuildPlan(spec, artifact.Manifest.PlannedScope)
	if err != nil {
		return nil, err
	}
	sourceKinds, err := pinSourceKinds(artifact, ws.WorkspaceRoot.Abs, plan)
	if err != nil {
		return nil, err
	}
	if err := store.Save(artifact); err != nil {
		return nil, err
	}
	graph, err := BuildGraph(spec, artifact.Manifest.PlannedScope, sourceKinds)
	if err != nil {
		return nil, err
	}
	logs, err := openRunLogs(runAbs)
	if err != nil {
		return nil, err
	}
	defer logs.Close()
	logs.event("run-started", map[string]any{
		"run":         artifact.Manifest.Run.Label,
		"evaluator":   selection.Name,
		"concurrency": artifact.Manifest.Concurrency,
		"reason":      selection.Reason,
		"workUnits":   len(graph.Units),
	})
	eng := &engine{
		store:         store,
		artifact:      artifact,
		graph:         graph,
		spec:          spec,
		selection:     selection,
		logs:          logs,
		stderr:        opts.Stderr,
		workspaceRoot: ws.WorkspaceRoot.Abs,
		runAbs:        runAbs,
		displayPath:   displayPath,
		now:           time.Now,
		sleep:         sleepWithContext,
		sourceBundles: map[string]*SourceBundle{},
		harnessResult: harnessResult,
	}
	eng.progress("Evaluating %s with %s (concurrency %d, %d work units)",
		displayPath, selection.Name, artifact.Manifest.Concurrency, len(graph.Units))
	status, err := eng.execute(ctx)
	if err != nil {
		return nil, err
	}
	result := buildResult(artifact, graph, displayPath, status)
	if status == StatusAwaitingEvaluator {
		result.EvaluatorRequest = eng.awaitingRequest
		result.Failure = eng.awaitingFailure
	}
	return result, nil
}

// pinSourceKinds detects and pins each in-scope area's selector kind on a
// fresh run, or reads the pinned kinds back on resume. Detection runs once,
// at run creation: resume reads the pinned kind instead of re-detecting, so a
// file appearing or vanishing mid-run cannot silently re-dispatch a selector
// to a different resolver or change the graph shape.
func pinSourceKinds(artifact *Artifact, workspaceRoot string, plan *ModelPlan) (map[string]SourceKind, error) {
	kinds := map[string]SourceKind{}
	if artifact.Sources == nil {
		artifact.Sources = map[string]*SourceRecord{}
		for _, area := range plan.Areas {
			kind := detectSourceKind(workspaceRoot, area.Source)
			artifact.Sources[area.Ref] = &SourceRecord{
				Selector: area.Source,
				Kind:     string(kind),
				Resolver: resolverForKind(kind),
			}
			kinds[area.Ref] = kind
		}
		return kinds, nil
	}
	for _, area := range plan.Areas {
		record := artifact.Sources[area.Ref]
		if record == nil {
			return nil, &RunError{
				Category: evaluator.FailureRunStateInvalid,
				Detail:   fmt.Sprintf("the run artifact pins no source record for %s; start a new run", area.Ref),
			}
		}
		kinds[area.Ref] = SourceKind(record.Kind)
	}
	return kinds, nil
}

func buildResult(artifact *Artifact, graph *Graph, displayPath, status string) *Result {
	counts := WorkUnitCounts{Total: len(graph.Units), EvaluatorUnits: graph.EvaluatorUnits()}
	for _, unit := range graph.Units {
		state, ok := artifact.State.WorkUnits[unit.ID]
		if !ok {
			continue
		}
		switch state.Status {
		case UnitCompleted:
			counts.Completed++
		case UnitFailed:
			counts.Failed++
		}
	}
	sources := make([]SourcePlan, 0, len(graph.Plan.Areas))
	for _, area := range graph.Plan.Areas {
		record := artifact.Sources[area.Ref]
		if record == nil {
			continue
		}
		sources = append(sources, SourcePlan{
			Area:     area.Ref,
			Selector: record.Selector,
			Kind:     record.Kind,
			Resolver: record.Resolver,
		})
	}
	result := &Result{
		SchemaVersion: evaluation.SchemaVersion,
		Path:          displayPath,
		Status:        status,
		Evaluator:     artifact.Manifest.Evaluator,
		EvaluatorKind: artifact.Manifest.EvaluatorKind,
		Concurrency:   artifact.Manifest.Concurrency,
		WorkUnits:     counts,
		Sources:       sources,
		Failure:       artifact.State.Failure,
	}
	switch status {
	case StatusCompleted:
		if artifact.Outputs != nil {
			result.ReportMD = artifact.Outputs.ReportMD
			result.RatingResult = artifact.Outputs.Rating
		}
		result.NextActions = []receipt.Action{{
			ID:      "evaluation-report-read",
			Label:   "Read the evaluation report",
			Command: "cat " + result.ReportMD,
		}}
	case StatusAwaitingEvaluator:
		result.NextActions = []receipt.Action{{
			ID:      "evaluation-evaluator-result",
			Label:   "Submit the harness judgment result for the pending request",
			Command: "qualitymd evaluation run --resume " + displayPath + " --evaluator-result - --json",
		}, {
			ID:      "evaluation-run-reemit",
			Label:   "Recover the pending work request",
			Command: "qualitymd evaluation run --resume " + displayPath + " --json",
		}}
	case StatusCancelled, StatusFailed:
		result.NextActions = []receipt.Action{{
			ID:      "evaluation-run-resume",
			Label:   "Resume the run",
			Command: "qualitymd evaluation run --resume " + displayPath,
		}}
	}
	return result
}

func resolveConcurrency(configured *int, caps evaluator.Capabilities) (int, error) {
	requested := 0
	if configured == nil {
		requested = defaultConcurrency()
	} else {
		requested = *configured
	}
	if requested < 1 {
		return 0, &UsageError{Err: fmt.Errorf("evaluation.concurrency must be a positive integer")}
	}
	if requested > 1 && !caps.Concurrent {
		return 1, nil
	}
	return requested, nil
}

func defaultConcurrency() int {
	n := runtime.NumCPU() * 2
	if n < 2 {
		return 2
	}
	return n
}

func resolveRunnerWorkspace(opts Options) (*workspace.Workspace, error) {
	ws, err := workspace.Resolve(workspace.Options{
		RepoRoot:              opts.RepoRoot,
		Model:                 opts.Model,
		EvaluationDirOverride: opts.EvaluationDir,
	})
	if err != nil {
		return nil, err
	}
	if err := evaluator.ValidateProfiles(ws.Evaluators); err != nil {
		return nil, &UsageError{Err: err}
	}
	return ws, nil
}

// requestedEvaluator resolves the evaluator name: the --evaluator flag, then
// the workspace config, then auto.
func requestedEvaluator(opts Options, ws *workspace.Workspace) string {
	if opts.Evaluator != "" {
		return opts.Evaluator
	}
	if ws.Evaluation.Evaluator != "" {
		return ws.Evaluation.Evaluator
	}
	return "auto"
}

func selectEvaluator(opts Options, ws *workspace.Workspace, name string) (*evaluator.Selection, error) {
	sel := opts.SelectEvaluator
	if sel == nil {
		sel = evaluator.Select
	}
	selection, err := sel(evaluator.Options{Name: name, Profiles: ws.Evaluators})
	if err != nil {
		return nil, err
	}
	return selection, nil
}

func resolveResumePath(opts Options) (string, string, error) {
	resolved, err := evaluation.ResolveRunSelection(evaluation.RunSelection{
		Model:         opts.Model,
		EvaluationDir: opts.EvaluationDir,
		RunArg:        opts.Resume,
	})
	if err != nil {
		return "", "", wrapEvaluationError(err)
	}
	abs, err := filepath.Abs(resolved.Path)
	if err != nil {
		return "", "", err
	}
	display := resolved.DisplayPath
	if display == "" {
		display = resolved.Path
	}
	return abs, display, nil
}

func wrapEvaluationError(err error) error {
	var usage *evaluation.UsageError
	if errors.As(err, &usage) {
		return &UsageError{Err: err}
	}
	return err
}

func sleepWithContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
