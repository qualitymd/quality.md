package evaluation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	md "github.com/qualitymd/quality.md/internal/markdown"
	"github.com/qualitymd/quality.md/internal/model"
)

func evaluationRenderableGaps(runAbs string) []RunGap {
	manifest, err := loadEvaluationManifest(runAbs)
	if err != nil {
		return []RunGap{{Kind: GapMissingEvaluationData, Ref: evaluationManifestPath, Detail: err.Error()}}
	}
	spec, err := loadRunModel(runAbs)
	if err != nil {
		return []RunGap{{Kind: GapUnreadableEvaluationData, Ref: ModelSnapshotFile, Detail: err.Error()}}
	}
	raw, gap := readRequiredEvaluationPayload(runAbs, "data/frame/evaluation-frame.json", DataKindEvaluationFrame)
	if gap != nil {
		return []RunGap{*gap}
	}
	artifacts, collectErr := collectEvaluationArtifacts(runAbs)
	if collectErr != nil {
		return []RunGap{{Kind: GapMalformedEvaluationData, Ref: "data", Detail: collectErr.Error()}}
	}
	artifacts.Manifest = manifest
	artifacts.Frame = raw
	plan, gap := resolveEvaluationReportPlan(artifacts)
	if gap != nil {
		return []RunGap{*gap}
	}
	if gaps := plannedCoverageGaps(spec, artifacts, plan); len(gaps) > 0 {
		return gaps
	}
	if gaps := adviceCoverageGaps(artifacts); len(gaps) > 0 {
		return gaps
	}
	return nil
}

func readRequiredEvaluationPayload(runAbs, rel string, want DataKind) (map[string]any, *RunGap) {
	raw, err := os.ReadFile(filepath.Join(runAbs, rel))
	if os.IsNotExist(err) {
		return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: rel, Detail: "required evaluation payload is missing"}
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
	runRel := workspaceRelativeRunPath(runAbs, artifacts.Manifest)
	reports := renderEvaluationReportTree(spec, artifacts, plan, runRel)
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
		RatingResult:           evaluationReceiptRating(plan.ScopedAreaAnalysis),
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
	Manifest              *EvaluationManifest
	RunLabel              string
	Frame                 map[string]any
	FindingRanking        map[string]any
	Recommendations       map[string]map[string]any
	RecommendationRanking map[string]any
	Areas                 map[string]*evaluationAreaArtifacts
	Factors               map[string]*evaluationFactorArtifacts
	Requirements          map[string]*evaluationRequirementArtifacts
}

type evaluationRenderedReport struct {
	Kind             string
	Path             string
	AreaID           []string
	FactorID         *factorID
	RequirementID    *requirementID
	RecommendationID string
	Content          string
}

type rankedFinding struct {
	Rank        int
	Total       int
	Key         string
	Selector    string
	FindingID   string
	Requirement *evaluationRequirementArtifacts
	Finding     map[string]any
	Ranking     map[string]any
}

type evaluationReportPlan struct {
	Frame              map[string]any
	RequestedScope     RunScope
	ScopedAreaID       []string
	FactorFilter       []factorID
	ScopedAreaAnalysis map[string]any
	ScopedAreaReport   *evaluationRenderedReport
}

func collectEvaluationArtifacts(runAbs string) (*evaluationArtifacts, error) {
	out := &evaluationArtifacts{
		RunLabel:        filepath.Base(runAbs),
		Recommendations: map[string]map[string]any{},
		Areas:           map[string]*evaluationAreaArtifacts{},
		Factors:         map[string]*evaluationFactorArtifacts{},
		Requirements:    map[string]*evaluationRequirementArtifacts{},
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
		return nil, fmt.Errorf("collecting evaluation data: %w", err)
	}
	return out, nil
}

type evaluationPayloadCollector func(*evaluationArtifacts, map[string]any) error

var evaluationPayloadCollectors = map[DataKind]evaluationPayloadCollector{
	DataKindEvaluationManifest:         collectEvaluationManifest,
	DataKindEvaluationFrame:            collectEvaluationFrame,
	DataKindAreaEvaluationFrame:        collectEvaluationAreaFrame,
	DataKindAreaAnalysis:               collectEvaluationAreaAnalysis,
	DataKindRequirementEvaluationFrame: collectEvaluationRequirementFrame,
	DataKindRequirementAssessment:      collectEvaluationRequirementAssessment,
	DataKindRequirementRating:          collectEvaluationRequirementRating,
	DataKindFactorAnalysisFrame:        collectEvaluationFactorFrame,
	DataKindFactorAnalysis:             collectEvaluationFactorAnalysis,
	DataKindFindingRanking:             collectFindingRanking,
	DataKindRecommendation:             collectRecommendation,
	DataKindRecommendationRanking:      collectRecommendationRanking,
}

func collectEvaluationManifest(out *evaluationArtifacts, payload map[string]any) error {
	raw, err := canonicalJSON(payload)
	if err != nil {
		return err
	}
	var manifest EvaluationManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return err
	}
	out.Manifest = &manifest
	return nil
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

func collectFindingRanking(out *evaluationArtifacts, payload map[string]any) error {
	out.FindingRanking = payload
	return nil
}

func collectRecommendation(out *evaluationArtifacts, payload map[string]any) error {
	id := recommendationID(payload)
	if !validRecommendationID(id) {
		return usagef("RecommendationResult.id is missing")
	}
	out.Recommendations[id] = payload
	return nil
}

func collectRecommendationRanking(out *evaluationArtifacts, payload map[string]any) error {
	out.RecommendationRanking = payload
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

func renderEvaluationReportTree(spec *model.Spec, artifacts *evaluationArtifacts, plan *evaluationReportPlan, runRel string) []evaluationRenderedReport {
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
			Content: renderEvaluationAreaReport(spec, artifacts, area, path, runRel),
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
			Content:  renderEvaluationFactorReport(spec, artifacts, factor, path, runRel),
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
			Content:       renderEvaluationRequirementReport(spec, artifacts, requirement, path, runRel),
		})
	}
	findingsIndex := evaluationRenderedReport{
		Kind:    string(ReportKindFindings),
		Path:    "findings.md",
		Content: renderEvaluationFindingsIndex(spec, artifacts, runRel),
	}
	recommendationReports := renderEvaluationRecommendationReports(spec, artifacts, runRel)
	reports = append(reports, findingsIndex)
	reports = append(reports, recommendationReports...)
	linkEvaluationReportPlan(plan, reports)
	run := evaluationRenderedReport{
		Kind:    string(ReportKindRun),
		Path:    "report.md",
		Content: renderEvaluationRunReport(spec, artifacts, plan, reports, "report.md", runRel),
	}
	return append([]evaluationRenderedReport{run}, reports...)
}

func renderEvaluationFindingsIndex(spec *model.Spec, artifacts *evaluationArtifacts, runRel string) string {
	var b strings.Builder
	data := reportSourceData(append([]string{evaluationManifestPath}, rankedFindingSourceData(artifacts, 0)...)...)
	renderReportHeader(&b, reportHeader{
		Type:       reportTypeFindingIndex,
		Heading:    "Findings",
		ReportPath: "findings.md",
		RunPath:    runRel,
		Run:        artifacts.Manifest,
		SummaryHead: []string{
			"Findings",
			"Highest concern severity",
		},
		SummaryRow: []string{
			fmt.Sprintf("%d findings", len(artifacts.rankedFindings())),
			highestFindingSeverityTitle(artifacts),
		},
		Contents: []reportContentLink{
			{Label: "Ranked findings", Anchor: "#ranked-findings"},
			{Label: "Primary source data", Anchor: "#primary-source-data"},
		},
	})
	b.WriteString("## Ranked findings\n\n")
	writeRankedFindingsTable(&b, spec, artifacts, "findings.md", 0)
	writePrimarySourceDataSection(&b, "findings.md", data)
	return b.String()
}

func renderEvaluationRecommendationReports(spec *model.Spec, artifacts *evaluationArtifacts, runRel string) []evaluationRenderedReport {
	var reports []evaluationRenderedReport
	recommendations := artifacts.rankedRecommendations()
	var index strings.Builder
	indexData := reportSourceData(append([]string{evaluationManifestPath}, rankedRecommendationSourceData(artifacts, 0)...)...)
	renderReportHeader(&index, reportHeader{
		Type:       reportTypeRecommendationIndex,
		Heading:    "Recommendations",
		ReportPath: "recommendations.md",
		RunPath:    runRel,
		Run:        artifacts.Manifest,
		SummaryHead: []string{
			"Recommendations",
			"Highest impact",
			"Coverage",
		},
		SummaryRow: []string{
			fmt.Sprintf("%d recommendations", len(recommendations)),
			highestRecommendationImpactTitle(recommendations),
			recommendationCoverageSummary(artifacts),
		},
		Contents: []reportContentLink{
			{Label: "Ranked recommendations", Anchor: "#ranked-recommendations"},
			{Label: "Coverage", Anchor: "#coverage"},
			{Label: "Primary source data", Anchor: "#primary-source-data"},
		},
	})
	index.WriteString("## Ranked recommendations\n\n")
	writeRecommendationIndexTable(&index, spec, artifacts, "recommendations.md")
	writeAdviceCoverageSummary(&index, artifacts)
	writePrimarySourceDataSection(&index, "recommendations.md", indexData)
	reports = append(reports, evaluationRenderedReport{
		Kind:    string(ReportKindAdviceIndex),
		Path:    "recommendations.md",
		Content: index.String(),
	})
	for _, item := range recommendations {
		title := firstString(item.Recommendation, "title")
		path := recommendationReportPath(item.Rank, title)
		reports = append(reports, evaluationRenderedReport{
			Kind:             string(ReportKindAdvice),
			Path:             path,
			RecommendationID: recommendationID(item.Recommendation),
			Content:          renderEvaluationRecommendationReport(artifacts, item, runRel),
		})
	}
	return reports
}

type rankedRecommendation struct {
	Rank           int
	Recommendation map[string]any
	Ranking        map[string]any
}

const (
	reportTypeEvaluationOverview    = "Evaluation Overview Report"
	reportTypeAreaEvaluation        = "Area Evaluation Report"
	reportTypeFactorEvaluation      = "Factor Evaluation Report"
	reportTypeRequirementEvaluation = "Requirement Evaluation Report"
	reportTypeFindingIndex          = "Finding Index Report"
	reportTypeRecommendationIndex   = "Recommendation Index Report"
	reportTypeRecommendation        = "Recommendation Report"
)

