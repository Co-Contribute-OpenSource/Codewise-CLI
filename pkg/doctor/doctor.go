package doctor

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

type Status string

const (
	Pass Status = "PASS"
	Warn Status = "WARN"
	Fail Status = "FAIL"
)

type Result struct {
	Name    string
	Status  Status
	Message string
}

type Options struct {
	Cluster bool
	Context string
	Require []string
}

type System interface {
	LookPath(file string) (string, error)
	CombinedOutput(name string, args ...string) ([]byte, error)
}

type OSSystem struct{}

func (OSSystem) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (OSSystem) CombinedOutput(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

func Run(options Options) ([]Result, error) {
	return RunWithSystem(options, OSSystem{})
}

func RunWithSystem(options Options, system System) ([]Result, error) {
	required := make(map[string]bool)
	for _, name := range options.Require {
		required[strings.ToLower(strings.TrimSpace(name))] = true
	}

	type tool struct {
		name string
		args []string
	}
	tools := []tool{
		{name: "docker", args: []string{"--version"}},
		{name: "kubectl", args: []string{"version", "--client", "--output=yaml"}},
		{name: "helm", args: []string{"version", "--short"}},
		{name: "git", args: []string{"--version"}},
	}

	results := make([]Result, 0, len(tools)+1)
	failed := make([]string, 0)
	for _, candidate := range tools {
		if _, err := system.LookPath(candidate.name); err != nil {
			status := Warn
			if required[candidate.name] {
				status = Fail
				failed = append(failed, candidate.name)
			}
			results = append(results, Result{Name: candidate.name, Status: status, Message: "not found in PATH"})
			continue
		}

		output, err := system.CombinedOutput(candidate.name, candidate.args...)
		if err != nil {
			status := Warn
			if required[candidate.name] {
				status = Fail
				failed = append(failed, candidate.name)
			}
			message := strings.TrimSpace(string(output))
			if message == "" {
				message = err.Error()
			}
			results = append(results, Result{Name: candidate.name, Status: status, Message: firstLine(message)})
			continue
		}
		results = append(results, Result{Name: candidate.name, Status: Pass, Message: firstLine(string(output))})
	}

	if options.Cluster {
		args := []string{"cluster-info"}
		if options.Context != "" {
			args = append(args, "--context", options.Context)
		}
		output, err := system.CombinedOutput("kubectl", args...)
		if err != nil {
			message := strings.TrimSpace(string(output))
			if message == "" {
				message = err.Error()
			}
			results = append(results, Result{Name: "cluster", Status: Fail, Message: firstLine(message)})
			failed = append(failed, "cluster")
		} else {
			results = append(results, Result{Name: "cluster", Status: Pass, Message: "Kubernetes API is reachable"})
		}
	}

	if len(failed) > 0 {
		sort.Strings(failed)
		return results, fmt.Errorf("required checks failed: %s", strings.Join(failed, ", "))
	}
	return results, nil
}

func firstLine(value string) string {
	value = strings.TrimSpace(value)
	if before, _, found := strings.Cut(value, "\n"); found {
		return before
	}
	return value
}
