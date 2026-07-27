package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var validStrategies = map[string]bool{
	"":        true,
	"helm":    true,
	"kubectl": true,
	"gitops":  true,
}

// SetStrategy persists a deployment strategy. "auto" clears the explicit
// value and restores backwards-compatible automatic detection.
func SetStrategy(name, strategy string) error {
	strategy = strings.ToLower(strings.TrimSpace(strategy))
	if strategy == "auto" {
		strategy = ""
	}
	if !validStrategies[strategy] {
		return fmt.Errorf("invalid strategy %q: use auto, helm, kubectl, or gitops", strategy)
	}

	environment, err := LoadEnv(name)
	if err != nil {
		return err
	}
	environment.K8s.Strategy = strategy

	dir, err := envDir(name)
	if err != nil {
		return err
	}
	target := filepath.Join(dir, "k8s.yaml")
	temporary, err := os.CreateTemp(dir, ".k8s-*.yaml")
	if err != nil {
		return fmt.Errorf("create temporary strategy file: %w", err)
	}
	temporaryPath := temporary.Name()
	if err := temporary.Close(); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}
	if err := writeYAML(temporaryPath, environment.K8s); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}
	if err := os.Rename(temporaryPath, target); err != nil {
		_ = os.Remove(temporaryPath)
		return fmt.Errorf("replace strategy configuration: %w", err)
	}
	return nil
}
