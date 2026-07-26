package deploy

import (
	"testing"

	"github.com/aryansharma9917/codewise-cli/pkg/env"
)

func TestResolveStrategyHonorsExplicitConfiguration(t *testing.T) {
	for _, strategy := range []Strategy{StrategyHelm, StrategyKubectl, StrategyGitOps} {
		e := &env.Env{
			K8s:    env.K8sConfig{Strategy: string(strategy)},
			GitOps: env.GitOpsConfig{Repo: "https://example.invalid/repo"},
		}
		if got := ResolveStrategy(e); got != strategy {
			t.Errorf("ResolveStrategy(%q) = %q", strategy, got)
		}
	}
}
