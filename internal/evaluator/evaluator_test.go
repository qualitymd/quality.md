package evaluator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/workspace"
)

func TestExtractJSONObject(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantKey string
		wantErr bool
	}{
		{name: "bare object", text: `{"kind":"x"}`, wantKey: "kind"},
		{name: "fenced object", text: "```json\n{\"kind\":\"x\"}\n```", wantKey: "kind"},
		{name: "prose around object", text: "Here you go:\n{\"kind\":\"x\"}\nDone.", wantKey: "kind"},
		{name: "no object", text: "no json here", wantErr: true},
		{name: "truncated object", text: `{"kind":`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := ExtractJSONObject(tt.text)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ExtractJSONObject() = %v, want error", payload)
				}
				return
			}
			if err != nil {
				t.Fatalf("ExtractJSONObject() error = %v", err)
			}
			if _, ok := payload[tt.wantKey]; !ok {
				t.Fatalf("payload = %v, want key %q", payload, tt.wantKey)
			}
		})
	}
}

func TestBuildPromptSplitsStablePrefixFromDelta(t *testing.T) {
	req := WorkRequest{
		WorkUnitID:     "assessRateRequirement:requirement:root::has-tests",
		Kind:           "assessRateRequirement",
		Subject:        "requirement:root::has-tests",
		Instructions:   "Assess and rate it.",
		ExpectedSchema: []byte(`{"type":"object"}`),
		SharedContext:  map[string]any{"areaEvaluationFrame": map[string]any{"a": 1}},
		Context:        map[string]any{"requirementEvaluationFrame": map[string]any{"b": 2}},
		Source:         []SourceFile{{Path: "src/main.txt", Content: "hello", SHA256: "abc"}},
	}
	system, stablePrefix, delta := BuildPrompt(req)
	if !strings.Contains(system, "DATA under") {
		t.Errorf("system prompt lacks safety instructions")
	}
	// Everything stable across an area's work units renders in the stable
	// prefix, ordered task, schema, source, shared context.
	taskAt := strings.Index(stablePrefix, "## Task")
	schemaAt := strings.Index(stablePrefix, "## Expected result schema")
	sourceAt := strings.Index(stablePrefix, "## Source files")
	sharedAt := strings.Index(stablePrefix, "## Shared context")
	if taskAt < 0 || schemaAt < 0 || sourceAt < 0 || sharedAt < 0 ||
		taskAt >= schemaAt || schemaAt >= sourceAt || sourceAt >= sharedAt {
		t.Errorf("stable prefix layers out of order: task=%d schema=%d source=%d shared=%d",
			taskAt, schemaAt, sourceAt, sharedAt)
	}
	for _, varying := range []string{req.WorkUnitID, "requirementEvaluationFrame"} {
		if strings.Contains(stablePrefix, varying) {
			t.Errorf("stable prefix contains per-work-unit content %q", varying)
		}
	}
	if !strings.Contains(delta, req.WorkUnitID) || !strings.Contains(delta, "requirementEvaluationFrame") {
		t.Errorf("delta lacks the work-unit header or per-unit context: %q", delta)
	}
	if !strings.HasSuffix(delta, "Return ONLY the JSON object now.") {
		t.Errorf("delta must end with the JSON-only closing line")
	}

	// The stable prefix is byte-identical for a sibling work unit in the same
	// area — that identity is what provider prefix caching keys on.
	sibling := req
	sibling.WorkUnitID = "assessRateRequirement:requirement:root::has-docs"
	sibling.Subject = "requirement:root::has-docs"
	sibling.Context = map[string]any{"requirementEvaluationFrame": map[string]any{"b": 3}}
	_, siblingPrefix, _ := BuildPrompt(sibling)
	if siblingPrefix != stablePrefix {
		t.Errorf("stable prefix differs across the area's work units")
	}
}

func TestValidateProfilesRejectsReservedShadow(t *testing.T) {
	err := ValidateProfiles(map[string]workspace.EvaluatorProfile{
		"claude": {Kind: "anthropic"},
	})
	var selErr *SelectionError
	if err == nil || !errors.As(err, &selErr) {
		t.Fatalf("ValidateProfiles() error = %v, want SelectionError", err)
	}
}

func TestValidateProfilesRequiresKind(t *testing.T) {
	err := ValidateProfiles(map[string]workspace.EvaluatorProfile{
		"mine": {},
	})
	if err == nil {
		t.Fatal("ValidateProfiles() = nil, want missing-kind error")
	}
}

// fakeProbe simulates readiness probe subprocess output per command.
func fakeProbe(responses map[string]struct {
	out string
	err error
},
) func(string, ...string) (string, error) {
	return func(name string, args ...string) (string, error) {
		key := name + " " + strings.Join(args, " ")
		response, ok := responses[key]
		if !ok {
			return "", errors.New("unexpected probe: " + key)
		}
		return response.out, response.err
	}
}

