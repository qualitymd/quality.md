// Package lint validates QUALITY.md documents against the mechanical format
// rules and reports findings in the lint command's public result shape.
package lint

import (
	"errors"
	"fmt"
	"os"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
)

// Check parses and lints path, defaulting to QUALITY.md.
func Check(path string) (Result, error) {
	doc, early, err := parse(path)
	if err != nil {
		return Result{}, err
	}
	if early != nil {
		return *early, nil
	}
	state := newRunState(doc)
	state.run()
	return state.result(nil), nil
}

// Fix applies deterministic repairs, re-lints, and returns the post-repair
// result with repair records from the original document.
func Fix(path string) (Result, error) {
	doc, early, err := parse(path)
	if err != nil {
		return Result{}, err
	}
	if early != nil {
		return *early, nil
	}
	if info, err := os.Lstat(doc.Path); err != nil {
		return Result{}, fmt.Errorf("stat %s: %w", doc.Path, err)
	} else if info.Mode()&os.ModeSymlink != 0 {
		return Result{}, fmt.Errorf("%s is a symbolic link; refusing to repair it", doc.Path)
	}

	original := newRunState(doc)
	original.run()
	repairs := original.repairs
	repairRecords := make([]RepairRecord, 0, len(repairs))
	for _, repair := range repairs {
		if err := repair.apply(); err != nil {
			return Result{}, fmt.Errorf("applying repair %s at %s: %w", repair.record.RuleID, repair.record.Location.Label, err)
		}
		repairRecords = append(repairRecords, repair.record)
	}
	if len(repairs) > 0 {
		rendered, err := document.Render(doc)
		if err != nil {
			return Result{}, err
		}
		if err := document.WriteAtomic(doc.Path, rendered); err != nil {
			return Result{}, err
		}
		doc, early, err = parse(doc.Path)
		if err != nil {
			return Result{}, err
		}
		if early != nil {
			return *early, nil
		}
	}

	repaired := newRunState(doc)
	repaired.run()
	return repaired.result(repairRecords), nil
}

// Load returns the typed model only when the shared lint rule catalog finds no
// error-severity findings.
func Load(path string) (*model.Spec, error) {
	doc, early, err := parse(path)
	if err != nil {
		return nil, err
	}
	if early != nil {
		return nil, lintError{result: *early}
	}
	state := newRunState(doc)
	state.run()
	result := state.result(nil)
	if !result.Valid {
		return nil, lintError{result: result}
	}
	return model.Decode(doc)
}

type lintError struct {
	result Result
}

func (e lintError) Error() string {
	if e.result.Summary.Errors == 1 {
		return fmt.Sprintf("%s has 1 lint error", e.result.Path)
	}
	return fmt.Sprintf("%s has %d lint errors", e.result.Path, e.result.Summary.Errors)
}

// parse reads and parses path. A non-nil *Result means the document could not be
// parsed and the returned result should be emitted as-is without running rules.
func parse(path string) (*document.Document, *Result, error) {
	if path == "" {
		path = "QUALITY.md"
	}
	doc, err := document.Parse(path)
	if err == nil {
		return doc, nil, nil
	}

	var parseErr *document.ParseError
	if errors.As(err, &parseErr) {
		finding := Finding{
			RuleID:   RuleInvalidFrontmatter,
			Severity: SeverityError,
			Message:  "The file does not begin with valid QUALITY.md frontmatter; a QUALITY.md file requires a YAML frontmatter block matching the model shape.",
			Location: Location{
				Path:      path,
				ModelPath: []PathSegment{},
				Label:     "frontmatter",
			},
			Fixable: false,
		}
		result := Result{
			SchemaVersion: schemaVersion,
			Path:          path,
			Valid:         false,
			Summary:       Summary{Errors: 1},
			Findings:      []Finding{finding},
			Repairs:       []RepairRecord{},
			NextActions:   nextActions(path, Summary{Errors: 1}),
		}
		return nil, &result, nil
	}
	return nil, nil, err
}
