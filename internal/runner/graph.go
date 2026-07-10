package runner

import (
	"fmt"
	"strings"

	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/model"
)

// UnitKind names one evaluation work-unit kind. The vocabulary follows the
// evaluation protocol's routine order.
type UnitKind string

const (
	KindFrameEvaluation            UnitKind = "frameEvaluation"
	KindFrameAreaEvaluation        UnitKind = "frameAreaEvaluation"
	KindFrameRequirementEvaluation UnitKind = "frameRequirementEvaluation"
	// KindAssessRateRequirement executes the protocol's assessRequirement and
	// rateRequirement moves as one evaluator call; the unit persists both
	// payload kinds.
	KindAssessRateRequirement UnitKind = "assessRateRequirement"
	KindFrameFactorAnalysis   UnitKind = "frameFactorAnalysis"
	KindAnalyzeFactor         UnitKind = "analyzeFactor"
	KindFrameAreaAnalysis     UnitKind = "frameAreaAnalysis"
	KindAnalyzeArea           UnitKind = "analyzeArea"
	KindRankFindings          UnitKind = "rankFindings"
	KindRecommend             UnitKind = "recommend"
	KindRankRecommendations   UnitKind = "rankRecommendations"
	KindBuildReports          UnitKind = "buildReports"
)

// Unit is one node of the deterministic evaluation work graph.
type Unit struct {
	// ID is the deterministic work-unit identifier, derived from the routine
	// kind and the subject's canonical model reference.
	ID string
	// Kind is the work-unit kind.
	Kind UnitKind
	// Subject is the canonical model reference the unit addresses, empty for
	// run-wide units.
	Subject string
	// DependsOn lists work-unit IDs that must complete first.
	DependsOn []string
	// EvaluatorBacked marks bounded judgment work dispatched to the selected
	// evaluator; deterministic-only units are produced by the runner itself.
	EvaluatorBacked bool
	// DataKind is the evaluation data kind of the unit's result payload,
	// empty for composite or non-payload units.
	DataKind evaluation.DataKind
}

// Graph is the deterministic evaluation work graph in topological (model)
// order: executing units in slice order satisfies every dependency.
type Graph struct {
	Units []*Unit
	Plan  *ModelPlan
	byID  map[string]*Unit
}

// Unit returns the unit with the given ID, or nil.
func (g *Graph) Unit(id string) *Unit {
	return g.byID[id]
}

// EvaluatorUnits counts the evaluator-backed units.
func (g *Graph) EvaluatorUnits() int {
	count := 0
	for _, unit := range g.Units {
		if unit.EvaluatorBacked {
			count++
		}
	}
	return count
}

// ModelPlan is the planned-scope expansion of the model the work graph is
// built over.
type ModelPlan struct {
	ScopedArea   model.AreaPath
	FactorFilter []string
	Areas        []*PlannedArea
	Factors      []*PlannedFactor
	Requirements []*PlannedRequirement

	areasByRef        map[string]*PlannedArea
	factorsByRef      map[string]*PlannedFactor
	requirementsByRef map[string]*PlannedRequirement
}

// PlannedArea is one in-scope area with its local structure.
type PlannedArea struct {
	Path   model.AreaPath
	Ref    string
	Title  string
	Source string
	// LocalRequirements are refs of planned requirements declared in this
	// area (including under its factors).
	LocalRequirements []string
	// RootFactors are refs of planned top-level factors declared here.
	RootFactors []string
	// ChildAreas are refs of planned direct child areas.
	ChildAreas []string
}

// PlannedFactor is one in-scope factor node.
type PlannedFactor struct {
	Area  model.AreaPath
	Path  model.FactorPath
	Ref   string
	Title string
	// ChildFactors are refs of planned direct sub-factors.
	ChildFactors []string
	// Requirements are refs of planned requirements directly attached to
	// this factor (declared under it or linking it).
	Requirements []string
}

// PlannedRequirement is one in-scope requirement.
type PlannedRequirement struct {
	Area model.AreaPath
	Name string
	Ref  string
	// Factors are refs of planned factors the requirement is attached to.
	Factors     []string
	Title       string
	Description string
	Assessment  string
	// Ratings are per-rating-level criterion overrides.
	Ratings map[string]string
}

// Area returns the planned area with the given canonical reference, or nil.
func (p *ModelPlan) Area(ref string) *PlannedArea { return p.areasByRef[ref] }

// Factor returns the planned factor with the given canonical reference, or nil.
func (p *ModelPlan) Factor(ref string) *PlannedFactor { return p.factorsByRef[ref] }

