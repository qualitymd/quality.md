package evaluator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/qualitymd/quality.md/internal/workspace"
)

const apiCallTimeout = 10 * time.Minute

// apiEvaluator is the shared base for direct provider API adapters. The API
// key is read from the configured environment variable at call time and never
// persisted or logged.
type apiEvaluator struct {
	kind    string
	model   string
	baseURL string
	env     string
	getenv  func(string) string
	client  *http.Client
}

func (e *apiEvaluator) Kind() string { return e.kind }

func (e *apiEvaluator) Capabilities() Capabilities {
	return Capabilities{
		Concurrent:      true,
		ReusableContext: []string{"prompt-cache"},
		ReportsUsage:    true,
	}
}

func (e *apiEvaluator) apiKey() (string, error) {
	key := e.getenv(e.env)
	if key == "" {
		return "", &SelectionError{
			Category: FailureMissingAPIKey,
			Message:  fmt.Sprintf("the %s evaluator needs an API key in $%s", e.kind, e.env),
			Remedies: []string{"export " + e.env, "or select another evaluator with --evaluator"},
		}
	}
	return key, nil
}

func newAPIEvaluator(kind string, profile workspace.EvaluatorProfile, defaultModel, defaultBase, defaultEnv string, opts Options) (*apiEvaluator, error) {
	e := &apiEvaluator{
		kind:    kind,
		model:   profile.Model,
		baseURL: profile.BaseURL,
		env:     profileEnv(profile, defaultEnv),
		getenv:  opts.getenv,
		client:  &http.Client{Timeout: apiCallTimeout},
	}
	if e.model == "" {
		e.model = defaultModel
	}
	if e.baseURL == "" {
		e.baseURL = defaultBase
	}
	if _, err := e.apiKey(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *apiEvaluator) post(ctx context.Context, url string, headers map[string]string, body any) ([]byte, int, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close() //nolint:errcheck // Read errors surface via io.ReadAll.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return data, resp.StatusCode, nil
}

func classifyAPIFailure(ctx context.Context, status int, err error) (FailureCategory, string) {
	if ctx.Err() != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return FailureTimeout, "evaluator API call timed out"
		}
		return FailureCancelled, "evaluator API call was cancelled"
	}
	if err != nil {
		return FailureInvalidEvaluatorOutput, err.Error()
	}
	switch {
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return FailureEvaluatorUnauthenticated, fmt.Sprintf("provider rejected the API key (HTTP %d)", status)
	case status == http.StatusTooManyRequests:
		return FailureRateLimited, "provider rate limit reached (HTTP 429)"
	case status >= 500:
		return FailureInvalidEvaluatorOutput, fmt.Sprintf("provider error (HTTP %d)", status)
	default:
		return FailureInvalidEvaluatorOutput, fmt.Sprintf("unexpected provider response (HTTP %d)", status)
	}
}

// anthropicEvaluator calls the Anthropic Messages API directly.
type anthropicEvaluator struct{ apiEvaluator }

var _ Evaluator = (*anthropicEvaluator)(nil)

func newAnthropicEvaluator(profile workspace.EvaluatorProfile, opts Options) (*anthropicEvaluator, error) {
	base, err := newAPIEvaluator("anthropic", profile, defaultAnthropicModel, "https://api.anthropic.com", defaultAnthropicEnv, opts)
	if err != nil {
		return nil, err
	}
	return &anthropicEvaluator{apiEvaluator: *base}, nil
}

