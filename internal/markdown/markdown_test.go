package markdown

import "testing"

func TestCellEscapesTableHazards(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: Empty},
		{name: "pipe", in: "a | b", want: `a \| b`},
		{name: "already escaped pipe", in: `a \| b`, want: `a \| b`},
		{name: "multiline", in: "a\nb\r\nc\rd", want: "a<br>b<br>c<br>d"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Cell(tc.in); got != tc.want {
				t.Fatalf("Cell(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestCodeHandlesBackticks(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: "`" + Empty + "`"},
		{name: "plain", in: "QUALITY.md", want: "`QUALITY.md`"},
		{name: "backtick", in: "a`b", want: "`` a`b ``"},
		{name: "double backtick", in: "a``b", want: "``` a``b ```"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Code(tc.in); got != tc.want {
				t.Fatalf("Code(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestRelLinkAndDataLink(t *testing.T) {
	if got := RelLink("areas/api/api-area.md", "root-area.md", "Root [Area] | link"); got != `[Root \[Area\] \| link](../../root-area.md)` {
		t.Fatalf("RelLink() = %q", got)
	}
	if got := DataLink("report.md", "data/evaluation-output-result.json"); got != "[evaluation-output-result.json](data/evaluation-output-result.json)" {
		t.Fatalf("DataLink() = %q", got)
	}
}

func TestFrontmatter(t *testing.T) {
	got := Frontmatter(
		FrontmatterField{Name: "type", Value: "Requirement Evaluation Report"},
		FrontmatterField{Name: "title", Value: "needs: quoting"},
		FrontmatterField{Name: "data", Values: []string{"data/evaluation-output-result.json", "data/areas/api/report.json"}},
	)
	want := "---\n" +
		"type: Requirement Evaluation Report\n" +
		"title: \"needs: quoting\"\n" +
		"data:\n" +
		"  - data/evaluation-output-result.json\n" +
		"  - data/areas/api/report.json\n" +
		"---\n\n"
	if got != want {
		t.Fatalf("Frontmatter() =\n%s\nwant:\n%s", got, want)
	}
}

func TestWriterTable(t *testing.T) {
	var w Writer
	w.Heading(1, "Report")
	w.Table(
		[]string{"Name", "Detail"},
		[][]string{
			{"A", "has | pipe"},
			{"B", "two\nlines"},
			{"C", Link("Open", "target.md")},
		},
	)
	want := "# Report\n\n" +
		"| Name | Detail |\n" +
		"| --- | --- |\n" +
		"| A | has \\| pipe |\n" +
		"| B | two<br>lines |\n" +
		"| C | [Open](target.md) |\n\n"
	if got := w.String(); got != want {
		t.Fatalf("Writer output:\n%s\nwant:\n%s", got, want)
	}
}
