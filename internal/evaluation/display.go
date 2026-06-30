package evaluation

import (
	"strings"
	"unicode"
)

type enumValue[T ~string] struct {
	Value       T
	Label       string
	Marker      string
	Description string
	Rank        int
}

type enumCatalog[T ~string] struct {
	Label       string
	Description string
	Values      []enumValue[T]
}

func enumStrings[T ~string](catalog enumCatalog[T]) []string {
	values := catalog.Values
	out := make([]string, len(values))
	for i, item := range values {
		out[i] = string(item.Value)
	}
	return out
}

func enumTitle[T ~string](catalog enumCatalog[T], value T) string {
	return enumStringTitle(catalog, string(value))
}

func enumStringTitle[T ~string](catalog enumCatalog[T], raw string) string {
	if display, ok := enumDisplay(catalog, raw); ok {
		return display
	}
	return humanizeEnum(raw)
}

func enumDisplay[T ~string](catalog enumCatalog[T], raw string) (string, bool) {
	for _, item := range catalog.Values {
		if string(item.Value) == raw {
			return markerLabel(item.Marker, item.Label), true
		}
	}
	return "", false
}

func enumKeyLabels[T ~string](catalog enumCatalog[T]) []string {
	out := make([]string, 0, len(catalog.Values))
	for _, item := range catalog.Values {
		out = append(out, markerLabel(item.Marker, item.Label))
	}
	return out
}

func enumRank[T ~string](catalog enumCatalog[T], raw string) (int, bool) {
	for _, item := range catalog.Values {
		if string(item.Value) == raw {
			return item.Rank, true
		}
	}
	return 0, false
}

func markerLabel(marker, label string) string {
	if marker == "" {
		return label
	}
	return marker + " " + label
}

type AnalysisStatus string

const (
	AnalysisStatusAnalyzed    AnalysisStatus = "analyzed"
	AnalysisStatusEmpty       AnalysisStatus = "empty"
	AnalysisStatusNotAnalyzed AnalysisStatus = "not_analyzed"
	AnalysisStatusBlocked     AnalysisStatus = "blocked"
)

type AssessmentStatus string

const (
	AssessmentStatusAssessed          AssessmentStatus = "assessed"
	AssessmentStatusPartiallyAssessed AssessmentStatus = "partially_assessed"
	AssessmentStatusNotAssessed       AssessmentStatus = "not_assessed"
	AssessmentStatusBlocked           AssessmentStatus = "blocked"
)

type RatingStatus string

const (
	RatingStatusRated    RatingStatus = "rated"
	RatingStatusNotRated RatingStatus = "not_rated"
	RatingStatusBlocked  RatingStatus = "blocked"
)

type ConfidenceLevel string

const (
	ConfidenceHigh   ConfidenceLevel = "high"
	ConfidenceMedium ConfidenceLevel = "medium"
	ConfidenceLow    ConfidenceLevel = "low"
	ConfidenceNone   ConfidenceLevel = "none"
)

type FindingType string

const (
	FindingTypeStrength FindingType = "strength"
	FindingTypeGap      FindingType = "gap"
	FindingTypeRisk     FindingType = "risk"
	FindingTypeNote     FindingType = "note"
)

type FindingSeverity string

const (
	FindingSeverityCritical FindingSeverity = "critical"
	FindingSeverityHigh     FindingSeverity = "high"
	FindingSeverityMedium   FindingSeverity = "medium"
	FindingSeverityLow      FindingSeverity = "low"
)

type FindingBasisStatus string

const (
	FindingBasisStatusVerified      FindingBasisStatus = "verified"
	FindingBasisStatusPlausible     FindingBasisStatus = "plausible"
	FindingBasisStatusNotAssessed   FindingBasisStatus = "not_assessed"
	FindingBasisStatusNotApplicable FindingBasisStatus = "not_applicable"
)

type RecommendationImpact string

const (
	RecommendationImpactVeryHigh RecommendationImpact = "very_high"
	RecommendationImpactHigh     RecommendationImpact = "high"
	RecommendationImpactMedium   RecommendationImpact = "medium"
	RecommendationImpactLow      RecommendationImpact = "low"
)

