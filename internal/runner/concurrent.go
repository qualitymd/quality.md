package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/evaluator"
)

type EvaluationStepCompletion struct {
	Step     *Unit
	Accepted *AcceptedEvaluationStep
	Rejected *RejectedEvaluationStep
	Reused   bool
	Attempts int
	Err      error
}

type AcceptedEvaluationStep struct {
	Payloads  []map[string]any
	InputHash string
}

type RejectedEvaluationStep struct {
	Failure *Failure
}

type evaluationStepCall struct {
	step           *Unit
	request        evaluator.WorkRequest
	inputHash      string
	attemptsBefore int
}

type evaluationScheduler struct {
	e           *engine
	ctx         context.Context
	cancel      context.CancelFunc
	jobs        chan evaluationStepCall
	completions chan EvaluationStepCompletion
	workers     sync.WaitGroup
	scheduled   map[string]bool
	done        map[string]bool
	active      int
	completed   int
}

func (e *engine) executeConcurrent(ctx context.Context) (string, error) {
	s := newEvaluationScheduler(ctx, e)
	defer s.close()
	return s.run()
}

func newEvaluationScheduler(ctx context.Context, e *engine) *evaluationScheduler {
	ctx, cancel := context.WithCancel(ctx)
	s := &evaluationScheduler{
		e:           e,
		ctx:         ctx,
		cancel:      cancel,
		jobs:        make(chan evaluationStepCall),
		completions: make(chan EvaluationStepCompletion, e.artifact.Manifest.Concurrency),
		scheduled:   map[string]bool{},
		done:        map[string]bool{},
	}
	for i := 0; i < e.artifact.Manifest.Concurrency; i++ {
		s.workers.Add(1)
		go func() {
			defer s.workers.Done()
			for call := range s.jobs {
				s.completions <- e.completeJudgmentStep(ctx, call)
			}
		}()
	}
	return s
}

func (s *evaluationScheduler) close() {
	s.cancel()
	close(s.jobs)
	s.workers.Wait()
}

func (s *evaluationScheduler) run() (string, error) {
	for s.completed < len(s.e.graph.Units) {
		if s.ctx.Err() != nil {
			return s.e.markCancelled()
		}
		progressed, status, err := s.scheduleReady()
		if err != nil || status != "" {
			s.cancel()
			return status, err
		}
		if s.completed == len(s.e.graph.Units) {
			break
		}
		if s.active == 0 {
			if progressed {
				continue
			}
			return StatusFailed, fmt.Errorf("evaluation work graph made no progress")
		}
		if status, err := s.receiveCompletion(); err != nil || status != "" {
			s.cancel()
			return status, err
		}
	}
	return s.completeRun()
}

func (s *evaluationScheduler) scheduleReady() (bool, string, error) {
	progressed := false
	for {
		unit := nextReadyEvaluationStep(s.e.graph.Units, s.scheduled, s.done)
		if unit == nil || (unit.EvaluatorBacked && s.active >= s.e.artifact.Manifest.Concurrency) {
			return progressed, "", nil
		}
		s.scheduled[unit.ID] = true
		progressed = true
		status, err := s.schedule(unit)
		if err != nil || status != "" {
			return progressed, status, err
		}
	}
}

func (s *evaluationScheduler) schedule(unit *Unit) (string, error) {
	switch {
	case unit.Kind == KindBuildReports:
		return s.runReportStep(unit)
	case !unit.EvaluatorBacked:
		return s.acceptCompletion(s.e.completeFrameStep(unit))
	default:
		return s.scheduleJudgmentStep(unit)
	}
}

func (s *evaluationScheduler) runReportStep(unit *Unit) (string, error) {
	failed, err := s.e.runBuildReports(unit)
	if err != nil {
		return StatusFailed, err
	}
	if failed {
		return StatusFailed, nil
	}
	s.markDone(unit.ID)
	return "", nil
}

