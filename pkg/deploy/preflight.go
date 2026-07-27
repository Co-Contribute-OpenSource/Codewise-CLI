package deploy

import (
	"fmt"

	"github.com/aryansharma9917/codewise-cli/pkg/env"
)

func runCheck(name string, args ...string) error {

	if err := commandRunner.Run(name, args, nil, nil, nil); err != nil {
		return err
	}

	return nil
}

func Preflight(environment *env.Env) error {

	fmt.Println("Running preflight checks...")

	strategy := ResolveStrategy(environment)

	////////////////////////////////////////////////////
	// Check binary availability
	////////////////////////////////////////////////////

	if strategy == StrategyHelm {
		if err := runCheck("helm", "version"); err != nil {
			return fmt.Errorf("helm not available or not functioning: %w", err)
		}
	}

	// Every strategy uses kubectl for connectivity, namespace management,
	// rollout monitoring, or applying the Argo CD Application.
	if err := runCheck("kubectl", "version", "--client"); err != nil {
		return fmt.Errorf("kubectl not available or not functioning: %w", err)
	}

	////////////////////////////////////////////////////
	// Cluster connectivity check
	////////////////////////////////////////////////////

	args := []string{"cluster-info"}

	if environment.K8s.Context != "" {
		args = append(args, "--context", environment.K8s.Context)
	}

	if err := runCheck("kubectl", args...); err != nil {
		return fmt.Errorf("cannot reach kubernetes cluster; check kube-context: %w", err)
	}

	fmt.Println("Cluster reachable")
	return nil
}
