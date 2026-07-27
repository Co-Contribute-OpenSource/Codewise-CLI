package deploy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aryansharma9917/codewise-cli/pkg/env"
)

type Strategy string

const (
	StrategyHelm    Strategy = "helm"
	StrategyKubectl Strategy = "kubectl"
	StrategyGitOps  Strategy = "gitops"
)

func ResolveStrategy(e *env.Env) Strategy {
	switch Strategy(strings.ToLower(e.K8s.Strategy)) {
	case StrategyHelm:
		return StrategyHelm
	case StrategyKubectl:
		return StrategyKubectl
	case StrategyGitOps:
		return StrategyGitOps
	}

	// check if gitops repo is configured
	if e.GitOps.Repo != "" {
		return StrategyGitOps
	}

	// check if the configured Helm chart exists
	chart := e.Helm.Chart
	if chart == "" {
		chart = filepath.Join(".", "helm", "chart")
	}
	if _, err := os.Stat(chart); err == nil {
		return StrategyHelm
	}

	return StrategyKubectl
}
