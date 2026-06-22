package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestVersionJSONCommand(t *testing.T) {
	savedVersion, savedCommit := version, commit
	t.Cleanup(func() { version, commit = savedVersion, savedCommit })
	version, commit = "v1.2.3", "abc1234"

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"version", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{`"version": "v1.2.3"`, `"commit": "abc1234"`, `"developmentBuild": false`, `"specificationVersion": "0.3 (Draft)"`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %s, want %s", out.String(), want)
		}
	}
}

func TestCurrentVersionInfoMarksPseudoVersionDevelopmentBuild(t *testing.T) {
	withVersion(t, "v0.4.2-0.20260620161924-c48da3ec9d7f+dirty", "c48da3e")

	info := currentVersionInfo()
	if !info.DevelopmentBuild {
		t.Fatalf("DevelopmentBuild = false, want true for pseudo-version local build")
	}
}

func TestUpdateCheckJSONUsesNPMMarker(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "npm")
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	withLatestVersion(t, func(_ context.Context, method installMethod) (latestRelease, error) {
		if method != installNPM {
			t.Fatalf("install method = %s, want npm", method)
		}
		return latestRelease{Version: "v9.9.9", Ready: true}, nil
	})

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update", "--check", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{`"installMethod": "npm"`, `"latestVersion": "v9.9.9"`, `"latestVersionReady": true`, `"updateAvailable": true`, `"applySupported": true`, `"applied": false`, `"recommendedCommand": "npm install -g quality.md@latest"`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %s, want %s", out.String(), want)
		}
	}
}

func TestUpdateDefaultAppliesSupportedOwnerCommand(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "npm")
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (latestRelease, error) {
		return latestRelease{Version: "v9.9.9", Ready: true, ReleaseNotesURL: "https://example.test/release"}, nil
	})
	var ran string
	withUpdateRunner(t, func(_ context.Context, name string, args []string) error {
		ran = name + " " + strings.Join(args, " ")
		return nil
	})
	withVisibleVersion(t, func(context.Context) (string, error) {
		return "qualitymd v9.9.9", nil
	})

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if ran != "npm install -g quality.md@latest" {
		t.Fatalf("ran = %q, want npm install", ran)
	}
	for _, want := range []string{`"currentVersion": "qualitymd v9.9.9"`, `"applied": true`, `"releaseNotesURL": "https://example.test/release"`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %s, want %s", out.String(), want)
		}
	}
}

func TestUpdateCheckDoesNotApply(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "npm")
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (latestRelease, error) {
		return latestRelease{Version: "v9.9.9", Ready: true}, nil
	})
	withUpdateRunner(t, func(_ context.Context, _ string, _ []string) error {
		t.Fatal("update runner called during --check")
		return nil
	})

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update", "--check"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestUpdateOptOutDisablesExplicitCheck(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "npm")
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	t.Setenv("QUALITYMD_NO_UPDATE_CHECK", "1")
	withVersion(t, "v1.2.3", "abc1234")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (latestRelease, error) {
		t.Fatal("latest-version provider called with update checks disabled")
		return latestRelease{}, nil
	})

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update", "--check", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{`"latestVersionReady": false`, `"updateAvailable": false`, `"applied": false`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %s, want %s", out.String(), want)
		}
	}
}

func TestUpdateDefaultRefusesUnknownInstallWhenUpdateAvailable(t *testing.T) {
	if err := os.Unsetenv("QUALITYMD_INSTALL_METHOD"); err != nil {
		t.Fatalf("Unsetenv() error = %v", err)
	}
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (latestRelease, error) {
		return latestRelease{Version: "v9.9.9", Ready: true}, nil
	})

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want unsupported apply")
	}
}

func TestUpdateDefaultDoesNotApplyDevelopmentBuild(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "npm")
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "dev (abc1234)", "")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (latestRelease, error) {
		return latestRelease{Version: "v9.9.9", Ready: true}, nil
	})
	withUpdateRunner(t, func(_ context.Context, _ string, _ []string) error {
		t.Fatal("update runner called for development build")
		return nil
	})

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{`"developmentBuild": true`, `"updateAvailable": false`, `"applied": false`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %s, want %s", out.String(), want)
		}
	}
}

func TestUpdateRefusesNewerReleaseThatIsNotReadyBeforeApply(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "npm")
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (latestRelease, error) {
		return latestRelease{Version: "v9.9.9", Ready: false}, nil
	})
	withUpdateRunner(t, func(_ context.Context, _ string, _ []string) error {
		t.Fatal("update runner called for not-ready release")
		return nil
	})

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want not-ready refusal")
	}
	if !strings.Contains(err.Error(), "not yet available") {
		t.Fatalf("error = %v, want not-ready diagnostic", err)
	}
}

func TestUpdateManagedStandaloneAppliesThroughInstaller(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "managed-standalone")
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (latestRelease, error) {
		return latestRelease{Version: "v9.9.9", Ready: true}, nil
	})
	var ran string
	withUpdateRunner(t, func(_ context.Context, name string, args []string) error {
		ran = name + " " + strings.Join(args, " ")
		return nil
	})
	withVisibleVersion(t, func(context.Context) (string, error) {
		return "v9.9.9", nil
	})

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(ran, "install.sh") && !strings.Contains(ran, "install.ps1") {
		t.Fatalf("ran = %q, want managed installer", ran)
	}
	if !strings.Contains(ran, "QUALITYMD_NO_INPUT") {
		t.Fatalf("ran = %q, want non-interactive installer env", ran)
	}
}

