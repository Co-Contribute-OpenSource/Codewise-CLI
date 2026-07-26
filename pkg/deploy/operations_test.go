package deploy

import (
	"errors"
	"strings"
	"testing"

	"github.com/aryansharma9917/codewise-cli/pkg/env"
)

func useFakeRunner(t *testing.T, fake *fakeRunner) {
	t.Helper()
	previous := commandRunner
	commandRunner = fake
	t.Cleanup(func() { commandRunner = previous })
}

func createTestEnvironment(t *testing.T, name string, strategy Strategy) *env.Env {
	t.Helper()
	t.Setenv("CODEWISE_HOME", t.TempDir())
	e := &env.Env{
		Name: name,
		K8s: env.K8sConfig{
			Namespace: "codewise-test",
			Context:   "kind-codewise",
			Strategy:  string(strategy),
		},
		Helm: env.HelmConfig{
			Release: name,
			Chart:   "./helm/chart",
		},
	}
	if strategy == StrategyGitOps {
		e.GitOps = env.GitOpsConfig{Repo: "https://example.invalid/repo.git", Branch: "main"}
	}
	if err := env.CreateEnvFromParts(name, e.K8s, e.Helm, e.GitOps, e.Values); err != nil {
		t.Fatal(err)
	}
	return e
}

func TestBuildCommandForEveryStrategy(t *testing.T) {
	tests := []struct {
		strategy Strategy
		binary   string
		contains string
	}{
		{StrategyKubectl, "kubectl", "apply"},
		{StrategyHelm, "helm", "upgrade"},
		{StrategyGitOps, "kubectl", "-"},
	}
	for _, tt := range tests {
		e := &env.Env{
			Name:   "test",
			K8s:    env.K8sConfig{Namespace: "test", Context: "kind-test", Strategy: string(tt.strategy)},
			Helm:   env.HelmConfig{Release: "test", Chart: "./chart"},
			GitOps: env.GitOpsConfig{Repo: "https://example.invalid/repo"},
		}
		command, strategy, err := BuildCommand(e)
		if err != nil {
			t.Fatal(err)
		}
		if strategy != tt.strategy || command.Name != tt.binary || !strings.Contains(strings.Join(command.Args, " "), tt.contains) {
			t.Errorf("%s command = %#v, strategy = %s", tt.strategy, command, strategy)
		}
	}
}

func TestHistoryAndRollbackRejectNonHelmStrategies(t *testing.T) {
	fake := &fakeRunner{}
	useFakeRunner(t, fake)
	createTestEnvironment(t, "raw", StrategyKubectl)

	if err := History("raw"); err == nil || !strings.Contains(err.Error(), "only available for Helm") {
		t.Fatalf("History error = %v", err)
	}
	if err := Rollback("raw", 1); err == nil || !strings.Contains(err.Error(), "only supported for Helm") {
		t.Fatalf("Rollback error = %v", err)
	}
	if len(fake.calls) != 0 {
		t.Fatalf("unsupported operations executed commands: %#v", fake.calls)
	}
}

func TestNamespaceLookupDoesNotHideAuthorizationFailure(t *testing.T) {
	fake := &fakeRunner{combinedOutput: []byte("Error from server (Forbidden)"), err: errors.New("exit status 1")}
	useFakeRunner(t, fake)

	_, err := namespaceExists("production", "kind-test")
	if err == nil || !strings.Contains(err.Error(), "Forbidden") {
		t.Fatalf("namespaceExists error = %v", err)
	}
	if len(fake.calls) != 1 {
		t.Fatalf("namespace lookup unexpectedly attempted creation: %#v", fake.calls)
	}
}

func TestNamespaceNotFoundIsRecognized(t *testing.T) {
	fake := &fakeRunner{combinedOutput: []byte(`Error from server (NotFound): namespaces "missing" not found`), err: errors.New("exit status 1")}
	useFakeRunner(t, fake)

	exists, err := namespaceExists("missing", "")
	if err != nil || exists {
		t.Fatalf("namespaceExists = %v, %v", exists, err)
	}
}

func TestLogsFetchesPodThenStreamsLogs(t *testing.T) {
	fake := &fakeRunner{output: []byte("pod/api-123\n")}
	useFakeRunner(t, fake)
	createTestEnvironment(t, "logs", StrategyKubectl)

	if err := Logs("logs", true); err != nil {
		t.Fatal(err)
	}
	if len(fake.calls) != 2 {
		t.Fatalf("calls = %#v", fake.calls)
	}
	logArgs := strings.Join(fake.calls[1].args, " ")
	if !strings.Contains(logArgs, "logs pod/api-123") || !strings.Contains(logArgs, "-f") {
		t.Fatalf("log args = %q", logArgs)
	}
}

func TestHelmStatusUsesHelmAndKubectl(t *testing.T) {
	fake := &fakeRunner{combinedOutput: []byte("ok")}
	useFakeRunner(t, fake)
	createTestEnvironment(t, "release", StrategyHelm)

	if err := Status("release"); err != nil {
		t.Fatal(err)
	}
	if len(fake.calls) != 2 || fake.calls[0].name != "helm" || fake.calls[1].name != "kubectl" {
		t.Fatalf("calls = %#v", fake.calls)
	}
}
