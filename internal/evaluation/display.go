package evaluation

import (
	"strings"
	"unicode"
)

type displayCatalog[T ~string] map[T]string

func (c displayCatalog[T]) title(value T) string {
	if title, ok := c[value]; ok {
		return title
	}
	return humanizeEnum(string(value))
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

type ReportKind string

const (
	ReportKindRun         ReportKind = "run"
	ReportKindArea        ReportKind = "area"
	ReportKindFactor      ReportKind = "factor"
	ReportKindRequirement ReportKind = "requirement"
	ReportKindAdviceIndex ReportKind = "recommendations"
	ReportKindAdvice      ReportKind = "recommendation"
)

// reportKinds is the single typed source for the report-reference kind
// vocabulary: every report kind a report reference object may name.
var reportKinds = []ReportKind{
	ReportKindRun,
	ReportKindArea,
	ReportKindFactor,
	ReportKindRequirement,
	ReportKindAdviceIndex,
	ReportKindAdvice,
}

var dataKindTitles = displayCatalog[DataKind]{
	DataKindRunManifest:                "📋 Run Manifest",
	DataKindEvaluationFrame:            "🧭 Evaluation Frame",
	DataKindAreaEvaluationFrame:        "🗺️ Area Evaluation Frame",
	DataKindRequirementEvaluationFrame: "📋 Requirement Evaluation Frame",
	DataKindRequirementAssessment:      "🔎 Requirement Assessment",
	DataKindRequirementRating:          "🎚️ Requirement Rating",
	DataKindFactorAnalysisFrame:        "🧩 Factor Analysis Frame",
	DataKindFactorAnalysis:             "📊 Factor Analysis",
	DataKindAreaAnalysisFrame:          "🏗️ Area Analysis Frame",
	DataKindAreaAnalysis:               "📈 Area Analysis",
	DataKindFindingRanking:             "🔝 Finding Ranking",
	DataKindRecommendation:             "💡 Recommendation",
	DataKindRecommendationRanking:      "🏁 Recommendation Ranking",
	DataKindEvaluationOutput:           "📦 Evaluation Output",
}

var runGapKindTitles = displayCatalog[RunGapKind]{
	GapMissingEvaluationData:    "📭 Missing Evaluation Data",
	GapMalformedEvaluationData:  "⚠️ Malformed Evaluation Data",
	GapUnreadableEvaluationData: "🚫 Unreadable Evaluation Data",
	GapIncompleteEvaluationData: "🧩 Incomplete Evaluation Data",
}

var analysisStatusTitles = displayCatalog[AnalysisStatus]{
	AnalysisStatusAnalyzed:    "✅ Analyzed",
	AnalysisStatusEmpty:       "⬜ Empty",
	AnalysisStatusNotAnalyzed: "⚪ Not Analyzed",
	AnalysisStatusBlocked:     "⛔ Blocked",
}

var assessmentStatusTitles = displayCatalog[AssessmentStatus]{
	AssessmentStatusAssessed:          "✅ Assessed",
	AssessmentStatusPartiallyAssessed: "🟡 Partially Assessed",
	AssessmentStatusNotAssessed:       "⚪ Not Assessed",
	AssessmentStatusBlocked:           "⛔ Blocked",
}

var ratingStatusTitles = displayCatalog[RatingStatus]{
	RatingStatusRated:    "✅ Rated",
	RatingStatusNotRated: "⚪ Not Rated",
	RatingStatusBlocked:  "⛔ Blocked",
}

var confidenceTitles = displayCatalog[ConfidenceLevel]{
	ConfidenceHigh:   "🟢 High",
	ConfidenceMedium: "🔵 Medium",
	ConfidenceLow:    "🟡 Low",
	ConfidenceNone:   "⚪ None",
}

var ratingResultKindTitles = displayCatalog[RatingResultKind]{
	RatingResultRated:       "✅ Rated",
	RatingResultNotAssessed: "⚪ Not Assessed",
}

var reportKindTitles = displayCatalog[ReportKind]{
	ReportKindRun:         "📄 Run",
	ReportKindArea:        "🗺️ Area",
	ReportKindFactor:      "🧩 Factor",
	ReportKindRequirement: "📋 Requirement",
}

var limitTypeTitles = map[string]string{
	"incompleteInputs": "🧩 Incomplete Inputs",
	"evaluationLimits": "⚠️ Evaluation Limits",
}

var unknownTypeTitles = map[string]string{
	"unknowns":        "❓ Unknowns",
	"missingEvidence": "🔎 Missing Evidence",
}

var findingSeverityTitles = map[string]string{
	"critical": "🔴 Critical",
	"high":     "🔴 High",
	"medium":   "🟡 Medium",
	"low":      "🔵 Low",
}

var findingTypeTitles = map[string]string{
	"gap":            "⚠️ Gap",
	"risk":           "⚠️ Risk",
	"strength":       "✅ Strength",
	"evidence":       "🔎 Evidence",
	"unknown":        "❓ Unknown",
	"note":           "ℹ️ Note",
	"recommendation": "➡️ Recommendation",
}

var basisStatusTitles = map[string]string{
	"verified":       "Verified",
	"plausible":      "Plausible",
	"not_assessed":   "Not Assessed",
	"not_applicable": "Not Applicable",
}

func dataKindTitle(kind DataKind) string {
	return dataKindTitles.title(kind)
}

func runGapKindTitle(kind RunGapKind) string {
	return runGapKindTitles.title(kind)
}

func analysisStatusTitle(status string) string {
	return analysisStatusTitles.title(AnalysisStatus(status))
}

func assessmentStatusTitle(status string) string {
	return assessmentStatusTitles.title(AssessmentStatus(status))
}

func ratingStatusTitle(status string) string {
	return ratingStatusTitles.title(RatingStatus(status))
}

func confidenceTitle(confidence string) string {
	return confidenceTitles.title(ConfidenceLevel(confidence))
}

func ratingResultKindTitle(kind RatingResultKind) string {
	return ratingResultKindTitles.title(kind)
}

func reportKindTitle(kind string) string {
	return reportKindTitles.title(ReportKind(kind))
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
	return stringTitle(value, findingSeverityTitles)
}

func findingTypeTitle(value string) string {
	return stringTitle(value, findingTypeTitles)
}

func basisStatusTitle(value string) string {
	return stringTitle(value, basisStatusTitles)
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
