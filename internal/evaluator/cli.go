package evaluator

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// cliEvaluator invokes a coding-agent CLI non-interactively as a subprocess
// transport for one work request at a time. The CLI runs in its
// machine-readable print/exec mode; an installation that cannot honor
// non-interactive structured invocation fails selection with
// evaluator_incompatible rather than degrading into unparseable runs.
type cliEvaluator struct {
	kind    string
	command string
	// runCommand overrides subprocess execution in tests.
	runCommand func(ctx context.Context, name string, args []string, stdin string) (string, string, error)
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

	var args []string
	switch e.kind {
	case "claude":
		args = []string{"-p", "--output-format", "json"}
	case "codex":
		args = []string{"exec", "--json", "-"}
	default:
		result.Failure = FailureInternal
		result.FailureDetail = fmt.Sprintf("unsupported CLI evaluator kind %q", e.kind)
		return result, nil
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
