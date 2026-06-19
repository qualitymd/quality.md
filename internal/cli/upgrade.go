package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type installMethod string

const (
	installManagedStandalone installMethod = "managed-standalone"
	installNPM               installMethod = "npm"
	installHomebrew          installMethod = "homebrew"
	installGoSource          installMethod = "go-source"
	installArchive           installMethod = "archive"
	installUnknown           installMethod = "unknown"
)

type upgradeResult struct {
	SchemaVersion        int    `json:"schemaVersion"`
	CurrentVersion       string `json:"currentVersion"`
	Commit               string `json:"commit,omitempty"`
	SpecificationVersion string `json:"specificationVersion"`
	DevelopmentBuild     bool   `json:"developmentBuild"`
	InstallMethod        string `json:"installMethod"`
	LatestVersion        string `json:"latestVersion,omitempty"`
	UpdateAvailable      bool   `json:"updateAvailable"`
	ApplySupported       bool   `json:"applySupported"`
	RecommendedAction    string `json:"recommendedAction"`
	RecommendedCommand   string `json:"recommendedCommand,omitempty"`
	Applied              bool   `json:"applied"`
}

type latestVersionProvider func(context.Context, installMethod) (string, error)
type upgradeCommandRunner func(context.Context, string, []string) error

var (
	fetchLatestVersion latestVersionProvider = defaultLatestVersion
	runUpgradeCommand  upgradeCommandRunner  = defaultRunUpgradeCommand
)

func newUpgradeCmd() *cobra.Command {
	var checkOnly bool
	var apply bool
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Check for and advise on qualitymd CLI upgrades",
		Args:  usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if checkOnly && apply {
				return usageError(fmt.Errorf("--check and --apply cannot be used together"))
			}
			result, err := checkUpgrade(cmd.Context())
			if err != nil {
				return err
			}
			if apply {
				if !result.ApplySupported {
					return fmt.Errorf("upgrade apply is not supported for %s installs; %s", result.InstallMethod, result.RecommendedAction)
				}
				name, args := upgradeCommand(installMethod(result.InstallMethod))
				if name == "" {
					return fmt.Errorf("upgrade apply is not supported for %s installs; %s", result.InstallMethod, result.RecommendedAction)
				}
				if err := runUpgradeCommand(cmd.Context(), name, args); err != nil {
					return err
				}
				result.Applied = true
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			return renderUpgradeResult(cmd.OutOrStdout(), result, apply)
		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check", false, "check whether a newer CLI release is available")
	cmd.Flags().BoolVar(&apply, "apply", false, "apply an upgrade when the detected install method supports it")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable upgrade result")
	return cmd
}

func checkUpgrade(ctx context.Context) (upgradeResult, error) {
	info := currentVersionInfo()
	method := detectInstallMethod()
	latest, err := fetchLatestVersion(ctx, method)
	if err != nil {
		return upgradeResult{}, err
	}
	command := recommendedCommand(method)
	result := upgradeResult{
		SchemaVersion:        1,
		CurrentVersion:       info.Version,
		Commit:               info.Commit,
		SpecificationVersion: info.SpecificationVersion,
		DevelopmentBuild:     info.DevelopmentBuild,
		InstallMethod:        string(method),
		LatestVersion:        latest,
		UpdateAvailable:      updateAvailable(info.Version, latest),
		ApplySupported:       applySupported(method),
		RecommendedAction:    recommendedAction(method, command),
		RecommendedCommand:   command,
	}
	return result, nil
}

func renderUpgradeResult(w interface{ Write([]byte) (int, error) }, result upgradeResult, applied bool) error {
	if _, err := fmt.Fprintf(w, "Current version: %s\nLatest version: %s\nInstall method: %s\nUpdate available: %t\n",
		result.CurrentVersion, emptyDisplay(result.LatestVersion), result.InstallMethod, result.UpdateAvailable); err != nil {
		return err
	}
	if result.ApplySupported {
		if applied {
			_, err := fmt.Fprintln(w, "Upgrade applied.")
			return err
		}
		_, err := fmt.Fprintf(w, "Recommended command: %s\n", result.RecommendedCommand)
		return err
	}
	_, err := fmt.Fprintf(w, "Recommended action: %s\n", result.RecommendedAction)
	return err
}

func detectInstallMethod() installMethod {
	if strings.EqualFold(os.Getenv("QUALITYMD_INSTALL_METHOD"), "npm") {
		return installNPM
	}
	exe, err := os.Executable()
	if err != nil {
		return installUnknown
	}
	exe, _ = filepath.EvalSymlinks(exe)
	lower := strings.ToLower(filepath.ToSlash(exe))
	if managedRootOwns(exe) {
		return installManagedStandalone
	}
	if strings.Contains(lower, "/cellar/qualitymd/") || strings.Contains(lower, "/caskroom/qualitymd/") || strings.Contains(lower, "/homebrew/") && strings.Contains(lower, "/qualitymd") {
		return installHomebrew
	}
	if strings.Contains(lower, "/go-build/") || strings.Contains(lower, "/pkg/mod/") || strings.Contains(lower, "/go/bin/qualitymd") {
		return installGoSource
	}
	if strings.Contains(lower, "/qualitymd") {
		return installArchive
	}
	return installUnknown
}

func managedRootOwns(exe string) bool {
	root := os.Getenv("QUALITYMD_HOME")
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		root = filepath.Join(home, ".qualitymd")
	}
	marker := filepath.Join(root, ".qualitymd-managed-install")
	if _, err := os.Stat(marker); err != nil {
		return false
	}
	rel, err := filepath.Rel(root, exe)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func defaultLatestVersion(ctx context.Context, method installMethod) (string, error) {
	switch method {
	case installNPM:
		return latestNPMVersion(ctx)
	default:
		return latestGitHubVersion(ctx)
	}
}

func latestGitHubVersion(ctx context.Context) (string, error) {
	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := getJSON(ctx, "https://api.github.com/repos/qualitymd/quality.md/releases/latest", &payload); err != nil {
		return "", err
	}
	return payload.TagName, nil
}

func latestNPMVersion(ctx context.Context) (string, error) {
	var payload struct {
		Version string `json:"version"`
	}
	if err := getJSON(ctx, "https://registry.npmjs.org/quality.md/latest", &payload); err != nil {
		return "", err
	}
	if payload.Version == "" {
		return "", nil
	}
	return "v" + strings.TrimPrefix(payload.Version, "v"), nil
}

func getJSON(ctx context.Context, url string, target any) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "qualitymd")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("latest version check failed: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func updateAvailable(current, latest string) bool {
	current = normalizeVersion(current)
	latest = normalizeVersion(latest)
	if current == "" || latest == "" || strings.HasPrefix(current, "dev") {
		return false
	}
	// Compare by SemVer precedence when both parse, so a downgrade or an older
	// prerelease is not reported as an upgrade. semver wants a leading "v".
	cv, lv := "v"+current, "v"+latest
	if semver.IsValid(cv) && semver.IsValid(lv) {
		return semver.Compare(lv, cv) > 0
	}
	// Fall back to a plain difference for version strings SemVer can't parse.
	return current != latest
}