func TestUpdateReportsFailedPostApplyVerification(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "homebrew")
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (latestRelease, error) {
		return latestRelease{Version: "v9.9.9", Ready: true}, nil
	})
	withUpdateRunner(t, func(_ context.Context, _ string, _ []string) error {
		return nil
	})
	withVisibleVersion(t, func(context.Context) (string, error) {
		return "v1.2.3", nil
	})

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update", "--json"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want verification failure")
	}
	if !strings.Contains(err.Error(), "not v9.9.9") {
		t.Fatalf("error = %v, want version verification detail", err)
	}
}

func TestUpgradeAliasRemoved(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"upgrade"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want unknown command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("error = %v, want unknown command", err)
	}
}

func TestUpdateAvailableSemVerOnly(t *testing.T) {
	cases := []struct {
		name    string
		current string
		latest  string
		want    bool
	}{
		{name: "newer patch", current: "v0.4.0", latest: "v0.4.1", want: true},
		{name: "latest older", current: "v0.5.0", latest: "v0.4.1", want: false},
		{name: "equal", current: "v0.4.1", latest: "v0.4.1", want: false},
		{name: "current prerelease", current: "v0.4.1-rc.1", latest: "v0.4.1", want: false},
		{name: "latest prerelease", current: "v0.4.1", latest: "v0.4.2-rc.1", want: false},
		{name: "development build", current: "dev (abc1234)", latest: "v9.9.9", want: false},
		{name: "prefixed current", current: "qualitymd v0.4.0", latest: "v0.4.1", want: true},
		{name: "empty latest", current: "v0.4.0", latest: "", want: false},
		{name: "non-semver differ", current: "weird", latest: "weirder", want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := updateAvailable(tc.current, tc.latest); got != tc.want {
				t.Errorf("updateAvailable(%q, %q) = %v, want %v", tc.current, tc.latest, got, tc.want)
			}
		})
	}
}

func TestAmbientNoticeUsesCacheAndSkipsJSON(t *testing.T) {
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	if err := writeUpdateCache(updateCacheRecord{
		LatestVersion:   "v9.9.9",
		ReleaseNotesURL: "https://example.test/release",
		Ready:           true,
		CheckedAt:       time.Now().UTC(),
	}); err != nil {
		t.Fatalf("writeUpdateCache() error = %v", err)
	}
	withAmbientAllowed(t, true)

	cmd := newRootCmd()
	var stderr bytes.Buffer
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"version", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want no notice for --json", stderr.String())
	}

	cmd = newRootCmd()
	stderr.Reset()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"version"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"update available: v1.2.3 -> v9.9.9", "qualitymd update", "https://example.test/release"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("stderr = %q, want %q", stderr.String(), want)
		}
	}
}

func TestAmbientRefreshSpawnsOnlyWhenAllowedAndCacheStale(t *testing.T) {
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	withAmbientAllowed(t, true)
	var starts int
	withStartUpdateRefresh(t, func() error {
		starts++
		return nil
	})

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"version"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if starts != 1 {
		t.Fatalf("starts = %d, want 1 for missing cache", starts)
	}

	if err := writeUpdateCache(updateCacheRecord{LatestVersion: "v9.9.9", Ready: true, CheckedAt: time.Now().UTC()}); err != nil {
		t.Fatalf("writeUpdateCache() error = %v", err)
	}
	cmd = newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"version"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if starts != 1 {
		t.Fatalf("starts = %d, want no second start for fresh cache", starts)
	}
}

func TestUpdateCheckRefreshesCache(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "npm")
	t.Setenv("QUALITYMD_HOME", t.TempDir())
	withVersion(t, "v1.2.3", "abc1234")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (latestRelease, error) {
		return latestRelease{Version: "v9.9.9", Ready: true, ReleaseNotesURL: "https://example.test/release"}, nil
	})

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"update", "--check", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	record, err := readUpdateCache()
	if err != nil {
		t.Fatalf("readUpdateCache() error = %v", err)
	}
	if record.LatestVersion != "v9.9.9" || !record.Ready || record.ReleaseNotesURL != "https://example.test/release" {
		data, _ := json.Marshal(record)
		t.Fatalf("cache = %s, want latest release facts", data)
	}
}

func withLatestVersion(t *testing.T, provider latestVersionProvider) {
	t.Helper()
	saved := fetchLatestVersion
	fetchLatestVersion = provider
	t.Cleanup(func() { fetchLatestVersion = saved })
}

func withUpdateRunner(t *testing.T, runner updateCommandRunner) {
	t.Helper()
	saved := runUpdateCommand
	runUpdateCommand = runner
	t.Cleanup(func() { runUpdateCommand = saved })
}

func withVisibleVersion(t *testing.T, reader installedVersionReader) {
	t.Helper()
	saved := readVisibleVersion
	readVisibleVersion = reader
	t.Cleanup(func() { readVisibleVersion = saved })
}

func withVersion(t *testing.T, v, c string) {
	t.Helper()
	savedVersion, savedCommit := version, commit
	version, commit = v, c
	t.Cleanup(func() { version, commit = savedVersion, savedCommit })
}

func withAmbientAllowed(t *testing.T, allowed bool) {
	t.Helper()
	saved := allowAmbientUpdates
	allowAmbientUpdates = func(io.Writer) bool { return allowed }
	t.Cleanup(func() { allowAmbientUpdates = saved })
}

func withStartUpdateRefresh(t *testing.T, starter func() error) {
	t.Helper()
	saved := startUpdateRefresh
	startUpdateRefresh = starter
	t.Cleanup(func() { startUpdateRefresh = saved })
}
