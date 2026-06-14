// Package eval scores the requirements in a quality.md spec.
//
// STUB: This evaluator does not implement the evaluation model described in
// specs/cli-evaluate.md. The spec defines an agentic, judgment-based audit that
// reads the subject codebase and scores each requirement against a ratings
// scale, producing an evaluation bundle. What lives here instead is a flat
// pass/fail/skip runner over inline `bash`/`cel`/`rules` evaluators: `bash`
// shells out on exit code, `cel` runs against an empty environment, and `rules`
// (the LLM tier) is unimplemented. None of this matches the spec's two-tier
// design and is expected to be replaced wholesale.
package eval

import (
	"context"
	"os/exec"
	"strings"

	"github.com/qualitymd/quality.md/internal/spec"
)

// Status is the outcome of evaluating a requirement.
type Status int

const (
	StatusPass Status = iota
	StatusFail
	StatusSkip
)

// Result is the outcome for a single requirement.
type Result struct {
	Factor      string
	Requirement string
	Status      Status
	Detail      string
}

// Results is the full set of requirement outcomes.
type Results struct {
	Items []Result
}

// Failed reports whether any requirement failed.
func (r Results) Failed() bool {
	for _, it := range r.Items {
		if it.Status == StatusFail {
			return true
		}
	}
	return false
}

// Run evaluates every requirement in the spec.
func Run(ctx context.Context, s *spec.Spec) Results {
	var res Results
	for factorName, factor := range s.Factors {
		for reqName, req := range factor {
			res.Items = append(res.Items, evalRequirement(ctx, factorName, reqName, req))
		}
	}
	return res
}

func evalRequirement(ctx context.Context, factor, name string, req spec.Requirement) Result {
	switch {
	case req.Bash != "":
		return evalBash(ctx, factor, name, req.Bash)
	case req.CEL != "":
		return evalCEL(factor, name, req.CEL)
	case req.Prompt != "":
		return Result{
			Factor:      factor,
			Requirement: name,
			Status:      StatusSkip,
			Detail:      "prompt (LLM) evaluation not yet implemented",
		}
	default:
		return Result{
			Factor:      factor,
			Requirement: name,
			Status:      StatusSkip,
			Detail:      "no evaluator specified",
		}
	}
}

func evalBash(ctx context.Context, factor, name, script string) Result {
	out, err := exec.CommandContext(ctx, "bash", "-c", script).CombinedOutput()
	r := Result{Factor: factor, Requirement: name, Detail: strings.TrimSpace(string(out))}
	if err != nil {
		r.Status = StatusFail
		if r.Detail == "" {
			r.Detail = err.Error()
		}
		return r
	}
	r.Status = StatusPass
	return r
}
