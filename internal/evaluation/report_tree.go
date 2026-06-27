package evaluation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
)

func evaluationRenderableGaps(runAbs string) []RunGap {
	raw, err := readRequiredEvaluationPayload(runAbs, "data/frame/evaluation-frame.json", DataKindEvaluationFrame)
	if err != nil {
		return []RunGap{*err}
	}
	artifacts, collectErr := collectEvaluationArtifacts(runAbs)
	if collectErr != nil {
		return []RunGap{{Kind: GapMalformedEvaluationData, Ref: "data", Detail: collectErr.Error()}}
	}
	artifacts.Frame = raw
	if _, gap := resolveEvaluationReportPlan(artifacts); gap != nil {
		return []RunGap{*gap}
	}
	return nil
}

func readRequiredEvaluationPayload(runAbs, rel string, want DataKind) (map[string]any, *RunGap) {
	raw, err := os.ReadFile(filepath.Join(runAbs, rel))
	if os.IsNotExist(err) {
		return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: rel, Detail: "required Evaluation evaluation payload is missing"}
	}
	if err != nil {
		return nil, &RunGap{Kind: GapUnreadableEvaluationData, Ref: rel, Detail: err.Error()}
	}
	payload, err := decodeDataPayload(raw)
	if err != nil {
		return nil, &RunGap{Kind: GapMalformedEvaluationData, Ref: rel, Detail: err.Error()}
	}
	kind, err := payloadKind(payload)
	if err != nil {
		return nil, &RunGap{Kind: GapIncompleteEvaluationData, Ref: rel, Detail: err.Error()}
	}
	if kind != want {
		return nil, &RunGap{Kind: GapIncompleteEvaluationData, Ref: rel, Detail: fmt.Sprintf("kind = %s, want %s", kind, want)}
	}
	if err := validateDataPayload(kind, payload); err != nil {
		return nil, &RunGap{Kind: GapIncompleteEvaluationData, Ref: rel, Detail: err.Error()}
	}
	return payload, nil
}

func buildEvaluationReport(path, displayPath string) (*BuildReportReceipt, error) {
	runAbs, err := verifyRun(path)
	if err != nil {
		return nil, err
	}
	if displayPath == "" {
		displayPath = displayRunPath(runAbs)
	}
	if gaps := evaluationRenderableGaps(runAbs); len(gaps) > 0 {
		return nil, nonReportableRunError(displayPath, gaps[0])
	}
	spec, err := loadRunModel(runAbs)
	if err != nil {
		return nil, err
	}
	artifacts, err := collectEvaluationArtifacts(runAbs)
	if err != nil {
		return nil, err
	}
	plan, gap := resolveEvaluationReportPlan(artifacts)
	if gap != nil {
		return nil, nonReportableRunError(displayPath, *gap)
	}
	reports := renderEvaluationReportTree(spec, artifacts, plan)
	for _, report := range reports {
		reportAbs := filepath.Join(runAbs, report.Path)
		if err := os.MkdirAll(filepath.Dir(reportAbs), 0o755); err != nil {
			return nil, fmt.Errorf("creating report directory: %w", err)
		}
		if err := writeReportFile(reportAbs, []byte(report.Content)); err != nil {
			return nil, err
		}
	}
	output := evaluationOutputResult(artifacts, plan, reports)
	outputRaw, err := canonicalJSON(output)
	if err != nil {
		return nil, err
	}
	outputRel := filepath.Join("data", "evaluation-output-result.json")
	if err := os.WriteFile(filepath.Join(runAbs, outputRel), outputRaw, 0o644); err != nil {
		return nil, fmt.Errorf("writing %s: %w", filepath.ToSlash(outputRel), err)
	}
	reportRel := "report.md"
	receipt := &BuildReportReceipt{
		SchemaVersion:          SchemaVersion,
		Path:                   displayPath,
		ReportMD:               filepath.ToSlash(filepath.Join(displayPath, reportRel)),
		EvaluationOutputResult: filepath.ToSlash(filepath.Join(displayPath, outputRel)),
		RatingResult:           evaluationReceiptRating(plan.HeadlineAnalysis),
	}
	if plan.HeadlineReport != nil {
		receipt.HeadlineReportMD = filepath.ToSlash(filepath.Join(displayPath, plan.HeadlineReport.Path))
	}
	if rootReport := reportForRootArea(reports); rootReport != nil {
		receipt.RootAreaReportMD = filepath.ToSlash(filepath.Join(displayPath, rootReport.Path))
	}
	return receipt, nil
}

func loadRunModel(runAbs string) (*model.Spec, error) {
	doc, err := document.Parse(filepath.Join(runAbs, ModelSnapshotFile))
	if err != nil {
		return nil, err
	}
	return model.Decode(doc)
}

func readJSONMap(path string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	return result, nil
}

type evaluationAreaArtifacts struct {
	ID       []string
	Frame    map[string]any
	Analysis map[string]any
}

type evaluationFactorArtifacts struct {
	ID       factorID
	Frame    map[string]any
	Analysis map[string]any
}

type evaluationRequirementArtifacts struct {
	ID         requirementID
	Frame      map[string]any
	Assessment map[string]any
	Rating     map[string]any
}

type evaluationArtifacts struct {
	Frame        map[string]any
	Areas        map[string]*evaluationAreaArtifacts
	Factors      map[string]*evaluationFactorArtifacts
	Requirements map[string]*evaluationRequirementArtifacts
}

type evaluationRenderedReport struct {
	Kind          string
	Path          string
	AreaID        []string
	FactorID      *factorID
	RequirementID *requirementID
	Content       string
}

type evaluationHeadlineKind string

const (
	evaluationHeadlineArea   evaluationHeadlineKind = "area"
	evaluationHeadlineFactor evaluationHeadlineKind = "factor"
)

type evaluationReportPlan struct {
	Frame            map[string]any
	HeadlineKind     evaluationHeadlineKind
	HeadlineAreaID   []string
	HeadlineFactorID *factorID
	HeadlineAnalysis map[string]any
	HeadlineReport   *evaluationRenderedReport
}

func collectEvaluationArtifacts(runAbs string) (*evaluationArtifacts, error) {
	out := &evaluationArtifacts{
		Areas:        map[string]*evaluationAreaArtifacts{},
		Factors:      map[string]*evaluationFactorArtifacts{},
		Requirements: map[string]*evaluationRequirementArtifacts{},
	}
	root := filepath.Join(runAbs, "data")
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		payload, err := readJSONMap(path)
		if err != nil {
			return nil
		}
		kind, err := payloadKind(payload)
		if err != nil {
			return nil
		}
		if collector := evaluationPayloadCollectors[kind]; collector != nil {
			return collector(out, payload)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("collecting Evaluation evaluation data: %w", err)
	}
	return out, nil
}

type evaluationPayloadCollector func(*evaluationArtifacts, map[string]any) error

var evaluationPayloadCollectors = map[DataKind]evaluationPayloadCollector{
	DataKindEvaluationFrame:            collectEvaluationFrame,
	DataKindAreaEvaluationFrame:        collectEvaluationAreaFrame,
	DataKindAreaAnalysis:               collectEvaluationAreaAnalysis,
	DataKindRequirementEvaluationFrame: collectEvaluationRequirementFrame,
	DataKindRequirementAssessment:      collectEvaluationRequirementAssessment,
	DataKindRequirementRating:          collectEvaluationRequirementRating,
	DataKindFactorAnalysisFrame:        collectEvaluationFactorFrame,
	DataKindFactorAnalysis:             collectEvaluationFactorAnalysis,
}

func collectEvaluationFrame(out *evaluationArtifacts, payload map[string]any) error {
	out.Frame = payload
	return nil
}

func collectEvaluationAreaFrame(out *evaluationArtifacts, payload map[string]any) error {
	id, err := subjectAreaID(payload)
	if err != nil {
		return err
	}
	out.area(areaKey(id)).Frame = payload
	return nil
}

func collectEvaluationAreaAnalysis(out *evaluationArtifacts, payload map[string]any) error {
	id, err := topAreaID(payload)
	if err != nil {
		return err
	}
	out.area(areaKey(id)).Analysis = payload
	return nil
}

func collectEvaluationRequirementFrame(out *evaluationArtifacts, payload map[string]any) error {
	id, err := subjectRequirementID(payload)
	if err != nil {
		return err
	}
	out.requirement(requirementKey(id)).Frame = payload
	return nil
}

func collectEvaluationRequirementAssessment(out *evaluationArtifacts, payload map[string]any) error {
	id, err := topRequirementID(payload)
	if err != nil {
		return err
	}
	out.requirement(requirementKey(id)).Assessment = payload
	return nil
}

func collectEvaluationRequirementRating(out *evaluationArtifacts, payload map[string]any) error {
	id, err := topRequirementID(payload)
	if err != nil {
		return err
	}
	out.requirement(requirementKey(id)).Rating = payload
	return nil
}

