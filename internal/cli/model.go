package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
)

// newModelCmd builds the read-only `model` group: a structure-and-identity
// projection of a quality model. It owns logical structure, canonical reference
// IDs, labels, and containment — never source coverage, counts, readiness, or
// evaluation results, which stay with `status` and `evaluation`.
func newModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "model",
		Short: "Query a quality model's structure and canonical reference IDs",
		Long: "model projects a QUALITY.md model as a read-only structure: its " +
			"areas, factors, and requirements, each with the canonical reference " +
			"ID used in evaluation payloads, and how they contain one another.\n\n" +
			"It never emits source coverage, counts, readiness, or evaluation " +
			"results — those stay with status and evaluation.",
		Args: usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newModelTreeCmd())
	cmd.AddCommand(newModelListCmd())
	cmd.AddCommand(newModelGetCmd())
	return cmd
}

func newModelTreeCmd() *cobra.Command {
	var jsonOutput bool
	var areaRef string
	var depth int
	cmd := &cobra.Command{
		Use:   "tree [path]",
		Short: "Render the model as a containment hierarchy",
		Example: "  qualitymd model tree\n" +
			"  qualitymd model tree --area area:agent-harness\n" +
			"  qualitymd model tree --depth 1 --json",
		Args: usage(cobra.MaximumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			spec, err := loadModelSpec(modelPathArg(args, 0))
			if err != nil {
				return err
			}
			root, err := rootElement(spec, areaRef)
			if err != nil {
				return err
			}
			if !cmd.Flags().Changed("depth") {
				depth = -1
			} else if depth < 0 {
				return usageError(fmt.Errorf("--depth must be zero or greater"))
			}
			view := truncateDepth(root, depth)
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), view)
			}
			return renderModelTree(cmd.OutOrStdout(), view)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit the tree as nested JSON")
	cmd.Flags().StringVar(&areaRef, "area", "", "root the tree at a canonical area reference (area:<path>)")
	cmd.Flags().IntVar(&depth, "depth", -1, "limit nesting depth; 0 emits only the rooted node")
	return cmd
}

func newModelListCmd() *cobra.Command {
	var jsonOutput bool
	var areaRef string
	var types []string
	cmd := &cobra.Command{
		Use:   "list [path]",
		Short: "Enumerate model elements with their canonical IDs",
		Example: "  qualitymd model list\n" +
			"  qualitymd model list --type factor\n" +
			"  qualitymd model list --area area:agent-harness --type requirement --json",
		Args: usage(cobra.MaximumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			kinds, err := parseKindFilter(types)
			if err != nil {
				return err
			}
			spec, err := loadModelSpec(modelPathArg(args, 0))
			if err != nil {
				return err
			}
			root, err := rootElement(spec, areaRef)
			if err != nil {
				return err
			}
			rows := listRows(root, kinds)
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), rows)
			}
			return renderModelList(cmd.OutOrStdout(), rows)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit the enumeration as a JSON array")
	cmd.Flags().StringVar(&areaRef, "area", "", "restrict output to one area's subtree (area:<path>)")
	cmd.Flags().StringSliceVar(&types, "type", nil, "restrict output to kinds: area, factor, requirement")
	return cmd
}

func newModelGetCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "get <id> [path]",
		Short: "Show one element's detail and immediate relations",
		Example: "  qualitymd model get area:root\n" +
			"  qualitymd model get factor:client-app::performance\n" +
			"  qualitymd model get requirement:agent-harness::has-tests --json",
		Args: usage(cobra.RangeArgs(1, 2)),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			spec, err := loadModelSpec(modelPathArg(args, 1))
			if err != nil {
				return err
			}
			root := model.Project(spec)
			element := model.Find(root, id)
			if element == nil {
				return usageError(unknownIDError(root, id))
			}
			detail := newElementDetail(element)
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), detail)
			}
			return renderModelGet(cmd.OutOrStdout(), detail)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit the element detail as a JSON object")
	return cmd
}

