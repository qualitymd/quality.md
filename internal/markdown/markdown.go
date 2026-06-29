// Package markdown provides small Markdown authoring primitives for generated
// project artifacts.
package markdown

import (
	"path/filepath"
	"strconv"
	"strings"
)

const Empty = "—"

// Writer accumulates generated Markdown.
type Writer struct {
	b strings.Builder
}

// FrontmatterField is one YAML frontmatter field.
type FrontmatterField struct {
	Name   string
	Value  string
	Values []string
}

// WriteString appends raw Markdown.
func (w *Writer) WriteString(s string) {
	w.b.WriteString(s)
}

// String returns the generated Markdown.
func (w *Writer) String() string {
	return w.b.String()
}

// Frontmatter renders a small YAML frontmatter block. It supports scalar string
// fields and string-list fields, which is enough for generated report metadata.
func Frontmatter(fields ...FrontmatterField) string {
	var b strings.Builder
	b.WriteString("---\n")
	for _, field := range fields {
		if len(field.Values) > 0 {
			b.WriteString(field.Name)
			b.WriteString(":\n")
			for _, value := range field.Values {
				b.WriteString("  - ")
				b.WriteString(yamlScalar(value))
				b.WriteString("\n")
			}
			continue
		}
		b.WriteString(field.Name)
		b.WriteString(": ")
		b.WriteString(yamlScalar(field.Value))
		b.WriteString("\n")
	}
	b.WriteString("---\n\n")
	return b.String()
}

// Heading writes a Markdown ATX heading followed by a blank line.
func (w *Writer) Heading(level int, title string) {
	if level < 1 {
		level = 1
	}
	w.b.WriteString(strings.Repeat("#", level))
	w.b.WriteString(" ")
	w.b.WriteString(title)
	w.b.WriteString("\n\n")
}

// Paragraph writes one paragraph followed by a blank line.
func (w *Writer) Paragraph(text string) {
	if text == "" {
		text = Empty
	}
	w.b.WriteString(text)
	w.b.WriteString("\n\n")
}

// RawBlock writes Markdown content followed by a blank line.
func (w *Writer) RawBlock(markdown string) {
	w.b.WriteString(markdown)
	if !strings.HasSuffix(markdown, "\n") {
		w.b.WriteString("\n")
	}
	w.b.WriteString("\n")
}

// Bullet writes an unordered-list item.
func (w *Writer) Bullet(markdown string) {
	w.b.WriteString("- ")
	w.b.WriteString(markdown)
	w.b.WriteString("\n")
}

// Table writes a pipe table. Cell values may contain inline Markdown.
func (w *Writer) Table(headers []string, rows [][]string) {
	w.b.WriteString(TableRow(headers...))
	w.b.WriteString(separatorRow(len(headers)))
	for _, row := range rows {
		w.b.WriteString(TableRow(row...))
	}
	w.b.WriteString("\n")
}

// TableRow renders one Markdown pipe-table row. Cell values may contain inline
// Markdown; unsafe table separators and multiline content are normalized.
func TableRow(cells ...string) string {
	escaped := make([]string, 0, len(cells))
	for _, cell := range cells {
		escaped = append(escaped, Cell(cell))
	}
	return "| " + strings.Join(escaped, " | ") + " |\n"
}

func separatorRow(columns int) string {
	cells := make([]string, columns)
	for i := range cells {
		cells[i] = "---"
	}
	return "| " + strings.Join(cells, " | ") + " |\n"
}

// Cell renders a value for a Markdown table cell.
func Cell(s string) string {
	if s == "" {
		return Empty
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.ReplaceAll(s, "\n", "<br>")
	return escapeTablePipes(s)
}

func escapeTablePipes(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r == '|' && (i == 0 || s[i-1] != '\\') {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Code renders an inline code span.
func Code(s string) string {
	if s == "" {
		s = Empty
	}
	fence := "`"
	longest := 0
	current := 0
	for _, r := range s {
		if r == '`' {
			current++
			if current > longest {
				longest = current
			}
			continue
		}
		current = 0
	}
	if longest > 0 {
		fence = strings.Repeat("`", longest+1)
		return fence + " " + s + " " + fence
	}
	return fence + s + fence
}

// Link renders an inline Markdown link.
func Link(label, target string) string {
	return "[" + linkLabel(label) + "](" + filepath.ToSlash(target) + ")"
}

// RelLink renders an inline link from one generated Markdown file to another.
func RelLink(fromPath, toPath, label string) string {
	fromDir := filepath.Dir(fromPath)
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
	return Link(label, rel)
}

// DataLink renders a generated report link to a source data artifact.
func DataLink(fromPath, toPath string) string {
	return RelLink(fromPath, toPath, filepath.Base(toPath))
}

func yamlScalar(s string) string {
	if isPlainYAMLScalar(s) {
		return s
	}
	return strconv.Quote(s)
}

func isPlainYAMLScalar(s string) bool {
	if s == "" || strings.TrimSpace(s) != s || strings.ContainsAny(s, "\r\n:#[]{}&,*!|>@`\"'") {
		return false
	}
	switch strings.ToLower(s) {
	case "null", "true", "false", "~":
		return false
	}
	switch s[0] {
	case '-', '?':
		return false
	}
	return true
}

func linkLabel(label string) string {
	label = Cell(label)
	label = strings.ReplaceAll(label, `[`, `\[`)
	label = strings.ReplaceAll(label, `]`, `\]`)
	return label
}
