package deploy

import (
	"fmt"
	"strings"
)

func Status(envName string) error {

	environment, err := LoadEnvironment(envName)
	if err != nil {
		return err
	}

	ns := environment.K8s.Namespace
	ctx := environment.K8s.Context
	release := environment.Helm.Release
	strategy := ResolveStrategy(environment)

	fmt.Println("Deployment Status")
	fmt.Println("-----------------")
	fmt.Println("Environment:", envName)
	fmt.Println("Namespace:", ns)
	fmt.Println("Strategy:", strategy)
	fmt.Println()

	switch strategy {
	case StrategyHelm:
		fmt.Println("Release:", release)
		args := []string{"status", release, "-n", ns}
		if ctx != "" {
			args = append(args, "--kube-context", ctx)
		}
		out, err := commandRunner.CombinedOutput("helm", args...)
		if err != nil {
			return outputError("failed to fetch Helm release status", out, err)
		}
		fmt.Println(string(out))
	case StrategyGitOps:
		args := []string{"get", "application", environment.Name, "-n", "argocd", "-o", "wide"}
		if ctx != "" {
			args = append(args, "--context", ctx)
		}
		out, err := commandRunner.CombinedOutput("kubectl", args...)
		if err != nil {
			return outputError("failed to fetch Argo CD Application status", out, err)
		}
		fmt.Println(string(out))
	case StrategyKubectl:
		fmt.Println("kubectl-managed resources do not have revision metadata; showing live workloads.")
	}

	// Pods
	fmt.Println("Pods:")
	podArgs := []string{
		"get",
		"pods",
		"-n",
		ns,
		"-o",
		"wide",
	}

	if ctx != "" {
		podArgs = append(podArgs, "--context", ctx)
	}

	podsOut, err := commandRunner.CombinedOutput("kubectl", podArgs...)
	if err != nil {
		return outputError("failed to fetch pods", podsOut, err)
	}
	fmt.Println(strings.TrimSpace(string(podsOut)))

	return nil
}
