package deploy

import (
	"fmt"
	"os"
	"strings"
)

func Logs(envName string, follow bool) error {

	environment, err := LoadEnvironment(envName)
	if err != nil {
		return err
	}

	ns := environment.K8s.Namespace
	ctx := environment.K8s.Context

	// Get pods
	args := []string{
		"get",
		"pods",
		"-n",
		ns,
		"-o",
		"name",
	}

	if ctx != "" {
		args = append(args, "--context", ctx)
	}

	out, err := commandRunner.Output("kubectl", args...)
	if err != nil {
		return fmt.Errorf("failed to fetch pods: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return fmt.Errorf("no pods found")
	}

	// Take first pod
	pod := lines[0]

	fmt.Println("Fetching logs for:", pod)
	fmt.Println()

	logArgs := []string{
		"logs",
		pod,
		"-n",
		ns,
	}

	if follow {
		logArgs = append(logArgs, "-f")
	}

	if ctx != "" {
		logArgs = append(logArgs, "--context", ctx)
	}

	return commandRunner.Run("kubectl", logArgs, nil, os.Stdout, os.Stderr)
}