func (s *evaluationScheduler) scheduleJudgmentStep(unit *Unit) (string, error) {
	call, completion, err := s.e.prepareJudgmentStep(unit)
	if err != nil {
		return StatusFailed, err
	}
	if completion != nil {
		return s.acceptCompletion(*completion)
	}
	s.active++
	s.jobs <- *call
	return "", nil
}

func (s *evaluationScheduler) receiveCompletion() (string, error) {
	completion := <-s.completions
	s.active--
	if s.ctx.Err() != nil {
		return s.e.markCancelled()
	}
	return s.acceptCompletion(completion)
}

func (s *evaluationScheduler) acceptCompletion(completion EvaluationStepCompletion) (string, error) {
	if completion.Err != nil {
		return StatusFailed, completion.Err
	}
	failed, err := s.e.acceptEvaluationStepCompletion(completion)
	if err != nil {
		return StatusFailed, err
	}
	if failed {
		return StatusFailed, nil
	}
	s.markDone(completion.Step.ID)
	return "", nil
}

func (s *evaluationScheduler) markDone(unitID string) {
	s.done[unitID] = true
	s.completed++
}

func (s *evaluationScheduler) completeRun() (string, error) {
	s.e.artifact.State.Status = StatusCompleted
	s.e.artifact.State.CompletedAt = s.e.timestamp()
	if err := s.e.save(); err != nil {
		return StatusFailed, err
	}
	s.e.logs.event("run-completed", map[string]any{"run": s.e.artifact.Manifest.Run.Label})
	return StatusCompleted, nil
}

func nextReadyEvaluationStep(units []*Unit, scheduled, done map[string]bool) *Unit {
	for _, unit := range units {
		if scheduled[unit.ID] {
			continue
		}
		ready := true
		for _, dep := range unit.DependsOn {
			if !done[dep] {
				ready = false
				break
			}
		}
		if ready {
			return unit
		}
	}
	return nil
}

func (e *engine) completeFrameStep(unit *Unit) EvaluationStepCompletion {
	payload, err := e.deterministicPayload(unit)
	if err != nil {
		return EvaluationStepCompletion{Step: unit, Err: err}
	}
	hash := hashJSON(payload)
	state := e.artifact.State.unit(unit.ID)
	if state.Status == UnitCompleted && state.InputHash == hash && e.payloadFor(unit.ID) != nil {
		e.logs.event("work-unit-reused", map[string]any{"workUnit": unit.ID})
		return EvaluationStepCompletion{Step: unit, Reused: true}
	}
	if err := validateEvaluationStepPayload(unit, payload, e); err != nil {
		return EvaluationStepCompletion{Step: unit, Err: err}
	}
	return EvaluationStepCompletion{
		Step: unit,
		Accepted: &AcceptedEvaluationStep{
			Payloads:  []map[string]any{payload},
			InputHash: hash,
		},
	}
}

func (e *engine) prepareJudgmentStep(unit *Unit) (*evaluationStepCall, *EvaluationStepCompletion, error) {
	if failure := e.resolveSourceGuard(unit); failure != nil {
		return nil, &EvaluationStepCompletion{
			Step:     unit,
			Rejected: &RejectedEvaluationStep{Failure: failure},
		}, nil
	}
	req, err := e.buildWorkRequest(unit)
	if err != nil {
		if failure := sourceUnavailableFailure(err); failure != nil {
			return nil, &EvaluationStepCompletion{
				Step:     unit,
				Rejected: &RejectedEvaluationStep{Failure: failure},
			}, nil
		}
		return nil, nil, err
	}
	inputHash := workUnitInputHash(req)
	state := e.artifact.State.unit(unit.ID)
	if state.Status == UnitCompleted && state.InputHash == inputHash && e.unitResultPresent(unit) {
		e.logs.event("work-unit-reused", map[string]any{"workUnit": unit.ID})
		return nil, &EvaluationStepCompletion{Step: unit, Reused: true}, nil
	}
	e.progress("%s...", unit.ID)
	e.logs.event("work-unit-started", map[string]any{"workUnit": unit.ID, "evaluator": e.selection.Name})
	return &evaluationStepCall{
		step:           unit,
		request:        req,
		inputHash:      inputHash,
		attemptsBefore: state.Attempts,
	}, nil, nil
}

