package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func findCommandByName(parent *cobra.Command, use string) *cobra.Command {
	for _, c := range parent.Commands() {
		if c.Name() == use {
			return c
		}
	}
	return nil
}

func TestVersionCommandIsRegistered(t *testing.T) {
	if got := findCommandByName(rootCmd, "version"); got == nil {
		t.Fatalf("expected version command to be registered on root")
	}
}

func TestRootVersionMatchesBuildVersion(t *testing.T) {
	if rootCmd.Version != Version {
		t.Fatalf("expected root version %q to match build version %q", rootCmd.Version, Version)
	}
}

func TestK8sInitCommandIsRegistered(t *testing.T) {
	k8s := findCommandByName(rootCmd, "k8s")
	if k8s == nil {
		t.Fatalf("expected k8s command to be registered on root")
	}

	if got := findCommandByName(k8s, "init"); got == nil {
		t.Fatalf("expected k8s init subcommand to be registered")
	}
}

func TestDockerPushTagFlagIsRegistered(t *testing.T) {
	docker := findCommandByName(rootCmd, "docker")
	if docker == nil {
		t.Fatalf("expected docker command to be registered on root")
	}

	push := findCommandByName(docker, "push")
	if push == nil {
		t.Fatalf("expected docker push subcommand to be registered")
	}

	if got := push.Flags().Lookup("tag"); got == nil {
		t.Fatalf("expected docker push to register the --tag flag")
	}
}

func TestDoctorCommandIsRegistered(t *testing.T) {
	if got := findCommandByName(rootCmd, "doctor"); got == nil {
		t.Fatal("expected doctor command to be registered")
	}
}

func TestEnvironmentControlCommandsAreRegistered(t *testing.T) {
	envCommand := findCommandByName(rootCmd, "env")
	for _, name := range []string{"show", "strategy"} {
		if got := findCommandByName(envCommand, name); got == nil {
			t.Fatalf("expected env %s command to be registered", name)
		}
	}
}
