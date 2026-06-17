package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/receipt"
	"github.com/qualitymd/quality.md/internal/scaffold"
)

func newInitCmd() *cobra.Command {
	var force bool
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Scaffold a starter QUALITY.md",
		Args:  usage(cobra.MaximumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "QUALITY.md"
			if len(args) == 1 {
				path = args[0]
			}

			if jsonOutput && path == "-" {
				return usageError(fmt.Errorf("--json cannot be combined with path -"))
			}

			if path == "-" {
				_, err := cmd.OutOrStdout().Write(scaffold.Bytes())
				return err
			}

			created := true
			if _, err := os.Stat(path); err == nil {
				created = false
			} else if !os.IsNotExist(err) {
				return err
			}

			if err := scaffold.Create(path, force); err != nil {
				if jsonOutput {
					writeInitError(cmd.ErrOrStderr(), path, err)
				}
				return err
			}
			actions := initActions(path)
			if jsonOutput {
				return writeInitReceipt(cmd.OutOrStdout(), InitReceipt{
					SchemaVersion: initSchemaVersion,
					Path:          path,
					Created:       created,
					NextActions:   actions,
				})
			}
			_, err := fmt.Fprintf(cmd.ErrOrStderr(), "Created %s\n\nNext: %s\n", path, actions[0].Command)
			return err
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing file")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable JSON init receipt")
	return cmd
}

const initSchemaVersion = 1

type InitReceipt struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	Created       bool             `json:"created"`
	NextActions   []receipt.Action `json:"nextActions"`
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

func writeInitReceipt(w interface{ Write([]byte) (int, error) }, receipt InitReceipt) error {
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}

func writeInitError(w interface{ Write([]byte) (int, error) }, path string, err error) {
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