// modelPathArg returns the model file path positioned at index, defaulting to
// the conventional QUALITY.md.
func modelPathArg(args []string, index int) string {
	if len(args) > index {
		return args[index]
	}
	return "QUALITY.md"
}

// loadModelSpec reads and decodes a model file. It runs no lint rules — `model`
// assumes a parseable model and leaves diagnostics to `lint` — so a read or
// parse failure surfaces as an internal error (exit 70) pointing at lint.
func loadModelSpec(path string) (*model.Spec, error) {
	if path == "-" {
		return nil, usageError(fmt.Errorf("model does not read from stdin; pass a file path"))
	}
	doc, err := document.Parse(path)
	if err != nil {
		return nil, modelReadError(path, err)
	}
	spec, err := model.Decode(doc)
	if err != nil {
		return nil, modelReadError(path, err)
	}
	return spec, nil
}

func modelReadError(path string, err error) error {
	return fmt.Errorf("%w\n\nmodel needs a parseable model; run `qualitymd lint %s` for diagnostics", err, path)
}

// rootElement projects spec and, when areaRef is set, roots the projection at
// that area's subtree. A non-canonical or unresolved areaRef is a usage error.
func rootElement(spec *model.Spec, areaRef string) (*model.Element, error) {
	root := model.Project(spec)
	if areaRef == "" {
		return root, nil
	}
	areaPath, err := model.ParseAreaReference(spec, areaRef)
	if err != nil {
		return nil, usageError(err)
	}
	element := model.Find(root, areaPath.Reference())
	if element == nil {
		return nil, usageError(fmt.Errorf("--area %q does not resolve to an area in the model", areaRef))
	}
	return element, nil
}

func parseKindFilter(types []string) (map[model.Kind]bool, error) {
	if len(types) == 0 {
		return nil, nil
	}
	kinds := map[model.Kind]bool{}
	for _, raw := range types {
		switch model.Kind(raw) {
		case model.KindArea, model.KindFactor, model.KindRequirement:
			kinds[model.Kind(raw)] = true
		default:
			return nil, usageError(fmt.Errorf("--type %q is not one of: area, factor, requirement", raw))
		}
	}
	return kinds, nil
}

// listElement is the flat per-element row for `model list`.
type listElement struct {
	ID       string     `json:"id"`
	Kind     model.Kind `json:"kind"`
	Label    string     `json:"label"`
	ParentID string     `json:"parentId,omitempty"`
}

func listRows(root *model.Element, kinds map[model.Kind]bool) []listElement {
	rows := []listElement{}
	for _, e := range model.Flatten(root) {
		if kinds != nil && !kinds[e.Kind] {
			continue
		}
		rows = append(rows, listElement{ID: e.ID, Kind: e.Kind, Label: e.Label, ParentID: e.ParentID})
	}
	return rows
}

// elementDetail is the `model get` projection of one element and the IDs of its
// immediate relations, grouped by kind.
type elementDetail struct {
	ID           string     `json:"id"`
	Kind         model.Kind `json:"kind"`
	Label        string     `json:"label"`
	ParentID     string     `json:"parentId,omitempty"`
	Factors      []string   `json:"factors,omitempty"`
	Requirements []string   `json:"requirements,omitempty"`
	Areas        []string   `json:"areas,omitempty"`
}

func newElementDetail(e *model.Element) elementDetail {
	detail := elementDetail{ID: e.ID, Kind: e.Kind, Label: e.Label, ParentID: e.ParentID}
	for _, child := range e.Children {
		switch child.Kind {
		case model.KindFactor:
			detail.Factors = append(detail.Factors, child.ID)
		case model.KindRequirement:
			detail.Requirements = append(detail.Requirements, child.ID)
		case model.KindArea:
			detail.Areas = append(detail.Areas, child.ID)
		}
	}
	return detail
}

