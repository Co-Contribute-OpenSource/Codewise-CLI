package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aryansharma9917/codewise-cli/pkg/env"
)

type Explanation struct {
	Environment   string
	Strategy      Strategy
	Reason        string
	Prerequisites []string
	Command       string
	Steps         []string
	Recovery      string
}

func BuildExplanation(environment *env.Env) (*Explanation, error) {
	command, strategy, err := BuildCommand(environment)
	if err != nil {
		return nil, err
	}

	explanation := &Explanation{
		Environment: environment.Name,
		Strategy:    strategy,
		Reason:      strategyReason(environment, strategy),
		Command:     command.Name + " " + strings.Join(command.Args, " "),
		Steps: []string{
			"Load and validate the selected environment",
			"Run local tool and Kubernetes connectivity preflight checks",
			fmt.Sprintf("Ensure namespace %q exists", environment.K8s.Namespace),
			"Execute the deployment command and stream its output",
			"Wait for Kubernetes Deployments to finish rolling out",
			"Collect pod diagnostics if rollout monitoring fails",
		},
	}

	switch strategy {
	case StrategyHelm:
		explanation.Prerequisites = []string{"helm", "kubectl", "kubeconfig and cluster access"}
		explanation.Recovery = fmt.Sprintf(
			"codewise deploy history --env %s, then codewise deploy rollback --env %s --revision REVISION",
			environment.Name, environment.Name,
		)
	case StrategyGitOps:
		explanation.Prerequisites = []string{"kubectl", "Argo CD Application CRD", "argocd namespace permissions", "repository access configured in Argo CD"}
		explanation.Recovery = "Revert the desired-state Git commit and allow Argo CD to reconcile it"
	case StrategyKubectl:
		explanation.Prerequisites = []string{"kubectl", "kubeconfig and cluster access", "Kubernetes manifests under k8s/"}
		explanation.Recovery = "Reapply known-good manifests or use kubectl rollout undo for a Deployment"
	}
	return explanation, nil
}

func Explain(envName string) error {
	environment, err := LoadEnvironment(envName)
	if err != nil {
		return err
	}
	explanation, err := BuildExplanation(environment)
	if err != nil {
		return err
	}

	fmt.Println("Deployment Explanation")
	fmt.Println("----------------------")
	fmt.Println("Environment:", explanation.Environment)
	fmt.Println("Strategy:", explanation.Strategy)
	fmt.Println("Selected because:", explanation.Reason)
	fmt.Println()
	fmt.Println("Prerequisites:")
	for _, prerequisite := range explanation.Prerequisites {
		fmt.Println(" -", prerequisite)
	}
	fmt.Println()
	fmt.Println("Command:")
	fmt.Println(" ", explanation.Command)
	fmt.Println()
	fmt.Println("Execution:")
	for index, step := range explanation.Steps {
		fmt.Printf(" %d. %s\n", index+1, step)
	}
	fmt.Println()
	fmt.Println("Recovery:")
	fmt.Println(" ", explanation.Recovery)
	return nil
}

func strategyReason(environment *env.Env, strategy Strategy) string {
	if environment.K8s.Strategy != "" {
		return fmt.Sprintf("k8s.strategy explicitly specifies %q", environment.K8s.Strategy)
	}
	if strategy == StrategyGitOps {
		return "gitops.repo is configured"
	}
	if strategy == StrategyHelm {
		chart := filepath.Join(".", "helm", "chart")
		if environment.Helm.Chart != "" {
			chart = environment.Helm.Chart
		}
		if _, err := os.Stat(chart); err == nil {
			return fmt.Sprintf("Helm chart %q exists", chart)
		}
		return "a Helm chart was detected"
	}
	return "no GitOps repository or Helm chart was detected, so raw Kubernetes manifests are used"
}