type FindingRankingTier string

const (
	FindingRankingTierP1 FindingRankingTier = "P1"
	FindingRankingTierP2 FindingRankingTier = "P2"
	FindingRankingTierP3 FindingRankingTier = "P3"
	FindingRankingTierP4 FindingRankingTier = "P4"
)

type FindingCoverageDisposition string

const (
	FindingCoverageAddressedByRecommendation FindingCoverageDisposition = "addressed_by_recommendation"
	FindingCoverageNotAdviceDriving          FindingCoverageDisposition = "not_advice_driving"
)

type ReportKind string

const (
	ReportKindRun         ReportKind = "run"
	ReportKindArea        ReportKind = "area"
	ReportKindFactor      ReportKind = "factor"
	ReportKindRequirement ReportKind = "requirement"
	ReportKindFindings    ReportKind = "findings"
	ReportKindAdviceIndex ReportKind = "recommendations"
	ReportKindAdvice      ReportKind = "recommendation"
)

var dataKindValues = enumCatalog[DataKind]{
	Label:       "Data kind",
	Description: "Kind of structured Evaluation payload stored for a run.",
	Values: []enumValue[DataKind]{
		{Value: DataKindEvaluationManifest, Label: "Evaluation Manifest", Marker: "📋", Description: "Evaluation metadata, scope, and run context."},
		{Value: DataKindEvaluationFrame, Label: "Evaluation Frame", Marker: "🧭", Description: "Top-level evaluation planning frame."},
		{Value: DataKindAreaEvaluationFrame, Label: "Area Evaluation Frame", Marker: "🗺️", Description: "Planned evaluation frame for an Area."},
		{Value: DataKindRequirementEvaluationFrame, Label: "Requirement Evaluation Frame", Marker: "📋", Description: "Planned assessment frame for a Requirement."},
		{Value: DataKindRequirementAssessment, Label: "Requirement Assessment", Marker: "🔎", Description: "Judgment evidence and findings for a Requirement."},
		{Value: DataKindRequirementRating, Label: "Requirement Rating", Marker: "🎚️", Description: "Rating assigned to a Requirement."},
		{Value: DataKindFactorAnalysisFrame, Label: "Factor Analysis Frame", Marker: "🧩", Description: "Planned analysis frame for a Factor."},
		{Value: DataKindFactorAnalysis, Label: "Factor Analysis", Marker: "📊", Description: "Synthesized judgment for a Factor."},
		{Value: DataKindAreaAnalysisFrame, Label: "Area Analysis Frame", Marker: "🏗️", Description: "Planned analysis frame for an Area."},
		{Value: DataKindAreaAnalysis, Label: "Area Analysis", Marker: "📈", Description: "Synthesized judgment for an Area."},
		{Value: DataKindFindingRanking, Label: "Finding Ranking", Marker: "🔝", Description: "Ordered finding priority set."},
		{Value: DataKindRecommendation, Label: "Recommendation", Marker: "💡", Description: "Proposed improvement action."},
		{Value: DataKindRecommendationRanking, Label: "Recommendation Ranking", Marker: "🏁", Description: "Ordered recommendation priority set."},
		{Value: DataKindEvaluationOutput, Label: "Evaluation Output", Marker: "📦", Description: "Generated report-output index."},
	},
}

var runGapKindValues = enumCatalog[RunGapKind]{
	Label:       "Run gap kind",
	Description: "Kind of missing, unreadable, malformed, or incomplete run data blocking reportability.",
	Values: []enumValue[RunGapKind]{
		{Value: GapMissingEvaluationData, Label: "Missing Evaluation Data", Marker: "📭", Description: "Required payload is absent."},
		{Value: GapMalformedEvaluationData, Label: "Malformed Evaluation Data", Marker: "⚠️", Description: "Payload cannot be parsed or has the wrong structure."},
		{Value: GapUnreadableEvaluationData, Label: "Unreadable Evaluation Data", Marker: "🚫", Description: "Payload exists but cannot be read."},
		{Value: GapIncompleteEvaluationData, Label: "Incomplete Evaluation Data", Marker: "🧩", Description: "Payload is readable but lacks required usable content."},
	},
}

