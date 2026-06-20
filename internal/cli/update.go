package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/x/term"
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

const (
	updateCacheName     = ".qualitymd-update-cache"
	updateCacheTTL      = 20 * time.Hour
	updateRefreshMarker = "QUALITYMD_UPDATE_REFRESH"
)

type updateResult struct {
	SchemaVersion        int    `json:"schemaVersion"`
	CurrentVersion       string `json:"currentVersion"`
	Commit               string `json:"commit,omitempty"`
	SpecificationVersion string `json:"specificationVersion"`
	DevelopmentBuild     bool   `json:"developmentBuild"`
	InstallMethod        string `json:"installMethod"`
	LatestVersion        string `json:"latestVersion,omitempty"`
	LatestVersionReady   bool   `json:"latestVersionReady"`
	UpdateAvailable      bool   `json:"updateAvailable"`
	ApplySupported       bool   `json:"applySupported"`
	RecommendedAction    string `json:"recommendedAction,omitempty"`
	RecommendedCommand   string `json:"recommendedCommand,omitempty"`
	Applied              bool   `json:"applied"`
	ReleaseNotesURL      string `json:"releaseNotesURL,omitempty"`
}

type latestRelease struct {
	Version         string
	Ready           bool
	ReleaseNotesURL string
}

type updateCacheRecord struct {
	LatestVersion   string    `json:"latestVersion,omitempty"`
	ReleaseNotesURL string    `json:"releaseNotesURL,omitempty"`
	Ready           bool      `json:"ready"`
	CheckedAt       time.Time `json:"checkedAt"`
}

type latestVersionProvider func(context.Context, installMethod) (latestRelease, error)
type updateCommandRunner func(context.Context, string, []string) error
type installedVersionReader func(context.Context) (string, error)

var (
	fetchLatestVersion  latestVersionProvider  = defaultLatestVersion
	runUpdateCommand    updateCommandRunner    = defaultRunUpdateCommand
	readVisibleVersion  installedVersionReader = defaultVisibleVersion
	now                                        = time.Now
	allowAmbientUpdates                        = defaultAllowAmbientUpdates
	startUpdateRefresh                         = defaultStartUpdateRefresh
)

func newUpdateCmd() *cobra.Command {
	var checkOnly bool
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the qualitymd CLI through its owning install channel",
		Args:  usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, newerButNotReady, err := checkUpdate(cmd.Context())
			if err != nil {
				return err
			}
			if !checkOnly {
				if err := applyUpdate(cmd.Context(), &result, newerButNotReady); err != nil {
					return err
				}
			}
			writeUpdateCacheFromResult(result)
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			return renderUpdateResult(cmd.OutOrStdout(), result, checkOnly)
		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check", false, "check whether a newer CLI release is available without applying it")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable update result")
	return cmd
}

func newUpdateRefreshCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "__update-refresh",
		Hidden: true,
		Args:   usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			_ = refreshUpdateCache(cmd.Context())
			return nil
		},
	}
}

func checkUpdate(ctx context.Context) (updateResult, bool, error) {
	info := currentVersionInfo()
	method := detectInstallMethod()
	command := recommendedCommand(method)
	if updateChecksDisabled() {
		return updateResult{
			SchemaVersion:        1,
			CurrentVersion:       info.Version,
			Commit:               info.Commit,
			SpecificationVersion: info.SpecificationVersion,
			DevelopmentBuild:     info.DevelopmentBuild,
			InstallMethod:        string(method),
			ApplySupported:       applySupported(method),
			RecommendedAction:    recommendedAction(method, command),
			RecommendedCommand:   command,
		}, false, nil
	}
	latest, err := fetchLatestVersion(ctx, method)
	if err != nil {
		return updateResult{}, false, err
	}
	newer := updateAvailable(info.Version, latest.Version)
	available := !info.DevelopmentBuild && newer && latest.Ready
	result := updateResult{
		SchemaVersion:        1,
		CurrentVersion:       info.Version,
		Commit:               info.Commit,
		SpecificationVersion: info.SpecificationVersion,
		DevelopmentBuild:     info.DevelopmentBuild,
		InstallMethod:        string(method),
		LatestVersion:        latest.Version,
		LatestVersionReady:   latest.Ready,
		UpdateAvailable:      available,
		ApplySupported:       applySupported(method),
		RecommendedAction:    recommendedAction(method, command),
		RecommendedCommand:   command,
		ReleaseNotesURL:      latest.ReleaseNotesURL,
	}
	return result, !info.DevelopmentBuild && newer && !latest.Ready, nil
}

