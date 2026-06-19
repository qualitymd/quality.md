package cli

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
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
	for _, want := range []string{`"version": "v1.2.3"`, `"commit": "abc1234"`, `"developmentBuild": false`, `"specificationVersion": "0.1 (Draft)"`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %s, want %s", out.String(), want)
		}
	}
}

func TestUpgradeCheckJSONUsesNPMMarker(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "npm")
	withLatestVersion(t, func(_ context.Context, method installMethod) (string, error) {
		if method != installNPM {
			t.Fatalf("install method = %s, want npm", method)
		}
		return "v9.9.9", nil
	})

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"upgrade", "--check", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{`"installMethod": "npm"`, `"latestVersion": "v9.9.9"`, `"applySupported": true`, `"recommendedCommand": "npm install -g quality.md@latest"`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %s, want %s", out.String(), want)
		}
	}
}

func TestUpgradeApplyRunsSupportedOwnerCommand(t *testing.T) {
	t.Setenv("QUALITYMD_INSTALL_METHOD", "npm")
	withLatestVersion(t, func(_ context.Context, _ installMethod) (string, error) {
		return "v9.9.9", nil
	})
	var ran string
	withUpgradeRunner(t, func(_ context.Context, name string, args []string) error {
		ran = name + " " + strings.Join(args, " ")
		return nil
	})

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"upgrade", "--apply", "--json"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if ran != "npm install -g quality.md@latest" {
		t.Fatalf("ran = %q, want npm install", ran)
	}
	if !strings.Contains(out.String(), `"applied": true`) {
		t.Fatalf("stdout = %s, want applied true", out.String())
	}
}

func TestUpgradeApplyRefusesUnknownInstall(t *testing.T) {
	if err := os.Unsetenv("QUALITYMD_INSTALL_METHOD"); err != nil {
		t.Fatalf("Unsetenv() error = %v", err)
	}
	withLatestVersion(t, func(_ context.Context, _ installMethod) (string, error) {
		return "v9.9.9", nil
	})

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"upgrade", "--apply"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want unsupported apply")
	}
}

func TestUpdateAvailableSemVer(t *testing.T) {
	cases := []struct {
		name            string
		current, latest string
		want            bool
	}{
		{"newer patch", "v0.4.0", "v0.4.1", true},
		{"latest older", "v0.5.0", "v0.4.1", false},
		{"equal", "v0.4.1", "v0.4.1", false},
		{"release beats prerelease", "v0.4.1-rc.1", "v0.4.1", true},
		{"prerelease loses to release", "v0.4.1", "v0.4.1-rc.1", false},
		{"development build", "dev (abc1234)", "v9.9.9", false},
		{"prefixed current", "qualitymd v0.4.0", "v0.4.1", true},
		{"empty latest", "v0.4.0", "", false},
		{"non-semver differ", "weird", "weirder", true},
	}
	for _, tc := range cases {
		if got := updateAvailable(tc.current, tc.latest); got != tc.want {
			t.Errorf("%s: updateAvailable(%q, %q) = %v, want %v", tc.name, tc.current, tc.latest, got, tc.want)
		}
	}
}

func withLatestVersion(t *testing.T, provider latestVersionProvider) {
	t.Helper()
	saved := fetchLatestVersion
	fetchLatestVersion = provider
	t.Cleanup(func() { fetchLatestVersion = saved })
}

func withUpgradeRunner(t *testing.T, runner upgradeCommandRunner) {
	t.Helper()
	saved := runUpgradeCommand
	runUpgradeCommand = runner
	t.Cleanup(func() { runUpgradeCommand = saved })
}