var analysisStatusValues = enumCatalog[AnalysisStatus]{
	Label:       "Analysis status",
	Description: "Whether an Area or Factor analysis has been completed, skipped as empty, not analyzed, or blocked.",
	Values: []enumValue[AnalysisStatus]{
		{Value: AnalysisStatusAnalyzed, Label: "Analyzed", Marker: "✅", Description: "Analysis completed with a recorded result."},
		{Value: AnalysisStatusEmpty, Label: "Empty", Marker: "⬜", Description: "Scope had nothing applicable to analyze."},
		{Value: AnalysisStatusNotAnalyzed, Label: "Not Analyzed", Marker: "⚪", Description: "Analysis has not been performed."},
		{Value: AnalysisStatusBlocked, Label: "Blocked", Marker: "⛔", Description: "Analysis could not be completed."},
	},
}

var assessmentStatusValues = enumCatalog[AssessmentStatus]{
	Label:       "Assessment status",
	Description: "Whether a Requirement assessment has been completed, partially completed, not assessed, or blocked.",
	Values: []enumValue[AssessmentStatus]{
		{Value: AssessmentStatusAssessed, Label: "Assessed", Marker: "✅", Description: "Requirement was assessed with usable findings."},
		{Value: AssessmentStatusPartiallyAssessed, Label: "Partially Assessed", Marker: "🟡", Description: "Assessment is incomplete but contains usable judgment."},
		{Value: AssessmentStatusNotAssessed, Label: "Not Assessed", Marker: "⚪", Description: "No assessment was completed."},
		{Value: AssessmentStatusBlocked, Label: "Blocked", Marker: "⛔", Description: "Assessment could not be completed."},
	},
}

var ratingStatusValues = enumCatalog[RatingStatus]{
	Label:       "Rating status",
	Description: "Whether a rating result has been assigned, not assigned, or blocked.",
	Values: []enumValue[RatingStatus]{
		{Value: RatingStatusRated, Label: "Rated", Marker: "✅", Description: "A rating level was assigned."},
		{Value: RatingStatusNotRated, Label: "Not Rated", Marker: "⚪", Description: "No rating level was assigned."},
		{Value: RatingStatusBlocked, Label: "Blocked", Marker: "⛔", Description: "Rating could not be assigned."},
	},
}

var confidenceValues = enumCatalog[ConfidenceLevel]{
	Label:       "Confidence",
	Description: "Confidence in the recorded judgment based on available evidence.",
	Values: []enumValue[ConfidenceLevel]{
		{Value: ConfidenceHigh, Label: "High", Marker: "🟢", Description: "Strong evidence supports the judgment.", Rank: 0},
		{Value: ConfidenceMedium, Label: "Medium", Marker: "🔵", Description: "Adequate evidence supports the judgment with some uncertainty.", Rank: 1},
		{Value: ConfidenceLow, Label: "Low", Marker: "🟡", Description: "Limited evidence supports the judgment.", Rank: 2},
		{Value: ConfidenceNone, Label: "None", Marker: "⚪", Description: "No meaningful confidence is available.", Rank: 3},
	},
}

var ratingResultKindValues = enumCatalog[RatingResultKind]{
	Label:       "Rating result",
	Description: "Whether a Rating Result contains an assigned rating or records that the subject was not assessed.",
	Values: []enumValue[RatingResultKind]{
		{Value: RatingResultRated, Label: "Rated", Marker: "✅", Description: "The result contains an assigned rating."},
		{Value: RatingResultNotAssessed, Label: "Not Assessed", Marker: "⚪", Description: "The subject was not assessed."},
	},
}

