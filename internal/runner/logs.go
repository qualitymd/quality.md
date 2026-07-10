package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Run-local structured logs. Execution telemetry stays out of the
// authoritative evaluation.json: events.jsonl records runner lifecycle
// events; evaluator-calls.jsonl records evaluator-call metadata. Neither
// records raw prompts, raw source bundles, raw model responses, API keys,
// auth tokens, or environment values.

const (
	logsDir            = "logs"
	eventsLogFile      = "events.jsonl"
	evaluatorCallsFile = "evaluator-calls.jsonl"
)

type runLogs struct {
	events *os.File
	calls  *os.File
	now    func() time.Time
}

func openRunLogs(runAbs string) (*runLogs, error) {
	dir := filepath.Join(runAbs, logsDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating %s: %w", logsDir, err)
	}
	events, err := os.OpenFile(filepath.Join(dir, eventsLogFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", eventsLogFile, err)
	}
	calls, err := os.OpenFile(filepath.Join(dir, evaluatorCallsFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		_ = events.Close()
		return nil, fmt.Errorf("opening %s: %w", evaluatorCallsFile, err)
	}
	return &runLogs{events: events, calls: calls, now: time.Now}, nil
}

func (l *runLogs) Close() {
	if l == nil {
		return
	}
	_ = l.events.Close()
	_ = l.calls.Close()
}

func (l *runLogs) timestamp() string {
	return l.now().UTC().Format(time.RFC3339)
}

// event appends one runner lifecycle event.
func (l *runLogs) event(name string, fields map[string]any) {
	if l == nil {
		return
	}
	entry := map[string]any{"ts": l.timestamp(), "event": name}
	for key, value := range fields {
		entry[key] = value
	}
	l.append(l.events, entry)
}

// call appends one evaluator-call metadata entry.
func (l *runLogs) call(fields map[string]any) {
	if l == nil {
		return
	}
	entry := map[string]any{"ts": l.timestamp()}
	for key, value := range fields {
		entry[key] = value
	}
	l.append(l.calls, entry)
}

func (l *runLogs) append(file *os.File, entry map[string]any) {
	raw, err := json.Marshal(entry)
	if err != nil {
		return
	}
	raw = append(raw, '\n')
	_, _ = file.Write(raw)
}
