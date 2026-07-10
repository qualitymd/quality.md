package evaluator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/qualitymd/quality.md/internal/workspace"
)

// ReservedNames are the built-in evaluator names. Custom evaluator profiles
// must not shadow them.
var ReservedNames = []string{"auto", "harness", "codex", "claude", "openai", "anthropic", "shell", "manual"}

// Default provider models for the built-in API-backed evaluators. Profiles
// override these per workspace.
const (
	defaultAnthropicModel = "claude-sonnet-5"
	defaultOpenAIModel    = "gpt-5.1"
	defaultAnthropicEnv   = "ANTHROPIC_API_KEY"
	defaultOpenAIEnv      = "OPENAI_API_KEY"
)

// SelectionError is a typed evaluator selection failure carrying the stable
// failure category and remediation guidance.
type SelectionError struct {
	Category FailureCategory
	Message  string
	Remedies []string
}

func (e *SelectionError) Error() string {
	if len(e.Remedies) == 0 {
		return e.Message
	}
	return e.Message + " (" + strings.Join(e.Remedies, "; ") + ")"
}

// Selection is one resolved evaluator choice.
type Selection struct {
	// Name is the selected evaluator or profile name.
	Name string
	// Evaluator is the ready evaluator runtime.
	Evaluator Evaluator
	// Reason explains how the selection was made (for logs and receipts).
	Reason string
	// Candidates carries the readiness evidence for every CLI candidate auto
	// discovery considered, in consideration order. Empty for explicit
	// selections.
	Candidates []CLIReadiness
}

// CLIReadiness is the readiness evidence for one CLI evaluator candidate:
// executable presence, non-interactive structured-output capability, and
// authentication state where a documented non-interactive probe exists. It
// never carries credential values.
type CLIReadiness struct {
	Name       string `json:"name"`
	Executable bool   `json:"executable"`
	// StructuredOutput reports whether the installed CLI advertises the
	// non-interactive structured-output invocation the runner requires.
	StructuredOutput bool `json:"structuredOutput"`
	// Authenticated is nil when the CLI documents no non-interactive
	// authentication probe; readiness then assumes authentication and the
	// evidence says so.
	Authenticated *bool `json:"authenticated,omitempty"`
	// Usable reports whether auto discovery may select this candidate.
	Usable bool `json:"usable"`
	// Evidence lists human-readable probe observations.
	Evidence []string `json:"evidence,omitempty"`
}

// Options configures evaluator selection.
type Options struct {
	// Name is the requested evaluator: a built-in name, a configured profile
	// name, or "auto" (also the default when empty).
	Name string
	// Profiles are the configured evaluator profiles from workspace config.
	Profiles map[string]workspace.EvaluatorProfile
	// LookPath overrides executable discovery in tests.
	LookPath func(file string) (string, error)
	// Getenv overrides environment lookup in tests.
	Getenv func(key string) string
	// RunProbe overrides readiness probe subprocess execution in tests. It
	// returns the probe command's combined output.
	RunProbe func(name string, args ...string) (string, error)
}

func (o Options) lookPath(file string) (string, error) {
	if o.LookPath != nil {
		return o.LookPath(file)
	}
	return exec.LookPath(file)
}

func (o Options) getenv(key string) string {
	if o.Getenv != nil {
		return o.Getenv(key)
	}
	return os.Getenv(key)
}

// probeTimeout bounds one readiness probe subprocess.
const probeTimeout = 10 * time.Second