func collectEvaluationFactorFrame(out *evaluationArtifacts, payload map[string]any) error {
	id, err := subjectFactorID(payload)
	if err != nil {
		return err
	}
	out.factor(factorKey(id)).Frame = payload
	return nil
}

func collectEvaluationFactorAnalysis(out *evaluationArtifacts, payload map[string]any) error {
	id, err := topFactorID(payload)
	if err != nil {
		return err
	}
	out.factor(factorKey(id)).Analysis = payload
	return nil
}

func (a *evaluationArtifacts) area(key string) *evaluationAreaArtifacts {
	if existing, ok := a.Areas[key]; ok {
		return existing
	}
	id := areaIDFromKey(key)
	created := &evaluationAreaArtifacts{ID: id}
	a.Areas[key] = created
	return created
}

func (a *evaluationArtifacts) factor(key string) *evaluationFactorArtifacts {
	if existing, ok := a.Factors[key]; ok {
		return existing
	}
	id := factorIDFromKey(key)
	created := &evaluationFactorArtifacts{ID: id}
	a.Factors[key] = created
	return created
}

func (a *evaluationArtifacts) requirement(key string) *evaluationRequirementArtifacts {
	if existing, ok := a.Requirements[key]; ok {
		return existing
	}
	id := requirementIDFromKey(key)
	created := &evaluationRequirementArtifacts{ID: id}
	a.Requirements[key] = created
	return created
}

func renderEvaluationReportTree(spec *model.Spec, artifacts *evaluationArtifacts, plan *evaluationReportPlan) []evaluationRenderedReport {
	var reports []evaluationRenderedReport
	for _, area := range artifacts.sortedAreas() {
		if area.Analysis == nil {
			continue
		}
		path := areaReportPath(area.ID)
		reports = append(reports, evaluationRenderedReport{
			Kind:    string(ReportKindArea),
			Path:    path,
			AreaID:  copyStrings(area.ID),
			Content: renderEvaluationAreaReport(spec, artifacts, area, path),
		})
	}
	for _, factor := range artifacts.sortedFactors() {
		if factor.Analysis == nil {
			continue
		}
		id := factor.ID
		path := factorReportPath(id)
		reports = append(reports, evaluationRenderedReport{
			Kind:     string(ReportKindFactor),
			Path:     path,
			AreaID:   copyStrings(id.DeclaringArea),
			FactorID: &id,
			Content:  renderEvaluationFactorReport(spec, artifacts, factor, path),
		})
	}
	for _, requirement := range artifacts.sortedRequirements() {
		if requirement.Assessment == nil && requirement.Rating == nil {
			continue
		}
		id := requirement.ID
		path := requirementReportPath(id)
		reports = append(reports, evaluationRenderedReport{
			Kind:          string(ReportKindRequirement),
			Path:          path,
			AreaID:        copyStrings(id.DeclaringArea),
			RequirementID: &id,
			Content:       renderEvaluationRequirementReport(spec, artifacts, requirement, path),
		})
	}
	linkEvaluationReportPlan(plan, reports)
	run := evaluationRenderedReport{
		Kind:    string(ReportKindRun),
		Path:    "report.md",
		Content: renderEvaluationRunReport(spec, artifacts, plan, reports, "report.md"),
	}
	return append([]evaluationRenderedReport{run}, reports...)
}

func renderEvaluationRunReport(spec *model.Spec, artifacts *evaluationArtifacts, plan *evaluationReportPlan, reports []evaluationRenderedReport, reportPath string) string {
	headlineLabel := evaluationHeadlineLabel(spec, plan)
	headlineScope := scopedMap(plan.HeadlineAnalysis, "localAndDescendantAnalysis")
	var b strings.Builder
	b.WriteString("# Evaluation Report: " + headlineLabel + "\n\n")
	b.WriteString("| Headline Rating | Headline Subject | Data |\n")
	b.WriteString("| --- | --- | --- |\n")
	b.WriteString("| " + markdownCell(evaluationRatingLabel(spec, headlineScope)) + " | " + markdownCell(evaluationHeadlineReportLink(spec, reportPath, plan)) + " | " + reportDataLink(reportPath, "data/evaluation-output-result.json") + " |\n\n")
	b.WriteString("Summary:\n\n")
	b.WriteString(evaluationSummary(headlineScope))
	b.WriteString("\n\n## Scope\n\n")
	writeEvaluationRunScope(&b, artifacts)
	b.WriteString("## Subject Reports\n\n")
	writeEvaluationRunReportsTable(&b, spec, artifacts, reports, reportPath)
	b.WriteString("## Coverage\n\n")
	writeEvaluationRunCoverage(&b, artifacts, reports, reportPath)
	b.WriteString("## Limits & Incomplete Inputs\n\n")
	writeEvaluationLimitsTable(&b, scopedMap(plan.HeadlineAnalysis, "localAnalysis"), headlineScope)
	writeEvaluationLegend(&b)
	return b.String()
}

