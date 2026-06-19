package evaluation

import (
	"fmt"
	"slices"
	"strings"
)

func validatePlannedCoverage(coverage PlannedCoverage) error {
	if coverage.Assessments == nil {
		return usagef("assessments is required")
	}
	if coverage.Analyses == nil {
		return usagef("analyses is required")
	}
	assessmentKeys := map[string]bool{}
	for i, assessment := range coverage.Assessments {
		if assessment.TargetPath == nil {
			return usagef("assessments[%d].targetPath is required", i)
		}
		if strings.TrimSpace(assessment.Requirement) == "" {
			return usagef("assessments[%d].requirement is required", i)
		}
		key := plannedAssessmentIdentity(assessment)
		if assessmentKeys[key] {
			return usagef("assessments[%d] duplicates an earlier planned assessment", i)
		}
		assessmentKeys[key] = true
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
	slices.SortFunc(coverage.Assessments, func(a, b PlannedCoverageAssessment) int {
		return strings.Compare(plannedAssessmentIdentity(a), plannedAssessmentIdentity(b))
	})
	slices.SortFunc(coverage.Analyses, func(a, b PlannedCoverageAnalysis) int {
		return strings.Compare(plannedAnalysisIdentity(a), plannedAnalysisIdentity(b))
	})
}

func plannedAssessmentIdentity(assessment PlannedCoverageAssessment) string {
	return strings.Join(assessment.TargetPath, "\x00") + "\x00" + assessment.Requirement
}

func plannedAnalysisIdentity(analysis PlannedCoverageAnalysis) string {
	return strings.Join(analysis.TargetPath, "\x00")
}

func describePlannedAssessment(assessment PlannedCoverageAssessment) string {
	return fmt.Sprintf("targetPath=%s requirement=%q", describeTargetPath(assessment.TargetPath), assessment.Requirement)
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