func applyUpdate(ctx context.Context, result *updateResult, newerButNotReady bool) error {
	if result.DevelopmentBuild {
		return nil
	}
	if newerButNotReady {
		return fmt.Errorf("latest qualitymd release %s is not yet available for %s installs", result.LatestVersion, result.InstallMethod)
	}
	if !result.UpdateAvailable {
		return nil
	}
	if !result.ApplySupported {
		return fmt.Errorf("update apply is not supported for %s installs; %s", result.InstallMethod, result.RecommendedAction)
	}
	name, args := updateCommand(installMethod(result.InstallMethod))
	if name == "" {
		return fmt.Errorf("update apply is not supported for %s installs; %s", result.InstallMethod, result.RecommendedAction)
	}
	if err := runUpdateCommand(ctx, name, args); err != nil {
		return err
	}
	visibleVersion, err := readVisibleVersion(ctx)
	if err != nil {
		return fmt.Errorf("update command ran, but qualitymd --version could not be verified: %w", err)
	}
	if !sameRelease(visibleVersion, result.LatestVersion) {
		return fmt.Errorf("update command ran, but visible qualitymd version is %s, not %s", emptyDisplay(visibleVersion), result.LatestVersion)
	}
	result.CurrentVersion = visibleVersion
	result.UpdateAvailable = false
	result.Applied = true
	result.RecommendedAction = ""
	result.RecommendedCommand = ""
	return nil
}

