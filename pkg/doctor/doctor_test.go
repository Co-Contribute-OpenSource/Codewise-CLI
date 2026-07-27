package doctor

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type fakeSystem struct {
	missing map[string]bool
	failing map[string]bool
}

func (f fakeSystem) LookPath(file string) (string, error) {
	if f.missing[file] {
		return "", errors.New("missing")
	}
	return "/usr/bin/" + file, nil
}

func (f fakeSystem) CombinedOutput(name string, args ...string) ([]byte, error) {
	if f.failing[name] {
		return []byte("connection refused"), errors.New("failed")
	}
	if len(args) > 0 && args[0] == "cluster-info" {
		return []byte("Kubernetes control plane is running"), nil
	}
	return []byte(fmt.Sprintf("%s version\nmore details", name)), nil
}

func TestRunWarnsForOptionalMissingTool(t *testing.T) {
	results, err := RunWithSystem(Options{}, fakeSystem{missing: map[string]bool{"docker": true}})
	if err != nil {
		t.Fatal(err)
	}
	if results[0].Status != Warn {
		t.Fatalf("docker status = %s", results[0].Status)
	}
}

func TestRunFailsForRequiredMissingTool(t *testing.T) {
	_, err := RunWithSystem(
		Options{Require: []string{"docker"}},
		fakeSystem{missing: map[string]bool{"docker": true}},
	)
	if err == nil || !strings.Contains(err.Error(), "docker") {
		t.Fatalf("error = %v", err)
	}
}

func TestRunChecksCluster(t *testing.T) {
	results, err := RunWithSystem(Options{Cluster: true}, fakeSystem{})
	if err != nil {
		t.Fatal(err)
	}
	last := results[len(results)-1]
	if last.Name != "cluster" || last.Status != Pass {
		t.Fatalf("cluster result = %#v", last)
	}
}

func TestRunFailsWhenClusterIsUnreachable(t *testing.T) {
	_, err := RunWithSystem(Options{Cluster: true}, fakeSystem{failing: map[string]bool{"kubectl": true}})
	if err == nil || !strings.Contains(err.Error(), "cluster") {
		t.Fatalf("error = %v", err)
	}
}