func renderEvaluationAreaReport(spec *model.Spec, artifacts *evaluationArtifacts, area *evaluationAreaArtifacts, reportPath string) string {
	title := areaTitle(spec, area.ID)
	local := scopedMap(area.Analysis, "localAnalysis")
	overall := scopedMap(area.Analysis, "localAndDescendantAnalysis")
	var b strings.Builder
	b.WriteString("# Area: " + title + "\n\n")
	writeEvaluationAreaTrail(&b, spec, artifacts, area.ID, reportPath)
	b.WriteString("| Overall Rating | Local Rating | Confidence | Data |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	b.WriteString("| " + markdownCell(evaluationRatingLabel(spec, overall)) + " | " + markdownCell(evaluationRatingLabel(spec, local)) + " | " + markdownCell(evaluationConfidencePair(overall, local)) + " | " + reportDataLink(reportPath, areaDataPath(area.ID, "area-analysis-result.json")) + " |\n\n")
	b.WriteString("Summary:\n\n")
	b.WriteString(evaluationSummary(overall))
	b.WriteString("\n\n## Findings\n\n")
	writeEvaluationAreaFindingsTable(&b, spec, area.Analysis, area.ID, reportPath)
	b.WriteString("## Rating Drivers\n\n")
	writeEvaluationDriversTable(&b, spec, overall)
	b.WriteString("## Factors\n\n")
	b.WriteString("| Factor | Path | Local Rating | + Sub-Factors Rating | Sub-Factors |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	if factors := artifacts.rootFactorsForArea(area.ID); len(factors) == 0 {
		b.WriteString("| (no local Factors) |  |  |  |  |\n\n")
	} else {
		for _, factor := range factors {
			factorLocal := scopedMap(factor.Analysis, "localAnalysis")
			factorOverall := scopedMap(factor.Analysis, "localAndDescendantAnalysis")
			children := artifacts.childFactors(factor.ID)
			b.WriteString("| " + reportLink(reportPath, factorReportPath(factor.ID), factorTitle(spec, factor.ID)) + " | `" + factorDisplayPath(factor.ID) + "` | " + markdownCell(evaluationRatingLabel(spec, factorLocal)) + " | " + markdownCell(evaluationSubRatingCell(spec, factorOverall, len(children) > 0)) + " | " + markdownCell(factorList(children, spec, reportPath)) + " |\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("## Sub-Areas\n\n")
	b.WriteString("| Area | Path | Local Rating | + Sub-Areas Rating | Factors |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	if children := artifacts.childAreas(area.ID); len(children) == 0 {
		b.WriteString("| (no child Areas) |  |  |  |  |\n\n")
	} else {
		for _, child := range children {
			childLocal := scopedMap(child.Analysis, "localAnalysis")
			childOverall := scopedMap(child.Analysis, "localAndDescendantAnalysis")
			b.WriteString("| " + reportLink(reportPath, areaReportPath(child.ID), areaTitle(spec, child.ID)) + " | `" + areaDisplayPath(child.ID) + "` | " + markdownCell(evaluationRatingLabel(spec, childLocal)) + " | " + markdownCell(evaluationSubRatingCell(spec, childOverall, len(artifacts.childAreas(child.ID)) > 0)) + " | " + markdownCell(factorList(artifacts.rootFactorsForArea(child.ID), spec, reportPath)) + " |\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("## Requirements\n\n")
	b.WriteString("| Requirement | Rating | Status | Factors |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	if requirements := artifacts.requirementsForArea(area.ID); len(requirements) == 0 {
		b.WriteString("| (no local Requirements) |  |  |  |\n\n")
	} else {
		for _, req := range requirements {
			b.WriteString("| " + reportLink(reportPath, requirementReportPath(req.ID), requirementTitle(spec, req.ID)) + " | " + markdownCell(evaluationRequirementRatingLabel(spec, req.Rating)) + " | " + markdownCell(assessmentStatusTitle(evaluationString(req.Assessment, "status"))) + " | " + markdownCell(requirementFactorLinks(req, reportPath)) + " |\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("## Limits & Incomplete Inputs\n\n")
	writeEvaluationLimitsTable(&b, local, overall)
	writeEvaluationLegend(&b)
	return b.String()
}

func renderEvaluationFactorReport(spec *model.Spec, artifacts *evaluationArtifacts, factor *evaluationFactorArtifacts, reportPath string) string {
	local := scopedMap(factor.Analysis, "localAnalysis")
	overall := scopedMap(factor.Analysis, "localAndDescendantAnalysis")
	title := factorTitle(spec, factor.ID)
	var b strings.Builder
	b.WriteString("# Factor: " + title + "\n\n")
	writeEvaluationAreaTrail(&b, spec, artifacts, factor.ID.DeclaringArea, reportPath)
	writeEvaluationFactorTrail(&b, spec, factor.ID, reportPath)
	b.WriteString("| Overall Rating | Local Rating | Status | Confidence | Data |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	b.WriteString("| " + markdownCell(evaluationRatingLabel(spec, overall)) + " | " + markdownCell(evaluationRatingLabel(spec, local)) + " | " + markdownCell(evaluationAnalysisStatusPair(overall, local)) + " | " + markdownCell(evaluationConfidencePair(overall, local)) + " | " + reportDataLink(reportPath, factorDataPath(factor.ID, "factor-analysis-result.json")) + " |\n\n")
	b.WriteString("Summary:\n\n")
	b.WriteString(evaluationSummary(overall))
	b.WriteString("\n\n## Findings\n\n")
	writeEvaluationFactorFindingsTable(&b, spec, artifacts, factor.ID, reportPath)
	b.WriteString("## Rating Drivers\n\n")
	writeEvaluationDriversTable(&b, spec, overall)
	b.WriteString("## Requirements\n\n")
	b.WriteString("| Requirement | Rating | Status |\n")
	b.WriteString("| --- | --- | --- |\n")
	if requirements := artifacts.requirementsForFactor(factor.ID); len(requirements) == 0 {
		b.WriteString("| (no direct Requirements) |  |  |\n\n")
	} else {
		for _, req := range requirements {
			b.WriteString("| " + reportLink(reportPath, requirementReportPath(req.ID), requirementTitle(spec, req.ID)) + " | " + markdownCell(evaluationRequirementRatingLabel(spec, req.Rating)) + " | " + markdownCell(assessmentStatusTitle(evaluationString(req.Assessment, "status"))) + " |\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("## Child Factors\n\n")
	b.WriteString("| Factor | Path | Local Rating | + Sub-Factors Rating |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	if children := artifacts.childFactors(factor.ID); len(children) == 0 {
		b.WriteString("| (no child Factors) |  |  |  |\n\n")
	} else {
		for _, child := range children {
			childLocal := scopedMap(child.Analysis, "localAnalysis")
			childOverall := scopedMap(child.Analysis, "localAndDescendantAnalysis")
			b.WriteString("| " + reportLink(reportPath, factorReportPath(child.ID), factorTitle(spec, child.ID)) + " | `" + factorDisplayPath(child.ID) + "` | " + markdownCell(evaluationRatingLabel(spec, childLocal)) + " | " + markdownCell(evaluationSubRatingCell(spec, childOverall, len(artifacts.childFactors(child.ID)) > 0)) + " |\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("## Limits & Incomplete Inputs\n\n")
	writeEvaluationLimitsTable(&b, local, overall)
	writeEvaluationLegend(&b)
	return b.String()
}

func renderEvaluationRequirementReport(spec *model.Spec, artifacts *evaluationArtifacts, req *evaluationRequirementArtifacts, reportPath string) string {
	title := requirementTitle(spec, req.ID)
	var b strings.Builder
	b.WriteString("# Requirement: " + title + "\n\n")
	writeEvaluationAreaTrail(&b, spec, artifacts, req.ID.DeclaringArea, reportPath)
	writeEvaluationRequirementFactorsLine(&b, req, reportPath)
	b.WriteString("| Rating | Assessment | Confidence | Data |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	b.WriteString("| " + markdownCell(evaluationRequirementRatingLabel(spec, req.Rating)) + " | " + markdownCell(assessmentStatusTitle(evaluationString(req.Assessment, "status"))) + " | " + markdownCell(evaluationRequirementConfidencePair(req)) + " | " + reportDataLink(reportPath, requirementDataPath(req.ID, "requirement-assessment-result.json")) + ", " + reportDataLink(reportPath, requirementDataPath(req.ID, "requirement-rating-result.json")) + " |\n\n")
	b.WriteString("Summary:\n\n")
	if summary := evaluationString(req.Assessment, "evidenceSummary"); summary != "" {
		b.WriteString(summary)
	} else if rationale := evaluationString(req.Rating, "rationale"); rationale != "" {
		b.WriteString(rationale)
	} else {
		b.WriteString("No assessment summary was recorded.")
	}
	b.WriteString("\n\n## Findings Summary\n\n")
	writeEvaluationFindingsTable(&b, req.Assessment)
	b.WriteString("## Finding Details\n\n")
	writeEvaluationFindingDetails(&b, req.Assessment)
	b.WriteString("## Unknowns & Missing Evidence\n\n")
	writeEvaluationUnknownsTable(&b, req.Assessment, req.Rating)
	writeEvaluationLegend(&b)
	_ = artifacts
	return b.String()
}

func evaluationOutputResult(artifacts *evaluationArtifacts, plan *evaluationReportPlan, reports []evaluationRenderedReport) map[string]any {
	reportOutputs := make([]any, 0, len(reports))
	reportsByArea := map[string][]any{}
	for _, report := range reports {
		ref := evaluationReportRef(report)
		reportOutputs = append(reportOutputs, ref)
		if report.Kind != string(ReportKindRun) {
			reportsByArea[areaKey(report.AreaID)] = append(reportsByArea[areaKey(report.AreaID)], ref)
		}
	}
	areaOutputs := []any{}
	for _, area := range artifacts.sortedAreas() {
		if area.Analysis == nil {
			continue
		}
		areaID := model.AreaPath(area.ID).Reference()
		areaOutputs = append(areaOutputs, map[string]any{
			"areaId":                    areaID,
			"areaEvaluationFrameRef":    routineRef(DataKindAreaEvaluationFrame, map[string]any{"areaId": areaID}, ""),
			"areaAnalysisResultRef":     routineRef(DataKindAreaAnalysis, map[string]any{"areaId": areaID}, ""),
			"factorAnalysisRefs":        factorAnalysisRefs(artifacts.rootFactorsForArea(area.ID)),
			"requirementAssessmentRefs": requirementAssessmentRefs(artifacts.requirementsForArea(area.ID)),
			"requirementRatingRefs":     requirementRatingRefs(artifacts.requirementsForArea(area.ID)),
			"reportRefs":                reportsByArea[areaKey(area.ID)],
		})
	}
	output := map[string]any{
		"schemaVersion":     SchemaVersion,
		"kind":              string(DataKindEvaluationOutput),
		"runReportRef":      evaluationReportRef(evaluationRenderedReport{Kind: string(ReportKindRun), Path: "report.md"}),
		"headlineResultRef": evaluationHeadlineResultRef(plan),
		"headlineReportRef": evaluationReportRef(*plan.HeadlineReport),
		"areaOutputs":       areaOutputs,
		"reportOutputs":     reportOutputs,
	}
	if root := artifacts.Areas[areaKey(nil)]; root != nil && root.Analysis != nil {
		output["rootAreaAnalysisRef"] = routineRef(DataKindAreaAnalysis, map[string]any{"areaId": model.AreaPath{}.Reference()}, "localAndDescendantAnalysis")
	}
	return output
}

func evaluationReportRef(report evaluationRenderedReport) map[string]any {
	ref := map[string]any{"kind": report.Kind, "path": filepath.ToSlash(report.Path)}
	if report.Kind != string(ReportKindRun) {
		ref["areaId"] = model.AreaPath(report.AreaID).Reference()
	}
	if report.FactorID != nil {
		ref["factorId"] = factorIDJSON(*report.FactorID)
	}
	if report.RequirementID != nil {
		ref["requirementId"] = requirementIDJSON(*report.RequirementID)
	}
	return ref
}

func resolveEvaluationReportPlan(artifacts *evaluationArtifacts) (*evaluationReportPlan, *RunGap) {
	if artifacts.Frame == nil {
		return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: "data/frame/evaluation-frame.json", Detail: "required Evaluation evaluation payload is missing"}
	}
	inputs := objectMap(artifacts.Frame["inputs"])
	if plan, gap, scoped := resolveScopedFactorHeadline(artifacts, stringValues(inputs["factorIds"])); scoped {
		return plan, gap
	}
	if plan, gap, scoped := resolveScopedAreaHeadline(artifacts, stringValues(inputs["areaIds"])); scoped {
		return plan, gap
	}
	return resolveRootHeadline(artifacts)
}

func resolveScopedFactorHeadline(artifacts *evaluationArtifacts, refs []string) (*evaluationReportPlan, *RunGap, bool) {
	var firstMissingArea string
	for _, ref := range refs {
		id, err := factorIDFrom(ref)
		if err != nil {
			return nil, &RunGap{Kind: GapIncompleteEvaluationData, Ref: "data/frame/evaluation-frame.json", Detail: err.Error()}, true
		}
		path := factorDataPath(id, "factor-analysis-result.json")
		if firstMissingArea == "" {
			firstMissingArea = path
		}
		if factor := artifacts.Factors[factorKey(id)]; factor != nil && factor.Analysis != nil {
			return &evaluationReportPlan{
				Frame:            artifacts.Frame,
				HeadlineKind:     evaluationHeadlineFactor,
				HeadlineAreaID:   copyStrings(id.DeclaringArea),
				HeadlineFactorID: &id,
				HeadlineAnalysis: factor.Analysis,
			}, nil, true
		}
	}
	if firstMissingArea == "" {
		return nil, nil, false
	}
	return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: firstMissingArea, Detail: "required scoped Factor analysis payload is missing"}, true
}

func resolveScopedAreaHeadline(artifacts *evaluationArtifacts, refs []string) (*evaluationReportPlan, *RunGap, bool) {
	var firstMissingArea string
	for _, ref := range refs {
		id, err := areaIDFrom(ref)
		if err != nil {
			return nil, &RunGap{Kind: GapIncompleteEvaluationData, Ref: "data/frame/evaluation-frame.json", Detail: err.Error()}, true
		}
		path := areaDataPath(id, "area-analysis-result.json")
		if firstMissingArea == "" {
			firstMissingArea = path
		}
		if area := artifacts.Areas[areaKey(id)]; area != nil && area.Analysis != nil {
			return &evaluationReportPlan{
				Frame:            artifacts.Frame,
				HeadlineKind:     evaluationHeadlineArea,
				HeadlineAreaID:   copyStrings(id),
				HeadlineAnalysis: area.Analysis,
			}, nil, true
		}
	}
	if firstMissingArea == "" {
		return nil, nil, false
	}
	return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: firstMissingArea, Detail: "required scoped Area analysis payload is missing"}, true
}

func resolveRootHeadline(artifacts *evaluationArtifacts) (*evaluationReportPlan, *RunGap) {
	root := artifacts.Areas[areaKey(nil)]
	if root == nil || root.Analysis == nil {
		return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: "data/areas/root/area-analysis-result.json", Detail: "required Evaluation evaluation payload is missing"}
	}
	return &evaluationReportPlan{
		Frame:            artifacts.Frame,
		HeadlineKind:     evaluationHeadlineArea,
		HeadlineAreaID:   nil,
		HeadlineAnalysis: root.Analysis,
	}, nil
}

func linkEvaluationReportPlan(plan *evaluationReportPlan, reports []evaluationRenderedReport) {
	for i := range reports {
		report := &reports[i]
		switch plan.HeadlineKind {
		case evaluationHeadlineArea:
			if report.Kind == string(ReportKindArea) && sameStrings(report.AreaID, plan.HeadlineAreaID) {
				plan.HeadlineReport = report
				return
			}
		case evaluationHeadlineFactor:
			if report.Kind == string(ReportKindFactor) && plan.HeadlineFactorID != nil && report.FactorID != nil && sameStrings(report.FactorID.DeclaringArea, plan.HeadlineFactorID.DeclaringArea) && sameStrings(report.FactorID.Path, plan.HeadlineFactorID.Path) {
				plan.HeadlineReport = report
				return
			}
		}
	}
}

func evaluationHeadlineResultRef(plan *evaluationReportPlan) map[string]any {
	switch plan.HeadlineKind {
	case evaluationHeadlineFactor:
		return routineRef(DataKindFactorAnalysis, map[string]any{"factorId": factorIDJSON(*plan.HeadlineFactorID)}, "localAndDescendantAnalysis")
	default:
		return routineRef(DataKindAreaAnalysis, map[string]any{"areaId": model.AreaPath(plan.HeadlineAreaID).Reference()}, "localAndDescendantAnalysis")
	}
}

func reportForRootArea(reports []evaluationRenderedReport) *evaluationRenderedReport {
	for i := range reports {
		if reports[i].Kind == string(ReportKindArea) && len(reports[i].AreaID) == 0 {
			return &reports[i]
		}
	}
	return nil
}

func evaluationHeadlineLabel(spec *model.Spec, plan *evaluationReportPlan) string {
	switch plan.HeadlineKind {
	case evaluationHeadlineFactor:
		return "Factor: " + factorTitle(spec, *plan.HeadlineFactorID)
	default:
		return "Area: " + areaTitle(spec, plan.HeadlineAreaID)
	}
}

func evaluationHeadlineReportLink(spec *model.Spec, reportPath string, plan *evaluationReportPlan) string {
	label := evaluationHeadlineLabel(spec, plan)
	if plan.HeadlineReport == nil {
		return label
	}
	return reportLink(reportPath, plan.HeadlineReport.Path, label)
}

func writeEvaluationRunScope(b *strings.Builder, artifacts *evaluationArtifacts) {
	inputs := objectMap(artifacts.Frame["inputs"])
	derived := objectMap(artifacts.Frame["derivedContext"])
	b.WriteString("| Field | Value |\n")
	b.WriteString("| --- | --- |\n")
	b.WriteString("| Requested Scope | " + markdownCell(firstString(inputs, "requestedScope")) + " |\n")
	b.WriteString("| Resolved Scope | " + markdownCell(firstString(derived, "resolvedScope")) + " |\n")
	b.WriteString("| Areas | " + markdownCell(strings.Join(stringValues(inputs["areaIds"]), "; ")) + " |\n")
	b.WriteString("| Factors | " + markdownCell(strings.Join(stringValues(inputs["factorIds"]), "; ")) + " |\n\n")
}

func writeEvaluationRunReportsTable(b *strings.Builder, spec *model.Spec, artifacts *evaluationArtifacts, reports []evaluationRenderedReport, reportPath string) {
	b.WriteString("| Subject | Kind | Rating | Report |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	if len(reports) == 0 {
		b.WriteString("| (no subject reports) |  |  |  |\n\n")
		return
	}
	for _, report := range reports {
		if report.Kind == string(ReportKindRun) {
			continue
		}
		b.WriteString("| " + reportLink(reportPath, report.Path, reportSubjectTitle(spec, report)) + " | " + markdownCell(reportKindTitle(report.Kind)) + " | " + markdownCell(reportSubjectRating(spec, artifacts, report)) + " | " + reportLink(reportPath, report.Path, filepath.Base(report.Path)) + " |\n")
	}
	b.WriteString("\n")
}

func writeEvaluationRunCoverage(b *strings.Builder, artifacts *evaluationArtifacts, reports []evaluationRenderedReport, reportPath string) {
	if root := reportForRootArea(reports); root != nil {
		b.WriteString("- Root Area report: " + reportLink(reportPath, root.Path, filepath.Base(root.Path)) + "\n")
	} else {
		b.WriteString("- Root Area was not evaluated in this run.\n")
	}
	fmt.Fprintf(b, "- Generated subject reports: %d\n\n", len(reports))
	_ = artifacts
}

func reportSubjectTitle(spec *model.Spec, report evaluationRenderedReport) string {
	switch report.Kind {
	case string(ReportKindArea):
		return areaTitle(spec, report.AreaID)
	case string(ReportKindFactor):
		return factorTitle(spec, *report.FactorID)
	case string(ReportKindRequirement):
		return requirementTitle(spec, *report.RequirementID)
	default:
		return "Evaluation Report"
	}
}

func reportSubjectRating(spec *model.Spec, artifacts *evaluationArtifacts, report evaluationRenderedReport) string {
	switch report.Kind {
	case string(ReportKindArea):
		if area := artifacts.Areas[areaKey(report.AreaID)]; area != nil {
			return evaluationRatingLabel(spec, scopedMap(area.Analysis, "localAndDescendantAnalysis"))
		}
	case string(ReportKindFactor):
		if report.FactorID != nil {
			if factor := artifacts.Factors[factorKey(*report.FactorID)]; factor != nil {
				return evaluationRatingLabel(spec, scopedMap(factor.Analysis, "localAndDescendantAnalysis"))
			}
		}
	case string(ReportKindRequirement):
		if report.RequirementID != nil {
			if req := artifacts.Requirements[requirementKey(*report.RequirementID)]; req != nil {
				return evaluationRequirementRatingLabel(spec, req.Rating)
			}
		}
	default:
		return "—"
	}
	return "—"
}

func factorAnalysisRefs(factors []*evaluationFactorArtifacts) []any {
	refs := make([]any, 0, len(factors))
	for _, factor := range factors {
		refs = append(refs, routineRef(DataKindFactorAnalysis, map[string]any{"factorId": factorIDJSON(factor.ID)}, "localAndDescendantAnalysis"))
	}
	return refs
}

func requirementAssessmentRefs(requirements []*evaluationRequirementArtifacts) []any {
	refs := make([]any, 0, len(requirements))
	for _, requirement := range requirements {
		refs = append(refs, routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": requirementIDJSON(requirement.ID)}, ""))
	}
	return refs
}

func requirementRatingRefs(requirements []*evaluationRequirementArtifacts) []any {
	refs := make([]any, 0, len(requirements))
	for _, requirement := range requirements {
		refs = append(refs, routineRef(DataKindRequirementRating, map[string]any{"requirementId": requirementIDJSON(requirement.ID)}, ""))
	}
	return refs
}

func (a *evaluationArtifacts) sortedAreas() []*evaluationAreaArtifacts {
	keys := make([]string, 0, len(a.Areas))
	for key := range a.Areas {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i] == "" {
			return true
		}
		if keys[j] == "" {
			return false
		}
		return keys[i] < keys[j]
	})
	out := make([]*evaluationAreaArtifacts, 0, len(keys))
	for _, key := range keys {
		out = append(out, a.Areas[key])
	}
	return out
}

func (a *evaluationArtifacts) sortedFactors() []*evaluationFactorArtifacts {
	keys := make([]string, 0, len(a.Factors))
	for key := range a.Factors {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]*evaluationFactorArtifacts, 0, len(keys))
	for _, key := range keys {
		out = append(out, a.Factors[key])
	}
	return out
}

func (a *evaluationArtifacts) sortedRequirements() []*evaluationRequirementArtifacts {
	keys := make([]string, 0, len(a.Requirements))
	for key := range a.Requirements {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]*evaluationRequirementArtifacts, 0, len(keys))
	for _, key := range keys {
		out = append(out, a.Requirements[key])
	}
	return out
}

func (a *evaluationArtifacts) childAreas(parent []string) []*evaluationAreaArtifacts {
	var out []*evaluationAreaArtifacts
	for _, area := range a.sortedAreas() {
		if len(area.ID) == len(parent)+1 && sameStrings(area.ID[:len(parent)], parent) && area.Analysis != nil {
			out = append(out, area)
		}
	}
	return out
}

func (a *evaluationArtifacts) rootFactorsForArea(areaID []string) []*evaluationFactorArtifacts {
	var out []*evaluationFactorArtifacts
	for _, factor := range a.sortedFactors() {
		if sameStrings(factor.ID.DeclaringArea, areaID) && len(factor.ID.Path) == 1 && factor.Analysis != nil {
			out = append(out, factor)
		}
	}
	return out
}

func (a *evaluationArtifacts) childFactors(parent factorID) []*evaluationFactorArtifacts {
	var out []*evaluationFactorArtifacts
	for _, factor := range a.sortedFactors() {
		if sameStrings(factor.ID.DeclaringArea, parent.DeclaringArea) && len(factor.ID.Path) == len(parent.Path)+1 && sameStrings(factor.ID.Path[:len(parent.Path)], parent.Path) && factor.Analysis != nil {
			out = append(out, factor)
		}
	}
	return out
}

func (a *evaluationArtifacts) requirementsForArea(areaID []string) []*evaluationRequirementArtifacts {
	var out []*evaluationRequirementArtifacts
	for _, req := range a.sortedRequirements() {
		if sameStrings(req.ID.DeclaringArea, areaID) {
			out = append(out, req)
		}
	}
	return out
}

func (a *evaluationArtifacts) requirementsForFactor(factor factorID) []*evaluationRequirementArtifacts {
	var out []*evaluationRequirementArtifacts
	want := factorDisplayPath(factor)
	for _, req := range a.sortedRequirements() {
		if !sameStrings(req.ID.DeclaringArea, factor.DeclaringArea) {
			continue
		}
		for _, linked := range requirementFactorIDs(req) {
			parsed, err := parseRequirementFactorID(req.ID.DeclaringArea, linked)
			if err == nil && sameStrings(parsed.DeclaringArea, factor.DeclaringArea) && sameStrings(parsed.Path, factor.Path) {
				out = append(out, req)
				break
			}
			if linked == want || linked == strings.Join(factor.Path, "/") || linked == factor.Path[len(factor.Path)-1] {
				out = append(out, req)
				break
			}
		}
	}
	return out
}

func writeEvaluationAreaTrail(b *strings.Builder, spec *model.Spec, artifacts *evaluationArtifacts, areaID []string, reportPath string) {
	parts := []string{areaTrailPart(spec, artifacts, nil, reportPath)}
	for i := range areaID {
		id := areaID[:i+1]
		parts = append(parts, areaTrailPart(spec, artifacts, id, reportPath))
	}
	b.WriteString("Area: " + strings.Join(parts, " / ") + "\n\n")
}

func areaTrailPart(spec *model.Spec, artifacts *evaluationArtifacts, areaID []string, reportPath string) string {
	title := areaTitle(spec, areaID)
	if area := artifacts.Areas[areaKey(areaID)]; area != nil && area.Analysis != nil {
		return reportLink(reportPath, areaReportPath(areaID), title)
	}
	return markdownCell(title)
}

func writeEvaluationFactorTrail(b *strings.Builder, spec *model.Spec, factor factorID, reportPath string) {
	parts := make([]string, 0, len(factor.Path))
	for i := range factor.Path {
		id := factorID{DeclaringArea: factor.DeclaringArea, Path: factor.Path[:i+1]}
		parts = append(parts, reportLink(reportPath, factorReportPath(id), factorTitle(spec, id)))
	}
	b.WriteString("Factor: " + strings.Join(parts, " / ") + "\n\n")
}

func writeEvaluationRequirementFactorsLine(b *strings.Builder, req *evaluationRequirementArtifacts, reportPath string) {
	links := requirementFactorLinks(req, reportPath)
	if links == "" {
		links = "(none)"
	}
	b.WriteString("Factors: " + links + "\n\n")
}

func writeEvaluationDriversTable(b *strings.Builder, spec *model.Spec, scope map[string]any) {
	b.WriteString("| Driver | Effect | Inputs |\n")
	b.WriteString("| --- | --- | --- |\n")
	drivers := objectSlice(scope["ratingDrivers"])
	if len(drivers) == 0 {
		b.WriteString("| (no rating drivers) |  |  |\n\n")
		return
	}
	for _, driver := range drivers {
		effect := firstString(driver, "effect", "impact")
		if effect == "" {
			effect = ratingTitle(spec, firstString(driver, "ratingLevelId"))
		}
		b.WriteString("| " + markdownCell(firstString(driver, "description", "summary", "requirementRatingDriver")) + " | " + markdownCell(effect) + " | " + markdownCell(compactJSON(driver["inputRefs"])) + " |\n")
	}
	b.WriteString("\n")
}

func writeEvaluationLimitsTable(b *strings.Builder, scopes ...map[string]any) {
	b.WriteString("| Type | Scope | Impact |\n")
	b.WriteString("| --- | --- | --- |\n")
	wrote := false
	for _, scope := range scopes {
		for _, field := range []string{"incompleteInputs", "evaluationLimits"} {
			for _, item := range objectSlice(scope[field]) {
				wrote = true
				b.WriteString("| " + markdownCell(limitTypeTitle(field)) + " | " + markdownCell(firstString(item, "scope", "ref", "id")) + " | " + markdownCell(firstString(item, "impact", "description", "reason")) + " |\n")
			}
		}
	}
	if !wrote {
		b.WriteString("| (no limits or incomplete inputs) |  |  |\n")
	}
}

type indexedAreaFinding struct {
	Finding      map[string]any
	Index        int
	Relationship map[string]any
}

func writeEvaluationAreaFindingsTable(b *strings.Builder, spec *model.Spec, analysis map[string]any, areaID []string, reportPath string) {
	b.WriteString("| ID | Statement | Type | Severity | Confidence | Effect | Cause |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
	findings := sortedAreaFindings(analysis)
	if len(findings) == 0 {
		b.WriteString("| (no findings) |  |  |  |  |  |  |\n\n")
		return
	}
	for _, item := range findings {
		finding := item.Finding
		b.WriteString("| `" + markdownCell(areaFindingID(finding, item.Index)) + "` | " + markdownCell(firstString(finding, "statement")) + " | " + markdownCell(findingTypeTitle(firstString(finding, "type"))) + " | " + markdownCell(findingSeverityTitle(firstString(finding, "severity"))) + " | " + markdownCell(confidenceTitle(firstString(finding, "confidence"))) + " | " + markdownCell(findingEffectSummary(finding)) + " | " + markdownCell(findingCauseSummary(finding)) + " |\n")
	}
	b.WriteString("\n")
	writeAreaFindingDetails(b, spec, findings, areaID, reportPath)
}

func writeEvaluationFactorFindingsTable(b *strings.Builder, spec *model.Spec, artifacts *evaluationArtifacts, factor factorID, reportPath string) {
	b.WriteString("| ID | Statement | Type | Severity | Confidence | Effect | Cause |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
	area := artifacts.area(areaKey(factor.DeclaringArea))
	findings := sortedFactorFindings(area.Analysis, factor)
	if len(findings) == 0 {
		b.WriteString("| (no findings) |  |  |  |  |  |  |\n\n")
		return
	}
	for _, item := range findings {
		finding := item.Finding
		b.WriteString("| `" + markdownCell(areaFindingID(finding, item.Index)) + "` | " + markdownCell(firstString(finding, "statement")) + " | " + markdownCell(findingTypeTitle(firstString(finding, "type"))) + " | " + markdownCell(findingSeverityTitle(firstString(finding, "severity"))) + " | " + markdownCell(confidenceTitle(firstString(finding, "confidence"))) + " | " + markdownCell(findingEffectSummary(finding)) + " | " + markdownCell(findingCauseSummary(finding)) + " |\n")
	}
	b.WriteString("\n")
	writeAreaFindingDetails(b, spec, findings, factor.DeclaringArea, reportPath)
}

func writeAreaFindingDetails(b *strings.Builder, spec *model.Spec, findings []indexedAreaFinding, areaID []string, reportPath string) {
	b.WriteString("### Finding Details\n\n")
	if len(findings) == 0 {
		b.WriteString("(no finding details)\n\n")
		return
	}
	for _, item := range findings {
		finding := item.Finding
		writeFindingCoreDetails(b, 4, areaFindingID(finding, item.Index), finding)
		b.WriteString("##### Relationships\n\n")
		if relationships := areaFindingFactorLinks(spec, finding, areaID, reportPath); relationships != "(none)" {
			b.WriteString(relationships + "\n\n")
		} else {
			b.WriteString("(none)\n\n")
		}
		b.WriteString("##### Inputs\n\n")
		b.WriteString(compactJSON(finding["inputRefs"]) + "\n\n")
	}
}

func writeFindingCoreDetails(b *strings.Builder, headingLevel int, id string, finding map[string]any) {
	heading := strings.Repeat("#", headingLevel)
	title := id
	if statement := firstString(finding, "statement"); statement != "" {
		title += " " + statement
	}
	b.WriteString(heading + " " + title + "\n\n")
	writeFindingSection(b, headingLevel+1, "Condition", firstString(finding, "condition"))
	writeFindingCriteriaSection(b, headingLevel+1, finding)
	writeFindingCauseSection(b, headingLevel+1, finding)
	writeFindingEffectSection(b, headingLevel+1, finding)
	writeFindingEvidenceSection(b, headingLevel+1, "Evidence", objectSlice(finding["evidence"]))
}

func writeFindingSection(b *strings.Builder, headingLevel int, title, body string) {
	b.WriteString(strings.Repeat("#", headingLevel) + " " + title + "\n\n")
	if body == "" {
		b.WriteString("(not recorded)\n\n")
		return
	}
	b.WriteString(body + "\n\n")
}

func writeFindingCriteriaSection(b *strings.Builder, headingLevel int, finding map[string]any) {
	b.WriteString(strings.Repeat("#", headingLevel) + " Criteria\n\n")
	criteria := objectSlice(finding["criteria"])
	if len(criteria) == 0 {
		b.WriteString("(none recorded)\n\n")
		return
	}
	for _, criterion := range criteria {
		label := firstString(criterion, "requirementId")
		if rating := firstString(criterion, "ratingLevelId"); rating != "" {
			label += " / " + rating
		}
		b.WriteString("- `" + markdownCell(label) + "`: " + markdownCell(firstString(criterion, "criterion")) + "\n")
		if rationale := firstString(criterion, "rationale"); rationale != "" {
			b.WriteString("  Rationale: " + markdownCell(rationale) + "\n")
		}
	}
	b.WriteString("\n")
}

func writeFindingCauseSection(b *strings.Builder, headingLevel int, finding map[string]any) {
	b.WriteString(strings.Repeat("#", headingLevel) + " Cause\n\n")
	cause := objectMap(finding["cause"])
	if len(cause) == 0 {
		b.WriteString("(not recorded)\n\n")
		return
	}
	if status := firstString(cause, "status"); status != "" {
		b.WriteString("Status: " + causeStatusTitle(status) + "\n\n")
	}
	if statement := firstString(cause, "statement"); statement != "" {
		b.WriteString(statement + "\n\n")
	}
	if rationale := firstString(cause, "rationale"); rationale != "" {
		b.WriteString("Rationale: " + rationale + "\n\n")
	}
	writeFindingEvidenceSection(b, headingLevel+1, "Cause Evidence", objectSlice(cause["evidence"]))
}

func writeFindingEffectSection(b *strings.Builder, headingLevel int, finding map[string]any) {
	b.WriteString(strings.Repeat("#", headingLevel) + " Effect\n\n")
	effect := objectMap(finding["effect"])
	if len(effect) == 0 {
		b.WriteString("(not recorded)\n\n")
		return
	}
	if statement := firstString(effect, "statement"); statement != "" {
		b.WriteString(statement + "\n\n")
	}
	if ratingEffect := firstString(effect, "ratingEffect"); ratingEffect != "" {
		b.WriteString("Rating effect: " + ratingEffect + "\n\n")
	}
	if rationale := firstString(effect, "rationale"); rationale != "" {
		b.WriteString("Rationale: " + rationale + "\n\n")
	}
}

func writeFindingEvidenceSection(b *strings.Builder, headingLevel int, title string, evidence []map[string]any) {
	b.WriteString(strings.Repeat("#", headingLevel) + " " + title + "\n\n")
	if len(evidence) == 0 {
		b.WriteString("(none recorded)\n\n")
		return
	}
	for _, item := range evidence {
		b.WriteString("- `" + markdownCell(firstString(item, "sourceRef")) + "`: " + markdownCell(firstString(item, "statement")) + "\n")
		if rationale := firstString(item, "rationale"); rationale != "" {
			b.WriteString("  Rationale: " + markdownCell(rationale) + "\n")
		}
	}
	b.WriteString("\n")
}

func findingEffectSummary(finding map[string]any) string {
	effect := objectMap(finding["effect"])
	return firstString(effect, "statement", "ratingEffect")
}

func findingCauseSummary(finding map[string]any) string {
	cause := objectMap(finding["cause"])
	status := causeStatusTitle(firstString(cause, "status"))
	statement := firstString(cause, "statement")
	if status == "" {
		return statement
	}
	if statement == "" {
		return status
	}
	return status + ": " + statement
}

func sortedAreaFindings(analysis map[string]any) []indexedAreaFinding {
	findings := objectSlice(analysis["findings"])
	out := make([]indexedAreaFinding, 0, len(findings))
	for i, finding := range findings {
		out = append(out, indexedAreaFinding{Finding: finding, Index: i})
	}
	sort.SliceStable(out, func(i, j int) bool {
		return compareAreaFindings(out[i], out[j]) < 0
	})
	return out
}

func sortedFactorFindings(analysis map[string]any, factor factorID) []indexedAreaFinding {
	var out []indexedAreaFinding
	for i, finding := range objectSlice(analysis["findings"]) {
		for _, relationship := range objectSlice(finding["factorRelationships"]) {
			parsed, err := factorIDFrom(relationship["factorId"])
			if err == nil && sameStrings(parsed.DeclaringArea, factor.DeclaringArea) && sameStrings(parsed.Path, factor.Path) {
				out = append(out, indexedAreaFinding{Finding: finding, Index: i, Relationship: relationship})
				break
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		return compareFactorFindings(out[i], out[j]) < 0
	})
	return out
}

func compareAreaFindings(a, b indexedAreaFinding) int {
	if cmp := findingTypeRank(firstString(a.Finding, "type")) - findingTypeRank(firstString(b.Finding, "type")); cmp != 0 {
		return cmp
	}
	if cmp := findingSeverityRank(firstString(a.Finding, "severity")) - findingSeverityRank(firstString(b.Finding, "severity")); cmp != 0 {
		return cmp
	}
	if cmp := findingConfidenceRank(firstString(a.Finding, "confidence")) - findingConfidenceRank(firstString(b.Finding, "confidence")); cmp != 0 {
		return cmp
	}
	return a.Index - b.Index
}

func compareFactorFindings(a, b indexedAreaFinding) int {
	if cmp := findingTypeRank(firstString(a.Finding, "type")) - findingTypeRank(firstString(b.Finding, "type")); cmp != 0 {
		return cmp
	}
	if cmp := findingSeverityRank(firstString(a.Finding, "severity")) - findingSeverityRank(firstString(b.Finding, "severity")); cmp != 0 {
		return cmp
	}
	if cmp := factorRelationshipRank(firstString(a.Relationship, "relationship")) - factorRelationshipRank(firstString(b.Relationship, "relationship")); cmp != 0 {
		return cmp
	}
	if cmp := findingConfidenceRank(firstString(a.Finding, "confidence")) - findingConfidenceRank(firstString(b.Finding, "confidence")); cmp != 0 {
		return cmp
	}
	return a.Index - b.Index
}

func areaFindingID(finding map[string]any, index int) string {
	if id := firstString(finding, "id", "ID"); id != "" {
		return id
	}
	return fmt.Sprintf("finding-%d", index+1)
}

func areaFindingFactorLinks(spec *model.Spec, finding map[string]any, areaID []string, reportPath string) string {
	var links []string
	for _, relationship := range objectSlice(finding["factorRelationships"]) {
		factor, err := factorIDFrom(relationship["factorId"])
		if err != nil || !sameStrings(factor.DeclaringArea, areaID) {
			continue
		}
		label := factorTitle(spec, factor)
		rel := firstString(relationship, "relationship")
		if rel != "" {
			label += " " + factorRelationshipTitle(rel)
		}
		links = append(links, reportLink(reportPath, factorReportPath(factor), label))
	}
	if len(links) == 0 {
		return "(none)"
	}
	return strings.Join(links, "; ")
}

func findingTypeRank(value string) int {
	switch value {
	case "risk":
		return 0
	case "gap":
		return 1
	case "unknown":
		return 2
	case "note":
		return 3
	case "strength":
		return 4
	default:
		return 99
	}
}

func findingSeverityRank(value string) int {
	switch value {
	case "critical":
		return 0
	case "high":
		return 1
	case "medium":
		return 2
	case "low":
		return 3
	default:
		return 99
	}
}

func findingConfidenceRank(value string) int {
	switch value {
	case "high":
		return 0
	case "medium":
		return 1
	case "low":
		return 2
	case "none":
		return 3
	default:
		return 99
	}
}

func factorRelationshipRank(value string) int {
	switch value {
	case "primary-driver":
		return 0
	case "contributing-driver":
		return 1
	case "evidence-limit":
		return 2
	case "offsetting-strength":
		return 3
	case "related":
		return 4
	default:
		return 99
	}
}

func writeEvaluationFindingsTable(b *strings.Builder, assessment map[string]any) {
	b.WriteString("| ID | Statement | Type | Severity | Confidence | Effect | Cause |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
	findings := objectSlice(assessment["findings"])
	if len(findings) == 0 {
		b.WriteString("| (no findings) |  |  |  |  |  |  |\n\n")
		return
	}
	for i, finding := range findings {
		id := firstString(finding, "id", "ID")
		if id == "" {
			id = fmt.Sprintf("finding-%d", i+1)
		}
		b.WriteString("| `" + markdownCell(id) + "` | " + markdownCell(firstString(finding, "statement")) + " | " + markdownCell(findingTypeTitle(firstString(finding, "type"))) + " | " + markdownCell(findingSeverityTitle(firstString(finding, "severity"))) + " | " + markdownCell(confidenceTitle(firstString(finding, "confidence"))) + " | " + markdownCell(findingEffectSummary(finding)) + " | " + markdownCell(findingCauseSummary(finding)) + " |\n")
	}
	b.WriteString("\n")
}

func writeEvaluationFindingDetails(b *strings.Builder, assessment map[string]any) {
	findings := objectSlice(assessment["findings"])
	if len(findings) == 0 {
		b.WriteString("(no finding details)\n\n")
		return
	}
	for i, finding := range findings {
		id := firstString(finding, "id", "ID")
		if id == "" {
			id = fmt.Sprintf("finding-%d", i+1)
		}
		writeFindingCoreDetails(b, 3, id, finding)
	}
}

func writeEvaluationUnknownsTable(b *strings.Builder, assessment map[string]any, rating map[string]any) {
	b.WriteString("| Type | Detail |\n")
	b.WriteString("| --- | --- |\n")
	wrote := false
	for _, field := range []string{"unknowns", "missingEvidence"} {
		source := assessment
		if field == "missingEvidence" {
			source = rating
		}
		for _, item := range objectSlice(source[field]) {
			wrote = true
			b.WriteString("| " + markdownCell(unknownTypeTitle(field)) + " | " + markdownCell(firstString(item, "description", "reason", "ref", "id")) + " |\n")
		}
	}
	if !wrote {
		b.WriteString("| (none recorded) |  |\n")
	}
}

func areaReportPath(areaID []string) string {
	if len(areaID) == 0 {
		return "root-area.md"
	}
	parts := append([]string{"areas"}, areaID...)
	parts = append(parts, areaID[len(areaID)-1]+"-area.md")
	return filepath.ToSlash(filepath.Join(parts...))
}

func requirementReportPath(req requirementID) string {
	parts := reportAreaDirParts(req.DeclaringArea)
	parts = append(parts, "requirements", req.Name, req.Name+"-requirement.md")
	return filepath.ToSlash(filepath.Join(parts...))
}

func factorReportPath(factor factorID) string {
	parts := reportAreaDirParts(factor.DeclaringArea)
	for _, name := range factor.Path {
		parts = append(parts, "factors", name)
	}
	parts = append(parts, factor.Path[len(factor.Path)-1]+"-factor.md")
	return filepath.ToSlash(filepath.Join(parts...))
}

func reportAreaDirParts(areaID []string) []string {
	if len(areaID) == 0 {
		return nil
	}
	return append([]string{"areas"}, areaID...)
}

// reportDataLink renders a Data-column link to a machine-readable payload using
// the payload's base filename as the link text, so readers can tell which
// *-result.json a link opens.
func reportDataLink(fromReport, toPath string) string {
	return reportLink(fromReport, toPath, filepath.Base(toPath))
}

func reportLink(fromReport, toPath, label string) string {
	fromDir := filepath.Dir(fromReport)
	if fromDir == "." {
		fromDir = ""
	}
	rel, err := filepath.Rel(fromDir, toPath)
	if err != nil {
		rel = toPath
	}
	if rel == "." {
		rel = filepath.Base(toPath)
	}
	return "[" + markdownCell(label) + "](" + filepath.ToSlash(rel) + ")"
}

func areaTitle(spec *model.Spec, areaID []string) string {
	if len(areaID) == 0 {
		if spec.Title != "" {
			return spec.Title
		}
		return "Root Area"
	}
	area, ok := lookupArea(spec, areaID)
	if ok && area.Title != "" {
		return area.Title
	}
	return areaID[len(areaID)-1]
}

func factorTitle(spec *model.Spec, id factorID) string {
	factor, ok := lookupFactor(spec, id)
	if ok && factor.Title != "" {
		return factor.Title
	}
	return id.Path[len(id.Path)-1]
}

func requirementTitle(spec *model.Spec, id requirementID) string {
	req, ok := lookupRequirement(spec, id)
	if ok && req.Title != "" {
		return req.Title
	}
	return id.Name
}

func ratingTitle(spec *model.Spec, level string) string {
	if spec != nil {
		for _, candidate := range spec.RatingScale {
			if candidate.Level == level && candidate.Title != "" {
				return candidate.Title
			}
		}
	}
	return level
}

func lookupArea(spec *model.Spec, areaID []string) (model.Area, bool) {
	areas := spec.Areas
	var current model.Area
	for _, name := range areaID {
		next, ok := areas[name]
		if !ok {
			return model.Area{}, false
		}
		current = next
		areas = current.Areas
	}
	return current, len(areaID) > 0
}

func lookupFactor(spec *model.Spec, id factorID) (model.Factor, bool) {
	var factors map[string]model.Factor
	if len(id.DeclaringArea) == 0 {
		factors = spec.Factors
	} else {
		area, ok := lookupArea(spec, id.DeclaringArea)
		if !ok {
			return model.Factor{}, false
		}
		factors = area.Factors
	}
	var current model.Factor
	for _, name := range id.Path {
		next, ok := factors[name]
		if !ok {
			return model.Factor{}, false
		}
		current = next
		factors = current.Factors
	}
	return current, true
}

func lookupRequirement(spec *model.Spec, id requirementID) (model.Requirement, bool) {
	if len(id.DeclaringArea) == 0 {
		if req, ok := spec.Requirements[id.Name]; ok {
			return req, ok
		}
		return lookupRequirementInFactors(spec.Factors, id.Name)
	}
	area, ok := lookupArea(spec, id.DeclaringArea)
	if !ok {
		return model.Requirement{}, false
	}
	if req, ok := area.Requirements[id.Name]; ok {
		return req, ok
	}
	return lookupRequirementInFactors(area.Factors, id.Name)
}

func lookupRequirementInFactors(factors map[string]model.Factor, name string) (model.Requirement, bool) {
	for _, factor := range factors {
		if req, ok := factor.Requirements[name]; ok {
			return req, true
		}
		if req, ok := lookupRequirementInFactors(factor.Factors, name); ok {
			return req, true
		}
	}
	return model.Requirement{}, false
}

func factorList(factors []*evaluationFactorArtifacts, spec *model.Spec, fromReport string) string {
	if len(factors) == 0 {
		return ""
	}
	parts := make([]string, 0, len(factors))
	for _, factor := range factors {
		parts = append(parts, reportLink(fromReport, factorReportPath(factor.ID), factorTitle(spec, factor.ID))+" "+evaluationRatingLabel(spec, scopedMap(factor.Analysis, "localAndDescendantAnalysis")))
	}
	return strings.Join(parts, "; ")
}

func requirementFactorLinks(req *evaluationRequirementArtifacts, fromReport string) string {
	ids := requirementFactorIDs(req)
	if len(ids) == 0 {
		return ""
	}
	for i, id := range ids {
		parsed, err := parseRequirementFactorID(req.ID.DeclaringArea, id)
		if err != nil {
			continue
		}
		ids[i] = reportLink(fromReport, factorReportPath(parsed), strings.Join(parsed.Path, "/"))
	}
	return strings.Join(ids, "; ")
}

func requirementFactorIDs(req *evaluationRequirementArtifacts) []string {
	if req.Frame != nil {
		if raw, ok := req.Frame["subject"].(map[string]any); ok {
			return stringValues(raw["factorIds"])
		}
	}
	if req.Assessment != nil {
		return stringValues(req.Assessment["factors"])
	}
	return nil
}

func parseRequirementFactorID(areaID []string, ref string) (factorID, error) {
	if strings.HasPrefix(ref, "factor:") {
		return parseFactorRef(ref)
	}
	path := strings.Split(ref, "/")
	for _, part := range path {
		if !safeModelName(part) {
			return factorID{}, usagef("invalid factor ref %q", ref)
		}
	}
	return factorID{DeclaringArea: areaID, Path: path}, nil
}

func evaluationRequirementRatingLabel(spec *model.Spec, rating map[string]any) string {
	if rating == nil {
		return ratingStatusTitle(string(RatingStatusNotRated))
	}
	if status := evaluationString(rating, "status"); status != "rated" {
		if status == "" {
			return ratingStatusTitle(string(RatingStatusNotRated))
		}
		return ratingStatusTitle(status)
	}
	if level := evaluationString(rating, "ratingLevelId"); level != "" {
		return ratingTitle(spec, ratingLevelID(level))
	}
	return ratingStatusTitle(string(RatingStatusRated))
}

// evaluationSubRatingCell renders the descendant-inclusive ("+ Sub-Factors Rating" /
// "+ Sub-Areas Rating") cell for a breakdown row. When the node has no
// descendants there is no roll-up distinct from its local rating, so it renders
// an em dash rather than repeating the local rating.
func evaluationSubRatingCell(spec *model.Spec, aggregate map[string]any, hasDescendants bool) string {
	if !hasDescendants {
		return "—"
	}
	return evaluationRatingLabel(spec, aggregate)
}

func areaKey(id []string) string {
	return strings.Join(id, "/")
}

func factorKey(id factorID) string {
	return areaKey(id.DeclaringArea) + "::" + strings.Join(id.Path, "/")
}

func requirementKey(id requirementID) string {
	return areaKey(id.DeclaringArea) + "::" + id.Name
}

func areaIDFromKey(key string) []string {
	if key == "" {
		return nil
	}
	return strings.Split(key, "/")
}

func factorIDFromKey(key string) factorID {
	areaPart, factorPart, _ := strings.Cut(key, "::")
	return factorID{DeclaringArea: areaIDFromKey(areaPart), Path: strings.Split(factorPart, "/")}
}

func requirementIDFromKey(key string) requirementID {
	areaPart, name, _ := strings.Cut(key, "::")
	return requirementID{DeclaringArea: areaIDFromKey(areaPart), Name: name}
}

func areaDisplayPath(areaID []string) string {
	if len(areaID) == 0 {
		return "/"
	}
	return "/" + strings.Join(areaID, "/")
}

func factorDisplayPath(id factorID) string {
	area := areaDisplayPath(id.DeclaringArea)
	if area == "/" {
		return strings.Join(id.Path, "/")
	}
	return strings.TrimPrefix(area, "/") + "::" + strings.Join(id.Path, "/")
}

func anyStrings(values []string) []any {
	out := make([]any, 0, len(values))
	for _, value := range values {
		out = append(out, value)
	}
	return out
}

func copyStrings(values []string) []string {
	return append([]string(nil), values...)
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func factorIDJSON(id factorID) string {
	return model.FactorReference(model.AreaPath(id.DeclaringArea), model.FactorPath(id.Path))
}

func requirementIDJSON(id requirementID) string {
	return model.RequirementReference(model.AreaPath(id.DeclaringArea), id.Name)
}

func objectSlice(v any) []map[string]any {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if mapped, ok := item.(map[string]any); ok {
			out = append(out, mapped)
		}
	}
	return out
}

func objectMap(v any) map[string]any {
	if mapped, ok := v.(map[string]any); ok {
		return mapped
	}
	return nil
}

func stringValues(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	var out []string
	for _, item := range raw {
		switch value := item.(type) {
		case string:
			out = append(out, value)
		case map[string]any:
			if factorPath, ok := value["factorPath"].([]any); ok {
				out = append(out, strings.Join(stringValues(factorPath), "/"))
			}
		}
	}
	return out
}

func firstString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := m[key].(string); ok {
			return value
		}
	}
	return ""
}

func compactJSON(v any) string {
	if v == nil {
		return ""
	}
	switch value := v.(type) {
	case string:
		return value
	default:
		raw, err := json.Marshal(value)
		if err != nil {
			return ""
		}
		return string(raw)
	}
}

func scopedMap(payload map[string]any, field string) map[string]any {
	raw, ok := payload[field].(map[string]any)
	if !ok {
		return map[string]any{}
	}
	return raw
}

func evaluationRatingLabel(spec *model.Spec, scope map[string]any) string {
	status := evaluationString(scope, "status")
	level := ratingLevelID(evaluationString(scope, "ratingLevelId"))
	if status != "analyzed" || level == "" {
		if status == "" {
			return analysisStatusTitle(string(AnalysisStatusNotAnalyzed))
		}
		return analysisStatusTitle(status)
	}
	return ratingTitle(spec, level)
}

func evaluationAnalysisStatusPair(overall, local map[string]any) string {
	return emDashIfEmpty(analysisStatusTitle(evaluationString(overall, "status"))) + " / " + emDashIfEmpty(analysisStatusTitle(evaluationString(local, "status")))
}

func evaluationConfidencePair(overall, local map[string]any) string {
	return emDashIfEmpty(confidenceTitle(evaluationString(overall, "confidence"))) + " / " + emDashIfEmpty(confidenceTitle(evaluationString(local, "confidence")))
}

func evaluationRequirementConfidencePair(req *evaluationRequirementArtifacts) string {
	return emDashIfEmpty(confidenceTitle(evaluationString(req.Rating, "confidence"))) + " / " + emDashIfEmpty(confidenceTitle(evaluationString(req.Assessment, "confidence")))
}

func emDashIfEmpty(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func ratingLevelID(ref string) string {
	return strings.TrimPrefix(ref, "rating:")
}

func evaluationSummary(scope map[string]any) string {
	rationale := evaluationString(scope, "rationale")
	if rationale == "" {
		return "No analysis rationale was recorded."
	}
	return rationale
}

func evaluationString(payload map[string]any, field string) string {
	value, _ := payload[field].(string)
	return value
}

func evaluationReceiptRating(analysis map[string]any) RatingResult {
	overall := scopedMap(analysis, "localAndDescendantAnalysis")
	if evaluationString(overall, "status") != "analyzed" {
		return RatingResult{Kind: "not-assessed", Rationale: evaluationString(overall, "statusReason")}
	}
	return RatingResult{Kind: "rated", Level: ratingLevelID(evaluationString(overall, "ratingLevelId")), Rationale: evaluationString(overall, "rationale")}
}

func markdownCell(s string) string {
	if s == "" {
		return "—"
	}
	return strings.ReplaceAll(s, "|", "\\|")
}

func writeEvaluationLegend(b *strings.Builder) {
	b.WriteString("\n## Legend\n\n")
	b.WriteString("- `—` - not applicable or not recorded.\n")
}