func (o Options) runProbe(name string, args ...string) (string, error) {
	if o.RunProbe != nil {
		return o.RunProbe(name, args...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	return string(out), err
}

// ValidateProfiles rejects configured profiles that shadow reserved names or
// omit a kind.
func ValidateProfiles(profiles map[string]workspace.EvaluatorProfile) error {
	for name, profile := range profiles {
		if slices.Contains(ReservedNames, name) {
			return &SelectionError{
				Category: FailureRunStateInvalid,
				Message:  fmt.Sprintf("evaluator profile %q shadows a reserved built-in evaluator name", name),
				Remedies: []string{"rename the profile in .quality/config.yaml"},
			}
		}
		if strings.TrimSpace(profile.Kind) == "" {
			return &SelectionError{
				Category: FailureRunStateInvalid,
				Message:  fmt.Sprintf("evaluator profile %q must declare a kind", name),
				Remedies: []string{"set evaluators." + name + ".kind in .quality/config.yaml"},
			}
		}
	}
	return nil
}

// Select resolves an evaluator name to a ready evaluator runtime. The
// resolution never prompts: a missing or ambiguous evaluator fails with a
// typed SelectionError listing remedies.
func Select(opts Options) (*Selection, error) {
	if err := ValidateProfiles(opts.Profiles); err != nil {
		return nil, err
	}
	name := opts.Name
	if name == "" {
		name = "auto"
	}
	if name == "auto" {
		return selectAuto(opts)
	}
	if profile, ok := opts.Profiles[name]; ok {
		return selectProfile(name, profile, opts)
	}
	if slices.Contains(ReservedNames, name) {
		return selectBuiltin(name, opts)
	}
	return nil, &SelectionError{
		Category: FailureMissingEvaluator,
		Message:  fmt.Sprintf("unknown evaluator %q", name),
		Remedies: availableEvaluatorRemedies(opts),
	}
}

func selectBuiltin(name string, opts Options) (*Selection, error) {
	switch name {
	case "harness":
		return &Selection{
			Name:      "harness",
			Evaluator: &harnessEvaluator{},
			Reason:    "requested harness evaluator: judgment is supplied by the invoking agent harness through checkpoints",
		}, nil
	case "codex", "claude":
		ev, err := newCLIEvaluator(name, "", opts)
		if err != nil {
			return nil, err
		}
		return &Selection{Name: name, Evaluator: ev, Reason: "requested built-in CLI evaluator"}, nil
	case "openai":
		ev, err := newOpenAIEvaluator(workspace.EvaluatorProfile{Kind: "openai"}, opts)
		if err != nil {
			return nil, err
		}
		return &Selection{Name: name, Evaluator: ev, Reason: "requested built-in API evaluator"}, nil
	case "anthropic":
		ev, err := newAnthropicEvaluator(workspace.EvaluatorProfile{Kind: "anthropic"}, opts)
		if err != nil {
			return nil, err
		}
		return &Selection{Name: name, Evaluator: ev, Reason: "requested built-in API evaluator"}, nil
	case "shell", "manual":
		return nil, &SelectionError{
			Category: FailureMissingEvaluator,
			Message:  fmt.Sprintf("evaluator %q is a reserved name without an implementation yet", name),
			Remedies: availableEvaluatorRemedies(opts),
		}
	default:
		return nil, &SelectionError{Category: FailureMissingEvaluator, Message: fmt.Sprintf("unknown evaluator %q", name)}
	}
}

func selectProfile(name string, profile workspace.EvaluatorProfile, opts Options) (*Selection, error) {
	var ev Evaluator
	var err error
	switch profile.Kind {
	case "codex", "claude":
		ev, err = newCLIEvaluator(profile.Kind, profile.Command, opts)
	case "openai":
		ev, err = newOpenAIEvaluator(profile, opts)
	case "anthropic":
		ev, err = newAnthropicEvaluator(profile, opts)
	case "shell", "manual":
		return nil, &SelectionError{
			Category: FailureMissingEvaluator,
			Message:  fmt.Sprintf("evaluator profile %q uses reserved kind %q, which has no implementation yet", name, profile.Kind),
		}
	default:
		return nil, &SelectionError{
			Category: FailureMissingEvaluator,
			Message:  fmt.Sprintf("evaluator profile %q declares unsupported kind %q", name, profile.Kind),
			Remedies: []string{"use one of: codex, claude, openai, anthropic"},
		}
	}
	if err != nil {
		return nil, err
	}
	return &Selection{Name: name, Evaluator: ev, Reason: "configured evaluator profile"}, nil
}

// selectAuto performs deterministic local discovery: a ready Codex CLI, then
// a ready Claude CLI, then configured API profiles whose key environment
// variable is present, then a clear non-interactive failure. A CLI candidate
// is usable only after its readiness probe verifies executable presence,
// authentication where non-interactively verifiable, and the required
// structured-output invocation; auto never infers a parent agent harness —
// harness-backed runs are selected explicitly.
func selectAuto(opts Options) (*Selection, error) {
	var candidates []CLIReadiness
	for _, name := range []string{"codex", "claude"} {
		readiness := probeCLIReadiness(name, opts)
		candidates = append(candidates, readiness)
		if !readiness.Usable {
			continue
		}
		ev, err := newCLIEvaluator(name, "", opts)
		if err != nil {
			continue
		}
		return &Selection{
			Name:       name,
			Evaluator:  ev,
			Reason:     "auto: " + name + " CLI is installed and ready (" + strings.Join(readiness.Evidence, "; ") + ")",
			Candidates: candidates,
		}, nil
	}
	for _, name := range sortedProfileNames(opts.Profiles) {
		profile := opts.Profiles[name]
		if !profileAPIKeyPresent(profile, opts) {
			continue
		}
		selection, err := selectProfile(name, profile, opts)
		if err != nil {
			continue
		}
		selection.Reason = "auto: configured profile with API key present"
		selection.Candidates = candidates
		return selection, nil
	}
	return nil, &SelectionError{
		Category: FailureMissingEvaluator,
		Message:  "no evaluator is available" + readinessSummary(candidates),
		Remedies: []string{
			"install and authenticate the codex or claude CLI",
			"configure an evaluator profile in .quality/config.yaml and export its API key environment variable",
			"pass --evaluator <name> to select one explicitly",
		},
	}
}

// probeCLIReadiness verifies one CLI candidate beyond command presence:
// required structured-output flags from its help output, and authentication
// through the CLI's documented non-interactive probe where one exists.
func probeCLIReadiness(name string, opts Options) CLIReadiness {
	readiness := CLIReadiness{Name: name}
	command := cliCommandName(name, "")
	if _, err := opts.lookPath(command); err != nil {
		readiness.Evidence = append(readiness.Evidence, "executable not found on PATH")
		return readiness
	}
	readiness.Executable = true

	helpArgs := []string{"--help"}
	requiredFlag := "--output-format"
	if name == "codex" {
		helpArgs = []string{"exec", "--help"}
		requiredFlag = "--json"
	}
	help, err := opts.runProbe(command, helpArgs...)
	if err != nil {
		readiness.Evidence = append(readiness.Evidence, "help probe failed: "+firstLineOr(help, err.Error()))
		return readiness
	}
	if !strings.Contains(help, requiredFlag) {
		readiness.Evidence = append(readiness.Evidence,
			fmt.Sprintf("installed CLI does not advertise the required %s structured-output flag", requiredFlag))
		return readiness
	}
	readiness.StructuredOutput = true
	readiness.Evidence = append(readiness.Evidence, "non-interactive structured output available")

	if name == "codex" {
		out, err := opts.runProbe(command, "login", "status")
		authenticated := err == nil
		readiness.Authenticated = &authenticated
		if authenticated {
			readiness.Evidence = append(readiness.Evidence, "authenticated (codex login status)")
		} else {
			readiness.Evidence = append(readiness.Evidence, "not authenticated: "+firstLineOr(out, err.Error()))
		}
	} else {
		readiness.Evidence = append(readiness.Evidence,
			"authentication is not verifiable non-interactively; assumed available")
	}
	readiness.Usable = readiness.Executable && readiness.StructuredOutput &&
		(readiness.Authenticated == nil || *readiness.Authenticated)
	return readiness
}

func readinessSummary(candidates []CLIReadiness) string {
	parts := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		parts = append(parts, candidate.Name+": "+strings.Join(candidate.Evidence, ", "))
	}
	if len(parts) == 0 {
		return ""
	}
	return " (" + strings.Join(parts, "; ") + ")"
}

func profileAPIKeyPresent(profile workspace.EvaluatorProfile, opts Options) bool {
	switch profile.Kind {
	case "openai":
		return opts.getenv(profileEnv(profile, defaultOpenAIEnv)) != ""
	case "anthropic":
		return opts.getenv(profileEnv(profile, defaultAnthropicEnv)) != ""
	default:
		return false
	}
}

func profileEnv(profile workspace.EvaluatorProfile, fallback string) string {
	if profile.APIKeyEnv != "" {
		return profile.APIKeyEnv
	}
	return fallback
}

func sortedProfileNames(profiles map[string]workspace.EvaluatorProfile) []string {
	names := make([]string, 0, len(profiles))
	for name := range profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func availableEvaluatorRemedies(opts Options) []string {
	remedies := []string{"use a built-in evaluator: harness (agent-supplied judgment), codex, claude, openai, anthropic"}
	if names := sortedProfileNames(opts.Profiles); len(names) > 0 {
		remedies = append(remedies, "or a configured profile: "+strings.Join(names, ", "))
	}
	return remedies
}
