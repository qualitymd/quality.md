package evaluator

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// cliEvaluator invokes a coding-agent CLI non-interactively as a subprocess
// transport for one work request at a time. The CLI runs in its
// machine-readable print/exec mode; an installation that cannot honor
// non-interactive structured invocation fails selection with
// evaluator_incompatible rather than degrading into unparseable runs.
//
// Where the installed CLI advertises native JSON Schema output enforcement or
// no-persistence controls in its help output, the adapter uses them for
// bounded work requests; the runner still validates every returned payload
// independently, so an absent flag only loses the native enforcement.
type cliEvaluator struct {
	kind    string
	command string
	// runCommand overrides subprocess execution in tests.
	runCommand func(ctx context.Context, name string, args []string, stdin string) (string, string, error)

	// capabilities detected from the installed CLI's help output before the
	// first judgment call.
	probeOnce     sync.Once
	schemaFlag    string
	noPersistFlag string
}

// cliCapabilityFlags names the native structured-output and no-persistence
// flags each CLI kind may advertise. A flag is used only when the installed
// CLI's help output contains it verbatim.
var cliCapabilityFlags = map[string]struct{ schema, noPersist string }{
	"claude": {schema: "--json-schema", noPersist: "--no-session-persistence"},
	"codex":  {schema: "--output-schema", noPersist: "--ephemeral"},
}

// detectCapabilities inspects the installed CLI's help output once, before
// the first judgment call, so a version or capability mismatch never degrades
// mid-run.
func (e *cliEvaluator) detectCapabilities(ctx context.Context) {
	flags, ok := cliCapabilityFlags[e.kind]
	if !ok {
		return
	}
	helpArgs := []string{"--help"}
	if e.kind == "codex" {
		helpArgs = []string{"exec", "--help"}
	}
	help, stderr, err := e.run(ctx, helpArgs, "")
	if err != nil {
		help += stderr
	}
	if strings.Contains(help, flags.schema) {
		e.schemaFlag = flags.schema
	}
	if strings.Contains(help, flags.noPersist) {
		e.noPersistFlag = flags.noPersist
	}
}

var _ Evaluator = (*cliEvaluator)(nil)

func cliCommandName(kind, override string) string {
	if override != "" {
		return override
	}
	return kind
}

func newCLIEvaluator(kind, command string, opts Options) (*cliEvaluator, error) {
	name := cliCommandName(kind, command)
	if _, err := opts.lookPath(name); err != nil {
		return nil, &SelectionError{
			Category: FailureMissingEvaluator,
			Message:  fmt.Sprintf("the %s CLI (%q) is not installed or not on PATH", kind, name),
			Remedies: []string{"install the " + kind + " CLI", "or select another evaluator with --evaluator"},
		}
	}
	return &cliEvaluator{kind: kind, command: name}, nil
}

func (e *cliEvaluator) Kind() string { return e.kind }

func (e *cliEvaluator) Capabilities() Capabilities {
	return Capabilities{
		Strategies:   []Strategy{StrategySequential},
		Subagents:    false,
		ReportsUsage: e.kind == "claude",
	}
}

func (e *cliEvaluator) Evaluate(ctx context.Context, req WorkRequest) (WorkResult, error) {
	system, stablePrefix, delta := BuildPrompt(req)
	prompt := system + "\n\n" + stablePrefix + delta
	result := WorkResult{WorkUnitID: req.WorkUnitID, EvaluatorKind: e.kind, Strategy: StrategySequential}

	e.probeOnce.Do(func() { e.detectCapabilities(ctx) })

	var args []string
	switch e.kind {
	case "claude":
		args = []string{"-p", "--output-format", "json"}
	case "codex":
		args = []string{"exec", "--json"}
	default:
		result.Failure = FailureInternal
		result.FailureDetail = fmt.Sprintf("unsupported CLI evaluator kind %q", e.kind)
		return result, nil
	}
	if e.noPersistFlag != "" {
		args = append(args, e.noPersistFlag)
	}
	if e.schemaFlag != "" {
		schemaFile, cleanup, err := writeSchemaFile(req.ExpectedSchema)
		if err != nil {
			result.Failure = FailureInternal
			result.FailureDetail = "staging expected schema: " + err.Error()
			return result, nil
		}
		defer cleanup()
		args = append(args, e.schemaFlag, schemaFile)
	}
	if e.kind == "codex" {
		args = append(args, "-")
	}

	stdout, stderr, err := e.run(ctx, args, prompt)
	if err != nil {
		result.Failure, result.FailureDetail = classifyCLIError(ctx, err, stderr)
		return result, nil
	}

	var text string
	switch e.kind {
	case "claude":
		text = e.parseClaudeOutput(stdout, &result)
	case "codex":
		text = e.parseCodexOutput(stdout, &result)
	}
	if strings.TrimSpace(text) == "" {
		result.Failure = FailureInvalidEvaluatorOutput
		result.FailureDetail = "the " + e.kind + " CLI returned no structured result text"
		return result, nil
	}
	payload, err := ExtractJSONObject(text)
	if err != nil {
		result.Failure = FailureInvalidEvaluatorOutput
		result.FailureDetail = err.Error()
		return result, nil
	}
	result.Payload = payload
	return result, nil
}

