// Package spec parses, renders, and writes QUALITY.md documents.
package spec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Spec is the typed QUALITY.md frontmatter model used by callers that need a
// valid model after lint has accepted the document.
type Spec struct {
	Title        string                 `yaml:"title,omitempty"`
	RatingScale  []RatingLevel          `yaml:"ratingScale"`
	Factors      map[string]Factor      `yaml:"factors,omitempty"`
	Requirements map[string]Requirement `yaml:"requirements,omitempty"`
	Targets      map[string]Target      `yaml:"targets,omitempty"`
	Source       string                 `yaml:"source,omitempty"`

	// Path is the source file; not part of the YAML.
	Path string `yaml:"-"`
}

// RatingLevel is one level in a model's rating scale. Description states what
// the level means across the whole model and is never overridden; Criterion is
// the default rule for rating a requirement's findings and MAY be overridden per
// requirement via Requirement.Ratings.
type RatingLevel struct {
	Level       string `yaml:"level"`
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
	Criterion   string `yaml:"criterion"`
}

// Target is a recursive target node in the quality model.
type Target struct {
	Factors      map[string]Factor      `yaml:"factors,omitempty"`
	Requirements map[string]Requirement `yaml:"requirements,omitempty"`
	Targets      map[string]Target      `yaml:"targets,omitempty"`
	Source       string                 `yaml:"source,omitempty"`
}

// Factor is a quality lens scoped to the target where it is declared.
type Factor struct {
	Description  string                 `yaml:"description,omitempty"`
	Factors      map[string]Factor      `yaml:"factors,omitempty"`
	Requirements map[string]Requirement `yaml:"requirements,omitempty"`
}

// Requirement is one assessable expectation.
type Requirement struct {
	Assessment string            `yaml:"assessment"`
	Factors    []string          `yaml:"factors,omitempty"`
	Ratings    map[string]string `yaml:"ratings,omitempty"`
}

// Document is a parsed QUALITY.md document. Frontmatter is the YAML mapping
// node and Body is the original Markdown body, including its leading newline
// after the closing frontmatter fence when present.
type Document struct {
	Path        string
	Frontmatter *yaml.Node
	Body        []byte
}

// ParseError is returned when a file can be read but its frontmatter block
// cannot be parsed as a QUALITY.md YAML document.
type ParseError struct {
	Path string
	Err  error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: invalid frontmatter: %v", e.Path, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// Parse reads path, defaulting to QUALITY.md, and parses its required YAML
// frontmatter block. It does not run lint rules.
func Parse(path string) (*Document, error) {
	if path == "" {
		path = "QUALITY.md"
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	frontmatter, body, err := splitFrontmatter(raw)
	if err != nil {
		return nil, &ParseError{Path: path, Err: err}
	}

	var root yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(frontmatter))
	if err := decoder.Decode(&root); err != nil {
		return nil, &ParseError{Path: path, Err: err}
	}
	if len(root.Content) == 0 {
		return nil, &ParseError{Path: path, Err: errors.New("frontmatter is empty")}
	}

	return &Document{
		Path:        path,
		Frontmatter: root.Content[0],
		Body:        body,
	}, nil
}

// Decode unmarshals a parsed document into the typed model. Callers should only
// use this after lint has accepted the document.
func Decode(doc *Document) (*Spec, error) {
	var model Spec
	if err := doc.Frontmatter.Decode(&model); err != nil {
		return nil, fmt.Errorf("%s: decoding model: %w", doc.Path, err)
	}
	model.Path = doc.Path
	return &model, nil
}

// Render returns a complete QUALITY.md document from the current frontmatter
// node and the original Markdown body bytes.
func Render(doc *Document) ([]byte, error) {
	var out bytes.Buffer
	out.WriteString("---\n")
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(2)
	if err := encoder.Encode(doc.Frontmatter); err != nil {
		return nil, fmt.Errorf("rendering frontmatter: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("closing YAML encoder: %w", err)
	}
	out.WriteString("---\n")
	out.Write(doc.Body)
	return out.Bytes(), nil
}

// WriteAtomic replaces path with content through a temporary file in the same
// directory and preserves the original file mode when possible.
func WriteAtomic(path string, content []byte) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("%s is a symbolic link; refusing to repair it", path)
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temporary file: %w", err)
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if err := tmp.Chmod(info.Mode().Perm()); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("setting temporary file mode: %w", err)
	}
	if _, err := tmp.Write(content); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("writing temporary file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temporary file: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replacing %s: %w", path, err)
	}
	return nil
}

// MapEntries yields each key/value pair of a YAML mapping in source order. It
// yields no pairs when mapping is nil or is not a MappingNode, so callers can
// range over it unconditionally.
func MapEntries(mapping *yaml.Node) iter.Seq2[*yaml.Node, *yaml.Node] {
	return func(yield func(*yaml.Node, *yaml.Node) bool) {
		if mapping == nil || mapping.Kind != yaml.MappingNode {
			return
		}
		for i := 0; i+1 < len(mapping.Content); i += 2 {
			if !yield(mapping.Content[i], mapping.Content[i+1]) {
				return
			}
		}
	}
}

// MapEntry returns a mapping entry's key node, value node, and key index.
func MapEntry(mapping *yaml.Node, key string) (*yaml.Node, *yaml.Node, int) {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil, nil, -1
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i], mapping.Content[i+1], i
		}
	}
	return nil, nil, -1
}

// RemoveMapEntry removes a key/value pair from a YAML mapping.
func RemoveMapEntry(mapping *yaml.Node, key string) bool {
	_, _, i := MapEntry(mapping, key)
	if i < 0 {
		return false
	}
	mapping.Content = append(mapping.Content[:i], mapping.Content[i+2:]...)
	return true
}

func splitFrontmatter(raw []byte) ([]byte, []byte, error) {
	if !bytes.HasPrefix(raw, []byte("---")) {
		return nil, nil, errors.New("file must begin with a YAML frontmatter block")
	}

	lineEnd := bytes.IndexByte(raw, '\n')
	if lineEnd < 0 {
		return nil, nil, errors.New("unterminated frontmatter: missing closing \"---\"")
	}
	if string(bytes.TrimSpace(raw[:lineEnd])) != "---" {
		return nil, nil, errors.New("opening frontmatter fence must be \"---\"")
	}

	start := lineEnd + 1
	for pos := start; pos <= len(raw); {
		next := bytes.IndexByte(raw[pos:], '\n')
		lineEnd := len(raw)
		if next >= 0 {
			lineEnd = pos + next
		}
		line := bytes.TrimRight(raw[pos:lineEnd], "\r")
		if string(bytes.TrimSpace(line)) == "---" {
			bodyStart := lineEnd
			if next >= 0 {
				bodyStart++
			}
			return raw[start:pos], raw[bodyStart:], nil
		}
		if next < 0 {
			break
		}
		pos = lineEnd + 1
	}
	return nil, nil, io.ErrUnexpectedEOF
}
