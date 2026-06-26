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
	required := []struct {
		rel  string
		kind DataKind
	}{
		{"data/frame/evaluation-frame.json", DataKindEvaluationFrame},
		{"data/areas/root/area-evaluation-frame.json", DataKindAreaEvaluationFrame},
		{"data/areas/root/area-analysis-result.json", DataKindAreaAnalysis},
	}
	var gaps []RunGap
	for _, req := range required {
		raw, err := os.ReadFile(filepath.Join(runAbs, req.rel))
		if os.IsNotExist(err) {
			gaps = append(gaps, RunGap{Kind: GapMissingEvaluationData, Ref: req.rel, Detail: "required Evaluation evaluation payload is missing"})
			continue
		}
		if err != nil {
			gaps = append(gaps, RunGap{Kind: GapUnreadableEvaluationData, Ref: req.rel, Detail: err.Error()})
			continue
		}
		payload, _, err := decodeDataPayload(raw)
		if err != nil {
			gaps = append(gaps, RunGap{Kind: GapMalformedEvaluationData, Ref: req.rel, Detail: err.Error()})
			continue
		}
		kind, err := payloadKind(payload)
		if err != nil {
			gaps = append(gaps, RunGap{Kind: GapIncompleteEvaluationData, Ref: req.rel, Detail: err.Error()})
			continue
		}
		if kind != req.kind {
			gaps = append(gaps, RunGap{Kind: GapIncompleteEvaluationData, Ref: req.rel, Detail: fmt.Sprintf("kind = %s, want %s", kind, req.kind)})
			continue
		}
		if err := validateDataPayload(kind, payload); err != nil {
			gaps = append(gaps, RunGap{Kind: GapIncompleteEvaluationData, Ref: req.rel, Detail: err.Error()})
		}
	}
	return gaps
}

