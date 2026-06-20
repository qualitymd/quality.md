package evaluation

import (
	"fmt"
	"slices"
	"strings"
)

func validatePlannedCoverage(coverage PlannedCoverage) error {
	if coverage.AssessmentResults == nil {
		return usagef("assessmentResults is required")
	}
	if coverage.Analyses == nil {
		return usagef("analyses is required")
	}
	assessmentResultKeys := map[string]bool{}
	for i, assessmentResult := range coverage.AssessmentResults {
		if assessmentResult.TargetPath == nil {
			return usagef("assessmentResults[%d].targetPath is required", i)
		}
		if strings.TrimSpace(assessmentResult.Requirement) == "" {
			return usagef("assessmentResults[%d].requirement is required", i)
		}
		key := plannedAssessmentResultIdentity(assessmentResult)
		if assessmentResultKeys[key] {
			return usagef("assessmentResults[%d] duplicates an earlier planned assessment result", i)
		}
		assessmentResultKeys[key] = true
	}
	analysisKeys := map[string]bool{}
	for i, analysis := range coverage.Analyses {
		if analysis.TargetPath == nil {
			return usagef("analyses[%d].targetPath is required", i)
		}
		key := plannedAnalysisIdentity(analysis)
		if analysisKeys[key] {
			return usagef("analyses[%d] duplicates an earlier planned analysis", i)
		}
		analysisKeys[key] = true
	}
	return nil
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
	return strings.Join(assessmentResult.TargetPath, "\x00") + "\x00" + assessmentResult.Requirement
}

func plannedAnalysisIdentity(analysis PlannedCoverageAnalysis) string {
	return strings.Join(analysis.TargetPath, "\x00")
}

func describePlannedAssessmentResult(assessmentResult PlannedAssessmentResult) string {
	return fmt.Sprintf("targetPath=%s requirement=%q", describeTargetPath(assessmentResult.TargetPath), assessmentResult.Requirement)
}

func describePlannedAnalysis(analysis PlannedCoverageAnalysis) string {
	return "targetPath=" + describeTargetPath(analysis.TargetPath)
}

func describeTargetPath(targetPath []string) string {
	if len(targetPath) == 0 {
		return "[]"
	}
	return fmt.Sprintf("%q", targetPath)
}
