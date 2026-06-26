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

func v2RenderableGaps(runAbs string) []RunGap {
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
			gaps = append(gaps, RunGap{Kind: GapMissingEvaluationData, Ref: req.rel, Detail: "required Evaluation v2 payload is missing"})
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

func buildV2Report(path string) (*BuildReportReceipt, error) {
	runAbs, err := verifyRun(path)
	if err != nil {
		return nil, err
	}
	displayPath := displayRunPath(runAbs)
	if gaps := v2RenderableGaps(runAbs); len(gaps) > 0 {
		return nil, nonReportableRunError(displayPath, gaps[0])
	}
	spec, err := loadRunModel(runAbs)
	if err != nil {
		return nil, err
	}
	artifacts, err := collectV2Artifacts(runAbs)
	if err != nil {
		return nil, err
	}
	rootArea := artifacts.area(areaKey(nil))
	if rootArea == nil || rootArea.Analysis == nil {
		return nil, nonReportableRunError(displayPath, RunGap{Kind: GapMissingEvaluationData, Ref: "data/areas/root/area-analysis-result.json", Detail: "required Evaluation v2 payload is missing"})
	}
	reports := renderV2ReportTree(spec, artifacts)
	for _, report := range reports {
		reportAbs := filepath.Join(runAbs, report.Path)
		if err := os.MkdirAll(filepath.Dir(reportAbs), 0o755); err != nil {
			return nil, fmt.Errorf("creating report directory: %w", err)
		}
		if err := writeReportFile(reportAbs, []byte(report.Content)); err != nil {
			return nil, err
		}
	}
	output := v2OutputResult(artifacts, reports)
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
		RatingResult:           v2ReceiptRating(rootArea.Analysis),
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

type v2AreaArtifacts struct {
	ID       []string
	Frame    map[string]any
	Analysis map[string]any
}

type v2FactorArtifacts struct {
	ID       factorID
	Frame    map[string]any
	Analysis map[string]any
}

type v2RequirementArtifacts struct {
	ID         requirementID
	Frame      map[string]any
	Assessment map[string]any
	Rating     map[string]any
}

type v2Artifacts struct {
	Areas        map[string]*v2AreaArtifacts
	Factors      map[string]*v2FactorArtifacts
	Requirements map[string]*v2RequirementArtifacts
}

type v2RenderedReport struct {
	Kind          string
	Path          string
	AreaID        []string
	FactorID      *factorID
	RequirementID *requirementID
	Content       string
}

func collectV2Artifacts(runAbs string) (*v2Artifacts, error) {
	out := &v2Artifacts{
		Areas:        map[string]*v2AreaArtifacts{},
		Factors:      map[string]*v2FactorArtifacts{},
		Requirements: map[string]*v2RequirementArtifacts{},
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
		if collector := v2PayloadCollectors[kind]; collector != nil {
			return collector(out, payload)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("collecting Evaluation v2 data: %w", err)
	}
	return out, nil
}

type v2PayloadCollector func(*v2Artifacts, map[string]any) error

var v2PayloadCollectors = map[DataKind]v2PayloadCollector{
	DataKindAreaEvaluationFrame:        collectV2AreaFrame,
	DataKindAreaAnalysis:               collectV2AreaAnalysis,
	DataKindRequirementEvaluationFrame: collectV2RequirementFrame,
	DataKindRequirementAssessment:      collectV2RequirementAssessment,
	DataKindRequirementRating:          collectV2RequirementRating,
	DataKindFactorAnalysisFrame:        collectV2FactorFrame,
	DataKindFactorAnalysis:             collectV2FactorAnalysis,
}

func collectV2AreaFrame(out *v2Artifacts, payload map[string]any) error {
	id, err := subjectAreaID(payload)
	if err != nil {
		return err
	}
	out.area(areaKey(id)).Frame = payload
	return nil
}

func collectV2AreaAnalysis(out *v2Artifacts, payload map[string]any) error {
	id, err := topAreaID(payload)
	if err != nil {
		return err
	}
	out.area(areaKey(id)).Analysis = payload
	return nil
}

func collectV2RequirementFrame(out *v2Artifacts, payload map[string]any) error {
	id, err := subjectRequirementID(payload)
	if err != nil {
		return err
	}
	out.requirement(requirementKey(id)).Frame = payload
	return nil
}

func collectV2RequirementAssessment(out *v2Artifacts, payload map[string]any) error {
	id, err := topRequirementID(payload)
	if err != nil {
		return err
	}
	out.requirement(requirementKey(id)).Assessment = payload
	return nil
}

func collectV2RequirementRating(out *v2Artifacts, payload map[string]any) error {
	id, err := topRequirementID(payload)
	if err != nil {
		return err
	}
	out.requirement(requirementKey(id)).Rating = payload
	return nil
}

func collectV2FactorFrame(out *v2Artifacts, payload map[string]any) error {
	id, err := subjectFactorID(payload)
	if err != nil {
		return err
	}
	out.factor(factorKey(id)).Frame = payload
	return nil
}

func collectV2FactorAnalysis(out *v2Artifacts, payload map[string]any) error {
	id, err := topFactorID(payload)
	if err != nil {
		return err
	}
	out.factor(factorKey(id)).Analysis = payload
	return nil
}

func (a *v2Artifacts) area(key string) *v2AreaArtifacts {
	if existing, ok := a.Areas[key]; ok {
		return existing
	}
	id := areaIDFromKey(key)
	created := &v2AreaArtifacts{ID: id}
	a.Areas[key] = created
	return created
}

func (a *v2Artifacts) factor(key string) *v2FactorArtifacts {
	if existing, ok := a.Factors[key]; ok {
		return existing
	}
	id := factorIDFromKey(key)
	created := &v2FactorArtifacts{ID: id}
	a.Factors[key] = created
	return created
}

func (a *v2Artifacts) requirement(key string) *v2RequirementArtifacts {
	if existing, ok := a.Requirements[key]; ok {
		return existing
	}
	id := requirementIDFromKey(key)
	created := &v2RequirementArtifacts{ID: id}
	a.Requirements[key] = created
	return created
}

func renderV2ReportTree(spec *model.Spec, artifacts *v2Artifacts) []v2RenderedReport {
	var reports []v2RenderedReport
	for _, area := range artifacts.sortedAreas() {
		if area.Analysis == nil {
			continue
		}
		path := areaReportPath(area.ID)
		reports = append(reports, v2RenderedReport{
			Kind:    string(ReportKindArea),
			Path:    path,
			AreaID:  copyStrings(area.ID),
			Content: renderV2AreaReport(spec, artifacts, area, path),
		})
	}
	for _, factor := range artifacts.sortedFactors() {
		if factor.Analysis == nil {
			continue
		}
		id := factor.ID
		path := factorReportPath(id)
		reports = append(reports, v2RenderedReport{
			Kind:     string(ReportKindFactor),
			Path:     path,
			AreaID:   copyStrings(id.DeclaringArea),
			FactorID: &id,
			Content:  renderV2FactorReport(spec, artifacts, factor, path),
		})
	}
	for _, requirement := range artifacts.sortedRequirements() {
		if requirement.Assessment == nil && requirement.Rating == nil {
			continue
		}
		id := requirement.ID
		path := requirementReportPath(id)
		reports = append(reports, v2RenderedReport{
			Kind:          string(ReportKindRequirement),
			Path:          path,
			AreaID:        copyStrings(id.DeclaringArea),
			RequirementID: &id,
			Content:       renderV2RequirementReport(spec, artifacts, requirement, path),
		})
	}
	return reports
}

func renderV2AreaReport(spec *model.Spec, artifacts *v2Artifacts, area *v2AreaArtifacts, reportPath string) string {
	title := areaTitle(spec, area.ID)
	local := scopedMap(area.Analysis, "localAnalysis")
	overall := scopedMap(area.Analysis, "localAndDescendantAnalysis")
	var b strings.Builder
	writeV2AreaTrail(&b, spec, area.ID, reportPath)
	b.WriteString("# " + title + "\n\n")
	b.WriteString("Path: `" + areaDisplayPath(area.ID) + "`\n\n")
	b.WriteString("| Overall Rating | Local Rating | Confidence | Data |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	b.WriteString("| " + markdownCell(v2RatingLabel(spec, overall)) + " | " + markdownCell(v2RatingLabel(spec, local)) + " | " + markdownCell(v2ConfidencePair(overall, local)) + " | " + reportDataLink(reportPath, areaDataPath(area.ID, "area-analysis-result.json")) + " |\n\n")
	b.WriteString("Summary:\n\n")
	b.WriteString(v2Summary(overall))
	b.WriteString("\n\n## Rating Drivers\n\n")
	writeV2DriversTable(&b, spec, overall)
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
			b.WriteString("| " + reportLink(reportPath, factorReportPath(factor.ID), factorTitle(spec, factor.ID)) + " | `" + factorDisplayPath(factor.ID) + "` | " + markdownCell(v2RatingLabel(spec, factorLocal)) + " | " + markdownCell(v2SubRatingCell(spec, factorOverall, len(children) > 0)) + " | " + markdownCell(factorList(children, spec, reportPath)) + " |\n")
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
			b.WriteString("| " + reportLink(reportPath, areaReportPath(child.ID), areaTitle(spec, child.ID)) + " | `" + areaDisplayPath(child.ID) + "` | " + markdownCell(v2RatingLabel(spec, childLocal)) + " | " + markdownCell(v2SubRatingCell(spec, childOverall, len(artifacts.childAreas(child.ID)) > 0)) + " | " + markdownCell(factorList(artifacts.rootFactorsForArea(child.ID), spec, reportPath)) + " |\n")
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
			b.WriteString("| " + reportLink(reportPath, requirementReportPath(req.ID), requirementTitle(spec, req.ID)) + " | " + markdownCell(v2RequirementRatingLabel(spec, req.Rating)) + " | " + markdownCell(assessmentStatusTitle(v2String(req.Assessment, "status"))) + " | " + markdownCell(requirementFactorLinks(req, reportPath)) + " |\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("## Limits & Incomplete Inputs\n\n")
	writeV2LimitsTable(&b, local, overall)
	return b.String()
}

func renderV2FactorReport(spec *model.Spec, artifacts *v2Artifacts, factor *v2FactorArtifacts, reportPath string) string {
	local := scopedMap(factor.Analysis, "localAnalysis")
	overall := scopedMap(factor.Analysis, "localAndDescendantAnalysis")
	title := factorTitle(spec, factor.ID)
	var b strings.Builder
	writeV2AreaTrail(&b, spec, factor.ID.DeclaringArea, reportPath)
	writeV2FactorTrail(&b, spec, factor.ID, reportPath)
	b.WriteString("# " + title + "\n\n")
	b.WriteString("Path: `" + factorDisplayPath(factor.ID) + "`\n\n")
	b.WriteString("| Overall Rating | Local Rating | Status | Confidence | Data |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	b.WriteString("| " + markdownCell(v2RatingLabel(spec, overall)) + " | " + markdownCell(v2RatingLabel(spec, local)) + " | " + markdownCell(v2AnalysisStatusPair(overall, local)) + " | " + markdownCell(v2ConfidencePair(overall, local)) + " | " + reportDataLink(reportPath, factorDataPath(factor.ID, "factor-analysis-result.json")) + " |\n\n")
	b.WriteString("Summary:\n\n")
	b.WriteString(v2Summary(overall))
	b.WriteString("\n\n## Rating Drivers\n\n")
	writeV2DriversTable(&b, spec, overall)
	b.WriteString("## Requirements\n\n")
	b.WriteString("| Requirement | Rating | Status |\n")
	b.WriteString("| --- | --- | --- |\n")
	if requirements := artifacts.requirementsForFactor(factor.ID); len(requirements) == 0 {
		b.WriteString("| (no direct Requirements) |  |  |\n\n")
	} else {
		for _, req := range requirements {
			b.WriteString("| " + reportLink(reportPath, requirementReportPath(req.ID), requirementTitle(spec, req.ID)) + " | " + markdownCell(v2RequirementRatingLabel(spec, req.Rating)) + " | " + markdownCell(assessmentStatusTitle(v2String(req.Assessment, "status"))) + " |\n")
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
			b.WriteString("| " + reportLink(reportPath, factorReportPath(child.ID), factorTitle(spec, child.ID)) + " | `" + factorDisplayPath(child.ID) + "` | " + markdownCell(v2RatingLabel(spec, childLocal)) + " | " + markdownCell(v2SubRatingCell(spec, childOverall, len(artifacts.childFactors(child.ID)) > 0)) + " |\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("## Limits & Incomplete Inputs\n\n")
	writeV2LimitsTable(&b, local, overall)
	return b.String()
}

func renderV2RequirementReport(spec *model.Spec, artifacts *v2Artifacts, req *v2RequirementArtifacts, reportPath string) string {
	title := requirementTitle(spec, req.ID)
	var b strings.Builder
	writeV2AreaTrail(&b, spec, req.ID.DeclaringArea, reportPath)
	b.WriteString("# " + title + "\n\n")
	b.WriteString("Name: `" + markdownCell(req.ID.Name) + "`\n\n")
	b.WriteString("| Rating | Assessment | Factors | Confidence | Data |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	b.WriteString("| " + markdownCell(v2RequirementRatingLabel(spec, req.Rating)) + " | " + markdownCell(assessmentStatusTitle(v2String(req.Assessment, "status"))) + " | " + markdownCell(requirementFactorLinks(req, reportPath)) + " | " + markdownCell(confidenceTitle(v2String(req.Rating, "confidence"))+" / "+confidenceTitle(v2String(req.Assessment, "confidence"))) + " | " + reportDataLink(reportPath, requirementDataPath(req.ID, "requirement-assessment-result.json")) + ", " + reportDataLink(reportPath, requirementDataPath(req.ID, "requirement-rating-result.json")) + " |\n\n")
	b.WriteString("Summary:\n\n")
	if summary := v2String(req.Assessment, "evidenceSummary"); summary != "" {
		b.WriteString(summary)
	} else if rationale := v2String(req.Rating, "rationale"); rationale != "" {
		b.WriteString(rationale)
	} else {
		b.WriteString("No assessment summary was recorded.")
	}
	b.WriteString("\n\n## Findings Summary\n\n")
	writeV2FindingsTable(&b, req.Assessment)
	b.WriteString("## Finding Details\n\n")
	writeV2FindingDetails(&b, req.Assessment)
	b.WriteString("## Unknowns & Missing Evidence\n\n")
	writeV2UnknownsTable(&b, req.Assessment, req.Rating)
	_ = artifacts
	return b.String()
}

func v2OutputResult(artifacts *v2Artifacts, reports []v2RenderedReport) map[string]any {
	reportOutputs := make([]any, 0, len(reports))
	reportsByArea := map[string][]any{}
	for _, report := range reports {
		ref := v2ReportRef(report)
		reportOutputs = append(reportOutputs, ref)
		reportsByArea[areaKey(report.AreaID)] = append(reportsByArea[areaKey(report.AreaID)], ref)
	}
	var areaOutputs []any
	for _, area := range artifacts.sortedAreas() {
		if area.Analysis == nil {
			continue
		}
		areaID := anyStrings(area.ID)
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
		"rootAreaAnalysisRef": routineRef(DataKindAreaAnalysis, map[string]any{"areaId": []any{}}, "localAndDescendantAnalysis"),
		"areaOutputs":         areaOutputs,
		"reportOutputs":       reportOutputs,
	}
}

func v2ReportRef(report v2RenderedReport) map[string]any {
	ref := map[string]any{"kind": report.Kind, "areaId": anyStrings(report.AreaID), "path": filepath.ToSlash(report.Path)}
	if report.FactorID != nil {
		ref["factorId"] = factorIDJSON(*report.FactorID)
	}
	if report.RequirementID != nil {
		ref["requirementId"] = requirementIDJSON(*report.RequirementID)
	}
	return ref
}

func factorAnalysisRefs(factors []*v2FactorArtifacts) []any {
	refs := make([]any, 0, len(factors))
	for _, factor := range factors {
		refs = append(refs, routineRef(DataKindFactorAnalysis, map[string]any{"factorId": factorIDJSON(factor.ID)}, "localAndDescendantAnalysis"))
	}
	return refs
}

func requirementAssessmentRefs(requirements []*v2RequirementArtifacts) []any {
	refs := make([]any, 0, len(requirements))
	for _, requirement := range requirements {
		refs = append(refs, routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": requirementIDJSON(requirement.ID)}, ""))
	}
	return refs
}

func requirementRatingRefs(requirements []*v2RequirementArtifacts) []any {
	refs := make([]any, 0, len(requirements))
	for _, requirement := range requirements {
		refs = append(refs, routineRef(DataKindRequirementRating, map[string]any{"requirementId": requirementIDJSON(requirement.ID)}, ""))
	}
	return refs
}

func (a *v2Artifacts) sortedAreas() []*v2AreaArtifacts {
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
	out := make([]*v2AreaArtifacts, 0, len(keys))
	for _, key := range keys {
		out = append(out, a.Areas[key])
	}
	return out
}

func (a *v2Artifacts) sortedFactors() []*v2FactorArtifacts {
	keys := make([]string, 0, len(a.Factors))
	for key := range a.Factors {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]*v2FactorArtifacts, 0, len(keys))
	for _, key := range keys {
		out = append(out, a.Factors[key])
	}
	return out
}

func (a *v2Artifacts) sortedRequirements() []*v2RequirementArtifacts {
	keys := make([]string, 0, len(a.Requirements))
	for key := range a.Requirements {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]*v2RequirementArtifacts, 0, len(keys))
	for _, key := range keys {
		out = append(out, a.Requirements[key])
	}
	return out
}

func (a *v2Artifacts) childAreas(parent []string) []*v2AreaArtifacts {
	var out []*v2AreaArtifacts
	for _, area := range a.sortedAreas() {
		if len(area.ID) == len(parent)+1 && sameStrings(area.ID[:len(parent)], parent) && area.Analysis != nil {
			out = append(out, area)
		}
	}
	return out
}

func (a *v2Artifacts) rootFactorsForArea(areaID []string) []*v2FactorArtifacts {
	var out []*v2FactorArtifacts
	for _, factor := range a.sortedFactors() {
		if sameStrings(factor.ID.DeclaringArea, areaID) && len(factor.ID.Path) == 1 && factor.Analysis != nil {
			out = append(out, factor)
		}
	}
	return out
}

func (a *v2Artifacts) childFactors(parent factorID) []*v2FactorArtifacts {
	var out []*v2FactorArtifacts
	for _, factor := range a.sortedFactors() {
		if sameStrings(factor.ID.DeclaringArea, parent.DeclaringArea) && len(factor.ID.Path) == len(parent.Path)+1 && sameStrings(factor.ID.Path[:len(parent.Path)], parent.Path) && factor.Analysis != nil {
			out = append(out, factor)
		}
	}
	return out
}

func (a *v2Artifacts) requirementsForArea(areaID []string) []*v2RequirementArtifacts {
	var out []*v2RequirementArtifacts
	for _, req := range a.sortedRequirements() {
		if sameStrings(req.ID.DeclaringArea, areaID) {
			out = append(out, req)
		}
	}
	return out
}

func (a *v2Artifacts) requirementsForFactor(factor factorID) []*v2RequirementArtifacts {
	var out []*v2RequirementArtifacts
	want := factorDisplayPath(factor)
	for _, req := range a.sortedRequirements() {
		if !sameStrings(req.ID.DeclaringArea, factor.DeclaringArea) {
			continue
		}
		for _, linked := range requirementFactorIDs(req) {
			if linked == want || linked == strings.Join(factor.Path, "/") || linked == factor.Path[len(factor.Path)-1] {
				out = append(out, req)
				break
			}
		}
	}
	return out
}

func writeV2AreaTrail(b *strings.Builder, spec *model.Spec, areaID []string, reportPath string) {
	parts := []string{reportLink(reportPath, areaReportPath(nil), areaTitle(spec, nil))}
	for i := range areaID {
		id := areaID[:i+1]
		parts = append(parts, reportLink(reportPath, areaReportPath(id), areaTitle(spec, id)))
	}
	b.WriteString("Area: " + strings.Join(parts, " / ") + "\n\n")
}

func writeV2FactorTrail(b *strings.Builder, spec *model.Spec, factor factorID, reportPath string) {
	parts := make([]string, 0, len(factor.Path))
	for i := range factor.Path {
		id := factorID{DeclaringArea: factor.DeclaringArea, Path: factor.Path[:i+1]}
		parts = append(parts, reportLink(reportPath, factorReportPath(id), factorTitle(spec, id)))
	}
	b.WriteString("Factor: " + strings.Join(parts, " / ") + "\n\n")
}

func writeV2DriversTable(b *strings.Builder, spec *model.Spec, scope map[string]any) {
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

func writeV2LimitsTable(b *strings.Builder, scopes ...map[string]any) {
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

func writeV2FindingsTable(b *strings.Builder, assessment map[string]any) {
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

func writeV2FindingDetails(b *strings.Builder, assessment map[string]any) {
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
		b.WriteString("| Rationale | " + markdownCell(firstString(finding, "rationale", "Rationale")) + " |\n")
		b.WriteString("| Actions | " + markdownCell(compactJSON(firstPresent(finding, "actions", "Actions"))) + " |\n\n")
	}
}

func writeV2UnknownsTable(b *strings.Builder, assessment map[string]any, rating map[string]any) {
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

func factorList(factors []*v2FactorArtifacts, spec *model.Spec, fromReport string) string {
	if len(factors) == 0 {
		return ""
	}
	parts := make([]string, 0, len(factors))
	for _, factor := range factors {
		parts = append(parts, reportLink(fromReport, factorReportPath(factor.ID), factorTitle(spec, factor.ID))+" "+v2RatingLabel(spec, scopedMap(factor.Analysis, "localAndDescendantAnalysis")))
	}
	return strings.Join(parts, "; ")
}

func requirementFactorLinks(req *v2RequirementArtifacts, fromReport string) string {
	ids := requirementFactorIDs(req)
	if len(ids) == 0 {
		return ""
	}
	for i, id := range ids {
		parsed, err := parseRequirementFactorID(req.ID.DeclaringArea, id)
		if err != nil {
			continue
		}
		ids[i] = reportLink(fromReport, factorReportPath(parsed), id)
	}
	return strings.Join(ids, "; ")
}

func requirementFactorIDs(req *v2RequirementArtifacts) []string {
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

func v2RequirementRatingLabel(spec *model.Spec, rating map[string]any) string {
	if rating == nil {
		return ratingStatusTitle(string(RatingStatusNotRated))
	}
	if status := v2String(rating, "status"); status != "rated" {
		if status == "" {
			return ratingStatusTitle(string(RatingStatusNotRated))
		}
		return ratingStatusTitle(status)
	}
	if level := v2String(rating, "ratingLevelId"); level != "" {
		return ratingTitle(spec, level)
	}
	return ratingStatusTitle(string(RatingStatusRated))
}

// v2SubRatingCell renders the descendant-inclusive ("+ Sub-Factors Rating" /
// "+ Sub-Areas Rating") cell for a breakdown row. When the node has no
// descendants there is no roll-up distinct from its local rating, so it renders
// an em dash rather than repeating the local rating.
func v2SubRatingCell(spec *model.Spec, aggregate map[string]any, hasDescendants bool) string {
	if !hasDescendants {
		return "—"
	}
	return v2RatingLabel(spec, aggregate)
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

func factorIDJSON(id factorID) map[string]any {
	return map[string]any{"declaringAreaId": anyStrings(id.DeclaringArea), "factorPath": anyStrings(id.Path)}
}

func requirementIDJSON(id requirementID) map[string]any {
	return map[string]any{"declaringAreaId": anyStrings(id.DeclaringArea), "requirementName": id.Name}
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

func v2RatingLabel(spec *model.Spec, scope map[string]any) string {
	status := v2String(scope, "status")
	level := v2String(scope, "ratingLevelId")
	if status != "analyzed" || level == "" {
		if status == "" {
			return analysisStatusTitle(string(AnalysisStatusNotAnalyzed))
		}
		return analysisStatusTitle(status)
	}
	return ratingTitle(spec, level)
}

func v2AnalysisStatusPair(overall, local map[string]any) string {
	return analysisStatusTitle(v2String(overall, "status")) + " / " + analysisStatusTitle(v2String(local, "status"))
}

func v2ConfidencePair(overall, local map[string]any) string {
	return confidenceTitle(v2String(overall, "confidence")) + " / " + confidenceTitle(v2String(local, "confidence"))
}

func v2Summary(scope map[string]any) string {
	rationale := v2String(scope, "rationale")
	if rationale == "" {
		return "No analysis rationale was recorded."
	}
	return rationale
}

func v2String(payload map[string]any, field string) string {
	value, _ := payload[field].(string)
	return value
}

func v2ReceiptRating(analysis map[string]any) RatingResult {
	overall := scopedMap(analysis, "localAndDescendantAnalysis")
	if v2String(overall, "status") != "analyzed" {
		return RatingResult{Kind: "not-assessed", Rationale: v2String(overall, "statusReason")}
	}
	return RatingResult{Kind: "rated", Level: v2String(overall, "ratingLevelId"), Rationale: v2String(overall, "rationale")}
}

func markdownCell(s string) string {
	if s == "" {
		return ""
	}
	return strings.ReplaceAll(s, "|", "\\|")
}
