package deploy

import (
	"strings"
	"testing"

	"github.com/aryansharma9917/codewise-cli/pkg/env"
)

func TestBuildExplanationForEveryStrategy(t *testing.T) {
	tests := []struct {
		strategy Strategy
		recovery string
	}{
		{StrategyHelm, "deploy rollback"},
		{StrategyKubectl, "kubectl rollout undo"},
		{StrategyGitOps, "Revert"},
	}

	for _, tt := range tests {
		environment := &env.Env{
			Name: "staging",
			K8s: env.K8sConfig{
				Namespace: "staging",
				Strategy:  string(tt.strategy),
			},
			Helm: env.HelmConfig{Release: "staging", Chart: "./helm/chart"},
			GitOps: env.GitOpsConfig{
				Repo: "https://example.invalid/repo.git",
			},
		}
		explanation, err := BuildExplanation(environment)
		if err != nil {
			t.Fatal(err)
		}
		if explanation.Strategy != tt.strategy {
			t.Errorf("strategy = %s", explanation.Strategy)
		}
		if !strings.Contains(explanation.Reason, "explicitly") {
			t.Errorf("reason = %q", explanation.Reason)
		}
		if !strings.Contains(explanation.Recovery, tt.recovery) {
			t.Errorf("recovery = %q", explanation.Recovery)
		}
		if len(explanation.Steps) < 5 || explanation.Command == "" {
			t.Errorf("incomplete explanation: %#v", explanation)
		}
	}
}