// writeSchemaFile stages the expected result schema as a temporary file for a
// CLI's native schema-enforcement flag, removed after the call.
func writeSchemaFile(schema []byte) (string, func(), error) {
	file, err := os.CreateTemp("", "qualitymd-schema-*.json")
	if err != nil {
		return "", nil, err
	}
	name := file.Name()
	cleanup := func() { _ = os.Remove(name) }
	if _, err := file.Write(schema); err != nil {
		_ = file.Close()
		cleanup()
		return "", nil, err
	}
	if err := file.Close(); err != nil {
		cleanup()
		return "", nil, err
	}
	return name, cleanup, nil
}

func (e *cliEvaluator) run(ctx context.Context, args []string, stdin string) (string, string, error) {
	if e.runCommand != nil {
		return e.runCommand(ctx, e.command, args, stdin)
	}
	cmd := exec.CommandContext(ctx, e.command, args...)
	cmd.Stdin = strings.NewReader(stdin)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// parseClaudeOutput reads Claude Code print-mode JSON output: one JSON object
// with a "result" text field plus usage and session metadata.
func (e *cliEvaluator) parseClaudeOutput(stdout string, result *WorkResult) string {
	var doc struct {
		Result       string   `json:"result"`
		IsError      bool     `json:"is_error"`
		SessionID    string   `json:"session_id"`
		TotalCostUSD *float64 `json:"total_cost_usd"`
		Usage        *struct {
			InputTokens              *int64 `json:"input_tokens"`
			OutputTokens             *int64 `json:"output_tokens"`
			CacheCreationInputTokens *int64 `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     *int64 `json:"cache_read_input_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal([]byte(stdout), &doc); err != nil {
		// Fall back to treating stdout as the result text itself.
		return stdout
	}
	if doc.SessionID != "" {
		result.ContextMeta = map[string]string{"sessionId": doc.SessionID}
	}
	switch {
	case doc.Usage != nil:
		result.Usage = anthropicUsage(doc.Usage.InputTokens, doc.Usage.OutputTokens,
			doc.Usage.CacheCreationInputTokens, doc.Usage.CacheReadInputTokens)
		result.Usage.CostUSD = doc.TotalCostUSD
	case doc.TotalCostUSD != nil:
		result.Usage = &Usage{CostUSD: doc.TotalCostUSD}
	}
	return doc.Result
}

// parseCodexOutput reads Codex CLI exec-mode JSONL output and returns the
// last completed agent message text.
func (e *cliEvaluator) parseCodexOutput(stdout string, result *WorkResult) string {
	var text string
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	scanner.Buffer(make([]byte, 0, 1024*1024), 16*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var event struct {
			Type string `json:"type"`
			Item *struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"item"`
			ThreadID string `json:"thread_id"`
			Usage    *struct {
				InputTokens       *int64 `json:"input_tokens"`
				CachedInputTokens *int64 `json:"cached_input_tokens"`
				OutputTokens      *int64 `json:"output_tokens"`
			} `json:"usage"`
		}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if event.ThreadID != "" {
			result.ContextMeta = map[string]string{"threadId": event.ThreadID}
		}
		if event.Usage != nil {
			result.Usage = &Usage{
				InputTokens:       event.Usage.InputTokens,
				CachedInputTokens: event.Usage.CachedInputTokens,
				OutputTokens:      event.Usage.OutputTokens,
			}
		}
		if event.Item != nil && event.Item.Type == "agent_message" && event.Item.Text != "" {
			text = event.Item.Text
		}
	}
	return text
}

func classifyCLIError(ctx context.Context, err error, stderr string) (FailureCategory, string) {
	if ctx.Err() != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return FailureTimeout, "evaluator call timed out"
		}
		return FailureCancelled, "evaluator call was cancelled"
	}
	lower := strings.ToLower(stderr)
	switch {
	case strings.Contains(lower, "not logged in"),
		strings.Contains(lower, "login"),
		strings.Contains(lower, "unauthorized"),
		strings.Contains(lower, "authentication"):
		return FailureEvaluatorUnauthenticated, firstLineOr(stderr, err.Error())
	case strings.Contains(lower, "rate limit"), strings.Contains(lower, "429"):
		return FailureRateLimited, firstLineOr(stderr, err.Error())
	case strings.Contains(lower, "unknown flag"),
		strings.Contains(lower, "unknown option"),
		strings.Contains(lower, "unknown command"):
		return FailureEvaluatorIncompatible, firstLineOr(stderr, err.Error())
	default:
		return FailureInvalidEvaluatorOutput, firstLineOr(stderr, err.Error())
	}
}

func firstLineOr(text, fallback string) string {
	for _, line := range strings.Split(text, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			return trimmed
		}
	}
	return fallback
}
