package deploy

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type runnerCall struct {
	name string
	args []string
}

type fakeRunner struct {
	calls          []runnerCall
	output         []byte
	combinedOutput []byte
	err            error
}

func (f *fakeRunner) Run(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	f.calls = append(f.calls, runnerCall{name, append([]string(nil), args...)})
	if stdout != nil {
		_, _ = io.WriteString(stdout, "streamed output")
	}
	return f.err
}

func (f *fakeRunner) Output(name string, args ...string) ([]byte, error) {
	f.calls = append(f.calls, runnerCall{name, append([]string(nil), args...)})
	return f.output, f.err
}

func (f *fakeRunner) CombinedOutput(name string, args ...string) ([]byte, error) {
	f.calls = append(f.calls, runnerCall{name, append([]string(nil), args...)})
	return f.combinedOutput, f.err
}

func TestExecutorStreamsCommandOutput(t *testing.T) {
	fake := &fakeRunner{}
	var stdout bytes.Buffer
	executor := Executor{Runner: fake, Stdout: &stdout, Stderr: io.Discard}

	if err := executor.Run("kubectl", "logs", "pod/api"); err != nil {
		t.Fatal(err)
	}
	if stdout.String() != "streamed output" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if len(fake.calls) != 1 || fake.calls[0].name != "kubectl" {
		t.Fatalf("calls = %#v", fake.calls)
	}
}

func TestExecutorDryRunDoesNotExecute(t *testing.T) {
	fake := &fakeRunner{}
	if err := (&Executor{DryRun: true, Runner: fake}).Run("helm", "upgrade"); err != nil {
		t.Fatal(err)
	}
	if len(fake.calls) != 0 {
		t.Fatalf("dry run executed commands: %#v", fake.calls)
	}
}

func TestOutputErrorPreservesToolDiagnostics(t *testing.T) {
	err := outputError("deploy failed", []byte("forbidden"), errors.New("exit status 1"))
	if !strings.Contains(err.Error(), "forbidden") || !strings.Contains(err.Error(), "exit status 1") {
		t.Fatalf("error lost diagnostics: %v", err)
	}
}