var reportKindValues = enumCatalog[ReportKind]{
	Label:       "Report kind",
	Description: "Kind of generated Markdown report artifact.",
	Values: []enumValue[ReportKind]{
		{Value: ReportKindRun, Label: "Run", Marker: "📄", Description: "Run entrypoint report."},
		{Value: ReportKindArea, Label: "Area", Marker: "🗺️", Description: "Area report."},
		{Value: ReportKindFactor, Label: "Factor", Marker: "🧩", Description: "Factor report."},
		{Value: ReportKindRequirement, Label: "Requirement", Marker: "📋", Description: "Requirement report."},
		{Value: ReportKindFindings, Label: "Findings", Marker: "🔝", Description: "Findings index report."},
		{Value: ReportKindAdviceIndex, Label: "Recommendations", Marker: "📚", Description: "Recommendations index report."},
		{Value: ReportKindAdvice, Label: "Recommendation", Marker: "💡", Description: "Recommendation detail report."},
	},
}

var findingTypeValues = enumCatalog[FindingType]{
	Label:       "Finding type",
	Description: "Classification of what a finding contributes to the Evaluation judgment.",
	Values: []enumValue[FindingType]{
		{Value: FindingTypeStrength, Label: "Strength", Marker: "✅", Description: "Evidence of quality meeting or exceeding expectations."},
		{Value: FindingTypeGap, Label: "Gap", Marker: "🚩", Description: "Current shortfall against the quality bar."},
		{Value: FindingTypeRisk, Label: "Risk", Marker: "⚠️", Description: "Plausible future or conditional quality concern."},
		{Value: FindingTypeNote, Label: "Note", Marker: "ℹ️", Description: "Useful observation that is not itself a gap, risk, or strength."},
	},
}

var findingSeverityValues = enumCatalog[FindingSeverity]{
	Label:       "Finding severity",
	Description: "Severity of the finding's quality concern or evaluation significance.",
	Values: []enumValue[FindingSeverity]{
		{Value: FindingSeverityCritical, Label: "Critical", Marker: "🔴", Description: "Severe concern requiring urgent attention.", Rank: 0},
		{Value: FindingSeverityHigh, Label: "High", Marker: "🔴", Description: "Important concern with substantial quality impact.", Rank: 1},
		{Value: FindingSeverityMedium, Label: "Medium", Marker: "🟡", Description: "Meaningful concern worth addressing.", Rank: 2},
		{Value: FindingSeverityLow, Label: "Low", Marker: "🔵", Description: "Minor concern or low-impact observation.", Rank: 3},
	},
}

var findingBasisStatusValues = enumCatalog[FindingBasisStatus]{
	Label:       "Finding basis",
	Description: "Evidence support state for a finding's basis.",
	Values: []enumValue[FindingBasisStatus]{
		{Value: FindingBasisStatusVerified, Label: "Verified", Marker: "✅", Description: "Basis is directly supported by evidence."},
		{Value: FindingBasisStatusPlausible, Label: "Plausible", Marker: "🟡", Description: "Basis is reasonable but not fully verified."},
		{Value: FindingBasisStatusNotAssessed, Label: "Not Assessed", Marker: "⚪", Description: "Basis support was not assessed."},
		{Value: FindingBasisStatusNotApplicable, Label: "Not Applicable", Marker: "⬜", Description: "Basis support does not apply to this finding."},
	},
}

var recommendationImpactValues = enumCatalog[RecommendationImpact]{
	Label:       "Recommendation impact",
	Description: "Expected quality improvement from completing a recommendation.",
	Values: []enumValue[RecommendationImpact]{
		{Value: RecommendationImpactVeryHigh, Label: "Very high", Marker: "⬥⬥", Description: "Expected to materially improve important quality outcomes.", Rank: 0},
		{Value: RecommendationImpactHigh, Label: "High", Marker: "⬥", Description: "Expected to meaningfully improve quality.", Rank: 1},
		{Value: RecommendationImpactMedium, Label: "Medium", Marker: "●", Description: "Expected to provide useful but bounded improvement.", Rank: 2},
		{Value: RecommendationImpactLow, Label: "Low", Marker: "○", Description: "Expected to provide small or localized improvement.", Rank: 3},
	},
}

