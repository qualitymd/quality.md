// Package cli wires the QUALITY.md command tree (Cobra) and runs it through
// Fang for styled help, errors, version, and shell completion.
package cli

import (
	"context"
	"errors"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

const (
	ExitOK       = 0
	ExitProblems = 1
	ExitUsage    = 2
	ExitInternal = 70
)

type codedError struct {
	code   int
	silent bool
	err    error
}

func (e *codedError) Error() string {
	return e.err.Error()
}

func (e *codedError) Unwrap() error {
	return e.err
}

// Build-time metadata. goreleaser overrides these via -ldflags.
var (
	version = "dev"
	commit  = "none"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "qualitymd",
		Short: "Work with QUALITY.md files",
		Long: "qualitymd works with QUALITY.md files: Markdown documents whose YAML " +
			"frontmatter declares a quality model with a ratingScale, targets, factors, " +
			"requirements, and one assessment per requirement.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return usageError(err)
	})
	root.AddCommand(newInitCmd())
	root.AddCommand(newLintCmd())
	root.AddCommand(newSpecCmd())
	return root
}

// Execute builds the command tree and runs it. It exits non-zero on error;
// Fang renders the error.
func Execute() {
	os.Exit(execute(context.Background(), newRootCmd()))
}

func execute(ctx context.Context, root *cobra.Command) int {
	v, c := buildInfo()
	err := fang.Execute(
		ctx,
		root,
		fang.WithVersion(v),
		fang.WithCommit(c),
		fang.WithColorSchemeFunc(brandColorScheme),
		fang.WithErrorHandler(errorHandler),
	)
	return codeFor(err)
}

// buildInfo resolves the version and commit Fang reports. goreleaser stamps
// both via -ldflags for releases; otherwise we recover what we can from the
// embedded module build info so a `go install module@v1.2.3` shows its tag and
// a local build inside the repo shows "dev (<short-sha>)" rather than a bare,
// uninformative placeholder.
func buildInfo() (string, string) {
	if version != "dev" {
		return version, commit
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version, commit
	}
	var revision string
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			revision = setting.Value
		}
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return v, revision
	}
	return "dev", revision
}

func codeFor(err error) int {
	if err == nil {
		return ExitOK
	}
	var coded *codedError
	if errors.As(err, &coded) {
		return coded.code
	}
	if isUsageError(err) {
		return ExitUsage
	}
	return ExitInternal
}

func errorHandler(w io.Writer, styles fang.Styles, err error) {
	var coded *codedError
	if errors.As(err, &coded) && coded.silent {
		return
	}
	fang.DefaultErrorHandler(w, styles, err)
}

func usage(validator cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := validator(cmd, args); err != nil {
			return usageError(err)
		}
		return nil
	}
}

func usageError(err error) error {
	return &codedError{code: ExitUsage, err: err}
}

func silentProblems(err error) error {
	return &codedError{code: ExitProblems, silent: true, err: err}
}

func isUsageError(err error) bool {
	s := err.Error()
	for _, prefix := range []string{
		"flag needs an argument:",
		"unknown flag:",
		"unknown shorthand flag:",
		"unknown command",
		"invalid argument",
	} {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