type reportHeader struct {
	Type        string
	Heading     string
	ReportPath  string
	RunPath     string
	Run         *EvaluationManifest
	Context     []string
	SummaryHead []string
	SummaryRow  []string
	Contents    []reportContentLink
}

type reportContentLink struct {
	Label  string
	Anchor string
}

func renderReportHeader(b *strings.Builder, header reportHeader) {
	b.WriteString(md.Frontmatter(
		md.FrontmatterField{Name: "type", Value: header.Type},
		md.FrontmatterField{Name: "title", Value: header.Heading},
	))
	b.WriteString("# " + header.Heading + "\n\n")
	writeEvaluationLinks(b, header.ReportPath, header.RunPath)
	if line := reportRunLine(header.ReportPath, header.Run); line != "" {
		b.WriteString(line + "\n\n")
	}
	for _, line := range header.Context {
		if line != "" {
			b.WriteString(line + "\n\n")
		}
	}
	if len(header.SummaryHead) > 0 {
		b.WriteString("## Key details\n\n")
		writeReportSummaryTable(b, header.SummaryHead, header.SummaryRow)
	}
	writeContentsSection(b, header.Contents)
}

func reportRunLine(reportPath string, manifest *EvaluationManifest) string {
	if manifest == nil {
		return ""
	}
	label := runLabel(manifest)
	if reportPath != "report.md" {
		label = reportLink(reportPath, "report.md", label)
	}
	evaluationID := manifest.EvaluationID
	if evaluationID == "" {
		evaluationID = "—"
	}
	created := manifest.CreatedAt
	if created == "" {
		created = "—"
	}
	return "Run: " + label + " - Evaluation ID: " + md.Code(evaluationID) + " - Created: " + created + " - Scope: " + requestedScopeLabel(manifest.RequestedScope)
}

func writeReportSummaryTable(b *strings.Builder, headers, row []string) {
	b.WriteString(md.TableRow(headers...))
	separator := make([]string, len(headers))
	for i := range separator {
		separator[i] = "---"
	}
	b.WriteString(md.TableRow(separator...))
	b.WriteString(md.TableRow(row...))
	b.WriteString("\n")
}

func reportSourceData(paths ...string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(paths))
	for _, path := range paths {
		path = filepath.ToSlash(path)
		if path == "" {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		out = append(out, path)
	}
	return out
}

func writePrimarySourceDataSection(b *strings.Builder, reportPath string, paths []string) {
	switch content := b.String(); {
	case strings.HasSuffix(content, "\n\n"):
	case strings.HasSuffix(content, "\n"):
		b.WriteString("\n")
	default:
		b.WriteString("\n\n")
	}
	b.WriteString("## Primary source data\n\n")
	for _, path := range reportSourceData(paths...) {
		b.WriteString("- " + reportLink(reportPath, path, path) + "\n")
	}
	b.WriteString("\n")
}

func runReportSourceData(plan *evaluationReportPlan) []string {
	paths := []string{
		evaluationManifestPath,
		areaDataPath(plan.ScopedAreaID, "area-analysis-result.json"),
		"data/advice/finding-ranking-result.json",
		"data/advice/recommendation-ranking-result.json",
	}
	return reportSourceData(paths...)
}

func rankedFindingSourceData(artifacts *evaluationArtifacts, limit int) []string {
	paths := []string{"data/advice/finding-ranking-result.json"}
	for i, row := range artifacts.rankedFindings() {
		if limit > 0 && i >= limit {
			break
		}
		path, _, ok := strings.Cut(row.Key, "#")
		if ok {
			paths = append(paths, path)
		}
	}
	return reportSourceData(paths...)
}

func rankedRecommendationSourceData(artifacts *evaluationArtifacts, limit int) []string {
	paths := []string{"data/advice/recommendation-ranking-result.json"}
	for i, item := range artifacts.rankedRecommendations() {
		if limit > 0 && i >= limit {
			break
		}
		paths = append(paths, recommendationDataPath(recommendationID(item.Recommendation)))
	}
	return reportSourceData(paths...)
}

func areaReportSourceData(artifacts *evaluationArtifacts, area *evaluationAreaArtifacts) []string {
	paths := []string{evaluationManifestPath}
	if area.Analysis != nil {
		paths = append(paths, areaDataPath(area.ID, "area-analysis-result.json"))
	}
	paths = append(paths, "data/advice/finding-ranking-result.json")
	paths = append(paths, "data/advice/recommendation-ranking-result.json")
	for _, req := range artifacts.requirementsForArea(area.ID) {
		if req.Rating != nil {
			paths = append(paths, requirementDataPath(req.ID, "requirement-rating-result.json"))
		}
		if req.Assessment != nil {
			paths = append(paths, requirementDataPath(req.ID, "requirement-assessment-result.json"))
		}
	}
	return reportSourceData(paths...)
}

func factorReportSourceData(artifacts *evaluationArtifacts, factor *evaluationFactorArtifacts) []string {
	paths := []string{evaluationManifestPath}
	if factor.Analysis != nil {
		paths = append(paths, factorDataPath(factor.ID, "factor-analysis-result.json"))
	}
	for _, req := range artifacts.requirementsForFactor(factor.ID) {
		if req.Rating != nil {
			paths = append(paths, requirementDataPath(req.ID, "requirement-rating-result.json"))
		}
		if req.Assessment != nil {
			paths = append(paths, requirementDataPath(req.ID, "requirement-assessment-result.json"))
		}
	}
	return reportSourceData(paths...)
}

func requirementReportSourceData(req *evaluationRequirementArtifacts) []string {
	paths := []string{evaluationManifestPath}
	if req.Assessment != nil {
		paths = append(paths, requirementDataPath(req.ID, "requirement-assessment-result.json"))
	}
	if req.Rating != nil {
		paths = append(paths, requirementDataPath(req.ID, "requirement-rating-result.json"))
	}
	paths = append(paths, "data/advice/finding-ranking-result.json")
	return reportSourceData(paths...)
}

func writeContentsSection(b *strings.Builder, links []reportContentLink) {
	if len(links) < 2 {
		return
	}
	b.WriteString("## Contents\n\n")
	for _, link := range links {
		b.WriteString("- " + md.Link(link.Label, link.Anchor) + "\n")
	}
	b.WriteString("\n")
}

func writeEvaluationLinks(b *strings.Builder, reportPath, runPath string) {
	links := []string{
		reportLink(reportPath, "report.md", "report.md"),
		reportLink(reportPath, "findings.md", "findings.md"),
		reportLink(reportPath, "recommendations.md", "recommendations.md"),
		md.Link("glossary.md", glossaryLinkTarget(reportPath, runPath)),
	}
	b.WriteString("> **Evaluation links:** " + strings.Join(links, " | ") + "\n\n")
}

func highestFindingSeverityTitle(artifacts *evaluationArtifacts) string {
	bestRank := len(findingSeverityValues.Values)
	best := ""
	for _, row := range artifacts.rankedFindings() {
		severity := findingConcernSeverity(row.Finding)
		rank, ok := enumRank(findingSeverityValues, severity)
		if !ok {
			continue
		}
		if rank < bestRank {
			bestRank = rank
			best = severity
		}
	}
	if best == "" {
		return "—"
	}
	return findingSeverityTitle(best)
}

func highestRecommendationImpactTitle(items []rankedRecommendation) string {
	bestRank := len(recommendationImpactValues.Values)
	best := ""
	for _, item := range items {
		impact := firstString(item.Recommendation, "impact")
		rank, ok := enumRank(recommendationImpactValues, impact)
		if !ok {
			continue
		}
		if rank < bestRank {
			bestRank = rank
			best = impact
		}
	}
	if best == "" {
		return "—"
	}
	return impactTitle(best)
}

func recommendationCoverageSummary(artifacts *evaluationArtifacts) string {
	coverage := objectSlice(artifacts.RecommendationRanking["findingCoverage"])
	addressed := 0
	notDriving := 0
	for _, item := range coverage {
		switch firstString(item, "disposition") {
		case string(FindingCoverageAddressedByRecommendation):
			addressed++
		case string(FindingCoverageNotAdviceDriving):
			notDriving++
		}
	}
	return fmt.Sprintf("%s: %d / %s: %d",
		findingCoverageDispositionTitle(string(FindingCoverageAddressedByRecommendation)),
		addressed,
		findingCoverageDispositionTitle(string(FindingCoverageNotAdviceDriving)),
		notDriving,
	)
}

func (a *evaluationArtifacts) rankedRecommendations() []rankedRecommendation {
	ranking := objectSlice(a.RecommendationRanking["orderedRecommendations"])
	out := make([]rankedRecommendation, 0, len(ranking))
	used := map[string]struct{}{}
	for _, item := range ranking {
		id := firstString(item, "recommendationRef")
		if !validRecommendationID(id) {
			continue
		}
		rec := a.Recommendations[id]
		if rec == nil {
			continue
		}
		rank, _ := rankField(item)
		out = append(out, rankedRecommendation{Rank: rank, Recommendation: rec, Ranking: item})
		used[id] = struct{}{}
	}
	for _, id := range sortedRecommendationIDs(a.Recommendations) {
		if _, ok := used[id]; ok {
			continue
		}
		out = append(out, rankedRecommendation{Rank: len(out) + 1, Recommendation: a.Recommendations[id]})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Rank != out[j].Rank {
			return out[i].Rank < out[j].Rank
		}
		return recommendationID(out[i].Recommendation) < recommendationID(out[j].Recommendation)
	})
	return out
}