func TestSelectAutoPrefersInstalledCLI(t *testing.T) {
	opts := Options{
		LookPath: func(file string) (string, error) {
			if file == "claude" {
				return "/usr/local/bin/claude", nil
			}
			return "", errors.New("not found")
		},
		Getenv: func(string) string { return "" },
		RunProbe: fakeProbe(map[string]struct {
			out string
			err error
		}{
			"claude --help": {out: "--print --output-format json"},
		}),
	}
	selection, err := Select(opts)
	if err != nil {
		t.Fatalf("Select(auto) error = %v", err)
	}
	if selection.Name != "claude" {
		t.Fatalf("selection = %q, want claude", selection.Name)
	}
	if len(selection.Candidates) != 2 || selection.Candidates[0].Name != "codex" || selection.Candidates[0].Usable {
		t.Fatalf("candidates = %+v, want unusable codex evidence then claude", selection.Candidates)
	}
}

func TestSelectAutoSkipsUnauthenticatedCodex(t *testing.T) {
	opts := Options{
		LookPath: func(string) (string, error) { return "/usr/local/bin/cli", nil },
		Getenv:   func(string) string { return "" },
		RunProbe: fakeProbe(map[string]struct {
			out string
			err error
		}{
			"codex exec --help":  {out: "--json --output-schema"},
			"codex login status": {out: "Not logged in", err: errors.New("exit status 1")},
			"claude --help":      {out: "--print --output-format json"},
		}),
	}
	selection, err := Select(opts)
	if err != nil {
		t.Fatalf("Select(auto) error = %v", err)
	}
	if selection.Name != "claude" {
		t.Fatalf("selection = %q, want claude after skipping unauthenticated codex", selection.Name)
	}
	codex := selection.Candidates[0]
	if codex.Usable || codex.Authenticated == nil || *codex.Authenticated {
		t.Fatalf("codex readiness = %+v, want unauthenticated and unusable", codex)
	}
}

func TestSelectAutoSkipsCLIWithoutStructuredOutput(t *testing.T) {
	opts := Options{
		LookPath: func(file string) (string, error) {
			if file == "claude" {
				return "/usr/local/bin/claude", nil
			}
			return "", errors.New("not found")
		},
		Getenv: func(string) string { return "" },
		RunProbe: fakeProbe(map[string]struct {
			out string
			err error
		}{
			"claude --help": {out: "an old interactive-only build"},
		}),
	}
	_, err := Select(opts)
	var selErr *SelectionError
	if err == nil || !errors.As(err, &selErr) || selErr.Category != FailureMissingEvaluator {
		t.Fatalf("Select(auto) error = %v, want missing_evaluator after capability probe", err)
	}
	if !strings.Contains(selErr.Message, "structured-output") {
		t.Fatalf("message = %q, want structured-output readiness evidence", selErr.Message)
	}
}

func TestSelectHarnessEvaluator(t *testing.T) {
	selection, err := Select(Options{Name: "harness"})
	if err != nil {
		t.Fatalf("Select(harness) error = %v", err)
	}
	if selection.Name != "harness" || selection.Evaluator.Kind() != "harness" {
		t.Fatalf("selection = %q/%q, want harness", selection.Name, selection.Evaluator.Kind())
	}
}

func TestSelectAutoFallsBackToConfiguredAPIProfile(t *testing.T) {
	opts := Options{
		Profiles: map[string]workspace.EvaluatorProfile{
			"team": {Kind: "anthropic", APIKeyEnv: "TEAM_KEY"},
		},
		LookPath: func(string) (string, error) { return "", errors.New("not found") },
		Getenv: func(key string) string {
			if key == "TEAM_KEY" {
				return "secret"
			}
			return ""
		},
	}
	selection, err := Select(opts)
	if err != nil {
		t.Fatalf("Select(auto) error = %v", err)
	}
	if selection.Name != "team" || selection.Evaluator.Kind() != "anthropic" {
		t.Fatalf("selection = %q/%q, want team/anthropic", selection.Name, selection.Evaluator.Kind())
	}
}

func TestSelectAutoFailsWithRemedies(t *testing.T) {
	opts := Options{
		LookPath: func(string) (string, error) { return "", errors.New("not found") },
		Getenv:   func(string) string { return "" },
	}
	_, err := Select(opts)
	var selErr *SelectionError
	if err == nil || !errors.As(err, &selErr) {
		t.Fatalf("Select(auto) error = %v, want SelectionError", err)
	}
	if selErr.Category != FailureMissingEvaluator || len(selErr.Remedies) == 0 {
		t.Fatalf("selection error = %+v, want missing_evaluator with remedies", selErr)
	}
}