func buildEvaluationReport(path string) (*BuildReportReceipt, error) {
	runAbs, err := verifyRun(path)
	if err != nil {
		return nil, err
	}
	displayPath := displayRunPath(runAbs)
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
	rootArea := artifacts.area(areaKey(nil))
	if rootArea == nil || rootArea.Analysis == nil {
		return nil, nonReportableRunError(displayPath, RunGap{Kind: GapMissingEvaluationData, Ref: "data/areas/root/area-analysis-result.json", Detail: "required Evaluation evaluation payload is missing"})
	}
	reports := renderEvaluationReportTree(spec, artifacts)
	for _, report := range reports {
		reportAbs := filepath.Join(runAbs, report.Path)
		if err := os.MkdirAll(filepath.Dir(reportAbs), 0o755); err != nil {
			return nil, fmt.Errorf("creating report directory: %w", err)
		}
		if err := writeReportFile(reportAbs, []byte(report.Content)); err != nil {
			return nil, err
		}
	}
	output := evaluationOutputResult(artifacts, reports)
	outputRaw, err := canonicalJSON(output)
	if err != nil {
		return nil, err
	}
	outputRel := filepath.Join("data", "evaluation-output-result.json")
	if err := os.WriteFile(filepath.Join(runAbs, outputRel), outputRaw, 0o644); err != nil {
		return nil, fmt.Errorf("writing %s: %w", filepath.ToSlash(outputRel), err)
	}
	reportRel := "report.md"
	return &BuildReportReceipt{
		SchemaVersion:          SchemaVersion,
		Path:                   displayPath,
		ReportMD:               filepath.ToSlash(filepath.Join(displayPath, reportRel)),
		EvaluationOutputResult: filepath.ToSlash(filepath.Join(displayPath, outputRel)),
		RatingResult:           evaluationReceiptRating(rootArea.Analysis),
	}, nil
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
	DataKindAreaEvaluationFrame:        collectEvaluationAreaFrame,
	DataKindAreaAnalysis:               collectEvaluationAreaAnalysis,
	DataKindRequirementEvaluationFrame: collectEvaluationRequirementFrame,
	DataKindRequirementAssessment:      collectEvaluationRequirementAssessment,
	DataKindRequirementRating:          collectEvaluationRequirementRating,
	DataKindFactorAnalysisFrame:        collectEvaluationFactorFrame,
	DataKindFactorAnalysis:             collectEvaluationFactorAnalysis,
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

func renderEvaluationReportTree(spec *model.Spec, artifacts *evaluationArtifacts) []evaluationRenderedReport {
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
	return reports
}

func renderEvaluationAreaReport(spec *model.Spec, artifacts *evaluationArtifacts, area *evaluationAreaArtifacts, reportPath string) string {
	title := areaTitle(spec, area.ID)
	local := scopedMap(area.Analysis, "localAnalysis")
	overall := scopedMap(area.Analysis, "localAndDescendantAnalysis")
	var b strings.Builder
	b.WriteString("# Area: " + title + "\n\n")
	writeEvaluationAreaTrail(&b, spec, area.ID, reportPath)
	b.WriteString("| Overall Rating | Local Rating | Confidence | Data |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	b.WriteString("| " + markdownCell(evaluationRatingLabel(spec, overall)) + " | " + markdownCell(evaluationRatingLabel(spec, local)) + " | " + markdownCell(evaluationConfidencePair(overall, local)) + " | " + reportDataLink(reportPath, areaDataPath(area.ID, "area-analysis-result.json")) + " |\n\n")
	b.WriteString("Summary:\n\n")
	b.WriteString(evaluationSummary(overall))
	b.WriteString("\n\n## Rating Drivers\n\n")
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
	writeEvaluationAreaTrail(&b, spec, factor.ID.DeclaringArea, reportPath)
	writeEvaluationFactorTrail(&b, spec, factor.ID, reportPath)
	b.WriteString("| Overall Rating | Local Rating | Status | Confidence | Data |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	b.WriteString("| " + markdownCell(evaluationRatingLabel(spec, overall)) + " | " + markdownCell(evaluationRatingLabel(spec, local)) + " | " + markdownCell(evaluationAnalysisStatusPair(overall, local)) + " | " + markdownCell(evaluationConfidencePair(overall, local)) + " | " + reportDataLink(reportPath, factorDataPath(factor.ID, "factor-analysis-result.json")) + " |\n\n")
	b.WriteString("Summary:\n\n")
	b.WriteString(evaluationSummary(overall))
	b.WriteString("\n\n## Rating Drivers\n\n")
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
	writeEvaluationAreaTrail(&b, spec, req.ID.DeclaringArea, reportPath)
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

func evaluationOutputResult(artifacts *evaluationArtifacts, reports []evaluationRenderedReport) map[string]any {
	reportOutputs := make([]any, 0, len(reports))
	reportsByArea := map[string][]any{}
	for _, report := range reports {
		ref := evaluationReportRef(report)
		reportOutputs = append(reportOutputs, ref)
		reportsByArea[areaKey(report.AreaID)] = append(reportsByArea[areaKey(report.AreaID)], ref)
	}
	var areaOutputs []any
	for _, area := range artifacts.sortedAreas() {
		if area.Analysis == nil {
			continue
		}
		areaID := AreaPath(area.ID).Reference()
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
	return map[string]any{
		"schemaVersion":       SchemaVersion,
		"kind":                string(DataKindEvaluationOutput),
		"rootAreaAnalysisRef": routineRef(DataKindAreaAnalysis, map[string]any{"areaId": AreaPath{}.Reference()}, "localAndDescendantAnalysis"),
		"areaOutputs":         areaOutputs,
		"reportOutputs":       reportOutputs,
	}
}

func evaluationReportRef(report evaluationRenderedReport) map[string]any {
	ref := map[string]any{"kind": report.Kind, "areaId": AreaPath(report.AreaID).Reference(), "path": filepath.ToSlash(report.Path)}
	if report.FactorID != nil {
		ref["factorId"] = factorIDJSON(*report.FactorID)
	}
	if report.RequirementID != nil {
		ref["requirementId"] = requirementIDJSON(*report.RequirementID)
	}
	return ref
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

func writeEvaluationAreaTrail(b *strings.Builder, spec *model.Spec, areaID []string, reportPath string) {
	parts := []string{reportLink(reportPath, areaReportPath(nil), areaTitle(spec, nil))}
	for i := range areaID {
		id := areaID[:i+1]
		parts = append(parts, reportLink(reportPath, areaReportPath(id), areaTitle(spec, id)))
	}
	b.WriteString("Area: " + strings.Join(parts, " / ") + "\n\n")
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

func writeEvaluationFindingsTable(b *strings.Builder, assessment map[string]any) {
	b.WriteString("| ID | Type | Severity | Description |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	findings := objectSlice(assessment["findings"])
	if len(findings) == 0 {
		b.WriteString("| (no findings) |  |  |  |\n\n")
		return
	}
	for i, finding := range findings {
		id := firstString(finding, "id", "ID")
		if id == "" {
			id = fmt.Sprintf("finding-%d", i+1)
		}
		b.WriteString("| `" + markdownCell(id) + "` | " + markdownCell(findingTypeTitle(firstString(finding, "type", "Type"))) + " | " + markdownCell(findingSeverityTitle(firstString(finding, "severity", "Severity"))) + " | " + markdownCell(firstString(finding, "description", "Description")) + " |\n")
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
		b.WriteString("### " + id + "\n\n")
		if description := firstString(finding, "description", "Description"); description != "" {
			b.WriteString(description + "\n\n")
		}
		b.WriteString("| Field | Value |\n")
		b.WriteString("| --- | --- |\n")
		b.WriteString("| Type | " + markdownCell(findingTypeTitle(firstString(finding, "type", "Type"))) + " |\n")
		b.WriteString("| Severity | " + markdownCell(findingSeverityTitle(firstString(finding, "severity", "Severity"))) + " |\n")
		b.WriteString("| Location | " + markdownCell(compactJSON(firstPresent(finding, "location", "Location"))) + " |\n")
		b.WriteString("| Evidence | " + markdownCell(compactJSON(firstPresent(finding, "evidence", "Evidence"))) + " |\n")
		b.WriteString("| Rationale | " + markdownCell(firstString(finding, "rationale", "Rationale")) + " |\n\n")
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
		return "report.md"
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
	return FactorReference(AreaPath(id.DeclaringArea), FactorPath(id.Path))
}

func requirementIDJSON(id requirementID) string {
	return RequirementReference(AreaPath(id.DeclaringArea), id.Name)
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

func firstPresent(m map[string]any, keys ...string) any {
	for _, key := range keys {
		if value, ok := m[key]; ok {
			return value
		}
	}
	return nil
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