func renderUpdateResult(w interface{ Write([]byte) (int, error) }, result updateResult, checkOnly bool) error {
	if _, err := fmt.Fprintf(w, "Current version: %s\nLatest version: %s\nLatest version ready: %t\nInstall method: %s\nUpdate available: %t\n",
		result.CurrentVersion, emptyDisplay(result.LatestVersion), result.LatestVersionReady, result.InstallMethod, result.UpdateAvailable); err != nil {
		return err
	}
	if result.ReleaseNotesURL != "" {
		if _, err := fmt.Fprintf(w, "Release notes: %s\n", result.ReleaseNotesURL); err != nil {
			return err
		}
	}
	if result.Applied {
		_, err := fmt.Fprintln(w, "Update applied.")
		return err
	}
	if !checkOnly && !result.UpdateAvailable {
		_, err := fmt.Fprintln(w, "Already up to date.")
		return err
	}
	if result.ApplySupported {
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
	if strings.EqualFold(os.Getenv("QUALITYMD_INSTALL_METHOD"), "homebrew") {
		return installHomebrew
	}
	if strings.EqualFold(os.Getenv("QUALITYMD_INSTALL_METHOD"), "managed-standalone") {
		return installManagedStandalone
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
	root := qualitymdHome()
	marker := filepath.Join(root, ".qualitymd-managed-install")
	if _, err := os.Stat(marker); err != nil {
		return false
	}
	rel, err := filepath.Rel(root, exe)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func defaultLatestVersion(ctx context.Context, method installMethod) (latestRelease, error) {
	switch method {
	case installNPM:
		return latestNPMVersion(ctx)
	case installHomebrew:
		return latestHomebrewVersion(ctx)
	default:
		return latestGitHubVersion(ctx)
	}
}

func latestGitHubVersion(ctx context.Context) (latestRelease, error) {
	var payload struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
		Assets  []struct {
			Name string `json:"name"`
		} `json:"assets"`
	}
	if err := getJSON(ctx, "https://api.github.com/repos/qualitymd/quality.md/releases/latest", &payload); err != nil {
		return latestRelease{}, err
	}
	return latestRelease{
		Version:         payload.TagName,
		Ready:           githubReleaseReady(payload.Assets),
		ReleaseNotesURL: payload.HTMLURL,
	}, nil
}

func githubReleaseReady(assets []struct {
	Name string `json:"name"`
}) bool {
	archive := platformArchiveName()
	if archive == "" {
		return false
	}
	var hasArchive, hasChecksums bool
	for _, asset := range assets {
		switch asset.Name {
		case archive:
			hasArchive = true
		case "checksums.txt":
			hasChecksums = true
		}
	}
	return hasArchive && hasChecksums
}

func platformArchiveName() string {
	arch := runtime.GOARCH
	if arch != "amd64" && arch != "arm64" {
		return ""
	}
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("qualitymd_windows_%s.zip", arch)
	}
	switch runtime.GOOS {
	case "darwin", "linux":
		return fmt.Sprintf("qualitymd_%s_%s.tar.gz", runtime.GOOS, arch)
	default:
		return ""
	}
}

func latestNPMVersion(ctx context.Context) (latestRelease, error) {
	var payload struct {
		Version string `json:"version"`
	}
	if err := getJSON(ctx, "https://registry.npmjs.org/quality.md/latest", &payload); err != nil {
		return latestRelease{}, err
	}
	if payload.Version == "" {
		return latestRelease{}, nil
	}
	return latestRelease{Version: "v" + strings.TrimPrefix(payload.Version, "v"), Ready: true}, nil
}

func latestHomebrewVersion(ctx context.Context) (latestRelease, error) {
	body, err := getText(ctx, "https://raw.githubusercontent.com/qualitymd/homebrew-tap/main/Casks/qualitymd.rb")
	if err != nil {
		return latestRelease{}, err
	}
	versionRe := regexp.MustCompile(`(?m)^\s*version\s+"([^"]+)"`)
	match := versionRe.FindStringSubmatch(body)
	if len(match) != 2 {
		return latestRelease{}, nil
	}
	version := "v" + strings.TrimPrefix(match[1], "v")
	return latestRelease{
		Version:         version,
		Ready:           true,
		ReleaseNotesURL: "https://github.com/qualitymd/quality.md/releases/tag/" + version,
	}, nil
}

func getJSON(ctx context.Context, url string, target any) error {
	body, err := getBytes(ctx, url, "application/json")
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

func getText(ctx context.Context, url string) (string, error) {
	body, err := getBytes(ctx, url, "text/plain")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getBytes(ctx context.Context, url, accept string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", "qualitymd")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("latest version check failed: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func updateAvailable(current, latest string) bool {
	current = normalizeVersion(current)
	latest = normalizeVersion(latest)
	if current == "" || latest == "" {
		return false
	}
	cv, lv := "v"+current, "v"+latest
	if !semver.IsValid(cv) || !semver.IsValid(lv) {
		return false
	}
	if semver.Prerelease(cv) != "" || semver.Prerelease(lv) != "" {
		return false
	}
	return semver.Compare(lv, cv) > 0
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

func sameRelease(got, want string) bool {
	got = normalizeVersion(got)
	want = normalizeVersion(want)
	return got != "" && want != "" && got == want
}

func recommendedCommand(method installMethod) string {
	switch method {
	case installNPM:
		return "npm install -g quality.md@latest"
	case installHomebrew:
		return "brew upgrade qualitymd/tap/qualitymd"
	case installManagedStandalone:
		if runtime.GOOS == "windows" {
			return `powershell -NoProfile -ExecutionPolicy Bypass -Command "$env:QUALITYMD_NO_INPUT='1'; iwr https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.ps1 -UseB | iex"`
		}
		return "curl -fsSL https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.sh | QUALITYMD_NO_INPUT=1 sh"
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
	return method == installNPM || method == installHomebrew || method == installManagedStandalone
}

func updateCommand(method installMethod) (string, []string) {
	switch method {
	case installNPM:
		return "npm", []string{"install", "-g", "quality.md@latest"}
	case installHomebrew:
		return "brew", []string{"upgrade", "qualitymd/tap/qualitymd"}
	case installManagedStandalone:
		if runtime.GOOS == "windows" {
			return "powershell", []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "$env:QUALITYMD_NO_INPUT='1'; iwr https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.ps1 -UseB | iex"}
		}
		return "sh", []string{"-c", "curl -fsSL https://raw.githubusercontent.com/qualitymd/quality.md/main/install/install.sh | QUALITYMD_NO_INPUT=1 sh"}
	default:
		return "", nil
	}
}

func defaultRunUpdateCommand(ctx context.Context, name string, args []string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func defaultVisibleVersion(ctx context.Context) (string, error) {
	name, err := exec.LookPath("qualitymd")
	if err != nil {
		name, err = os.Executable()
		if err != nil {
			return "", err
		}
	}
	out, err := exec.CommandContext(ctx, name, "--version").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func qualitymdHome() string {
	if root := os.Getenv("QUALITYMD_HOME"); root != "" {
		return root
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".qualitymd"
	}
	return filepath.Join(home, ".qualitymd")
}

func updateCachePath() string {
	return filepath.Join(qualitymdHome(), updateCacheName)
}

func readUpdateCache() (updateCacheRecord, error) {
	var record updateCacheRecord
	data, err := os.ReadFile(updateCachePath())
	if err != nil {
		return record, err
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return updateCacheRecord{}, err
	}
	return record, nil
}

func writeUpdateCache(record updateCacheRecord) error {
	path := updateCachePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o600)
}

func writeUpdateCacheFromResult(result updateResult) {
	if result.DevelopmentBuild || result.LatestVersion == "" {
		return
	}
	_ = writeUpdateCache(updateCacheRecord{
		LatestVersion:   result.LatestVersion,
		ReleaseNotesURL: result.ReleaseNotesURL,
		Ready:           result.LatestVersionReady,
		CheckedAt:       now().UTC(),
	})
}

func refreshUpdateCache(ctx context.Context) error {
	if updateChecksDisabled() {
		return nil
	}
	info := currentVersionInfo()
	if info.DevelopmentBuild {
		return nil
	}
	method := detectInstallMethod()
	latest, err := fetchLatestVersion(ctx, method)
	if err != nil {
		return err
	}
	if latest.Version == "" {
		return nil
	}
	return writeUpdateCache(updateCacheRecord{
		LatestVersion:   latest.Version,
		ReleaseNotesURL: latest.ReleaseNotesURL,
		Ready:           latest.Ready,
		CheckedAt:       now().UTC(),
	})
}

func maybeStartUpdateRefresh(cmd *cobra.Command) {
	if os.Getenv(updateRefreshMarker) == "1" || !allowAmbientUpdates(cmd.ErrOrStderr()) || commandIs(cmd, "__update-refresh") || commandIs(cmd, "update") {
		return
	}
	record, err := readUpdateCache()
	if err == nil && now().Sub(record.CheckedAt) < updateCacheTTL {
		return
	}
	_ = startUpdateRefresh()
}

func maybeEmitUpdateNotice(cmd *cobra.Command) {
	if os.Getenv(updateRefreshMarker) == "1" || !allowAmbientUpdates(cmd.ErrOrStderr()) || commandHasJSON(cmd) || commandIs(cmd, "__update-refresh") || commandIs(cmd, "update") {
		return
	}
	record, err := readUpdateCache()
	if err != nil || !record.Ready || !updateAvailable(currentVersionInfo().Version, record.LatestVersion) {
		return
	}
	message := fmt.Sprintf("update available: %s -> %s (run `qualitymd update`)", currentVersionInfo().Version, record.LatestVersion)
	if record.ReleaseNotesURL != "" {
		message += " " + record.ReleaseNotesURL
	}
	_, _ = fmt.Fprintln(cmd.ErrOrStderr(), message)
}

func defaultAllowAmbientUpdates(w io.Writer) bool {
	if updateChecksDisabled() || os.Getenv("CI") != "" || currentVersionInfo().DevelopmentBuild {
		return false
	}
	fdWriter, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return false
	}
	return term.IsTerminal(fdWriter.Fd())
}

func updateChecksDisabled() bool {
	return envTruthy("QUALITYMD_NO_UPDATE_CHECK")
}

func envTruthy(name string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(name))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func defaultStartUpdateRefresh() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, "__update-refresh")
	cmd.Env = append(os.Environ(), updateRefreshMarker+"=1")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	configureDetachedUpdateRefresh(cmd)
	return cmd.Start()
}

func commandHasJSON(cmd *cobra.Command) bool {
	for c := cmd; c != nil; c = c.Parent() {
		flag := c.Flags().Lookup("json")
		if flag != nil && flag.Changed {
			return true
		}
		pflag := c.PersistentFlags().Lookup("json")
		if pflag != nil && pflag.Changed {
			return true
		}
	}
	return false
}

func commandIs(cmd *cobra.Command, name string) bool {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Name() == name {
			return true
		}
	}
	return false
}
