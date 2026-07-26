package deploy

import (
	"fmt"
	"strings"

	"github.com/aryansharma9917/codewise-cli/pkg/env"
)

func namespaceExists(namespace string, context string) (bool, error) {

	args := []string{"get", "ns", namespace}

	if context != "" {
		args = append(args, "--context", context)
	}

	output, err := commandRunner.CombinedOutput("kubectl", args...)
	if err == nil {
		return true, nil
	}
	if strings.Contains(strings.ToLower(string(output)), "notfound") ||
		strings.Contains(strings.ToLower(string(output)), "not found") {
		return false, nil
	}
	return false, outputError("failed to check Kubernetes namespace", output, err)
}

func createNamespace(namespace string, context string) error {

	args := []string{"create", "ns", namespace}

	if context != "" {
		args = append(args, "--context", context)
	}

	output, err := commandRunner.CombinedOutput("kubectl", args...)
	if err != nil {
		return outputError("failed to create namespace", output, err)
	}

	return nil
}

func EnsureNamespace(environment *env.Env) error {

	ns := environment.K8s.Namespace
	ctx := environment.K8s.Context

	fmt.Printf("Checking namespace \"%s\"...\n", ns)

	exists, err := namespaceExists(ns, ctx)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("Namespace exists")
		return nil
	}

	fmt.Println("Namespace not found. Creating namespace...")

	if err := createNamespace(ns, ctx); err != nil {
		return err
	}

	fmt.Println("Namespace created")
	return nil
}
