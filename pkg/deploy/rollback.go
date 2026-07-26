package deploy

import (
	"fmt"
	"os"
)

func Rollback(envName string, revision int) error {

	environment, err := LoadEnvironment(envName)
	if err != nil {
		return err
	}

	ns := environment.K8s.Namespace
	ctx := environment.K8s.Context
	release := environment.Helm.Release
	strategy := ResolveStrategy(environment)
	if strategy != StrategyHelm {
		return fmt.Errorf("automatic rollback is only supported for Helm environments (resolved strategy: %s); revert desired state in Git for GitOps or apply a previous manifest for kubectl", strategy)
	}

	fmt.Println("Starting rollback...")
	fmt.Println("Environment:", envName)
	fmt.Println("Release:", release)
	fmt.Println("Revision:", revision)
	fmt.Println()

	args := []string{
		"rollback",
		release,
		fmt.Sprintf("%d", revision),
		"-n",
		ns,
	}

	if ctx != "" {
		args = append(args, "--kube-context", ctx)
	}

	if err := commandRunner.Run("helm", args, nil, os.Stdout, os.Stderr); err != nil {
		return fmt.Errorf("Helm rollback failed: %w", err)
	}

	fmt.Println("Rollback executed successfully.")
	fmt.Println("Verifying rollout...")

	return MonitorRollout(environment)
}
