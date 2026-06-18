// Package models exposes bundled QUALITY.md models shipped with the CLI.
package models

import (
	"embed"
	"fmt"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
	"gopkg.in/yaml.v3"
)

//go:embed quality-meta-model.md
var files embed.FS

const qualityMetaModelName = "quality-meta-model"

// Entry is one bundled model in the catalog.
type Entry struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// View is the structured representation emitted by `qualitymd models view --json`.
type View struct {
	SchemaVersion int         `json:"schemaVersion"`
	Name          string      `json:"name"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Model         *model.Spec `json:"model"`
	BodyMarkdown  string      `json:"bodyMarkdown"`
}

type bundledModel struct {
	entry Entry
	file  string
}

var catalog = []bundledModel{
	{
		entry: Entry{
			Name:        qualityMetaModelName,
			Title:       "Quality meta-model",
			Description: "Criteria for evaluating a QUALITY.md model.",
		},
		file: "quality-meta-model.md",
	},
}

// Catalog returns the bundled-model catalog in stable display order.
func Catalog() []Entry {
	out := make([]Entry, 0, len(catalog))
	for _, item := range catalog {
		out = append(out, item.entry)
	}
	return out
}

// Markdown returns a bundled model as Markdown, optionally rewriting its apex
// source before rendering.
func Markdown(name, source string) ([]byte, error) {
	item, ok := find(name)
	if !ok {
		return nil, fmt.Errorf("unknown bundled model %q", name)
	}
	raw, err := files.ReadFile(item.file)
	if err != nil {
		return nil, err
	}
	if source == "" {
		return raw, nil
	}
	doc, err := parse(item, raw)
	if err != nil {
		return nil, err
	}
	setSource(doc.Frontmatter, source)
	return document.Render(doc)
}

// Structured returns a bundled model's parsed model and Markdown body, optionally
// rewriting the apex source before decoding.
func Structured(name, source string) (*View, error) {
	item, ok := find(name)
	if !ok {
		return nil, fmt.Errorf("unknown bundled model %q", name)
	}
	raw, err := files.ReadFile(item.file)
	if err != nil {
		return nil, err
	}
	doc, err := parse(item, raw)
	if err != nil {
		return nil, err
	}
	if source != "" {
		setSource(doc.Frontmatter, source)
	}
	spec, err := model.Decode(doc)
	if err != nil {
		return nil, err
	}
	return &View{
		SchemaVersion: 1,
		Name:          item.entry.Name,
		Title:         item.entry.Title,
		Description:   item.entry.Description,
		Model:         spec,
		BodyMarkdown:  string(doc.Body),
	}, nil
}

func find(name string) (bundledModel, bool) {
	for _, item := range catalog {
		if item.entry.Name == name {
			return item, true
		}
	}
	return bundledModel{}, false
}

func parse(item bundledModel, raw []byte) (*document.Document, error) {
	return document.ParseBytes(item.entry.Name, raw)
}

func setSource(mapping *yaml.Node, source string) {
	_, value, _ := document.MapEntry(mapping, "source")
	if value != nil {
		value.Kind = yaml.ScalarNode
		value.Tag = "!!str"
		value.Value = source
		value.Style = 0
		return
	}
	mapping.Content = append(mapping.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "source"},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: source},
	)
}
