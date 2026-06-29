package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/agentinstructions"
	"github.com/qualitymd/quality.md/internal/receipt"
	"github.com/qualitymd/quality.md/internal/scaffold"
)

// renderInitHuman writes the post-scaffold confirmation and next-action footer
// to w (stderr), styled when w is a terminal and plain otherwise.
func renderInitHuman(w io.Writer, path string, agentInstructionFiles []string, next string) error {
	agentLine := ""
	if len(agentInstructionFiles) > 0 {
		agentLine = "Agent instructions: " + strings.Join(agentInstructionFiles, ", ") + "\n"
	}
	if !colorEnabled(w) {
		_, err := fmt.Fprintf(w, "Created %s\n%s\nNext: %s\n", path, agentLine, next)
		return err
	}
	_, err := fmt.Fprintf(w, "%s Created %s\n%s\nNext: %s\n",
		styleSuccess.Render(glyphSuccess), path, agentLine, styleCommand.Render(next))
	return err
}

func newInitCmd() *cobra.Command {
	opts := initOptions{}
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Scaffold a starter QUALITY.md",
		Example: "  qualitymd init\n" +
			"  qualitymd init docs/QUALITY.md\n" +
			"  qualitymd init --minimal\n" +
			"  qualitymd init --no-agent-instructions\n" +
			"  qualitymd init - > QUALITY.md\n" +
			"  qualitymd init --force",
		Args: usage(cobra.MaximumNArgs(1)),
		RunE: opts.run,
	}
	cmd.Flags().BoolVar(&opts.force, "force", false, "overwrite an existing file")
	cmd.Flags().BoolVar(&opts.jsonOutput, "json", false, "emit a machine-readable JSON init receipt")
	cmd.Flags().BoolVar(&opts.minimal, "minimal", false, "write a minimal valid skeleton without the guided template prose")
	cmd.Flags().BoolVar(&opts.noAgentInstructions, "no-agent-instructions", false, "do not create or update agent instruction files")
	return cmd
}

type initOptions struct {
	force               bool
	jsonOutput          bool
	minimal             bool
	noAgentInstructions bool
}

func (opts *initOptions) run(cmd *cobra.Command, args []string) error {
	path := "QUALITY.md"
	if len(args) == 1 {
		path = args[0]
	}
	if opts.jsonOutput && path == "-" {
		return usageError(fmt.Errorf("--json cannot be combined with path -"))
	}
	if path == "-" {
		_, err := cmd.OutOrStdout().Write(scaffoldBytes(opts.minimal))
		return err
	}

	created, err := createScaffold(path, *opts, cmd.ErrOrStderr())
	if err != nil {
		return err
	}
	agentInstructionFiles, err := opts.updateAgentInstructions(path)
	if err != nil {
		return err
	}
	actions := initActions(path)
	if opts.jsonOutput {
		return writeInitReceipt(cmd.OutOrStdout(), InitReceipt{
			SchemaVersion:         initSchemaVersion,
			Path:                  path,
			Created:               created,
			AgentInstructionFiles: agentInstructionFiles,
			NextActions:           actions,
		})
	}
	return renderInitHuman(cmd.ErrOrStderr(), path, agentInstructionPaths(agentInstructionFiles), actions[0].Command)
}

func createScaffold(path string, opts initOptions, errOut io.Writer) (bool, error) {
	created := true
	if _, err := os.Stat(path); err == nil {
		created = false
	} else if !os.IsNotExist(err) {
		return false, err
	}

	if err := scaffold.Create(path, opts.force, opts.minimal); err != nil {
		if opts.jsonOutput {
			writeInitError(errOut, path, err)
			return false, silentInternal(err)
		}
		return false, err
	}
	return created, nil
}

func (opts *initOptions) updateAgentInstructions(path string) ([]agentinstructions.FileResult, error) {
	if opts.noAgentInstructions {
		return []agentinstructions.FileResult{}, nil
	}
	return agentinstructions.Update(agentinstructions.UpdateOptions{ModelPath: path})
}

// scaffoldBytes returns the scaffold variant selected by the minimal flag.
func scaffoldBytes(minimal bool) []byte {
	if minimal {
		return scaffold.MinimalBytes()
	}
	return scaffold.Bytes()
}

const initSchemaVersion = 1

// InitReceipt is the JSON contract emitted by `qualitymd init --json`.
type InitReceipt struct {
	SchemaVersion         int                            `json:"schemaVersion"`
	Path                  string                         `json:"path"`
	Created               bool                           `json:"created"`
	AgentInstructionFiles []agentinstructions.FileResult `json:"agentInstructionFiles"`
	NextActions           []receipt.Action               `json:"nextActions"`
}

type initError struct {
	SchemaVersion int    `json:"schemaVersion"`
	Path          string `json:"path"`
	Reason        string `json:"reason"`
}

func initActions(path string) []receipt.Action {
	return []receipt.Action{
		{
			ID:      "lint",
			Label:   "Validate the scaffolded file",
			Command: "qualitymd lint " + path,
		},
	}
}

func agentInstructionPaths(results []agentinstructions.FileResult) []string {
	paths := make([]string, 0, len(results))
	for _, result := range results {
		paths = append(paths, result.Path)
	}
	return paths
}

func writeInitReceipt(w io.Writer, receipt InitReceipt) error {
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}

func writeInitError(w io.Writer, path string, err error) {
	data, marshalErr := json.MarshalIndent(initError{
		SchemaVersion: initSchemaVersion,
		Path:          path,
		Reason:        err.Error(),
	}, "", "  ")
	if marshalErr != nil {
		return
	}
	_, _ = w.Write(append(data, '\n'))
}
