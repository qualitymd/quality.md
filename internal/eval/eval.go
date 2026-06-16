// Package eval traverses QUALITY.md requirements.
//
// STUB: The current specification describes judgment-based assessment and
// rating. This package does not perform that judgment yet; it records each
// requirement as not assessed so the placeholder CLI can exercise the current
// frontmatter shape without claiming a conforming evaluation.
package eval

import (
	"context"
	"strings"

	"github.com/qualitymd/quality.md/internal/spec"
)

// Status is the outcome of evaluating a requirement.
type Status int

const (
	StatusRated Status = iota
	StatusNotAssessed
)

// Result is the outcome for a single requirement.
type Result struct {
	Target      string
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
	return false
}

// Run traverses every requirement in the spec.
func Run(_ context.Context, s *spec.Spec) Results {
	var res Results
	for reqName, req := range s.Requirements {
		res.Items = append(res.Items, result("(root)", "", reqName, req))
	}
	for factorName, factor := range s.Factors {
		walkFactor(&res, "(root)", factorName, factor)
	}
	for targetName, target := range s.Targets {
		walkTarget(&res, targetName, target)
	}
	return res
}

func walkTarget(res *Results, targetName string, target spec.Target) {
	for reqName, req := range target.Requirements {
		res.Items = append(res.Items, result(targetName, "", reqName, req))
	}
	for factorName, factor := range target.Factors {
		walkFactor(res, targetName, factorName, factor)
	}
	for childName, child := range target.Targets {
		walkTarget(res, targetName+"/"+childName, child)
	}
}

func walkFactor(res *Results, targetName, factorName string, factor spec.Factor) {
	for reqName, req := range factor.Requirements {
		res.Items = append(res.Items, result(targetName, factorName, reqName, req))
	}
	for childName, child := range factor.Factors {
		walkFactor(res, targetName, factorName+"/"+childName, child)
	}
}

func result(targetName, factorName, reqName string, req spec.Requirement) Result {
	return Result{
		Target:      targetName,
		Factor:      factorName,
		Requirement: reqName,
		Status:      StatusNotAssessed,
		Detail:      summarizeAssessment(req.Assessment),
	}
}

func summarizeAssessment(assessment string) string {
	assessment = strings.Join(strings.Fields(assessment), " ")
	if len(assessment) <= 96 {
		return "assessment not run: " + assessment
	}
	return "assessment not run: " + assessment[:93] + "..."
}