func (e *anthropicEvaluator) Evaluate(ctx context.Context, req WorkRequest) (WorkResult, error) {
	result := WorkResult{WorkUnitID: req.WorkUnitID, EvaluatorKind: e.kind, Model: e.model}
	key, err := e.apiKey()
	if err != nil {
		var selErr *SelectionError
		if errors.As(err, &selErr) {
			result.Failure = selErr.Category
			result.FailureDetail = selErr.Message
			return result, nil
		}
		return result, err
	}
	system, stablePrefix, delta := BuildPrompt(req)
	// The stable prefix repeats verbatim across an area's work units, so it is
	// marked as a provider cache breakpoint; a cache miss only costs tokens.
	cacheControl := map[string]any{"type": "ephemeral"}
	body := map[string]any{
		"model":      e.model,
		"max_tokens": 16000,
		"system": []map[string]any{
			{"type": "text", "text": system, "cache_control": cacheControl},
		},
		"messages": []map[string]any{
			{"role": "user", "content": []map[string]any{
				{"type": "text", "text": stablePrefix, "cache_control": cacheControl},
				{"type": "text", "text": delta},
			}},
		},
	}
	headers := map[string]string{
		"x-api-key":         key,
		"anthropic-version": "2023-06-01",
	}
	data, status, err := e.post(ctx, e.baseURL+"/v1/messages", headers, body)
	if err != nil || status != http.StatusOK {
		result.Failure, result.FailureDetail = classifyAPIFailure(ctx, status, err)
		return result, nil
	}
	var doc struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage *struct {
			InputTokens              *int64 `json:"input_tokens"`
			OutputTokens             *int64 `json:"output_tokens"`
			CacheCreationInputTokens *int64 `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     *int64 `json:"cache_read_input_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		result.Failure = FailureInvalidEvaluatorOutput
		result.FailureDetail = "parsing provider response: " + err.Error()
		return result, nil
	}
	if doc.Usage != nil {
		result.Usage = anthropicUsage(doc.Usage.InputTokens, doc.Usage.OutputTokens,
			doc.Usage.CacheCreationInputTokens, doc.Usage.CacheReadInputTokens)
	}
	var text string
	for _, block := range doc.Content {
		if block.Type == "text" {
			text += block.Text
		}
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

// openAIEvaluator calls the OpenAI Chat Completions API directly.
type openAIEvaluator struct{ apiEvaluator }

var _ Evaluator = (*openAIEvaluator)(nil)

func newOpenAIEvaluator(profile workspace.EvaluatorProfile, opts Options) (*openAIEvaluator, error) {
	base, err := newAPIEvaluator("openai", profile, defaultOpenAIModel, "https://api.openai.com", defaultOpenAIEnv, opts)
	if err != nil {
		return nil, err
	}
	return &openAIEvaluator{apiEvaluator: *base}, nil
}

func (e *openAIEvaluator) Evaluate(ctx context.Context, req WorkRequest) (WorkResult, error) {
	result := WorkResult{WorkUnitID: req.WorkUnitID, EvaluatorKind: e.kind, Model: e.model}
	key, err := e.apiKey()
	if err != nil {
		var selErr *SelectionError
		if errors.As(err, &selErr) {
			result.Failure = selErr.Category
			result.FailureDetail = selErr.Message
			return result, nil
		}
		return result, err
	}
	// OpenAI prefix caching is automatic; the stable-first layout is enough.
	system, stablePrefix, delta := BuildPrompt(req)
	body := map[string]any{
		"model": e.model,
		"messages": []map[string]any{
			{"role": "system", "content": system},
			{"role": "user", "content": stablePrefix + delta},
		},
		"response_format": map[string]any{"type": "json_object"},
	}
	headers := map[string]string{"Authorization": "Bearer " + key}
	data, status, err := e.post(ctx, e.baseURL+"/v1/chat/completions", headers, body)
	if err != nil || status != http.StatusOK {
		result.Failure, result.FailureDetail = classifyAPIFailure(ctx, status, err)
		return result, nil
	}
	var doc struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage *struct {
			PromptTokens        *int64 `json:"prompt_tokens"`
			CompletionTokens    *int64 `json:"completion_tokens"`
			PromptTokensDetails *struct {
				CachedTokens *int64 `json:"cached_tokens"`
			} `json:"prompt_tokens_details"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		result.Failure = FailureInvalidEvaluatorOutput
		result.FailureDetail = "parsing provider response: " + err.Error()
		return result, nil
	}
	if doc.Usage != nil {
		result.Usage = &Usage{InputTokens: doc.Usage.PromptTokens, OutputTokens: doc.Usage.CompletionTokens}
		if doc.Usage.PromptTokensDetails != nil {
			result.Usage.CachedInputTokens = doc.Usage.PromptTokensDetails.CachedTokens
		}
	}
	if len(doc.Choices) == 0 {
		result.Failure = FailureInvalidEvaluatorOutput
		result.FailureDetail = "provider returned no choices"
		return result, nil
	}
	payload, err := ExtractJSONObject(doc.Choices[0].Message.Content)
	if err != nil {
		result.Failure = FailureInvalidEvaluatorOutput
		result.FailureDetail = err.Error()
		return result, nil
	}
	result.Payload = payload
	return result, nil
}