func TestSelectReservedUnimplementedNames(t *testing.T) {
	for _, name := range []string{"shell", "manual"} {
		opts := Options{
			Name:     name,
			LookPath: func(string) (string, error) { return "", errors.New("not found") },
			Getenv:   func(string) string { return "" },
		}
		_, err := Select(opts)
		var selErr *SelectionError
		if err == nil || !errors.As(err, &selErr) || selErr.Category != FailureMissingEvaluator {
			t.Fatalf("Select(%s) error = %v, want reserved-name failure", name, err)
		}
	}
}

func TestSelectMissingAPIKey(t *testing.T) {
	opts := Options{
		Name:     "anthropic",
		LookPath: func(string) (string, error) { return "", errors.New("not found") },
		Getenv:   func(string) string { return "" },
	}
	_, err := Select(opts)
	var selErr *SelectionError
	if err == nil || !errors.As(err, &selErr) || selErr.Category != FailureMissingAPIKey {
		t.Fatalf("Select(anthropic) error = %v, want missing_api_key", err)
	}
}

func TestCLIEvaluatorParsesClaudeOutput(t *testing.T) {
	ev := &cliEvaluator{kind: "claude", command: "claude"}
	ev.runCommand = func(_ context.Context, _ string, args []string, _ string) (string, string, error) {
		if len(args) == 0 || args[0] != "-p" {
			return "", "", fmt.Errorf("unexpected args %v", args)
		}
		return `{"result":"{\"kind\":\"RequirementRatingResult\"}","session_id":"s1",` +
			`"usage":{"input_tokens":10,"output_tokens":5,"cache_read_input_tokens":90}}`, "", nil
	}
	result, err := ev.Evaluate(context.Background(), WorkRequest{WorkUnitID: "u1", ExpectedSchema: []byte("{}")})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Failure != "" {
		t.Fatalf("failure = %s: %s", result.Failure, result.FailureDetail)
	}
	if result.Payload["kind"] != "RequirementRatingResult" {
		t.Fatalf("payload = %v", result.Payload)
	}
	// input_tokens excludes cache reads, so the total is 10+90.
	if result.Usage == nil || result.Usage.InputTokens == nil || *result.Usage.InputTokens != 100 {
		t.Fatalf("usage = %+v, want input tokens 100", result.Usage)
	}
	if result.Usage.CachedInputTokens == nil || *result.Usage.CachedInputTokens != 90 {
		t.Fatalf("usage = %+v, want cached input tokens 90", result.Usage)
	}
	if result.ContextMeta["sessionId"] != "s1" {
		t.Fatalf("contextMeta = %v", result.ContextMeta)
	}
}

func TestCLIEvaluatorParsesCodexOutput(t *testing.T) {
	ev := &cliEvaluator{kind: "codex", command: "codex"}
	ev.runCommand = func(_ context.Context, _ string, args []string, _ string) (string, string, error) {
		if len(args) == 0 || args[0] != "exec" {
			return "", "", fmt.Errorf("unexpected args %v", args)
		}
		lines := []string{
			`{"type":"thread.started","thread_id":"t1"}`,
			`{"type":"item.completed","item":{"type":"agent_message","text":"{\"kind\":\"RequirementRatingResult\"}"}}`,
			`{"type":"turn.completed","usage":{"input_tokens":7,"output_tokens":3}}`,
		}
		return strings.Join(lines, "\n"), "", nil
	}
	result, err := ev.Evaluate(context.Background(), WorkRequest{WorkUnitID: "u1", ExpectedSchema: []byte("{}")})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Failure != "" {
		t.Fatalf("failure = %s: %s", result.Failure, result.FailureDetail)
	}
	if result.Payload["kind"] != "RequirementRatingResult" {
		t.Fatalf("payload = %v", result.Payload)
	}
	if result.ContextMeta["threadId"] != "t1" {
		t.Fatalf("contextMeta = %v", result.ContextMeta)
	}
}