// Requirement returns the planned requirement with the given canonical
// reference, or nil.
func (p *ModelPlan) Requirement(ref string) *PlannedRequirement { return p.requirementsByRef[ref] }

// BuildPlan expands the planned scope over the model in deterministic
// projection order.
func BuildPlan(spec *model.Spec, planned evaluation.PlannedRunScope) (*ModelPlan, error) {
	scopedArea, err := model.ParseAreaReference(spec, planned.AreaID)
	if err != nil {
		return nil, fmt.Errorf("planned scope area: %w", err)
	}
	plan := &ModelPlan{
		ScopedArea:        scopedArea,
		FactorFilter:      planned.FactorFilter,
		areasByRef:        map[string]*PlannedArea{},
		factorsByRef:      map[string]*PlannedFactor{},
		requirementsByRef: map[string]*PlannedRequirement{},
	}
	filterFactors := make([]scopeFactor, 0, len(planned.FactorFilter))
	for _, ref := range planned.FactorFilter {
		area, factorPath, err := model.ParseFactorReference(spec, ref)
		if err != nil {
			return nil, fmt.Errorf("planned scope factor: %w", err)
		}
		filterFactors = append(filterFactors, scopeFactor{area: area, path: factorPath})
	}

	filtered := len(filterFactors) > 0
	plannedFactorIDs := map[string]struct{}{}
	requirementParents := map[string]string{}
	for _, element := range model.Flatten(model.Project(spec)) {
		switch element.Kind {
		case model.KindArea:
			plan.collectArea(spec, element)
		case model.KindFactor:
			plan.collectFactor(spec, element, filterFactors, plannedFactorIDs)
		case model.KindRequirement:
			plan.collectRequirement(spec, element, filtered, plannedFactorIDs, requirementParents)
		}
	}
	plan.linkStructure(requirementParents)
	return plan, nil
}

func (p *ModelPlan) collectArea(spec *model.Spec, element *model.Element) {
	area, err := model.ParseAreaReference(spec, element.ID)
	if err != nil || !pathInScope(area, p.ScopedArea) {
		return
	}
	p.addArea(spec, area, element.Label)
}

func (p *ModelPlan) collectFactor(spec *model.Spec, element *model.Element, filters []scopeFactor, plannedFactorIDs map[string]struct{}) {
	area, factorPath, err := model.ParseFactorReference(spec, element.ID)
	if err != nil || !pathInScope(area, p.ScopedArea) {
		return
	}
	if len(filters) > 0 && !factorMatchesFilter(area, factorPath, filters) {
		return
	}
	p.addFactor(area, factorPath, element.Label)
	plannedFactorIDs[element.ID] = struct{}{}
}

func (p *ModelPlan) collectRequirement(spec *model.Spec, element *model.Element, filtered bool, plannedFactorIDs map[string]struct{}, requirementParents map[string]string) {
	area, name, err := model.ParseRequirementReference(spec, element.ID)
	if err != nil || !pathInScope(area, p.ScopedArea) {
		return
	}
	if filtered {
		if _, ok := plannedFactorIDs[element.ParentID]; !ok {
			return
		}
	}
	if p.requirementsByRef[element.ID] != nil {
		return
	}
	requirementParents[element.ID] = element.ParentID
	p.addRequirement(spec, area, name)
}

func pathInScope(area model.AreaPath, scoped model.AreaPath) bool {
	if len(area) < len(scoped) {
		return false
	}
	for i := range scoped {
		if area[i] != scoped[i] {
			return false
		}
	}
	return true
}

// scopeFactor is one parsed factor-filter entry.
type scopeFactor struct {
	area model.AreaPath
	path model.FactorPath
}

func factorMatchesFilter(area model.AreaPath, path model.FactorPath, filters []scopeFactor) bool {
	for _, filter := range filters {
		if !samePath([]string(area), []string(filter.area)) {
			continue
		}
		if len(path) >= len(filter.path) && samePath([]string(path[:len(filter.path)]), []string(filter.path)) {
			return true
		}
	}
	return false
}

