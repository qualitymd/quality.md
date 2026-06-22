// Package cli wires the QUALITY.md command tree (Cobra) and runs it through
// Fang for styled help, errors, version, and shell completion.
package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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
		Short: "Evaluate and improve AI assistant projects with QUALITY.md",
		Long: "qualitymd is the companion CLI for the QUALITY.md file format for " +
			"evaluating and improving the quality of AI assistant projects and harnesses.\n\n" +
			"Designed to be used with the companion agent skill.\n\n" +
			"Learn more at https://getquality.md\n" +
			"Report issues at https://github.com/qualitymd/quality.md/issues",
		Example:       "npx skills add qualitymd/quality.md",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          usage(cobra.NoArgs),
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			maybeStartUpdateRefresh(cmd)
		},
		PersistentPostRun: func(cmd *cobra.Command, _ []string) {
			maybeEmitUpdateNotice(cmd)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return renderRootWelcome(cmd.OutOrStdout())
		},
	}
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return usageError(err)
	})
	// Keep commands in registration order so the help reference reads as a
	// workflow (init -> lint -> spec -> schema -> evaluation) rather than
	// alphabetically.
	cobra.EnableCommandSorting = false
	root.AddGroup(
		&cobra.Group{ID: groupCommon, Title: "Common Tasks"},
		&cobra.Group{ID: groupManage, Title: "Manage"},
	)
	addCommand(root, groupCommon, newInitCmd())
	addCommand(root, groupCommon, newLintCmd())
	addCommand(root, groupCommon, newSpecCmd())
	addCommand(root, groupCommon, newSchemaCmd())
	addCommand(root, groupCommon, newEvaluationCmd())
	addCommand(root, groupManage, newVersionCmd())
	addCommand(root, groupManage, newUpdateCmd())
	addCommand(root, groupManage, newStatusCmd())
	root.AddCommand(newUpdateRefreshCmd())
	// Cobra's auto-generated help and completion commands are housekeeping, not
	// part of the workflow, so file them under Manage rather than letting them
	// form a stray default group above the curated ones.
	root.SetHelpCommandGroupID(groupManage)
	root.SetCompletionCommandGroupID(groupManage)
	return root
}

// Command groups for the root help reference. Common Tasks is the authoring
// loop; Manage is housekeeping (state, updates, version, shell plumbing).
const (
	groupCommon = "common"
	groupManage = "manage"
)

func addCommand(root *cobra.Command, groupID string, cmd *cobra.Command) {
	cmd.GroupID = groupID
	root.AddCommand(cmd)
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
	if revision == "" {
		revision = gitRevision()
	}
	return fallbackBuildInfoVersion(info.Main.Version, revision)
}

func fallbackBuildInfoVersion(moduleVersion, revision string) (string, string) {
	short := shortRevision(revision)
	if moduleVersion != "" && moduleVersion != "(devel)" {
		return moduleVersion, short
	}
	if short != "" {
		return fmt.Sprintf("dev (%s)", short), ""
	}
	return "dev", ""
}

func shortRevision(revision string) string {
	if len(revision) > 7 {
		return revision[:7]
	}
	return revision
}

func gitRevision() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func codeFor(err error) int {
	if err == nil {
		return ExitOK
	}
	var coded *codedError
	if errors.As(err, &coded) {
		return coded.code
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

// silentInternal wraps err as an internal failure whose rendering the caller has
// already handled (e.g. an emitted JSON error receipt), so the top-level handler
// does not render it a second time. The exit code stays ExitInternal.
func silentInternal(err error) error {
	return &codedError{code: ExitInternal, silent: true, err: err}
}