func sortedRecommendationIDs(recommendations map[string]map[string]any) []string {
	ids := make([]string, 0, len(recommendations))
	for id := range recommendations {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func renderEvaluationRecommendationReport(artifacts *evaluationArtifacts, item rankedRecommendation, runRel string) string {
	rec := item.Recommendation
	var b strings.Builder
	reportPath := recommendationReportPath(item.Rank, firstString(rec, "title"))
	id := recommendationID(rec)
	recDataPath := recommendationDataPath(id)
	rankingPath := "data/advice/recommendation-ranking-result.json"
	data := reportSourceData(evaluationManifestPath, recDataPath, rankingPath)
	renderReportHeader(&b, reportHeader{
		Type:       reportTypeRecommendation,
		Heading:    "Recommendation: " + firstString(rec, "title"),
		ReportPath: reportPath,
		RunPath:    runRel,
		Run:        artifacts.Manifest,
		SummaryHead: []string{
			"#",
			"ID",
			"Impact",
			"Confidence",
			"Reference",
		},
		SummaryRow: []string{
			fmt.Sprintf("%d", item.Rank),
			md.Code(id),
			impactTitle(firstString(rec, "impact")),
			confidenceTitle(firstString(rec, "confidence")),
			md.Code(recommendationArtifactRef(artifacts.Manifest, id)),
		},
		Contents: []reportContentLink{
			{Label: "Description", Anchor: "#description"},
			{Label: "Background", Anchor: "#background"},
			{Label: "Expected value", Anchor: "#expected-value"},
			{Label: "Done criterion", Anchor: "#done-criterion"},
			{Label: "Ranking rationale", Anchor: "#ranking-rationale"},
			{Label: "Trace", Anchor: "#trace"},
			{Label: "Primary source data", Anchor: "#primary-source-data"},
		},
	})
	writeRecommendationSection(&b, "Description", firstString(rec, "description"))
	writeRecommendationSection(&b, "Background", firstString(rec, "background"))
	writeRecommendationSection(&b, "Expected value", firstString(rec, "expectedValue"))
	writeRecommendationSection(&b, "Done criterion", firstString(rec, "doneCriterion"))
	writeRecommendationSection(&b, "Ranking rationale", firstString(item.Ranking, "rationale"))
	b.WriteString("## Trace\n\n")
	if refs := objectSlice(rec["traceRefs"]); len(refs) == 0 {
		b.WriteString("(none recorded)\n\n")
	} else {
		for _, ref := range refs {
			b.WriteString("- " + md.Code(compactJSON(ref)) + "\n")
		}
		b.WriteString("\n")
	}
	writePrimarySourceDataSection(&b, reportPath, data)
	return b.String()
}

func recommendationArtifactRef(manifest *EvaluationManifest, recommendationID string) string {
	if manifest == nil || manifest.EvaluationID == "" || !validRecommendationID(recommendationID) {
		return ""
	}
	return fmt.Sprintf("evaluation:%s/recommendation/%s", manifest.EvaluationID, recommendationID)
}

func writeRecommendationSection(b *strings.Builder, title, body string) {
	b.WriteString("## " + title + "\n\n")
	if body == "" {
		b.WriteString("(not recorded)\n\n")
		return
	}
	b.WriteString(body + "\n\n")
}

func renderEvaluationRunReport(spec *model.Spec, artifacts *evaluationArtifacts, plan *evaluationReportPlan, reports []evaluationRenderedReport, reportPath, runRel string) string {
	title := evaluationRunReportTitle(spec, plan)
	scopedArea := scopedMap(plan.ScopedAreaAnalysis, "localAndDescendantAnalysis")
	localArea := scopedMap(plan.ScopedAreaAnalysis, "localAnalysis")
	var b strings.Builder
	data := runReportSourceData(plan)
	renderRunReportHeader(&b, artifacts, title)
	writeEvaluationLinks(&b, reportPath, runRel)
	b.WriteString("## Summary\n\n")
	b.WriteString(evaluationSummary(scopedArea))
	b.WriteString("\n\n## Key details\n\n")
	writeRunReportKeyDetails(&b, spec, artifacts, plan, scopedArea, localArea)
	writeContentsSection(&b, []reportContentLink{
		{Label: "Model evaluation", Anchor: "#model-evaluation"},
		{Label: "Top findings", Anchor: "#top-findings"},
		{Label: "Top recommendations", Anchor: "#top-recommendations"},
		{Label: "Primary source data", Anchor: "#primary-source-data"},
	})
	writeAreaFactorBreakdownSection(&b, "Model evaluation", spec, artifacts, plan.ScopedAreaID, reportPath)
	writeRunReportCoverageNote(&b, reports, reportPath)
	b.WriteString("## Top findings\n\n")
	writeFullFindingsReportLink(&b, reportPath, artifacts)
	writeRankedFindingsTable(&b, spec, artifacts, reportPath, 10)
	b.WriteString("## Top recommendations\n\n")
	writeFullRecommendationsReportLink(&b, reportPath, artifacts)
	writeTopRecommendationsTable(&b, spec, artifacts, reportPath, 10)
	writePrimarySourceDataSection(&b, reportPath, data)
	return b.String()
}

func renderRunReportHeader(b *strings.Builder, artifacts *evaluationArtifacts, title string) {
	b.WriteString(md.Frontmatter(
		md.FrontmatterField{Name: "type", Value: reportTypeEvaluationOverview},
		md.FrontmatterField{Name: "title", Value: title},
		md.FrontmatterField{Name: "evaluationId", Value: runReportEvaluationID(artifacts)},
		md.FrontmatterField{Name: "created", Value: runReportCreated(artifacts)},
		md.FrontmatterField{Name: "model", Value: runReportModel(artifacts)},
		md.FrontmatterField{Name: "run", Value: runReportRunLabel(artifacts)},
	))
	b.WriteString("# " + title + "\n\n")
}

func runReportRunLabel(artifacts *evaluationArtifacts) string {
	if artifacts != nil && artifacts.RunLabel != "" && artifacts.RunLabel != "." && artifacts.RunLabel != string(filepath.Separator) {
		return artifacts.RunLabel
	}
	if artifacts != nil && artifacts.Manifest != nil {
		return runLabel(artifacts.Manifest)
	}
	return ""
}

func runLabel(manifest *EvaluationManifest) string {
	if manifest == nil {
		return ""
	}
	if manifest.Run.Label != "" {
		return manifest.Run.Label
	}
	if manifest.Run.Number > 0 {
		return fmt.Sprintf("Run %04d", manifest.Run.Number)
	}
	return ""
}

func runReportEvaluationID(artifacts *evaluationArtifacts) string {
	if artifacts != nil && artifacts.Manifest != nil {
		return artifacts.Manifest.EvaluationID
	}
	return ""
}

func runReportModel(artifacts *evaluationArtifacts) string {
	if artifacts != nil && artifacts.Manifest != nil {
		return artifacts.Manifest.Model
	}
	return ""
}

func runReportCreated(artifacts *evaluationArtifacts) string {
	if artifacts != nil && artifacts.Manifest != nil {
		return artifacts.Manifest.CreatedAt
	}
	return ""
}

func writeRunReportKeyDetails(b *strings.Builder, spec *model.Spec, artifacts *evaluationArtifacts, plan *evaluationReportPlan, scopedArea, localArea map[string]any) {
	b.WriteString(md.TableRow("Overall rating", "Confidence", "Scope", "Findings", "Recommendations"))
	b.WriteString(md.TableRow("---", "---", "---", "---", "---"))
	b.WriteString(md.TableRow(
		evaluationRatingLabel(spec, scopedArea),
		emDashIfEmpty(confidenceTitle(evaluationString(scopedArea, "confidence"))),
		runReportScopeLabel(spec, plan),
		fmt.Sprintf("%d total", len(artifacts.rankedFindings())),
		fmt.Sprintf("%d total", len(artifacts.rankedRecommendations())),
	))
	b.WriteString("\n")
	_ = localArea
}

func runReportScopeLabel(spec *model.Spec, plan *evaluationReportPlan) string {
	if plan == nil {
		return "—"
	}
	area := areaTitle(spec, plan.ScopedAreaID)
	if len(plan.FactorFilter) == 0 {
		return "Full evaluation of " + area
	}
	factors := make([]string, 0, len(plan.FactorFilter))
	for _, factor := range plan.FactorFilter {
		factors = append(factors, factorTitle(spec, factor))
	}
	return "Evaluation of " + area + " for " + strings.Join(factors, ", ")
}

func writeFullFindingsReportLink(b *strings.Builder, fromReport string, artifacts *evaluationArtifacts) {
	total := len(artifacts.rankedFindings())
	summary := findingInlineCountSummary(artifacts)
	if summary == "" {
		fmt.Fprintf(b, "**Full findings report:** %s (%d total)\n\n", reportLink(fromReport, "findings.md", "findings.md"), total)
		return
	}
	fmt.Fprintf(b, "**Full findings report:** %s (%d total: %s)\n\n", reportLink(fromReport, "findings.md", "findings.md"), total, summary)
}

func findingInlineCountSummary(artifacts *evaluationArtifacts) string {
	counts := map[string]int{}
	severities := map[string]map[string]int{}
	for _, row := range artifacts.rankedFindings() {
		typ := firstString(row.Finding, "type")
		if typ == "" {
			continue
		}
		counts[typ]++
		sev := findingConcernSeverity(row.Finding)
		if sev == "" {
			continue
		}
		if severities[typ] == nil {
			severities[typ] = map[string]int{}
		}
		severities[typ][sev]++
	}
	parts := make([]string, 0, len(findingTypeValues.Values))
	for _, value := range findingTypeValues.Values {
		typ := string(value.Value)
		count := counts[typ]
		if count == 0 {
			continue
		}
		part := inlineMarkedCount(value.Marker, value.Label, count, true)
		if typ == string(FindingTypeGap) || typ == string(FindingTypeRisk) {
			if severitySummary := findingInlineSeveritySummary(severities[typ]); severitySummary != "" {
				part += ": " + severitySummary
			}
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "; ")
}

func findingInlineSeveritySummary(severities map[string]int) string {
	parts := make([]string, 0, len(findingSeverityValues.Values))
	for _, value := range findingSeverityValues.Values {
		count := severities[string(value.Value)]
		if count == 0 {
			continue
		}
		parts = append(parts, inlineMarkedCount(value.Marker, value.Label, count, false))
	}
	return strings.Join(parts, ", ")
}

func writeFullRecommendationsReportLink(b *strings.Builder, fromReport string, artifacts *evaluationArtifacts) {
	total := len(artifacts.rankedRecommendations())
	summary := recommendationImpactInlineCountSummary(artifacts)
	if summary == "" {
		fmt.Fprintf(b, "**Full recommendations report:** %s (%d total)\n\n", reportLink(fromReport, "recommendations.md", "recommendations.md"), total)
		return
	}
	fmt.Fprintf(b, "**Full recommendations report:** %s (%d total; impact: %s)\n\n", reportLink(fromReport, "recommendations.md", "recommendations.md"), total, summary)
}

func recommendationImpactInlineCountSummary(artifacts *evaluationArtifacts) string {
	counts := map[string]int{}
	for _, item := range artifacts.rankedRecommendations() {
		impact := firstString(item.Recommendation, "impact")
		if impact == "" {
			continue
		}
		counts[impact]++
	}
	parts := make([]string, 0, len(recommendationImpactValues.Values))
	for _, value := range recommendationImpactValues.Values {
		count := counts[string(value.Value)]
		if count == 0 {
			continue
		}
		parts = append(parts, inlineMarkedCount(value.Marker, value.Label, count, false))
	}
	return strings.Join(parts, ", ")
}

func inlineMarkedCount(marker, label string, count int, plural bool) string {
	if plural && count != 1 {
		label += "s"
	}
	if marker == "" {
		return fmt.Sprintf("%d %s", count, label)
	}
	return fmt.Sprintf("%s %d %s", marker, count, label)
}

func writeRunReportCoverageNote(b *strings.Builder, reports []evaluationRenderedReport, reportPath string) {
	if reportForRootArea(reports) != nil {
		return
	}
	b.WriteString("Root area was not evaluated in this run")
	if reportPath != "report.md" {
		b.WriteString(" for " + reportLink(reportPath, "report.md", "the run overview"))
	}
	b.WriteString(".\n\n")
}

func renderEvaluationAreaReport(spec *model.Spec, artifacts *evaluationArtifacts, area *evaluationAreaArtifacts, reportPath, runRel string) string {
	title := areaTitle(spec, area.ID)
	local := scopedMap(area.Analysis, "localAnalysis")
	overall := scopedMap(area.Analysis, "localAndDescendantAnalysis")
	var b strings.Builder
	data := areaReportSourceData(artifacts, area)
	renderReportHeader(&b, reportHeader{
		Type:       reportTypeAreaEvaluation,
		Heading:    "Area: " + title,
		ReportPath: reportPath,
		RunPath:    runRel,
		Run:        artifacts.Manifest,
		Context: []string{
			evaluationAreaTrailLine(spec, artifacts, area.ID, reportPath),
		},
		SummaryHead: []string{
			"Overall rating",
			"Local rating",
			"Confidence",
		},
		SummaryRow: []string{
			evaluationRatingLabel(spec, overall),
			evaluationRatingLabel(spec, local),
			evaluationConfidencePair(overall, local),
		},
		Contents: []reportContentLink{
			{Label: "Summary", Anchor: "#summary"},
			{Label: "Area / factor breakdown", Anchor: "#area--factor-breakdown"},
			{Label: "Requirements", Anchor: "#requirements"},
			{Label: "Limits and incomplete inputs", Anchor: "#limits-and-incomplete-inputs"},
			{Label: "Primary source data", Anchor: "#primary-source-data"},
		},
	})
	b.WriteString("## Summary\n\n")
	b.WriteString(evaluationSummary(overall))
	b.WriteString("\n\n")
	writeAreaFactorBreakdownSection(&b, "Area / factor breakdown", spec, artifacts, area.ID, reportPath)
	b.WriteString("## Requirements\n\n")
	b.WriteString("| Requirement | Rating | Status | Factors |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	if requirements := artifacts.requirementsForArea(area.ID); len(requirements) == 0 {
		b.WriteString("| (no local requirements) | — | — | — |\n\n")
	} else {
		for _, req := range requirements {
			b.WriteString(md.TableRow(reportLink(reportPath, requirementReportPath(req.ID), requirementTitle(spec, req.ID)), evaluationRequirementRatingLabel(spec, req.Rating), assessmentStatusTitle(evaluationString(req.Assessment, "status")), requirementFactorLinks(req, reportPath)))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Limits and incomplete inputs\n\n")
	writeEvaluationLimitsTable(&b, local, overall)
	writePrimarySourceDataSection(&b, reportPath, data)
	return b.String()
}

func renderEvaluationFactorReport(spec *model.Spec, artifacts *evaluationArtifacts, factor *evaluationFactorArtifacts, reportPath, runRel string) string {
	local := scopedMap(factor.Analysis, "localAnalysis")
	overall := scopedMap(factor.Analysis, "localAndDescendantAnalysis")
	title := factorTitle(spec, factor.ID)
	var b strings.Builder
	data := factorReportSourceData(artifacts, factor)
	renderReportHeader(&b, reportHeader{
		Type:       reportTypeFactorEvaluation,
		Heading:    "Factor: " + title,
		ReportPath: reportPath,
		RunPath:    runRel,
		Run:        artifacts.Manifest,
		Context: []string{
			evaluationAreaTrailLine(spec, artifacts, factor.ID.DeclaringArea, reportPath),
			evaluationFactorTrailLine(spec, factor.ID, reportPath),
		},
		SummaryHead: []string{
			"Overall rating",
			"Local rating",
			"Status",
			"Confidence",
		},
		SummaryRow: []string{
			evaluationRatingLabel(spec, overall),
			evaluationRatingLabel(spec, local),
			evaluationAnalysisStatusPair(overall, local),
			evaluationConfidencePair(overall, local),
		},
		Contents: []reportContentLink{
			{Label: "Summary", Anchor: "#summary"},
			{Label: "Requirements", Anchor: "#requirements"},
			{Label: "Sub-factors", Anchor: "#sub-factors"},
			{Label: "Limits and incomplete inputs", Anchor: "#limits-and-incomplete-inputs"},
			{Label: "Primary source data", Anchor: "#primary-source-data"},
		},
	})
	b.WriteString("## Summary\n\n")
	b.WriteString(evaluationSummary(overall))
	b.WriteString("\n\n## Requirements\n\n")
	b.WriteString("| Requirement | Rating | Status |\n")
	b.WriteString("| --- | --- | --- |\n")
	if requirements := artifacts.requirementsForFactor(factor.ID); len(requirements) == 0 {
		b.WriteString("| (no direct requirements) | — | — |\n\n")
	} else {
		for _, req := range requirements {
			b.WriteString(md.TableRow(reportLink(reportPath, requirementReportPath(req.ID), requirementTitle(spec, req.ID)), evaluationRequirementRatingLabel(spec, req.Rating), assessmentStatusTitle(evaluationString(req.Assessment, "status"))))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Sub-factors\n\n")
	b.WriteString("| Factor | Path | Local rating | + Sub-factors rating |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	if children := artifacts.childFactors(factor.ID); len(children) == 0 {
		b.WriteString("| (no sub-factors) | — | — | — |\n\n")
	} else {
		for _, child := range children {
			childLocal := scopedMap(child.Analysis, "localAnalysis")
			childOverall := scopedMap(child.Analysis, "localAndDescendantAnalysis")
			b.WriteString(md.TableRow(reportLink(reportPath, factorReportPath(child.ID), factorTitle(spec, child.ID)), md.Code(factorDisplayPath(child.ID)), evaluationRatingLabel(spec, childLocal), evaluationSubRatingCell(spec, childOverall, len(artifacts.childFactors(child.ID)) > 0)))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Limits and incomplete inputs\n\n")
	writeEvaluationLimitsTable(&b, local, overall)
	writePrimarySourceDataSection(&b, reportPath, data)
	return b.String()
}

func renderEvaluationRequirementReport(spec *model.Spec, artifacts *evaluationArtifacts, req *evaluationRequirementArtifacts, reportPath, runRel string) string {
	title := requirementTitle(spec, req.ID)
	var b strings.Builder
	data := requirementReportSourceData(req)
	renderReportHeader(&b, reportHeader{
		Type:       reportTypeRequirementEvaluation,
		Heading:    "Requirement: " + title,
		ReportPath: reportPath,
		RunPath:    runRel,
		Run:        artifacts.Manifest,
		Context: []string{
			evaluationAreaTrailLine(spec, artifacts, req.ID.DeclaringArea, reportPath),
			evaluationRequirementFactorsLine(req, reportPath),
		},
		SummaryHead: []string{
			"Rating",
			"Assessment",
			"Confidence",
		},
		SummaryRow: []string{
			evaluationRequirementRatingLabel(spec, req.Rating),
			assessmentStatusTitle(evaluationString(req.Assessment, "status")),
			evaluationRequirementConfidencePair(req),
		},
		Contents: []reportContentLink{
			{Label: "Summary", Anchor: "#summary"},
			{Label: "Findings summary", Anchor: "#findings-summary"},
			{Label: "Finding details", Anchor: "#finding-details"},
			{Label: "Unknowns and missing evidence", Anchor: "#unknowns-and-missing-evidence"},
			{Label: "Primary source data", Anchor: "#primary-source-data"},
		},
	})
	b.WriteString("## Summary\n\n")
	if summary := evaluationString(req.Assessment, "evidenceSummary"); summary != "" {
		b.WriteString(summary)
	} else if rationale := evaluationString(req.Rating, "rationale"); rationale != "" {
		b.WriteString(rationale)
	} else {
		b.WriteString("No assessment summary was recorded.")
	}
	b.WriteString("\n\n## Findings summary\n\n")
	writeEvaluationFindingsTable(&b, req.Assessment)
	b.WriteString("## Finding details\n\n")
	writeEvaluationFindingDetails(&b, artifacts, req)
	b.WriteString("## Unknowns and missing evidence\n\n")
	writeEvaluationUnknownsTable(&b, req.Assessment, req.Rating)
	writePrimarySourceDataSection(&b, reportPath, data)
	return b.String()
}

func evaluationOutputResult(artifacts *evaluationArtifacts, plan *evaluationReportPlan, reports []evaluationRenderedReport) map[string]any {
	reportOutputs := make([]any, 0, len(reports))
	reportsByArea := map[string][]any{}
	for _, report := range reports {
		ref := evaluationReportRef(report)
		reportOutputs = append(reportOutputs, ref)
		if report.Kind == string(ReportKindArea) || report.Kind == string(ReportKindFactor) || report.Kind == string(ReportKindRequirement) {
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
		"schemaVersion":         SchemaVersion,
		"kind":                  string(DataKindEvaluationOutput),
		"runReportRef":          evaluationReportRef(evaluationRenderedReport{Kind: string(ReportKindRun), Path: "report.md"}),
		"scopedAreaAnalysisRef": routineRef(DataKindAreaAnalysis, map[string]any{"areaId": model.AreaPath(plan.ScopedAreaID).Reference()}, "localAndDescendantAnalysis"),
		"areaOutputs":           areaOutputs,
		"reportOutputs":         reportOutputs,
	}
	if root := artifacts.Areas[areaKey(nil)]; root != nil && root.Analysis != nil {
		output["rootAreaAnalysisRef"] = routineRef(DataKindAreaAnalysis, map[string]any{"areaId": model.AreaPath{}.Reference()}, "localAndDescendantAnalysis")
	}
	return output
}

func evaluationReportRef(report evaluationRenderedReport) map[string]any {
	ref := map[string]any{"kind": report.Kind, "path": filepath.ToSlash(report.Path)}
	if report.Kind == string(ReportKindArea) || report.Kind == string(ReportKindFactor) || report.Kind == string(ReportKindRequirement) {
		ref["areaId"] = model.AreaPath(report.AreaID).Reference()
	}
	if report.FactorID != nil {
		ref["factorId"] = factorIDJSON(*report.FactorID)
	}
	if report.RequirementID != nil {
		ref["requirementId"] = requirementIDJSON(*report.RequirementID)
	}
	if report.RecommendationID != "" {
		ref["recommendationId"] = report.RecommendationID
	}
	return ref
}

func resolveEvaluationReportPlan(artifacts *evaluationArtifacts) (*evaluationReportPlan, *RunGap) {
	if artifacts.Manifest == nil {
		return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: evaluationManifestPath, Detail: "required EvaluationManifest payload is missing"}
	}
	if artifacts.Frame == nil {
		return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: "data/frame/evaluation-frame.json", Detail: "required evaluation payload is missing"}
	}
	areaID, err := areaIDFrom(artifacts.Manifest.PlannedScope.AreaID)
	if err != nil {
		return nil, &RunGap{Kind: GapIncompleteEvaluationData, Ref: evaluationManifestPath, Detail: err.Error()}
	}
	var factors []factorID
	for _, ref := range artifacts.Manifest.PlannedScope.FactorFilter {
		id, err := factorIDFrom(ref)
		if err != nil {
			return nil, &RunGap{Kind: GapIncompleteEvaluationData, Ref: evaluationManifestPath, Detail: err.Error()}
		}
		if !sameStrings(id.DeclaringArea, areaID) {
			return nil, &RunGap{Kind: GapIncompleteEvaluationData, Ref: evaluationManifestPath, Detail: fmt.Sprintf("factor %s does not belong to planned area %s", ref, artifacts.Manifest.PlannedScope.AreaID)}
		}
		if factor := artifacts.Factors[factorKey(id)]; factor == nil || factor.Analysis == nil {
			return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: factorDataPath(id, "factor-analysis-result.json"), Detail: "required planned factor analysis payload is missing"}
		}
		factors = append(factors, id)
	}
	area := artifacts.Areas[areaKey(areaID)]
	if area == nil || area.Analysis == nil {
		return nil, &RunGap{Kind: GapMissingEvaluationData, Ref: areaDataPath(areaID, "area-analysis-result.json"), Detail: "required scoped area analysis payload is missing"}
	}
	return &evaluationReportPlan{
		Frame:              artifacts.Frame,
		RequestedScope:     artifacts.Manifest.RequestedScope,
		ScopedAreaID:       areaID,
		FactorFilter:       factors,
		ScopedAreaAnalysis: area.Analysis,
	}, nil
}

func linkEvaluationReportPlan(plan *evaluationReportPlan, reports []evaluationRenderedReport) {
	for i := range reports {
		report := &reports[i]
		if report.Kind == string(ReportKindArea) && sameStrings(report.AreaID, plan.ScopedAreaID) {
			plan.ScopedAreaReport = report
			return
		}
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

func evaluationRunReportTitle(spec *model.Spec, plan *evaluationReportPlan) string {
	label := areaTitle(spec, plan.ScopedAreaID)
	if len(plan.FactorFilter) == 0 {
		return "Quality evaluation - " + label
	}
	factors := make([]string, 0, len(plan.FactorFilter))
	for _, factor := range plan.FactorFilter {
		factors = append(factors, factorTitle(spec, factor))
	}
	return "Quality evaluation - " + label + " (" + strings.Join(factors, ", ") + ")"
}

func requestedScopeLabel(scope RunScope) string {
	if scope.AreaID == "" && len(scope.FactorFilter) == 0 {
		return "full evaluation"
	}
	var parts []string
	if scope.AreaID != "" {
		parts = append(parts, scope.AreaID)
	}
	if len(scope.FactorFilter) > 0 {
		parts = append(parts, strings.Join(scope.FactorFilter, "; "))
	}
	return strings.Join(parts, " / ")
}

func plannedCoverageGaps(spec *model.Spec, artifacts *evaluationArtifacts, plan *evaluationReportPlan) []RunGap {
	plannedAreas, plannedFactors, plannedRequirements := plannedExpansion(spec, plan)
	for _, areaID := range plannedAreas {
		if area := artifacts.Areas[areaKey(areaID)]; area == nil || area.Analysis == nil {
			return []RunGap{{Kind: GapMissingEvaluationData, Ref: areaDataPath(areaID, "area-analysis-result.json"), Detail: "planned area analysis payload is missing"}}
		}
	}
	for _, factor := range plannedFactors {
		if item := artifacts.Factors[factorKey(factor)]; item == nil || item.Analysis == nil {
			return []RunGap{{Kind: GapMissingEvaluationData, Ref: factorDataPath(factor, "factor-analysis-result.json"), Detail: "planned factor analysis payload is missing"}}
		}
	}
	for _, req := range plannedRequirements {
		item := artifacts.Requirements[requirementKey(req)]
		if item == nil || item.Assessment == nil {
			return []RunGap{{Kind: GapMissingEvaluationData, Ref: requirementDataPath(req, "requirement-assessment-result.json"), Detail: "planned requirement assessment payload is missing"}}
		}
		if item.Rating == nil {
			return []RunGap{{Kind: GapMissingEvaluationData, Ref: requirementDataPath(req, "requirement-rating-result.json"), Detail: "planned requirement rating payload is missing"}}
		}
	}
	return nil
}

func plannedExpansion(spec *model.Spec, plan *evaluationReportPlan) ([][]string, []factorID, []requirementID) {
	var areas [][]string
	var factors []factorID
	var requirements []requirementID
	plannedFactorIDs := map[string]struct{}{}
	for _, element := range model.Flatten(model.Project(spec)) {
		switch element.Kind {
		case model.KindArea:
			if area, ok := plannedArea(spec, element, plan); ok {
				areas = append(areas, []string(area))
			}
		case model.KindFactor:
			if id, ok := plannedFactor(spec, element, plan); ok {
				factors = append(factors, id)
				plannedFactorIDs[element.ID] = struct{}{}
			}
		case model.KindRequirement:
			if req, ok := plannedRequirement(spec, element, plan, plannedFactorIDs); ok {
				requirements = append(requirements, req)
			}
		}
	}
	return areas, factors, requirements
}

func plannedArea(spec *model.Spec, element *model.Element, plan *evaluationReportPlan) (model.AreaPath, bool) {
	area, err := model.ParseAreaReference(spec, element.ID)
	if err != nil || !areaInScope(area, plan.ScopedAreaID) {
		return nil, false
	}
	return area, true
}

func plannedFactor(spec *model.Spec, element *model.Element, plan *evaluationReportPlan) (factorID, bool) {
	area, factorPath, err := model.ParseFactorReference(spec, element.ID)
	if err != nil {
		return factorID{}, false
	}
	id := factorID{DeclaringArea: []string(area), Path: []string(factorPath)}
	return id, factorInScope(id, plan)
}

func plannedRequirement(spec *model.Spec, element *model.Element, plan *evaluationReportPlan, plannedFactorIDs map[string]struct{}) (requirementID, bool) {
	area, name, err := model.ParseRequirementReference(spec, element.ID)
	if err != nil || !areaInScope(area, plan.ScopedAreaID) {
		return requirementID{}, false
	}
	if len(plan.FactorFilter) > 0 {
		if _, ok := plannedFactorIDs[element.ParentID]; !ok {
			return requirementID{}, false
		}
	}
	return requirementID{DeclaringArea: []string(area), Name: name}, true
}

func areaInScope(areaID model.AreaPath, scoped []string) bool {
	return len(areaID) >= len(scoped) && sameStrings([]string(areaID[:len(scoped)]), scoped)
}

func factorInScope(id factorID, plan *evaluationReportPlan) bool {
	if !areaInScope(model.AreaPath(id.DeclaringArea), plan.ScopedAreaID) {
		return false
	}
	if len(plan.FactorFilter) == 0 {
		return true
	}
	for _, filter := range plan.FactorFilter {
		if sameStrings(id.DeclaringArea, filter.DeclaringArea) && len(id.Path) >= len(filter.Path) && sameStrings(id.Path[:len(filter.Path)], filter.Path) {
			return true
		}
	}
	return false
}

func writeAreaFactorBreakdownSection(b *strings.Builder, heading string, spec *model.Spec, artifacts *evaluationArtifacts, rootAreaID []string, reportPath string) {
	b.WriteString("## " + heading + "\n\n")
	b.WriteString("| ▦ Area / □ Factor | Overall rating | Local rating | Findings | Recommendations |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	root := artifacts.Areas[areaKey(rootAreaID)]
	if root == nil || root.Analysis == nil {
		b.WriteString("| (no area / factor breakdown) | — | — | — | — |\n\n")
		return
	}
	writeAreaFactorBreakdownAreaRows(b, spec, artifacts, root, reportPath, 0, true)
	b.WriteString("\n")
}

func writeAreaFactorBreakdownAreaRows(b *strings.Builder, spec *model.Spec, artifacts *evaluationArtifacts, area *evaluationAreaArtifacts, reportPath string, depth int, root bool) {
	local := scopedMap(area.Analysis, "localAnalysis")
	overall := scopedMap(area.Analysis, "localAndDescendantAnalysis")
	title := areaFactorBreakdownTitle(reportPath, areaReportPath(area.ID), areaTitle(spec, area.ID), depth, root)
	b.WriteString(md.TableRow(
		title,
		evaluationRatingLabel(spec, overall),
		evaluationRatingLabel(spec, local),
		fmt.Sprintf("%d", artifacts.areaFindingCount(area.ID)),
		fmt.Sprintf("%d", artifacts.areaRecommendationCount(area.ID)),
	))
	for _, factor := range artifacts.rootFactorsForArea(area.ID) {
		writeAreaFactorBreakdownFactorRows(b, spec, artifacts, factor, reportPath, depth+1)
	}
	for _, child := range artifacts.childAreas(area.ID) {
		writeAreaFactorBreakdownAreaRows(b, spec, artifacts, child, reportPath, depth+1, false)
	}
}

func writeAreaFactorBreakdownFactorRows(b *strings.Builder, spec *model.Spec, artifacts *evaluationArtifacts, factor *evaluationFactorArtifacts, reportPath string, depth int) {
	local := scopedMap(factor.Analysis, "localAnalysis")
	overall := scopedMap(factor.Analysis, "localAndDescendantAnalysis")
	title := areaFactorBreakdownFactorTitle(spec, reportPath, factor.ID, depth)
	b.WriteString(md.TableRow(
		title,
		evaluationRatingLabel(spec, overall),
		evaluationRatingLabel(spec, local),
		fmt.Sprintf("%d", artifacts.factorFindingCount(factor.ID)),
		fmt.Sprintf("%d", artifacts.factorRecommendationCount(factor.ID)),
	))
	for _, child := range artifacts.childFactors(factor.ID) {
		writeAreaFactorBreakdownFactorRows(b, spec, artifacts, child, reportPath, depth+1)
	}
}

func areaFactorBreakdownTitle(reportPath, targetPath, title string, depth int, root bool) string {
	link := reportLink(reportPath, targetPath, areaBreakdownMarker()+" "+title)
	if root {
		return "**" + link + "**"
	}
	return areaFactorIndent(depth) + link
}

func areaFactorBreakdownFactorTitle(spec *model.Spec, reportPath string, id factorID, depth int) string {
	return areaFactorIndent(depth) + reportLink(reportPath, factorReportPath(id), factorBreakdownMarker()+" "+factorTitle(spec, id))
}

func areaFactorIndent(depth int) string {
	if depth <= 0 {
		return ""
	}
	return strings.Repeat("↳ ", depth)
}

func areaBreakdownMarker() string {
	return "▦"
}

func factorBreakdownMarker() string {
	return "□"
}

func (a *evaluationArtifacts) areaFindingCount(areaID []string) int {
	count := 0
	for _, row := range a.rankedFindings() {
		if row.Requirement == nil {
			continue
		}
		if areaContains(areaID, row.Requirement.ID.DeclaringArea) {
			count++
		}
	}
	return count
}

func (a *evaluationArtifacts) factorFindingCount(factor factorID) int {
	count := 0
	for _, row := range a.rankedFindings() {
		if row.Requirement == nil {
			continue
		}
		if requirementMatchesFactor(row.Requirement, factor) {
			count++
		}
	}
	return count
}

func (a *evaluationArtifacts) areaRecommendationCount(areaID []string) int {
	count := 0
	for _, item := range a.rankedRecommendations() {
		if a.recommendationMatchesArea(item.Recommendation, areaID) {
			count++
		}
	}
	return count
}

func (a *evaluationArtifacts) factorRecommendationCount(factor factorID) int {
	count := 0
	for _, item := range a.rankedRecommendations() {
		if a.recommendationMatchesFactor(item.Recommendation, factor) {
			count++
		}
	}
	return count
}

func (a *evaluationArtifacts) recommendationMatchesArea(rec map[string]any, areaID []string) bool {
	for _, ref := range objectSlice(rec["traceRefs"]) {
		for _, ctx := range a.recommendationTraceContexts(objectMap(ref)) {
			if areaContains(areaID, ctx.AreaID) {
				return true
			}
			for _, factor := range ctx.FactorIDs {
				if areaContains(areaID, factor.DeclaringArea) {
					return true
				}
			}
		}
	}
	return false
}

func (a *evaluationArtifacts) recommendationMatchesFactor(rec map[string]any, factor factorID) bool {
	for _, ref := range objectSlice(rec["traceRefs"]) {
		for _, ctx := range a.recommendationTraceContexts(objectMap(ref)) {
			for _, traced := range ctx.FactorIDs {
				if factorContains(factor, traced) {
					return true
				}
			}
		}
	}
	return false
}

func requirementMatchesFactor(req *evaluationRequirementArtifacts, factor factorID) bool {
	if !sameStrings(req.ID.DeclaringArea, factor.DeclaringArea) {
		return false
	}
	for _, linked := range requirementFactorIDs(req) {
		parsed, err := parseRequirementFactorID(req.ID.DeclaringArea, linked)
		if err == nil && factorContains(factor, parsed) {
			return true
		}
	}
	return false
}

func areaContains(parent, child []string) bool {
	return len(child) >= len(parent) && sameStrings(child[:len(parent)], parent)
}

func factorContains(parent, child factorID) bool {
	return sameStrings(parent.DeclaringArea, child.DeclaringArea) &&
		len(child.Path) >= len(parent.Path) &&
		sameStrings(child.Path[:len(parent.Path)], parent.Path)
}

func adviceCoverageGaps(artifacts *evaluationArtifacts) []RunGap {
	expected := artifactFindingRefs(artifacts)
	if artifacts.FindingRanking == nil {
		return []RunGap{{Kind: GapMissingEvaluationData, Ref: "data/advice/finding-ranking-result.json", Detail: "required FindingRankingResult payload is missing"}}
	}
	if len(artifacts.Recommendations) == 0 {
		return []RunGap{{Kind: GapMissingEvaluationData, Ref: "data/advice/recommendations", Detail: "at least one RecommendationResult payload is required"}}
	}
	if artifacts.RecommendationRanking == nil {
		return []RunGap{{Kind: GapMissingEvaluationData, Ref: "data/advice/recommendation-ranking-result.json", Detail: "required RecommendationRankingResult payload is missing"}}
	}
	if gap := missingRankedFindingGap(artifacts, expected); gap != nil {
		return []RunGap{*gap}
	}
	if gap := missingRankedRecommendationGap(artifacts); gap != nil {
		return []RunGap{*gap}
	}
	if gap := missingFindingCoverageGap(artifacts, expected); gap != nil {
		return []RunGap{*gap}
	}
	return nil
}

func missingRankedFindingGap(artifacts *evaluationArtifacts, expected map[string]struct{}) *RunGap {
	rankedFindings := map[string]struct{}{}
	for _, item := range objectSlice(artifacts.FindingRanking["orderedFindings"]) {
		if key := artifactFindingRefKey(item["findingRef"]); key != "" {
			rankedFindings[key] = struct{}{}
		}
	}
	for key := range expected {
		if _, ok := rankedFindings[key]; !ok {
			return &RunGap{Kind: GapIncompleteEvaluationData, Ref: "data/advice/finding-ranking-result.json", Detail: "FindingRankingResult is missing " + key}
		}
	}
	return nil
}

func missingRankedRecommendationGap(artifacts *evaluationArtifacts) *RunGap {
	rankedRecommendations := map[string]struct{}{}
	for _, item := range objectSlice(artifacts.RecommendationRanking["orderedRecommendations"]) {
		id := firstString(item, "recommendationRef")
		if validRecommendationID(id) {
			rankedRecommendations[id] = struct{}{}
		}
	}
	for id := range artifacts.Recommendations {
		if _, ok := rankedRecommendations[id]; !ok {
			return &RunGap{Kind: GapIncompleteEvaluationData, Ref: "data/advice/recommendation-ranking-result.json", Detail: fmt.Sprintf("RecommendationRankingResult is missing %s", id)}
		}
	}
	return nil
}

func missingFindingCoverageGap(artifacts *evaluationArtifacts, expected map[string]struct{}) *RunGap {
	covered := map[string]struct{}{}
	for _, item := range objectSlice(artifacts.RecommendationRanking["findingCoverage"]) {
		key := artifactFindingRefKey(item["findingRef"])
		if key == "" {
			continue
		}
		covered[key] = struct{}{}
	}
	for key := range expected {
		if _, ok := covered[key]; !ok {
			return &RunGap{Kind: GapIncompleteEvaluationData, Ref: "data/advice/recommendation-ranking-result.json", Detail: "findingCoverage is missing " + key}
		}
	}
	return nil
}

func artifactFindingRefs(artifacts *evaluationArtifacts) map[string]struct{} {
	out := map[string]struct{}{}
	for _, req := range artifacts.Requirements {
		if req.Assessment == nil {
			continue
		}
		path := requirementDataPath(req.ID, "requirement-assessment-result.json")
		for _, finding := range objectSlice(req.Assessment["findings"]) {
			if id := firstString(finding, "id"); id != "" {
				out[path+"#findings["+id+"]"] = struct{}{}
			}
		}
	}
	return out
}

func artifactFindingRefKey(v any) string {
	ref := objectMap(v)
	if ref == nil || DataKind(firstString(ref, "kind")) != DataKindRequirementAssessment {
		return ""
	}
	path, err := dataPathForRoutineRef(ref)
	if err != nil {
		return ""
	}
	selector := firstString(ref, "selector")
	if selector == "" {
		return ""
	}
	return path + "#" + selector
}

func writeRankedFindingsTable(b *strings.Builder, spec *model.Spec, artifacts *evaluationArtifacts, reportPath string, limit int) {
	b.WriteString("| Rank | Finding | Area | Factors | Type | Severity |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- |\n")
	rows := artifacts.rankedFindings()
	if len(rows) == 0 {
		b.WriteString("| (no ranked findings) | — | — | — | — | — |\n\n")
		return
	}
	wrote := 0
	for _, row := range rows {
		if limit > 0 && wrote >= limit {
			break
		}
		b.WriteString(md.TableRow(
			fmt.Sprintf("%d", row.Rank),
			rankedFindingLink(row, reportPath),
			rankedFindingAreaLink(spec, row, reportPath),
			rankedFindingFactorLinks(spec, row, reportPath),
			findingTypeTitle(firstString(row.Finding, "type")),
			findingSeverityCell(row.Finding),
		))
		wrote++
	}
	b.WriteString("\n")
}

func (a *evaluationArtifacts) rankedFindings() []rankedFinding {
	items := objectSlice(a.FindingRanking["orderedFindings"])
	rows := make([]rankedFinding, 0, len(items))
	for _, item := range items {
		ref := objectMap(item["findingRef"])
		path, _ := dataPathForRoutineRef(ref)
		selector := firstString(ref, "selector")
		key := path + "#" + selector
		reqID := requirementIDForAssessmentPath(path)
		req := a.Requirements[requirementKey(reqID)]
		rank, _ := rankField(item)
		rows = append(rows, rankedFinding{
			Rank:        rank,
			Key:         key,
			Selector:    selector,
			FindingID:   findingIDFromSelector(selector),
			Requirement: req,
			Finding:     a.findingByKey(key),
			Ranking:     item,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Rank != rows[j].Rank {
			return rows[i].Rank < rows[j].Rank
		}
		return rows[i].Key < rows[j].Key
	})
	for i := range rows {
		rows[i].Total = len(rows)
	}
	return rows
}

func (a *evaluationArtifacts) rankedFindingByKey(key string) (rankedFinding, bool) {
	for _, row := range a.rankedFindings() {
		if row.Key == key {
			return row, true
		}
	}
	return rankedFinding{}, false
}

func (a *evaluationArtifacts) findingByKey(key string) map[string]any {
	path, selector, ok := strings.Cut(key, "#")
	if !ok {
		return nil
	}
	req := requirementIDForAssessmentPath(path)
	item := a.Requirements[requirementKey(req)]
	if item == nil || item.Assessment == nil {
		return nil
	}
	id := strings.TrimSuffix(strings.TrimPrefix(selector, "findings["), "]")
	for _, finding := range objectSlice(item.Assessment["findings"]) {
		if firstString(finding, "id") == id {
			return finding
		}
	}
	return nil
}

func rankedFindingLink(row rankedFinding, reportPath string) string {
	label := row.Selector
	if statement := firstString(row.Finding, "statement"); statement != "" {
		label = statement
	}
	if row.Requirement == nil {
		return markdownCell(label)
	}
	target := requirementReportPath(row.Requirement.ID)
	if row.FindingID != "" {
		target += "#" + findingAnchorID(row.FindingID)
	}
	return reportLink(reportPath, target, label)
}

func rankedFindingAreaLink(spec *model.Spec, row rankedFinding, reportPath string) string {
	if row.Requirement == nil {
		return "—"
	}
	return reportLink(reportPath, areaReportPath(row.Requirement.ID.DeclaringArea), areaTitle(spec, row.Requirement.ID.DeclaringArea))
}

func rankedFindingFactorLinks(spec *model.Spec, row rankedFinding, reportPath string) string {
	if row.Requirement == nil {
		return "—"
	}
	ids := requirementFactorIDs(row.Requirement)
	links := make([]string, 0, len(ids))
	for _, ref := range ids {
		id, err := parseRequirementFactorID(row.Requirement.ID.DeclaringArea, ref)
		if err != nil || len(id.Path) == 0 {
			continue
		}
		links = append(links, reportLink(reportPath, factorReportPath(id), factorTitle(spec, id)))
	}
	if len(links) == 0 {
		return "—"
	}
	return strings.Join(links, ", ")
}

func findingIDFromSelector(selector string) string {
	return strings.TrimSuffix(strings.TrimPrefix(selector, "findings["), "]")
}

func findingAnchorID(id string) string {
	if id == "" {
		return ""
	}
	return "finding-" + id
}

func requirementIDForAssessmentPath(path string) requirementID {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) < 6 {
		return requirementID{}
	}
	name := parts[len(parts)-2]
	areaParts := parts[2 : len(parts)-3]
	if len(areaParts) == 1 && areaParts[0] == "root" {
		areaParts = nil
	}
	return requirementID{DeclaringArea: areaParts, Name: name}
}

func writeTopRecommendationsTable(b *strings.Builder, spec *model.Spec, artifacts *evaluationArtifacts, reportPath string, limit int) {
	b.WriteString("| # | Recommendation | Area / factors | Impact | Confidence | Reason |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- |\n")
	items := artifacts.rankedRecommendations()
	if len(items) == 0 {
		b.WriteString("| (no recommendations) | — | — | — | — | — |\n\n")
		return
	}
	for i, item := range items {
		if limit > 0 && i >= limit {
			break
		}
		title := firstString(item.Recommendation, "title")
		if title == "" {
			title = fmt.Sprintf("Recommendation %d", item.Rank)
		}
		path := recommendationReportPath(item.Rank, title)
		b.WriteString(md.TableRow(
			fmt.Sprintf("%d", item.Rank),
			reportLink(reportPath, path, title),
			recommendationAreaFactorLinks(spec, artifacts, item.Recommendation, reportPath),
			impactTitle(firstString(item.Recommendation, "impact")),
			emDashIfEmpty(confidenceTitle(firstString(item.Ranking, "confidence"))),
			firstString(item.Recommendation, "expectedValue"),
		))
	}
	b.WriteString("\n")
}

func writeRecommendationIndexTable(b *strings.Builder, spec *model.Spec, artifacts *evaluationArtifacts, reportPath string) {
	b.WriteString("| # | Recommendation | Area / factors | Impact | Confidence | Reason | Ranking rationale |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
	items := artifacts.rankedRecommendations()
	if len(items) == 0 {
		b.WriteString("| (no recommendations) | — | — | — | — | — | — |\n\n")
		return
	}
	for _, item := range items {
		title := firstString(item.Recommendation, "title")
		if title == "" {
			title = fmt.Sprintf("Recommendation %d", item.Rank)
		}
		path := recommendationReportPath(item.Rank, title)
		b.WriteString(md.TableRow(
			fmt.Sprintf("%d", item.Rank),
			reportLink(reportPath, path, title),
			recommendationAreaFactorLinks(spec, artifacts, item.Recommendation, reportPath),
			impactTitle(firstString(item.Recommendation, "impact")),
			confidenceTitle(firstString(item.Recommendation, "confidence")),
			firstString(item.Recommendation, "expectedValue"),
			firstString(item.Ranking, "rationale"),
		))
	}
	b.WriteString("\n")
}

type recommendationTraceContext struct {
	AreaID    []string
	FactorIDs []factorID
}

func recommendationAreaFactorLinks(spec *model.Spec, artifacts *evaluationArtifacts, rec map[string]any, reportPath string) string {
	groups := map[string][]factorID{}
	var order []string
	for _, ref := range objectSlice(rec["traceRefs"]) {
		for _, ctx := range artifacts.recommendationTraceContexts(objectMap(ref)) {
			key := areaKey(ctx.AreaID)
			if _, ok := groups[key]; !ok {
				order = append(order, key)
			}
			groups[key] = appendFactorIDs(groups[key], ctx.FactorIDs...)
		}
	}
	if len(order) == 0 {
		return "—"
	}
	parts := make([]string, 0, len(order))
	for _, key := range order {
		areaID := areaIDFromKey(key)
		area := reportLink(reportPath, areaReportPath(areaID), areaTitle(spec, areaID))
		factors := recommendationFactorLinks(spec, groups[key], reportPath)
		parts = append(parts, area+" / "+factors)
	}
	return strings.Join(parts, "; ")
}

func (a *evaluationArtifacts) recommendationTraceContexts(ref map[string]any) []recommendationTraceContext {
	if ref == nil {
		return nil
	}
	subject := objectMap(ref["subject"])
	if id, err := requirementIDFrom(subject["requirementId"]); err == nil {
		req := a.Requirements[requirementKey(id)]
		return []recommendationTraceContext{{AreaID: id.DeclaringArea, FactorIDs: recommendationRequirementFactorIDs(req, id)}}
	}
	if id, err := factorIDFrom(subject["factorId"]); err == nil {
		return []recommendationTraceContext{{AreaID: id.DeclaringArea, FactorIDs: []factorID{id}}}
	}
	if id, err := areaIDFrom(subject["areaId"]); err == nil {
		return []recommendationTraceContext{{AreaID: id}}
	}
	return nil
}

func recommendationRequirementFactorIDs(req *evaluationRequirementArtifacts, id requirementID) []factorID {
	if req == nil {
		return nil
	}
	var out []factorID
	for _, ref := range requirementFactorIDs(req) {
		parsed, err := parseRequirementFactorID(id.DeclaringArea, ref)
		if err != nil || len(parsed.Path) == 0 {
			continue
		}
		out = appendFactorIDs(out, parsed)
	}
	return out
}

func appendFactorIDs(ids []factorID, candidates ...factorID) []factorID {
	seen := map[string]struct{}{}
	for _, id := range ids {
		seen[factorKey(id)] = struct{}{}
	}
	for _, candidate := range candidates {
		if len(candidate.Path) == 0 {
			continue
		}
		key := factorKey(candidate)
		if _, ok := seen[key]; ok {
			continue
		}
		ids = append(ids, candidate)
		seen[key] = struct{}{}
	}
	return ids
}

func recommendationFactorLinks(spec *model.Spec, factors []factorID, reportPath string) string {
	if len(factors) == 0 {
		return "—"
	}
	links := make([]string, 0, len(factors))
	for _, id := range factors {
		links = append(links, reportLink(reportPath, factorReportPath(id), factorTitle(spec, id)))
	}
	return strings.Join(links, ", ")
}

func writeAdviceCoverageSummary(b *strings.Builder, artifacts *evaluationArtifacts) {
	coverage := objectSlice(artifacts.RecommendationRanking["findingCoverage"])
	addressed := 0
	notDriving := 0
	for _, item := range coverage {
		switch firstString(item, "disposition") {
		case string(FindingCoverageAddressedByRecommendation):
			addressed++
		case string(FindingCoverageNotAdviceDriving):
			notDriving++
		}
	}
	b.WriteString("## Coverage\n\n")
	fmt.Fprintf(b, "- %s: %d\n", findingCoverageDispositionTitle(string(FindingCoverageAddressedByRecommendation)), addressed)
	fmt.Fprintf(b, "- %s: %d\n\n", findingCoverageDispositionTitle(string(FindingCoverageNotAdviceDriving)), notDriving)
}

func recommendationReportPath(rank int, title string) string {
	slug := Slug(title)
	if slug == "" {
		slug = "recommendation"
	}
	return filepath.ToSlash(filepath.Join("recommendations", fmt.Sprintf("%03d-%s.md", rank, slug)))
}

func recommendationID(rec map[string]any) string {
	return firstString(rec, "id")
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

func evaluationAreaTrailLine(spec *model.Spec, artifacts *evaluationArtifacts, areaID []string, reportPath string) string {
	parts := []string{areaTrailPart(spec, artifacts, nil, reportPath)}
	for i := range areaID {
		id := areaID[:i+1]
		parts = append(parts, areaTrailPart(spec, artifacts, id, reportPath))
	}
	return "Area: " + strings.Join(parts, " / ")
}

func areaTrailPart(spec *model.Spec, artifacts *evaluationArtifacts, areaID []string, reportPath string) string {
	title := areaTitle(spec, areaID)
	if area := artifacts.Areas[areaKey(areaID)]; area != nil && area.Analysis != nil {
		return reportLink(reportPath, areaReportPath(areaID), title)
	}
	return markdownCell(title)
}

func evaluationFactorTrailLine(spec *model.Spec, factor factorID, reportPath string) string {
	parts := make([]string, 0, len(factor.Path))
	for i := range factor.Path {
		id := factorID{DeclaringArea: factor.DeclaringArea, Path: factor.Path[:i+1]}
		parts = append(parts, reportLink(reportPath, factorReportPath(id), factorTitle(spec, id)))
	}
	return "Factor: " + strings.Join(parts, " / ")
}

func evaluationRequirementFactorsLine(req *evaluationRequirementArtifacts, reportPath string) string {
	links := requirementFactorLinks(req, reportPath)
	if links == "" {
		links = "(none)"
	}
	return "Factors: " + links
}

func writeEvaluationLimitsTable(b *strings.Builder, scopes ...map[string]any) {
	b.WriteString("| Type | Scope | Impact |\n")
	b.WriteString("| --- | --- | --- |\n")
	wrote := false
	for _, scope := range scopes {
		for _, field := range []string{"incompleteInputs", "evaluationLimits"} {
			for _, item := range objectSlice(scope[field]) {
				wrote = true
				b.WriteString(md.TableRow(limitTypeTitle(field), firstString(item, "scope", "ref", "id"), firstString(item, "impact", "description", "reason")))
			}
		}
	}
	if !wrote {
		b.WriteString("| (no limits or incomplete inputs) | — | — |\n")
	}
}

func writeFindingCoreDetails(b *strings.Builder, headingLevel int, id string, finding map[string]any, ranking rankedFinding, ranked bool) {
	heading := strings.Repeat("#", headingLevel)
	title := id
	if statement := firstString(finding, "statement"); statement != "" {
		title += " " + statement
	}
	if anchor := findingAnchorID(id); anchor != "" {
		b.WriteString(`<a id="` + markdownCell(anchor) + `"></a>` + "\n\n")
	}
	b.WriteString(heading + " " + title + "\n\n")
	writeFindingRankingContext(b, ranking, ranked)
	writeFindingSection(b, headingLevel+1, "Condition", firstString(finding, "condition"))
	writeFindingCriteriaSection(b, headingLevel+1, finding)
	writeFindingBasisSection(b, headingLevel+1, finding)
	writeFindingEffectSection(b, headingLevel+1, finding)
	writeFindingEvidenceSection(b, headingLevel+1, "Evidence", objectSlice(finding["evidence"]))
}

func writeFindingRankingContext(b *strings.Builder, ranking rankedFinding, ranked bool) {
	b.WriteString("| Advice rank | Tier | Ranking rationale |\n")
	b.WriteString("| --- | --- | --- |\n")
	if !ranked {
		b.WriteString("| (not ranked) | — | — |\n\n")
		return
	}
	b.WriteString(md.TableRow(fmt.Sprintf("%d / %d", ranking.Rank, ranking.Total), findingRankingTierTitle(firstString(ranking.Ranking, "tier")), firstString(ranking.Ranking, "rationale")))
	b.WriteString("\n")
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
		b.WriteString("- " + md.Code(label) + ": " + markdownCell(firstString(criterion, "criterion")) + "\n")
		if rationale := firstString(criterion, "rationale"); rationale != "" {
			b.WriteString("  Rationale: " + markdownCell(rationale) + "\n")
		}
	}
	b.WriteString("\n")
}

func writeFindingBasisSection(b *strings.Builder, headingLevel int, finding map[string]any) {
	b.WriteString(strings.Repeat("#", headingLevel) + " Basis\n\n")
	basis := objectMap(finding["basis"])
	if len(basis) == 0 {
		b.WriteString("(not recorded)\n\n")
		return
	}
	if status := firstString(basis, "status"); status != "" {
		b.WriteString("Status: " + basisStatusTitle(status) + "\n\n")
	}
	if statement := firstString(basis, "statement"); statement != "" {
		b.WriteString(statement + "\n\n")
	}
	if rationale := firstString(basis, "rationale"); rationale != "" {
		b.WriteString("Rationale: " + rationale + "\n\n")
	}
	writeFindingEvidenceSection(b, headingLevel+1, "Basis evidence", objectSlice(basis["evidence"]))
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
		b.WriteString("- " + md.Code(firstString(item, "sourceRef")) + ": " + markdownCell(firstString(item, "statement")) + "\n")
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

func findingBasisSummary(finding map[string]any) string {
	basis := objectMap(finding["basis"])
	status := basisStatusTitle(firstString(basis, "status"))
	statement := firstString(basis, "statement")
	if status == "" {
		return statement
	}
	if statement == "" {
		return status
	}
	return status + ": " + statement
}

func writeEvaluationFindingsTable(b *strings.Builder, assessment map[string]any) {
	b.WriteString("| ID | Statement | Type | Severity | Confidence | Effect | Basis |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
	findings := objectSlice(assessment["findings"])
	if len(findings) == 0 {
		b.WriteString("| (no findings) | — | — | — | — | — | — |\n\n")
		return
	}
	for i, finding := range findings {
		id := firstString(finding, "id", "ID")
		if id == "" {
			id = fmt.Sprintf("finding-%d", i+1)
		}
		b.WriteString(md.TableRow(md.Code(id), firstString(finding, "statement"), findingTypeTitle(firstString(finding, "type")), findingSeverityCell(finding), confidenceTitle(firstString(finding, "confidence")), findingEffectSummary(finding), findingBasisSummary(finding)))
	}
	b.WriteString("\n")
}

func writeEvaluationFindingDetails(b *strings.Builder, artifacts *evaluationArtifacts, req *evaluationRequirementArtifacts) {
	findings := objectSlice(req.Assessment["findings"])
	if len(findings) == 0 {
		b.WriteString("(no finding details)\n\n")
		return
	}
	for i, finding := range findings {
		id := firstString(finding, "id", "ID")
		if id == "" {
			id = fmt.Sprintf("finding-%d", i+1)
		}
		key := requirementDataPath(req.ID, "requirement-assessment-result.json") + "#findings[" + id + "]"
		ranking, ranked := artifacts.rankedFindingByKey(key)
		writeFindingCoreDetails(b, 3, id, finding, ranking, ranked)
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
			b.WriteString(md.TableRow(unknownTypeTitle(field), firstString(item, "description", "reason", "ref", "id")))
		}
	}
	if !wrote {
		b.WriteString("| (none recorded) | — |\n")
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

func reportLink(fromReport, toPath, label string) string {
	return md.RelLink(fromReport, toPath, label)
}

func workspaceRelativeRunPath(runAbs string, manifest *EvaluationManifest) string {
	if manifest == nil || manifest.Model == "" {
		return ""
	}
	for dir := runAbs; dir != "" && dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, manifest.Model)); err == nil {
			if rel, relErr := filepath.Rel(dir, runAbs); relErr == nil {
				return filepath.ToSlash(rel)
			}
			return ""
		}
	}
	return ""
}

func findingConcernSeverity(finding map[string]any) string {
	switch firstString(finding, "type") {
	case string(FindingTypeGap), string(FindingTypeRisk):
		return firstString(finding, "severity")
	default:
		return ""
	}
}

func findingSeverityCell(finding map[string]any) string {
	if severity := findingConcernSeverity(finding); severity != "" {
		return findingSeverityTitle(severity)
	}
	return "—"
}

func glossaryLinkTarget(reportPath, runPath string) string {
	return filepath.ToSlash(filepath.Join(workspaceBackrefs(runPath, reportPath), "glossary.md"))
}

func workspaceBackrefs(runPath, reportPath string) string {
	depth := pathDepth(runPath) + pathDepth(filepath.Dir(reportPath))
	if depth == 0 {
		return "."
	}
	parts := make([]string, depth)
	for i := range parts {
		parts[i] = ".."
	}
	return filepath.Join(parts...)
}

func pathDepth(path string) int {
	path = filepath.ToSlash(filepath.Clean(path))
	if path == "" || path == "." {
		return 0
	}
	depth := 0
	for _, part := range strings.Split(path, "/") {
		if part != "" && part != "." {
			depth++
		}
	}
	return depth
}

func areaTitle(spec *model.Spec, areaID []string) string {
	if len(areaID) == 0 {
		if spec.Title != "" {
			return spec.Title
		}
		return "Root area"
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
	if status := evaluationString(rating, "status"); status != string(RatingStatusRated) {
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

// evaluationSubRatingCell renders the descendant-inclusive sub-factors rating
// cell. When the factor has no descendants there is no roll-up distinct from
// its local rating, so it renders an em dash rather than repeating the local
// rating.
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
	if status != string(AnalysisStatusAnalyzed) || level == "" {
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
	if evaluationString(overall, "status") != string(AnalysisStatusAnalyzed) {
		return RatingResult{Kind: RatingResultNotAssessed, Rationale: evaluationString(overall, "statusReason")}
	}
	return RatingResult{Kind: RatingResultRated, Level: ratingLevelID(evaluationString(overall, "ratingLevelId")), Rationale: evaluationString(overall, "rationale")}
}

func markdownCell(s string) string {
	return md.Cell(s)
}