// truncateDepth returns a copy of e whose nesting is capped at depth additional
// levels below the rooted node. A negative depth imposes no limit.
func truncateDepth(e *model.Element, depth int) *model.Element {
	clone := &model.Element{ID: e.ID, Kind: e.Kind, Label: e.Label, ParentID: e.ParentID}
	if depth == 0 {
		return clone
	}
	next := depth - 1
	for _, child := range e.Children {
		clone.Children = append(clone.Children, truncateDepth(child, next))
	}
	return clone
}

func renderModelTree(w io.Writer, root *model.Element) error {
	var walk func(e *model.Element, depth int) error
	walk = func(e *model.Element, depth int) error {
		indent := strings.Repeat("  ", depth)
		line := fmt.Sprintf("%s%s  %s\n", indent, modelID(w, e.ID), e.Label)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
		for _, child := range e.Children {
			if err := walk(child, depth+1); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(root, 0)
}

func renderModelList(w io.Writer, rows []listElement) error {
	if len(rows) == 0 {
		_, err := io.WriteString(w, "No elements.\n")
		return err
	}
	for _, row := range rows {
		if _, err := fmt.Fprintf(w, "%s  %s\n", modelID(w, row.ID), row.Label); err != nil {
			return err
		}
	}
	return nil
}

func renderModelGet(w io.Writer, detail elementDetail) error {
	if _, err := fmt.Fprintf(w, "%s\n", modelID(w, detail.ID)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  kind:   %s\n", detail.Kind); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "  label:  %s\n", detail.Label); err != nil {
		return err
	}
	if detail.ParentID != "" {
		if _, err := fmt.Fprintf(w, "  parent: %s\n", detail.ParentID); err != nil {
			return err
		}
	}
	if err := renderRelation(w, "factors", detail.Factors); err != nil {
		return err
	}
	if err := renderRelation(w, "requirements", detail.Requirements); err != nil {
		return err
	}
	return renderRelation(w, "areas", detail.Areas)
}

func renderRelation(w io.Writer, label string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if _, err := fmt.Fprintf(w, "  %s:\n", label); err != nil {
		return err
	}
	for _, id := range ids {
		if _, err := fmt.Fprintf(w, "    %s\n", id); err != nil {
			return err
		}
	}
	return nil
}

// modelID styles a canonical ID on a terminal and returns it verbatim otherwise,
// so piped and --json output stays byte-stable.
func modelID(w io.Writer, id string) string {
	if colorEnabled(w) {
		return styleCommand.Render(id)
	}
	return id
}

// unknownIDError reports an unresolved `get` id and suggests the nearest
// projected ids of the same kind by edit distance.
func unknownIDError(root *model.Element, id string) error {
	suggestions := nearestIDs(root, id, 3)
	if len(suggestions) == 0 {
		return fmt.Errorf("no element in the model has id %q", id)
	}
	return fmt.Errorf("no element in the model has id %q; did you mean: %s", id, strings.Join(suggestions, ", "))
}

func nearestIDs(root *model.Element, id string, limit int) []string {
	prefix := kindPrefix(id)
	type ranked struct {
		id   string
		dist int
	}
	var candidates []ranked
	for _, e := range model.Flatten(root) {
		if prefix != "" && kindPrefix(e.ID) != prefix {
			continue
		}
		candidates = append(candidates, ranked{id: e.ID, dist: editDistance(id, e.ID)})
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].dist != candidates[j].dist {
			return candidates[i].dist < candidates[j].dist
		}
		return candidates[i].id < candidates[j].id
	})
	out := []string{}
	for _, c := range candidates {
		out = append(out, c.id)
		if len(out) == limit {
			break
		}
	}
	return out
}

func kindPrefix(id string) string {
	if before, _, ok := strings.Cut(id, ":"); ok {
		return before
	}
	return ""
}

// editDistance is the Levenshtein distance between a and b, used only to rank
// near-match suggestions for an unresolved id.
func editDistance(a, b string) int {
	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(curr[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}
