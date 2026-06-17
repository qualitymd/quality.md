package cli

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestSpecWritesVerbatimSpecification(t *testing.T) {
	want, err := os.ReadFile("../../SPECIFICATION.md")
	if err != nil {
		t.Fatalf("os.ReadFile(SPECIFICATION.md) error = %v", err)
	}

	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"spec"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !bytes.Equal(out.Bytes(), want) {
		t.Fatal("stdout did not match SPECIFICATION.md byte-for-byte")
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestSpecDoesNotOfferJSON(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"spec", "--json"})
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

func TestSpecRejectsArguments(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"spec", "QUALITY.md"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if got := codeFor(err); got != ExitUsage {
		t.Fatalf("codeFor(error) = %d, want %d; err = %v", got, ExitUsage, err)
	}
}

func TestSpecNoColorForcesPlainOutput(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if shouldRenderSpec(panicFDWriter{}) {
		t.Fatal("shouldRenderSpec() = true, want false when NO_COLOR is set")
	}
}

func TestSpecWriteFailureMapsToInternal(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(failingWriter{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"spec"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want write error")
	}
	if got := codeFor(err); got != ExitInternal {
		t.Fatalf("codeFor(error) = %d, want %d; err = %v", got, ExitInternal, err)
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

type panicFDWriter struct{}

func (panicFDWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (panicFDWriter) Fd() uintptr {
	panic("Fd called despite NO_COLOR")
}
