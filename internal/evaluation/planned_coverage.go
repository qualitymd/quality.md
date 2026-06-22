package evaluation

import (
	"fmt"
	"slices"
	"strings"
)

func validatePlannedCoverage(coverage PlannedCoverage) error {
	var acc validationAccumulator
	if coverage.AssessmentResults == nil {
		acc.Add("assessmentResults", "is required")
	}
	if coverage.Analyses == nil {
		acc.Add("analyses", "is required")
	}
	assessmentResultKeys := map[string]bool{}
	for i, assessmentResult := range coverage.AssessmentResults {
		if assessmentResult.AreaPath == nil {
			acc.Add(fmt.Sprintf("assessmentResults[%d].areaPath", i), "is required")
		}
		if strings.TrimSpace(assessmentResult.Requirement) == "" {
			acc.Add(fmt.Sprintf("assessmentResults[%d].requirement", i), "is required")
		}
		key := plannedAssessmentResultIdentity(assessmentResult)
		if assessmentResultKeys[key] {
			acc.Add(fmt.Sprintf("assessmentResults[%d]", i), "duplicates an earlier planned assessment result")
		}
		assessmentResultKeys[key] = true
	}
	analysisKeys := map[string]bool{}
	for i, analysis := range coverage.Analyses {
		if analysis.AreaPath == nil {
			acc.Add(fmt.Sprintf("analyses[%d].areaPath", i), "is required")
		}
		key := plannedAnalysisIdentity(analysis)
		if analysisKeys[key] {
			acc.Add(fmt.Sprintf("analyses[%d]", i), "duplicates an earlier planned analysis")
		}
		analysisKeys[key] = true
	}
	return acc.Err()
}

func sortPlannedCoverage(coverage *PlannedCoverage) {
	slices.SortFunc(coverage.AssessmentResults, func(a, b PlannedAssessmentResult) int {
		return strings.Compare(plannedAssessmentResultIdentity(a), plannedAssessmentResultIdentity(b))
	})
	slices.SortFunc(coverage.Analyses, func(a, b PlannedCoverageAnalysis) int {
		return strings.Compare(plannedAnalysisIdentity(a), plannedAnalysisIdentity(b))
	})
}

func plannedAssessmentResultIdentity(assessmentResult PlannedAssessmentResult) string {
	return strings.Join(assessmentResult.AreaPath, "\x00") + "\x00" + assessmentResult.Requirement
}

func plannedAnalysisIdentity(analysis PlannedCoverageAnalysis) string {
	return strings.Join(analysis.AreaPath, "\x00")
}

func describePlannedAssessmentResult(assessmentResult PlannedAssessmentResult) string {
	return fmt.Sprintf("areaPath=%s requirement=%q", describeAreaPath(assessmentResult.AreaPath), assessmentResult.Requirement)
}

func describePlannedAnalysis(analysis PlannedCoverageAnalysis) string {
	return "areaPath=" + describeAreaPath(analysis.AreaPath)
}

func describeAreaPath(areaPath []string) string {
	if len(areaPath) == 0 {
		return "[]"
	}
	return fmt.Sprintf("%q", areaPath)
}
