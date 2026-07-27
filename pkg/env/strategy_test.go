package env

import "testing"

func TestSetStrategy(t *testing.T) {
	t.Setenv(codewiseHomeEnv, t.TempDir())
	if err := CreateEnv("staging", CreateOptions{}); err != nil {
		t.Fatal(err)
	}

	if err := SetStrategy("staging", "helm"); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadEnv("staging")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.K8s.Strategy != "helm" {
		t.Fatalf("strategy = %q", loaded.K8s.Strategy)
	}

	if err := SetStrategy("staging", "auto"); err != nil {
		t.Fatal(err)
	}
	loaded, err = LoadEnv("staging")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.K8s.Strategy != "" {
		t.Fatalf("auto strategy persisted as %q", loaded.K8s.Strategy)
	}
}

func TestSetStrategyRejectsUnknownValue(t *testing.T) {
	if err := SetStrategy("staging", "terraform"); err == nil {
		t.Fatal("expected invalid strategy error")
	}
}
