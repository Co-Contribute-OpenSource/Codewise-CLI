package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateName(t *testing.T) {
	valid := []string{"dev", "preview-pr-42", "team-blue", "test-env"}
	for _, name := range valid {
		if err := ValidateName(name); err != nil {
			t.Errorf("ValidateName(%q) unexpected error: %v", name, err)
		}
	}

	invalid := []string{"", ".", "..", "../prod", "team/dev", `team\dev`, "-dev", "dev-", "PROD", "has space", "team.blue", "test_env"}
	for _, name := range invalid {
		if err := ValidateName(name); err == nil {
			t.Errorf("ValidateName(%q) expected error", name)
		}
	}
}

func TestEnvironmentPathsStayUnderCodewiseHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv(codewiseHomeEnv, home)

	dir, err := EnvDir("staging")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, "envs", "staging")
	if dir != want {
		t.Fatalf("EnvDir() = %q, want %q", dir, want)
	}

	if _, err := EnvDir("../../outside"); err == nil {
		t.Fatal("expected traversal name to be rejected")
	}
}

func TestCreateLoadAndRemoveEnvironment(t *testing.T) {
	t.Setenv(codewiseHomeEnv, t.TempDir())
	parts := K8sConfig{Namespace: "staging", Strategy: "kubectl"}

	if err := CreateEnvFromParts("staging", parts, HelmConfig{}, GitOpsConfig{}, ValuesConfig{}); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadEnv("staging")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.K8s.Strategy != "kubectl" {
		t.Fatalf("strategy = %q", loaded.K8s.Strategy)
	}
	if err := RemoveEnv("staging"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(os.Getenv(codewiseHomeEnv), "envs", "staging")); !os.IsNotExist(err) {
		t.Fatalf("environment directory still exists: %v", err)
	}
}
