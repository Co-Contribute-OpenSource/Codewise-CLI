package deploy

import (
	"fmt"
)

func History(envName string) error {

	environment, err := LoadEnvironment(envName)
	if err != nil {
		return err
	}

	ns := environment.K8s.Namespace
	ctx := environment.K8s.Context
	release := environment.Helm.Release
	strategy := ResolveStrategy(environment)
	if strategy != StrategyHelm {
		return fmt.Errorf("deployment history is only available for Helm environments (resolved strategy: %s)", strategy)
	}

	fmt.Println("Release History")
	fmt.Println("----------------")
	fmt.Println("Environment:", envName)
	fmt.Println("Namespace:", ns)
	fmt.Println("Release:", release)
	fmt.Println()

	args := []string{
		"history",
		release,
		"-n",
		ns,
	}

	if ctx != "" {
		args = append(args, "--kube-context", ctx)
	}

	out, err := commandRunner.CombinedOutput("helm", args...)
	if err != nil {
		return outputError("failed to fetch Helm release history", out, err)
	}

	fmt.Println(string(out))
	return nil
}