func TestAnthropicEvaluatorCachesStablePrefix(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		_, _ = w.Write([]byte(`{
			"content":[{"type":"text","text":"{\"kind\":\"RequirementRatingResult\"}"}],
			"usage":{"input_tokens":3,"output_tokens":5,"cache_creation_input_tokens":7,"cache_read_input_tokens":90}
		}`))
	}))
	defer server.Close()
	ev := &anthropicEvaluator{apiEvaluator: apiEvaluator{
		kind:    "anthropic",
		model:   "test-model",
		baseURL: server.URL,
		env:     "TEST_KEY",
		getenv: func(key string) string {
			if key == "TEST_KEY" {
				return "secret"
			}
			return ""
		},
		client: server.Client(),
	}}
	result, err := ev.Evaluate(context.Background(), WorkRequest{
		WorkUnitID:     "u1",
		Kind:           "assessRateRequirement",
		Instructions:   "Assess and rate it.",
		ExpectedSchema: []byte(`{"type":"object"}`),
		Context:        map[string]any{"frame": map[string]any{"a": 1}},
	})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Failure != "" {
		t.Fatalf("failure = %s: %s", result.Failure, result.FailureDetail)
	}
	// The system block and the stable prefix carry cache_control; the trailing
	// delta block does not.
	systemBlocks, _ := body["system"].([]any)
	if len(systemBlocks) != 1 || !hasCacheControl(systemBlocks[0]) {
		t.Errorf("system = %v, want one cache_control text block", body["system"])
	}
	messages, _ := body["messages"].([]any)
	if len(messages) != 1 {
		t.Fatalf("messages = %v, want one user message", body["messages"])
	}
	message, _ := messages[0].(map[string]any)
	content, _ := message["content"].([]any)
	if len(content) != 2 || !hasCacheControl(content[0]) || hasCacheControl(content[1]) {
		t.Errorf("content = %v, want cached stable prefix then uncached delta", content)
	}
	// Anthropic input_tokens excludes cache reads and writes: total 3+7+90.
	if result.Usage == nil || result.Usage.InputTokens == nil || *result.Usage.InputTokens != 100 {
		t.Fatalf("usage = %+v, want input tokens 100", result.Usage)
	}
	if result.Usage.CachedInputTokens == nil || *result.Usage.CachedInputTokens != 90 {
		t.Fatalf("usage = %+v, want cached input tokens 90", result.Usage)
	}
}

func hasCacheControl(block any) bool {
	object, _ := block.(map[string]any)
	_, ok := object["cache_control"].(map[string]any)
	return ok
}

func TestCLIEvaluatorUsesNativeSchemaAndNoPersistFlags(t *testing.T) {
	ev := &cliEvaluator{kind: "codex", command: "codex"}
	var judgeArgs []string
	ev.runCommand = func(_ context.Context, _ string, args []string, _ string) (string, string, error) {
		if len(args) >= 2 && args[0] == "exec" && args[1] == "--help" {
			return "usage: codex exec --json --output-schema <file> --ephemeral", "", nil
		}
		judgeArgs = args
		return `{"type":"item.completed","item":{"type":"agent_message","text":"{\"kind\":\"x\"}"}}`, "", nil
	}
	result, err := ev.Evaluate(context.Background(), WorkRequest{WorkUnitID: "u1", ExpectedSchema: []byte(`{"type":"object"}`)})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Failure != "" {
		t.Fatalf("failure = %s: %s", result.Failure, result.FailureDetail)
	}
	joined := strings.Join(judgeArgs, " ")
	if !strings.Contains(joined, "--ephemeral") || !strings.Contains(joined, "--output-schema") {
		t.Fatalf("judge args = %v, want native --ephemeral and --output-schema flags", judgeArgs)
	}
	if judgeArgs[len(judgeArgs)-1] != "-" {
		t.Fatalf("judge args = %v, want trailing stdin marker", judgeArgs)
	}
}

func TestCLIEvaluatorSkipsUnsupportedNativeFlags(t *testing.T) {
	ev := &cliEvaluator{kind: "claude", command: "claude"}
	var judgeArgs []string
	ev.runCommand = func(_ context.Context, _ string, args []string, _ string) (string, string, error) {
		if len(args) >= 1 && args[0] == "--help" {
			return "usage: claude -p --output-format json", "", nil
		}
		judgeArgs = args
		return `{"result":"{\"kind\":\"x\"}"}`, "", nil
	}
	if _, err := ev.Evaluate(context.Background(), WorkRequest{WorkUnitID: "u1", ExpectedSchema: []byte(`{}`)}); err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	joined := strings.Join(judgeArgs, " ")
	if strings.Contains(joined, "--json-schema") || strings.Contains(joined, "--no-session-persistence") {
		t.Fatalf("judge args = %v, must not pass flags the installed CLI does not advertise", judgeArgs)
	}
}

func TestCLIEvaluatorClassifiesAuthFailure(t *testing.T) {
	ev := &cliEvaluator{kind: "claude", command: "claude"}
	ev.runCommand = func(context.Context, string, []string, string) (string, string, error) {
		return "", "Error: not logged in. Run claude login first.", errors.New("exit status 1")
	}
	result, err := ev.Evaluate(context.Background(), WorkRequest{WorkUnitID: "u1", ExpectedSchema: []byte("{}")})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Failure != FailureEvaluatorUnauthenticated {
		t.Fatalf("failure = %s, want evaluator_unauthenticated", result.Failure)
	}
}
