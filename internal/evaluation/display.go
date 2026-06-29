package evaluation

import (
	"strings"
	"unicode"
)

type enumValue[T ~string] struct {
	Value  T
	Label  string
	Marker string
	Rank   int
}

func enumStrings[T ~string](values []enumValue[T]) []string {
	out := make([]string, len(values))
	for i, item := range values {
		out[i] = string(item.Value)
	}
	return out
}

func enumTitle[T ~string](values []enumValue[T], value T) string {
	return enumStringTitle(values, string(value))
}

func enumStringTitle[T ~string](values []enumValue[T], raw string) string {
	if display, ok := enumDisplay(values, raw); ok {
		return display
	}
	return humanizeEnum(raw)
}

func enumDisplay[T ~string](values []enumValue[T], raw string) (string, bool) {
	for _, item := range values {
		if string(item.Value) == raw {
			return markerLabel(item.Marker, item.Label), true
		}
	}
	return "", false
}

func enumKeyLabels[T ~string](values []enumValue[T]) []string {
	out := make([]string, 0, len(values))
	for _, item := range values {
		out = append(out, markerLabel(item.Marker, item.Label))
	}
	return out
}

func enumRank[T ~string](values []enumValue[T], raw string) (int, bool) {
	for _, item := range values {
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
	FindingTypeUnknown  FindingType = "unknown"
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

var dataKindValues = []enumValue[DataKind]{
	{Value: DataKindRunManifest, Label: "Run Manifest", Marker: "📋"},
	{Value: DataKindEvaluationFrame, Label: "Evaluation Frame", Marker: "🧭"},
	{Value: DataKindAreaEvaluationFrame, Label: "Area Evaluation Frame", Marker: "🗺️"},
	{Value: DataKindRequirementEvaluationFrame, Label: "Requirement Evaluation Frame", Marker: "📋"},
	{Value: DataKindRequirementAssessment, Label: "Requirement Assessment", Marker: "🔎"},
	{Value: DataKindRequirementRating, Label: "Requirement Rating", Marker: "🎚️"},
	{Value: DataKindFactorAnalysisFrame, Label: "Factor Analysis Frame", Marker: "🧩"},
	{Value: DataKindFactorAnalysis, Label: "Factor Analysis", Marker: "📊"},
	{Value: DataKindAreaAnalysisFrame, Label: "Area Analysis Frame", Marker: "🏗️"},
	{Value: DataKindAreaAnalysis, Label: "Area Analysis", Marker: "📈"},
	{Value: DataKindFindingRanking, Label: "Finding Ranking", Marker: "🔝"},
	{Value: DataKindRecommendation, Label: "Recommendation", Marker: "💡"},
	{Value: DataKindRecommendationRanking, Label: "Recommendation Ranking", Marker: "🏁"},
	{Value: DataKindEvaluationOutput, Label: "Evaluation Output", Marker: "📦"},
}

var runGapKindValues = []enumValue[RunGapKind]{
	{Value: GapMissingEvaluationData, Label: "Missing Evaluation Data", Marker: "📭"},
	{Value: GapMalformedEvaluationData, Label: "Malformed Evaluation Data", Marker: "⚠️"},
	{Value: GapUnreadableEvaluationData, Label: "Unreadable Evaluation Data", Marker: "🚫"},
	{Value: GapIncompleteEvaluationData, Label: "Incomplete Evaluation Data", Marker: "🧩"},
}

var analysisStatusValues = []enumValue[AnalysisStatus]{
	{Value: AnalysisStatusAnalyzed, Label: "Analyzed", Marker: "✅"},
	{Value: AnalysisStatusEmpty, Label: "Empty", Marker: "⬜"},
	{Value: AnalysisStatusNotAnalyzed, Label: "Not Analyzed", Marker: "⚪"},
	{Value: AnalysisStatusBlocked, Label: "Blocked", Marker: "⛔"},
}

var assessmentStatusValues = []enumValue[AssessmentStatus]{
	{Value: AssessmentStatusAssessed, Label: "Assessed", Marker: "✅"},
	{Value: AssessmentStatusPartiallyAssessed, Label: "Partially Assessed", Marker: "🟡"},
	{Value: AssessmentStatusNotAssessed, Label: "Not Assessed", Marker: "⚪"},
	{Value: AssessmentStatusBlocked, Label: "Blocked", Marker: "⛔"},
}

var ratingStatusValues = []enumValue[RatingStatus]{
	{Value: RatingStatusRated, Label: "Rated", Marker: "✅"},
	{Value: RatingStatusNotRated, Label: "Not Rated", Marker: "⚪"},
	{Value: RatingStatusBlocked, Label: "Blocked", Marker: "⛔"},
}

var confidenceValues = []enumValue[ConfidenceLevel]{
	{Value: ConfidenceHigh, Label: "High", Marker: "🟢", Rank: 0},
	{Value: ConfidenceMedium, Label: "Medium", Marker: "🔵", Rank: 1},
	{Value: ConfidenceLow, Label: "Low", Marker: "🟡", Rank: 2},
	{Value: ConfidenceNone, Label: "None", Marker: "⚪", Rank: 3},
}

var ratingResultKindValues = []enumValue[RatingResultKind]{
	{Value: RatingResultRated, Label: "Rated", Marker: "✅"},
	{Value: RatingResultNotAssessed, Label: "Not Assessed", Marker: "⚪"},
}

var reportKindValues = []enumValue[ReportKind]{
	{Value: ReportKindRun, Label: "Run", Marker: "📄"},
	{Value: ReportKindArea, Label: "Area", Marker: "🗺️"},
	{Value: ReportKindFactor, Label: "Factor", Marker: "🧩"},
	{Value: ReportKindRequirement, Label: "Requirement", Marker: "📋"},
	{Value: ReportKindFindings, Label: "Findings", Marker: "🔝"},
	{Value: ReportKindAdviceIndex, Label: "Recommendations", Marker: "📚"},
	{Value: ReportKindAdvice, Label: "Recommendation", Marker: "💡"},
}

var findingTypeValues = []enumValue[FindingType]{
	{Value: FindingTypeStrength, Label: "Strength", Marker: "✅"},
	{Value: FindingTypeGap, Label: "Gap", Marker: "⚠️"},
	{Value: FindingTypeRisk, Label: "Risk", Marker: "⚠️"},
	{Value: FindingTypeUnknown, Label: "Unknown", Marker: "❓"},
	{Value: FindingTypeNote, Label: "Note", Marker: "ℹ️"},
}

var findingSeverityValues = []enumValue[FindingSeverity]{
	{Value: FindingSeverityCritical, Label: "Critical", Marker: "🔴", Rank: 0},
	{Value: FindingSeverityHigh, Label: "High", Marker: "🔴", Rank: 1},
	{Value: FindingSeverityMedium, Label: "Medium", Marker: "🟡", Rank: 2},
	{Value: FindingSeverityLow, Label: "Low", Marker: "🔵", Rank: 3},
}

var findingBasisStatusValues = []enumValue[FindingBasisStatus]{
	{Value: FindingBasisStatusVerified, Label: "Verified", Marker: "✅"},
	{Value: FindingBasisStatusPlausible, Label: "Plausible", Marker: "🟡"},
	{Value: FindingBasisStatusNotAssessed, Label: "Not Assessed", Marker: "⚪"},
	{Value: FindingBasisStatusNotApplicable, Label: "Not Applicable", Marker: "⬜"},
}

var recommendationImpactValues = []enumValue[RecommendationImpact]{
	{Value: RecommendationImpactVeryHigh, Label: "Very high", Marker: "◆", Rank: 0},
	{Value: RecommendationImpactHigh, Label: "High", Marker: "▲", Rank: 1},
	{Value: RecommendationImpactMedium, Label: "Medium", Marker: "●", Rank: 2},
	{Value: RecommendationImpactLow, Label: "Low", Marker: "○", Rank: 3},
}

var findingRankingTierValues = []enumValue[FindingRankingTier]{
	{Value: FindingRankingTierP1, Label: "P1 Highest", Marker: "🔴", Rank: 0},
	{Value: FindingRankingTierP2, Label: "P2 High", Marker: "🟠", Rank: 1},
	{Value: FindingRankingTierP3, Label: "P3 Medium", Marker: "🟡", Rank: 2},
	{Value: FindingRankingTierP4, Label: "P4 Low", Marker: "⚪", Rank: 3},
}

var findingCoverageDispositionValues = []enumValue[FindingCoverageDisposition]{
	{Value: FindingCoverageAddressedByRecommendation, Label: "Addressed by Recommendation", Marker: "✅"},
	{Value: FindingCoverageNotAdviceDriving, Label: "Not Advice Driving", Marker: "⬜"},
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