func samePath(a, b []string) bool {
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

func (p *ModelPlan) addArea(spec *model.Spec, path model.AreaPath, label string) {
	node := &PlannedArea{
		Path:   path,
		Ref:    path.Reference(),
		Title:  label,
		Source: areaSource(spec, path),
	}
	p.Areas = append(p.Areas, node)
	p.areasByRef[node.Ref] = node
}

func (p *ModelPlan) addFactor(area model.AreaPath, path model.FactorPath, label string) {
	node := &PlannedFactor{
		Area:  area,
		Path:  path,
		Ref:   model.FactorReference(area, path),
		Title: label,
	}
	p.Factors = append(p.Factors, node)
	p.factorsByRef[node.Ref] = node
}

func (p *ModelPlan) addRequirement(spec *model.Spec, area model.AreaPath, name string) {
	req := lookupRequirement(spec, area, name)
	node := &PlannedRequirement{
		Area: area,
		Name: name,
		Ref:  model.RequirementReference(area, name),
	}
	if req != nil {
		node.Title = req.Title
		node.Description = req.Description
		node.Assessment = req.Assessment
		node.Ratings = req.Ratings
		for _, linked := range req.Factors {
			if ref, ok := p.resolveLinkedFactor(area, linked); ok {
				node.Factors = append(node.Factors, ref)
			}
		}
	}
	p.Requirements = append(p.Requirements, node)
	p.requirementsByRef[node.Ref] = node
}

// resolveLinkedFactor resolves a requirement's factors entry — a canonical
// factor reference or an area-relative factor path — to a planned factor.
func (p *ModelPlan) resolveLinkedFactor(area model.AreaPath, linked string) (string, bool) {
	if strings.HasPrefix(linked, "factor:") {
		if _, ok := p.factorsByRef[linked]; ok {
			return linked, true
		}
		return "", false
	}
	ref := model.FactorReference(area, model.FactorPath(strings.Split(linked, "/")))
	if _, ok := p.factorsByRef[ref]; ok {
		return ref, true
	}
	return "", false
}

// linkStructure wires area/factor/requirement adjacency after the planned
// sets are complete.
func (p *ModelPlan) linkStructure(requirementParents map[string]string) {
	p.linkFactors()
	p.linkAreas()
	p.linkRequirements(requirementParents)
}

func (p *ModelPlan) linkFactors() {
	for _, factor := range p.Factors {
		if len(factor.Path) == 1 {
			if area := p.areasByRef[factor.Area.Reference()]; area != nil {
				area.RootFactors = append(area.RootFactors, factor.Ref)
			}
			continue
		}
		parentRef := model.FactorReference(factor.Area, factor.Path[:len(factor.Path)-1])
		if parent := p.factorsByRef[parentRef]; parent != nil {
			parent.ChildFactors = append(parent.ChildFactors, factor.Ref)
		}
	}
}

func (p *ModelPlan) linkAreas() {
	for _, area := range p.Areas {
		if len(area.Path) == 0 {
			continue
		}
		parentRef := area.Path[:len(area.Path)-1].Reference()
		if parent := p.areasByRef[parentRef]; parent != nil {
			parent.ChildAreas = append(parent.ChildAreas, area.Ref)
		}
	}
}

func (p *ModelPlan) linkRequirements(requirementParents map[string]string) {
	for _, req := range p.Requirements {
		if area := p.areasByRef[req.Area.Reference()]; area != nil {
			area.LocalRequirements = append(area.LocalRequirements, req.Ref)
		}
		parentID := requirementParents[req.Ref]
		if strings.HasPrefix(parentID, "factor:") && !contains(req.Factors, parentID) {
			if _, ok := p.factorsByRef[parentID]; ok {
				req.Factors = append(req.Factors, parentID)
			}
		}
		for _, factorRef := range req.Factors {
			if factor := p.factorsByRef[factorRef]; factor != nil && !contains(factor.Requirements, req.Ref) {
				factor.Requirements = append(factor.Requirements, req.Ref)
			}
		}
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func areaSource(spec *model.Spec, path model.AreaPath) string {
	if len(path) == 0 {
		return spec.Source
	}
	areas := spec.Areas
	var current model.Area
	for i, name := range path {
		area, ok := areas[name]
		if !ok {
			return ""
		}
		current = area
		if i < len(path)-1 {
			areas = area.Areas
		}
	}
	return current.Source
}

func lookupRequirement(spec *model.Spec, path model.AreaPath, name string) *model.Requirement {
	factors := spec.Factors
	requirements := spec.Requirements
	areas := spec.Areas
	for _, part := range path {
		area, ok := areas[part]
		if !ok {
			return nil
		}
		factors = area.Factors
		requirements = area.Requirements
		areas = area.Areas
	}
	if req, ok := requirements[name]; ok {
		return &req
	}
	return findFactorRequirement(factors, name)
}

func findFactorRequirement(factors map[string]model.Factor, name string) *model.Requirement {
	for _, factor := range factors {
		if req, ok := factor.Requirements[name]; ok {
			return &req
		}
		if req := findFactorRequirement(factor.Factors, name); req != nil {
			return req
		}
	}
	return nil
}

// BuildGraph builds the deterministic work graph for a planned scope. Units
// are emitted in topological model order: frames and requirement judgment in
// projection order, then factor and area analyses bottom-up, then advice,
// then report build.
func BuildGraph(spec *model.Spec, planned evaluation.PlannedRunScope) (*Graph, error) {
	plan, err := BuildPlan(spec, planned)
	if err != nil {
		return nil, err
	}
	g := &Graph{Plan: plan, byID: map[string]*Unit{}}

	frameEvaluation := g.add(&Unit{ID: string(KindFrameEvaluation), Kind: KindFrameEvaluation, DataKind: evaluation.DataKindEvaluationFrame})
	areaFrames := g.addAreaFrameUnits(plan, frameEvaluation.ID)
	requirementUnits := g.addRequirementUnits(plan, areaFrames)
	analyzeFactorUnits := g.addFactorAnalysisUnits(plan, requirementUnits)
	analyzeAreaUnits := g.addAreaAnalysisUnits(plan, requirementUnits, analyzeFactorUnits)

	assessDeps := make([]string, 0, len(plan.Requirements))
	for _, req := range plan.Requirements {
		assessDeps = append(assessDeps, unitID(KindAssessRateRequirement, req.Ref))
	}
	rankFindings := g.add(&Unit{
		ID:              string(KindRankFindings),
		Kind:            KindRankFindings,
		DependsOn:       assessDeps,
		EvaluatorBacked: true,
		DataKind:        evaluation.DataKindFindingRanking,
	})
	recommendDeps := []string{rankFindings.ID}
	for _, factor := range plan.Factors {
		recommendDeps = append(recommendDeps, analyzeFactorUnits[factor.Ref])
	}
	for _, area := range plan.Areas {
		recommendDeps = append(recommendDeps, analyzeAreaUnits[area.Ref])
	}
	recommend := g.add(&Unit{
		ID:              string(KindRecommend),
		Kind:            KindRecommend,
		DependsOn:       sortedUnique(recommendDeps),
		EvaluatorBacked: true,
		DataKind:        evaluation.DataKindRecommendation,
	})
	rankRecommendations := g.add(&Unit{
		ID:              string(KindRankRecommendations),
		Kind:            KindRankRecommendations,
		DependsOn:       []string{recommend.ID, rankFindings.ID},
		EvaluatorBacked: true,
		DataKind:        evaluation.DataKindRecommendationRanking,
	})

	buildDeps := make([]string, 0, len(g.Units))
	for _, unit := range g.Units {
		buildDeps = append(buildDeps, unit.ID)
	}
	g.add(&Unit{
		ID:        string(KindBuildReports),
		Kind:      KindBuildReports,
		DependsOn: buildDeps,
	})
	_ = rankRecommendations
	return g, nil
}

func (g *Graph) add(unit *Unit) *Unit {
	g.Units = append(g.Units, unit)
	g.byID[unit.ID] = unit
	return unit
}

func (g *Graph) addAreaFrameUnits(plan *ModelPlan, frameEvaluationID string) map[string]string {
	areaFrames := map[string]string{}
	for _, area := range plan.Areas {
		deps := []string{frameEvaluationID}
		if len(area.Path) > len(plan.ScopedArea) {
			parentRef := area.Path[:len(area.Path)-1].Reference()
			if parentFrame, ok := areaFrames[parentRef]; ok {
				deps = append(deps, parentFrame)
			}
		}
		unit := g.add(&Unit{
			ID:        unitID(KindFrameAreaEvaluation, area.Ref),
			Kind:      KindFrameAreaEvaluation,
			Subject:   area.Ref,
			DependsOn: deps,
			DataKind:  evaluation.DataKindAreaEvaluationFrame,
		})
		areaFrames[area.Ref] = unit.ID
	}
	return areaFrames
}

// addRequirementUnits emits one frame unit and one combined judgment unit per
// requirement, and returns the judgment unit IDs by requirement ref. The
// combined unit is composite (no single DataKind): it persists both the
// assessment and rating payloads, so every dependency the rating satisfied
// before the merge now targets the combined unit.
func (g *Graph) addRequirementUnits(plan *ModelPlan, areaFrames map[string]string) map[string]string {
	requirementUnits := map[string]string{}
	for _, req := range plan.Requirements {
		frame := g.add(&Unit{
			ID:        unitID(KindFrameRequirementEvaluation, req.Ref),
			Kind:      KindFrameRequirementEvaluation,
			Subject:   req.Ref,
			DependsOn: []string{areaFrames[req.Area.Reference()]},
			DataKind:  evaluation.DataKindRequirementEvaluationFrame,
		})
		judge := g.add(&Unit{
			ID:              unitID(KindAssessRateRequirement, req.Ref),
			Kind:            KindAssessRateRequirement,
			Subject:         req.Ref,
			DependsOn:       []string{frame.ID},
			EvaluatorBacked: true,
		})
		requirementUnits[req.Ref] = judge.ID
	}
	return requirementUnits
}

func (g *Graph) addFactorAnalysisUnits(plan *ModelPlan, requirementUnits map[string]string) map[string]string {
	analyzeFactorUnits := map[string]string{}
	for _, factor := range factorsBottomUp(plan) {
		deps := make([]string, 0, len(factor.Requirements)+len(factor.ChildFactors))
		for _, reqRef := range factor.Requirements {
			deps = append(deps, requirementUnits[reqRef])
		}
		for _, childRef := range factor.ChildFactors {
			deps = append(deps, analyzeFactorUnits[childRef])
		}
		frame := g.add(&Unit{
			ID:        unitID(KindFrameFactorAnalysis, factor.Ref),
			Kind:      KindFrameFactorAnalysis,
			Subject:   factor.Ref,
			DependsOn: deps,
			DataKind:  evaluation.DataKindFactorAnalysisFrame,
		})
		analyze := g.add(&Unit{
			ID:              unitID(KindAnalyzeFactor, factor.Ref),
			Kind:            KindAnalyzeFactor,
			Subject:         factor.Ref,
			DependsOn:       append([]string{frame.ID}, deps...),
			EvaluatorBacked: true,
			DataKind:        evaluation.DataKindFactorAnalysis,
		})
		analyzeFactorUnits[factor.Ref] = analyze.ID
	}
	return analyzeFactorUnits
}

func (g *Graph) addAreaAnalysisUnits(plan *ModelPlan, requirementUnits, analyzeFactorUnits map[string]string) map[string]string {
	analyzeAreaUnits := map[string]string{}
	for _, area := range areasBottomUp(plan) {
		deps := make([]string, 0, len(area.RootFactors)+len(area.ChildAreas)+len(area.LocalRequirements))
		for _, factorRef := range area.RootFactors {
			deps = append(deps, analyzeFactorUnits[factorRef])
		}
		for _, childRef := range area.ChildAreas {
			deps = append(deps, analyzeAreaUnits[childRef])
		}
		for _, reqRef := range area.LocalRequirements {
			deps = append(deps, requirementUnits[reqRef])
		}
		frame := g.add(&Unit{
			ID:        unitID(KindFrameAreaAnalysis, area.Ref),
			Kind:      KindFrameAreaAnalysis,
			Subject:   area.Ref,
			DependsOn: deps,
			DataKind:  evaluation.DataKindAreaAnalysisFrame,
		})
		analyze := g.add(&Unit{
			ID:              unitID(KindAnalyzeArea, area.Ref),
			Kind:            KindAnalyzeArea,
			Subject:         area.Ref,
			DependsOn:       append([]string{frame.ID}, deps...),
			EvaluatorBacked: true,
			DataKind:        evaluation.DataKindAreaAnalysis,
		})
		analyzeAreaUnits[area.Ref] = analyze.ID
	}
	return analyzeAreaUnits
}

func unitID(kind UnitKind, subject string) string {
	return string(kind) + ":" + subject
}

// factorsBottomUp orders planned factors children-first, preserving model
// order among siblings.
func factorsBottomUp(plan *ModelPlan) []*PlannedFactor {
	out := make([]*PlannedFactor, len(plan.Factors))
	copy(out, plan.Factors)
	// Deeper factor paths come first; the projection order breaks ties
	// deterministically because the sort is stable.
	stableSortBy(out, func(a, b *PlannedFactor) bool {
		return len(a.Path) > len(b.Path)
	})
	return out
}

// areasBottomUp orders planned areas children-first.
func areasBottomUp(plan *ModelPlan) []*PlannedArea {
	out := make([]*PlannedArea, len(plan.Areas))
	copy(out, plan.Areas)
	stableSortBy(out, func(a, b *PlannedArea) bool {
		return len(a.Path) > len(b.Path)
	})
	return out
}

func stableSortBy[T any](items []T, less func(a, b T) bool) {
	// Insertion sort keeps the implementation dependency-free and stable.
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && less(items[j], items[j-1]); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}

func sortedUnique(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	// Keep deterministic order without sorting: input order is already
	// deterministic; duplicates are simply dropped.
	return out
}
