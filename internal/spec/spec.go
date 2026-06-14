// Package spec loads and models a QUALITY.md specification. The spec lives in
// the YAML frontmatter of a Markdown file.
//
// STUB: This data model is a simplification that does not track the current
// format spec (SPECIFICATION.md and specs/). It models factors as a bare map of
// requirement name to a single `prompt`/`bash`/`cel` evaluator plus a pass/fail
// `rating`. The real format is richer: factors nest `subfactors` and carry
// `requirements`, each requirement declares exactly one assessment and is scored
// against a multi-level `ratings` scale (not a binary pass/fail), assessments
// reference `prompt`/`target` paths, and the Markdown body has structural
// requirements of its own. The loader here also ignores the body entirely. Do
// not treat these types as the canonical schema.
package spec

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Spec is a parsed QUALITY.md document.
type Spec struct {
	Factors map[string]Factor `yaml:"factors"`

	// Path is the source file; not part of the YAML.
	Path string `yaml:"-"`
}

// Factor groups related requirements (e.g. "reliability") by name.
type Factor map[string]Requirement

// Requirement is a single scored expectation. Exactly one evaluator
// (Prompt, Bash, or CEL) is expected; Rating supplies pass/fail criteria.
type Requirement struct {
	// Prompt holds inferential assessment conditions, either inline text or a
	// path to a Markdown guide. Evaluated by an LLM (not yet implemented).
	Prompt string `yaml:"prompt,omitempty"`
	// Bash is a shell command; a zero exit status means pass.
	Bash string `yaml:"bash,omitempty"`
	// CEL is a Common Expression Language predicate; a true result means pass.
	CEL string `yaml:"cel,omitempty"`
	// Rating describes the pass/fail criteria.
	Rating *Rating `yaml:"rating,omitempty"`
}

// Rating describes the conditions under which a requirement passes or fails.
type Rating struct {
	Pass string `yaml:"pass,omitempty"`
	Fail string `yaml:"fail,omitempty"`
}

// Load reads the spec at path (defaulting to QUALITY.md) and parses its
// frontmatter.
func Load(path string) (*Spec, error) {
	if path == "" {
		path = "QUALITY.md"
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading spec: %w", err)
	}

	fm, err := frontmatter(raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	var s Spec
	if err := yaml.Unmarshal(fm, &s); err != nil {
		return nil, fmt.Errorf("%s: parsing spec: %w", path, err)
	}
	s.Path = path
	return &s, nil
}

// frontmatter extracts the YAML block delimited by a leading and trailing
// `---` fence. A document without a fence is treated as YAML in its entirety.
func frontmatter(raw []byte) ([]byte, error) {
	const fence = "---"
	trimmed := bytes.TrimLeft(raw, " \t\r\n")
	if !bytes.HasPrefix(trimmed, []byte(fence)) {
		return raw, nil
	}
	rest := trimmed[len(fence):]
	if i := bytes.Index(rest, []byte("\n"+fence)); i >= 0 {
		return rest[:i], nil
	}
	return nil, fmt.Errorf("unterminated frontmatter: missing closing %q", fence)
}
