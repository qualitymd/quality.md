package evaluation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
	"gopkg.in/yaml.v3"
)

type Run struct {
	Path            string
	Assessments     []AssessmentRecord
	Analyses        []AnalysisRecord
	Recommendations []RecommendationRecord
	Scale           []model.RatingLevel
}

type Counts struct {
	Assessments     int `json:"assessments"`
	Analyses        int `json:"analyses"`
	Recommendations int `json:"recommendations"`
}

type Gap struct {
	Kind   string `json:"kind"`
	Ref    string `json:"ref"`
	Detail string `json:"detail"`
}

type Status struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	Reportable    bool             `json:"reportable"`
	Counts        Counts           `json:"counts"`
	Gaps          []Gap            `json:"gaps"`
	NextActions   []receipt.Action `json:"nextActions"`
}

func Load(path string) (*Run, error) {
	runAbs, err := verifyRun(path)
	if err != nil {
		return nil, err
	}
	doc, err := document.Parse(filepath.Join(runAbs, "model.md"))
	if err != nil {
		return nil, err
	}
	spec, err := model.Decode(doc)
	if err != nil {
		return nil, err
	}
	run := &Run{Path: filepath.ToSlash(runAbs), Scale: spec.RatingScale}
	if err := loadJSONRecords(filepath.Join(runAbs, "assessments"), "*.json", func(path string, raw []byte) error {
		var rec AssessmentRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			return err
		}
		if rec.SchemaVersion != SchemaVersion {
			return fmt.Errorf("schemaVersion = %d, want %d", rec.SchemaVersion, SchemaVersion)
		}
		rec.File = filepath.ToSlash(filepath.Join("assessments", filepath.Base(path)))
		run.Assessments = append(run.Assessments, rec)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := loadJSONRecords(filepath.Join(runAbs, "analysis"), "*.json", func(path string, raw []byte) error {
		var rec AnalysisRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			return err
		}
		if rec.SchemaVersion != SchemaVersion {
			return fmt.Errorf("schemaVersion = %d, want %d", rec.SchemaVersion, SchemaVersion)
		}
		rec.File = filepath.ToSlash(filepath.Join("analysis", filepath.Base(path)))
		run.Analyses = append(run.Analyses, rec)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := loadRecommendations(filepath.Join(runAbs, "recommendations"), &run.Recommendations); err != nil {
		return nil, err
	}
	return run, nil
}

func loadJSONRecords(dir, pattern string, decode func(string, []byte) error) error {
	paths, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return err
	}
	slices.Sort(paths)
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		if err := decode(path, raw); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	}
	return nil
}

func loadRecommendations(dir string, out *[]RecommendationRecord) error {
	paths, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return err
	}
	slices.Sort(paths)
	for _, path := range paths {
		rec, err := parseRecommendation(path)
		if err != nil {
			return err
		}
		rec.File = filepath.ToSlash(filepath.Join("recommendations", filepath.Base(path)))
		*out = append(*out, rec)
	}
	return nil
}

func parseRecommendation(path string) (RecommendationRecord, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return RecommendationRecord{}, fmt.Errorf("reading %s: %w", path, err)
	}
	front, body, err := splitMarkdownFrontmatter(raw)
	if err != nil {
		return RecommendationRecord{}, fmt.Errorf("%s: %w", path, err)
	}
	var rec RecommendationRecord
	if err := yaml.Unmarshal(front, &rec); err != nil {
		return RecommendationRecord{}, fmt.Errorf("%s: parsing frontmatter: %w", path, err)
	}
	if rec.SchemaVersion != SchemaVersion {
		return RecommendationRecord{}, fmt.Errorf("%s: schemaVersion = %d, want %d", path, rec.SchemaVersion, SchemaVersion)
	}
	rec.Body = string(body)
	return rec, nil
}

func splitMarkdownFrontmatter(raw []byte) ([]byte, []byte, error) {
	if !bytes.HasPrefix(raw, []byte("---\n")) {
		return nil, nil, fmt.Errorf("missing runtime frontmatter")
	}
	rest := raw[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return nil, nil, fmt.Errorf("unterminated runtime frontmatter")
	}
	return rest[:end], rest[end+len("\n---\n"):], nil
}

func (r *Run) Counts() Counts {
	return Counts{
		Assessments:     len(r.Assessments),
		Analyses:        len(r.Analyses),
		Recommendations: len(r.Recommendations),
	}
}

func (r *Run) Status() Status {
	gaps := r.Renderable()
	status := Status{
		SchemaVersion: SchemaVersion,
		Path:          r.Path,
		Reportable:    len(gaps) == 0,
		Counts:        r.Counts(),
		Gaps:          gaps,
	}
	if status.Reportable {
		status.NextActions = []receipt.Action{{
			ID:      "build-report",
			Label:   "Build the evaluation report",
			Command: "qualitymd evaluation build-report " + r.Path,
		}}
	} else {
		status.NextActions = []receipt.Action{{
			ID:      "add-record",
			Label:   "Add the missing evaluation records",
			Command: "qualitymd evaluation add-record <kind> " + r.Path,
		}}
	}
	return status
}

func (r *Run) Renderable() []Gap {
	var gaps []Gap
	assessments := map[string]bool{}
	for _, rec := range r.Assessments {
		assessments[rec.File] = true
	}
	analyses := map[string]AnalysisRecord{}
	for _, rec := range r.Analyses {
		analyses[rec.File] = rec
	}
	recs := map[string]bool{}
	for _, rec := range r.Recommendations {
		recs[rec.File] = true
		recs[strings.TrimSuffix(filepath.Base(rec.File), ".md")] = true
	}
	if len(r.Analyses) == 0 {
		gaps = append(gaps, Gap{Kind: "missing-analysis", Ref: "analysis/", Detail: "no analysis records are present"})
	}
	for _, analysis := range r.Analyses {
		for _, ref := range analysis.AssessmentRecords {
			if !assessments[ref] {
				gaps = append(gaps, Gap{Kind: "missing-assessment", Ref: ref, Detail: "referenced by " + analysis.File})
			}
		}
		for _, ref := range analysis.ChildAnalysisRecords {
			if _, ok := analyses[ref]; !ok {
				gaps = append(gaps, Gap{Kind: "missing-analysis", Ref: ref, Detail: "referenced by " + analysis.File})
			}
		}
	}
	for _, assessment := range r.Assessments {
		for _, ref := range assessment.Recommendations {
			if recs[ref] || recs["recommendations/"+ref+".md"] {
				continue
			}
			gaps = append(gaps, Gap{Kind: "missing-recommendation", Ref: ref, Detail: "referenced by " + assessment.File})
		}
	}
	return gaps
}
