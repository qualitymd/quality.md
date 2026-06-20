package cli

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	qualitymd "github.com/qualitymd/quality.md"
)

type versionInfo struct {
	SchemaVersion        int    `json:"schemaVersion"`
	Version              string `json:"version"`
	Commit               string `json:"commit,omitempty"`
	DevelopmentBuild     bool   `json:"developmentBuild"`
	SpecificationVersion string `json:"specificationVersion"`
}

func newVersionCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show structured qualitymd version metadata",
		Args:  usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			info := currentVersionInfo()
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), info)
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(),
				"Version: %s\nCommit: %s\nDevelopment build: %t\nSpecification version: %s\n",
				info.Version,
				emptyDisplay(info.Commit),
				info.DevelopmentBuild,
				emptyDisplay(info.SpecificationVersion),
			)
			return err
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit machine-readable version metadata")
	return cmd
}

func currentVersionInfo() versionInfo {
	v, c := buildInfo()
	return versionInfo{
		SchemaVersion:        1,
		Version:              v,
		Commit:               c,
		DevelopmentBuild:     isDevelopmentVersion(v),
		SpecificationVersion: bundledSpecificationVersion(),
	}
}

func isDevelopmentVersion(v string) bool {
	v = strings.TrimSpace(v)
	if strings.HasPrefix(v, "dev") || strings.Contains(v, "+dirty") {
		return true
	}
	normalized := "v" + normalizeVersion(v)
	return strings.Contains(normalized, "-0.")
}

func bundledSpecificationVersion() string {
	for _, line := range strings.Split(string(bytes.TrimSpace(qualitymd.Specification())), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "**Specification version:**") {
			return strings.TrimSpace(strings.TrimPrefix(line, "**Specification version:**"))
		}
	}
	return ""
}

func emptyDisplay(value string) string {
	if value == "" {
		return "not recorded"
	}
	return value
}