var findingRankingTierValues = enumCatalog[FindingRankingTier]{
	Label:       "Finding rank",
	Description: "Priority tier assigned when ranking findings for attention.",
	Values: []enumValue[FindingRankingTier]{
		{Value: FindingRankingTierP1, Label: "P1 Highest", Marker: "🔴", Description: "Top-priority finding for attention.", Rank: 0},
		{Value: FindingRankingTierP2, Label: "P2 High", Marker: "🟠", Description: "High-priority finding.", Rank: 1},
		{Value: FindingRankingTierP3, Label: "P3 Medium", Marker: "🟡", Description: "Medium-priority finding.", Rank: 2},
		{Value: FindingRankingTierP4, Label: "P4 Low", Marker: "⚪", Description: "Lower-priority finding.", Rank: 3},
	},
}

var findingCoverageDispositionValues = enumCatalog[FindingCoverageDisposition]{
	Label:       "Finding coverage",
	Description: "Whether a finding is addressed by recommendations or intentionally not advice-driving.",
	Values: []enumValue[FindingCoverageDisposition]{
		{Value: FindingCoverageAddressedByRecommendation, Label: "Addressed by Recommendation", Marker: "✅", Description: "Finding is covered by one or more recommendations."},
		{Value: FindingCoverageNotAdviceDriving, Label: "Not Advice Driving", Marker: "⬜", Description: "Finding is intentionally not driving recommendation work."},
	},
}

var limitTypeTitles = map[string]string{
	"incompleteInputs": "🧩 Incomplete Inputs",
	"evaluationLimits": "⚠️ Evaluation Limits",
}

var unknownTypeTitles = map[string]string{
	"unknowns":        "❓ Unknowns",
	"missingEvidence": "🔎 Missing Evidence",
}

func dataKindTitle(kind DataKind) string {
	return enumTitle(dataKindValues, kind)
}

func runGapKindTitle(kind RunGapKind) string {
	return enumTitle(runGapKindValues, kind)
}

func analysisStatusTitle(status string) string {
	return enumStringTitle(analysisStatusValues, status)
}

func assessmentStatusTitle(status string) string {
	return enumStringTitle(assessmentStatusValues, status)
}

func ratingStatusTitle(status string) string {
	return enumStringTitle(ratingStatusValues, status)
}

func confidenceTitle(confidence string) string {
	return enumStringTitle(confidenceValues, confidence)
}

func ratingResultKindTitle(kind RatingResultKind) string {
	return enumTitle(ratingResultKindValues, kind)
}

func reportKindTitle(kind string) string {
	return enumStringTitle(reportKindValues, kind)
}

func boolTitle(value bool) string {
	if value {
		return "✅ Yes"
	}
	return "⬜ No"
}

func limitTypeTitle(value string) string {
	return stringTitle(value, limitTypeTitles)
}

func unknownTypeTitle(value string) string {
	return stringTitle(value, unknownTypeTitles)
}

func findingSeverityTitle(value string) string {
	return enumStringTitle(findingSeverityValues, value)
}

func findingTypeTitle(value string) string {
	return enumStringTitle(findingTypeValues, value)
}

func impactTitle(value string) string {
	return enumStringTitle(recommendationImpactValues, value)
}

func basisStatusTitle(value string) string {
	return enumStringTitle(findingBasisStatusValues, value)
}

func findingRankingTierTitle(value string) string {
	return enumStringTitle(findingRankingTierValues, value)
}

func findingCoverageDispositionTitle(value string) string {
	return enumStringTitle(findingCoverageDispositionValues, value)
}

func stringTitle(value string, titles map[string]string) string {
	if title, ok := titles[value]; ok {
		return title
	}
	return humanizeEnum(value)
}

func humanizeEnum(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = strings.ReplaceAll(value, "_", " ")
	value = strings.ReplaceAll(value, "-", " ")
	value = splitCamel(value)
	words := strings.Fields(value)
	for i, word := range words {
		words[i] = titleWord(word)
	}
	return strings.Join(words, " ")
}

func splitCamel(value string) string {
	var b strings.Builder
	var previous rune
	for i, current := range value {
		if i > 0 && previous != ' ' && unicode.IsLower(previous) && unicode.IsUpper(current) {
			b.WriteRune(' ')
		}
		b.WriteRune(current)
		previous = current
	}
	return b.String()
}

func titleWord(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}
