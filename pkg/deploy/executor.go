package deploy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// CommandRunner is the single boundary between deployment logic and external
// programs. Tests can replace it with a fake without installing Docker,
// kubectl, or Helm.
type CommandRunner interface {
	Run(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) error
	Output(name string, args ...string) ([]byte, error)
	CombinedOutput(name string, args ...string) ([]byte, error)
}

type OSCommandRunner struct{}

func (OSCommandRunner) Run(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func (OSCommandRunner) Output(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func (OSCommandRunner) CombinedOutput(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

var commandRunner CommandRunner = OSCommandRunner{}

type Executor struct {
	DryRun bool
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Runner CommandRunner
}

func (e *Executor) Run(name string, args ...string) error {

	cmdStr := name + " " + strings.Join(args, " ")

	if e.DryRun {
		fmt.Println("[dry-run]", cmdStr)
		return nil
	}

	runner := e.Runner
	if runner == nil {
		runner = commandRunner
	}
	stdout := e.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := e.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	fmt.Println("Running:", cmdStr)

	return runner.Run(name, args, e.Stdin, stdout, stderr)
}

func outputError(action string, output []byte, err error) error {
	if err == nil {
		return nil
	}
	message := strings.TrimSpace(string(bytes.TrimSpace(output)))
	if message == "" {
		return fmt.Errorf("%s: %w", action, err)
	}
	return fmt.Errorf("%s: %s: %w", action, message, err)
}
