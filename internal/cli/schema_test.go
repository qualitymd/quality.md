package cli

import (
	"bytes"
	"strings"
	"testing"

	qualitymd "github.com/qualitymd/quality.md"
)

func TestSchemaWritesVerbatimSchema(t *testing.T) {
	want := qualitymd.Schema()

	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"schema"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !bytes.Equal(out.Bytes(), want) {
		t.Fatal("stdout did not match the bundled schema byte-for-byte")
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestSchemaOutputIsValidJSON(t *testing.T) {
	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"schema"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !bytes.HasPrefix(bytes.TrimSpace(out.Bytes()), []byte("{")) {
		t.Fatal("schema output is not a JSON object")
	}
}

func TestSchemaDoesNotOfferJSON(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"schema", "--json"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if got := codeFor(err); got != ExitUsage {
		t.Fatalf("codeFor(error) = %d, want %d; err = %v", got, ExitUsage, err)
	}
	if !strings.Contains(err.Error(), "unknown flag") {
		t.Fatalf("Execute() error = %v, want unknown flag", err)
	}
}

func TestSchemaRejectsArguments(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"schema", "QUALITY.md"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if got := codeFor(err); got != ExitUsage {
		t.Fatalf("codeFor(error) = %d, want %d; err = %v", got, ExitUsage, err)
	}
}

func TestSchemaWriteFailureMapsToInternal(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(failingWriter{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"schema"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want write error")
	}
	if got := codeFor(err); got != ExitInternal {
		t.Fatalf("codeFor(error) = %d, want %d; err = %v", got, ExitInternal, err)
	}
}