func normalizeVersion(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "qualitymd ")
	value = strings.TrimPrefix(value, "v")
	if i := strings.IndexAny(value, " +("); i >= 0 {
		value = value[:i]
	}
	return value
}

func recommendedCommand(method installMethod) string {
	switch method {
	case installNPM:
		return "npm install -g quality.md@latest"
	case installHomebrew:
		return "brew upgrade qualitymd/tap/qualitymd"
	case installManagedStandalone:
		if runtime.GOOS == "windows" {
			return `powershell -NoProfile -ExecutionPolicy Bypass -Command "iwr https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.ps1 -UseB | iex"`
		}
		return "curl -fsSL https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.sh | sh"
	default:
		return ""
	}
}

func recommendedAction(method installMethod, command string) string {
	if command != "" {
		return command
	}
	switch method {
	case installGoSource:
		return "reinstall with go install github.com/qualitymd/quality.md/cmd/qualitymd@latest"
	case installArchive:
		return "download and replace the archive-installed binary from the latest GitHub release"
	default:
		return "install or upgrade with a managed channel such as npm, Homebrew, or the GitHub-hosted installer"
	}
}

func applySupported(method installMethod) bool {
	return method == installNPM || method == installHomebrew
}

func upgradeCommand(method installMethod) (string, []string) {
	switch method {
	case installNPM:
		return "npm", []string{"install", "-g", "quality.md@latest"}
	case installHomebrew:
		return "brew", []string{"upgrade", "qualitymd/tap/qualitymd"}
	default:
		return "", nil
	}
}

func defaultRunUpgradeCommand(ctx context.Context, name string, args []string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
