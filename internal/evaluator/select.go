package evaluator

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"sort"
	"strings"

	"github.com/qualitymd/quality.md/internal/workspace"
)

// ReservedNames are the built-in evaluator names. Custom evaluator profiles
// must not shadow them.
var ReservedNames = []string{"auto", "codex", "claude", "openai", "anthropic", "shell", "manual"}

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

// selectAuto performs deterministic local discovery: an installed Codex CLI,
// then an installed Claude CLI, then configured API profiles whose key
// environment variable is present, then a clear non-interactive failure.
func selectAuto(opts Options) (*Selection, error) {
	for _, name := range []string{"codex", "claude"} {
		if _, err := opts.lookPath(cliCommandName(name, "")); err != nil {
			continue
		}
		ev, err := newCLIEvaluator(name, "", opts)
		if err != nil {
			continue
		}
		return &Selection{Name: name, Evaluator: ev, Reason: "auto: installed " + name + " CLI"}, nil
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
		return selection, nil
	}
	return nil, &SelectionError{
		Category: FailureMissingEvaluator,
		Message:  "no evaluator is available",
		Remedies: []string{
			"install and authenticate the codex or claude CLI",
			"configure an evaluator profile in .quality/config.yaml and export its API key environment variable",
			"pass --evaluator <name> to select one explicitly",
		},
	}
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
	remedies := []string{"use a built-in evaluator: codex, claude, openai, anthropic"}
	if names := sortedProfileNames(opts.Profiles); len(names) > 0 {
		remedies = append(remedies, "or a configured profile: "+strings.Join(names, ", "))
	}
	return remedies
}