func (e *engine) completeJudgmentStep(ctx context.Context, call evaluationStepCall) EvaluationStepCompletion {
	var lastFailure *Failure
	attemptsUsed := 0
	for attempt := 1; attempt <= maxUnitAttempts; attempt++ {
		if ctx.Err() != nil {
			return EvaluationStepCompletion{Step: call.step, Attempts: attemptsUsed}
		}
		attemptsUsed++
		attemptNumber := call.attemptsBefore + attemptsUsed
		started := e.now()
		callCtx, cancelCall := context.WithTimeout(ctx, evaluatorCallTimeout)
		result, err := e.selection.Evaluator.Evaluate(callCtx, call.request)
		cancelCall()
		duration := e.now().Sub(started)
		if err != nil {
			return EvaluationStepCompletion{
				Step:     call.step,
				Attempts: attemptsUsed,
				Err:      fmt.Errorf("evaluator %s: %w", e.selection.Name, err),
			}
		}

		failure := result.Failure
		detail := result.FailureDetail
		var payloads []map[string]any
		if failure == "" {
			payloads, failure, detail = e.acceptResultPayload(call.step, result.Payload)
		}
		e.logCall(call.step, call.request, result, attemptNumber, duration, failure, detail)

		if failure == "" {
			return EvaluationStepCompletion{
				Step:     call.step,
				Attempts: attemptsUsed,
				Accepted: &AcceptedEvaluationStep{
					Payloads:  payloads,
					InputHash: call.inputHash,
				},
			}
		}
		lastFailure = &Failure{Category: failure, Detail: detail}
		e.logs.event("work-unit-attempt-failed", map[string]any{
			"workUnit": call.step.ID,
			"attempt":  attemptNumber,
			"category": string(failure),
			"detail":   detail,
		})
		if failure == evaluator.FailureCancelled || ctx.Err() != nil {
			return EvaluationStepCompletion{Step: call.step, Attempts: attemptsUsed}
		}
		if _, retryable := retryableFailures[failure]; !retryable || attempt == maxUnitAttempts {
			break
		}
		if failure == evaluator.FailureRateLimited {
			e.sleep(ctx, time.Duration(attempt)*2*time.Second)
		}
		e.progress("%s: retrying after %s", call.step.ID, failure)
	}
	return EvaluationStepCompletion{
		Step:     call.step,
		Attempts: attemptsUsed,
		Rejected: &RejectedEvaluationStep{Failure: lastFailure},
	}
}

func (e *engine) acceptEvaluationStepCompletion(completion EvaluationStepCompletion) (bool, error) {
	state := e.artifact.State.unit(completion.Step.ID)
	state.Attempts += completion.Attempts
	if completion.Reused {
		return false, nil
	}
	if completion.Rejected != nil {
		state.Status = UnitFailed
		state.Failure = completion.Rejected.Failure
		e.artifact.State.Failure = completion.Rejected.Failure
		if err := e.save(); err != nil {
			return true, err
		}
		e.progress("%s failed: %s: %s", completion.Step.ID, completion.Rejected.Failure.Category, completion.Rejected.Failure.Detail)
		return true, nil
	}
	if completion.Accepted == nil {
		return false, nil
	}
	if err := e.acceptUnit(completion.Step, state, completion.Accepted.Payloads, completion.Accepted.InputHash); err != nil {
		return false, err
	}
	return false, nil
}

func validateEvaluationStepPayload(unit *Unit, payload map[string]any, e *engine) error {
	return evaluation.ValidatePayload(unit.DataKind, payload, e.spec)
}
