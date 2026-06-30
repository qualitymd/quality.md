// Command cli-docs generates the Mintlify CLI reference page from the qualitymd
// Cobra command tree, keeping the published reference in lockstep with the CLI's
// own command, flag, and example definitions.
//
// The command tree is the single source of truth. This generator introspects it
// and writes mintlify/cli.mdx; it never executes a command. Run via
// `mise run cli-docs`; the pre-commit hook and `mise run check` keep the
// generated page in sync.
//
// Do not edit mintlify/cli.mdx by hand. Edit the command definitions
// in internal/cli instead.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/qualitymd/quality.md/internal/cli"
	"github.com/qualitymd/quality.md/internal/workspace"
)

const outputRel = "mintlify/cli.mdx"

const frontmatter = `---
title: CLI reference
description: The qualitymd command-line interface that the /quality skill builds on.
---

{/* Generated from the Cobra command tree by scripts/cli-docs. Do not edit directly. Run ` + "`mise run cli-docs`" + `. */}

The ` + "`/quality`" + ` skill drives the ` + "`qualitymd`" + ` CLI for you. This reference is for
when you want to run commands directly. It is generated from the CLI's own
command definitions, so it always matches the installed binary.
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "cli-docs: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	repoRoot, err := workspace.FindRepoRoot("")
	if err != nil {
		return err
	}
	root := cli.NewRootForDocs()

	var b strings.Builder
	b.WriteString(frontmatter)
	writeOverview(&b, root)
	for _, cmd := range documented(root) {
		writeCommand(&b, cmd, 2)
	}

	outPath := filepath.Join(repoRoot, filepath.FromSlash(outputRel))
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(b.String()), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", outputRel, err)
	}
	return nil
}

// writeOverview emits the top-level command index table, linking each command to
// its detail section.
func writeOverview(b *strings.Builder, root *cobra.Command) {
	b.WriteString("\n## Commands\n\n")
	b.WriteString("| Command | Description |\n| --- | --- |\n")
	for _, cmd := range documented(root) {
		fmt.Fprintf(b, "| `%s` | %s |\n",
			cmd.CommandPath(), tableCell(cmd.Short))
	}
}

// writeCommand emits one command's section at the given Markdown heading level
// and recurses into its documented subcommands one level deeper.
func writeCommand(b *strings.Builder, cmd *cobra.Command, level int) {
	level = min(level, 6)
	fmt.Fprintf(b, "\n%s `%s`\n\n", strings.Repeat("#", level), cmd.CommandPath())
	if cmd.Short != "" {
		fmt.Fprintf(b, "%s\n\n", mdxEscape(cmd.Short))
	}

	fmt.Fprintf(b, "```bash\n%s\n```\n", cmd.UseLine())

	subs := documented(cmd)
	if len(subs) > 0 {
		b.WriteString("\n**Subcommands**\n\n")
		b.WriteString("| Command | Description |\n| --- | --- |\n")
		for _, sub := range subs {
			fmt.Fprintf(b, "| `%s` | %s |\n",
				sub.CommandPath(), tableCell(sub.Short))
		}
	}

	writeFlags(b, cmd)

	if example := strings.TrimRight(cmd.Example, "\n"); example != "" {
		b.WriteString("\n**Examples**\n\n")
		fmt.Fprintf(b, "```bash\n%s\n```\n", dedent(example))
	}

	for _, sub := range subs {
		writeCommand(b, sub, level+1)
	}
}

// writeFlags emits the command's local flags (excluding inherited persistent
// flags) as a table. Hidden flags are skipped.
func writeFlags(b *strings.Builder, cmd *cobra.Command) {
	var rows []string
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		rows = append(rows, fmt.Sprintf("| %s | %s |", flagName(f), flagDesc(f)))
	})
	if len(rows) == 0 {
		return
	}
	b.WriteString("\n**Flags**\n\n")
	b.WriteString("| Flag | Description |\n| --- | --- |\n")
	for _, row := range rows {
		b.WriteString(row + "\n")
	}
}

// documented returns a command's visible subcommands in registration order,
// dropping hidden commands and Cobra's auto-generated help and completion
// commands.
func documented(cmd *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Hidden || sub.IsAdditionalHelpTopicCommand() {
			continue
		}
		switch sub.Name() {
		case "help", "completion":
			continue
		}
		out = append(out, sub)
	}
	return out
}

func flagName(f *pflag.Flag) string {
	name := "`--" + f.Name + "`"
	if f.Shorthand != "" {
		name += ", `-" + f.Shorthand + "`"
	}
	return name
}

func flagDesc(f *pflag.Flag) string {
	desc := tableCell(f.Usage)
	if d := f.DefValue; d != "" && d != "false" && d != "[]" {
		desc += fmt.Sprintf(" (default: `%s`)", d)
	}
	return desc
}

// mdxEscape neutralizes the characters MDX would parse as JSX (`<` opening a tag,
// `{` opening an expression) so prose drawn from command metadata cannot break
// the Mintlify build. Flag usages carry placeholders such as `area:<path>`.
func mdxEscape(s string) string {
	r := strings.NewReplacer("<", "&lt;", "{", "&#123;", "}", "&#125;")
	return r.Replace(s)
}

// tableCell escapes a value for a Markdown table cell: MDX-safe, single-line,
// with pipes escaped.
func tableCell(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "|", "\\|")
	return mdxEscape(s)
}

// dedent strips the common leading-space indent shared by every non-blank line,
// since Example blocks in the command tree are indented for help rendering.
func dedent(s string) string {
	lines := strings.Split(s, "\n")
	indent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		n := len(line) - len(strings.TrimLeft(line, " "))
		if indent == -1 || n < indent {
			indent = n
		}
	}
	if indent <= 0 {
		return s
	}
	for i, line := range lines {
		if len(line) >= indent {
			lines[i] = line[indent:]
		}
	}
	return strings.Join(lines, "\n")
}
