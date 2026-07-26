package env

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	codewiseHomeEnv = "CODEWISE_HOME"
)

var environmentNamePattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

// ValidateName ensures an environment name is both predictable and safe to
// use as a directory and as a default Helm release name.
func ValidateName(name string) error {
	if !environmentNamePattern.MatchString(name) || name == "." || name == ".." {
		return fmt.Errorf("invalid environment name %q: use a 1-63 character DNS label containing lowercase letters, numbers, or hyphens; start and end with a letter or number", name)
	}
	if strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("invalid environment name %q: path separators are not allowed", name)
	}
	return nil
}

/////////////////////////////////////////////////////////
// PATH HELPERS
/////////////////////////////////////////////////////////

func baseEnvPath() (string, error) {

	if home := os.Getenv(codewiseHomeEnv); home != "" {
		return filepath.Join(home, "envs"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve user home: %w", err)
	}

	return filepath.Join(home, ".codewise", "envs"), nil
}

func envDir(name string) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}

	base, err := baseEnvPath()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(base, name)
	relative, err := filepath.Rel(base, dir)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("environment path escapes Codewise home")
	}
	return dir, nil
}

// EnvDir returns the full path to an environment directory.
func EnvDir(name string) (string, error) {
	return envDir(name)
}

// RemoveEnv deletes an environment directory and its contents.
func RemoveEnv(name string) error {
	dir, err := envDir(name)
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

func envExists(name string) bool {

	dir, err := envDir(name)
	if err != nil {
		return false
	}

	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

func ensureBaseDir() error {

	base, err := baseEnvPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(base, 0755)
}

/////////////////////////////////////////////////////////
// YAML READER
/////////////////////////////////////////////////////////

func readYAML(path string, out interface{}) error {

	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(raw, out)
}

/////////////////////////////////////////////////////////
// SINGLE SOURCE OF TRUTH
/////////////////////////////////////////////////////////

func LoadEnv(name string) (*Env, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}

	if !envExists(name) {
		return nil, fmt.Errorf("environment %q does not exist", name)
	}

	dir, err := envDir(name)
	if err != nil {
		return nil, err
	}

	k8s := K8sConfig{}
	helm := HelmConfig{}
	gitops := GitOpsConfig{}
	values := ValuesConfig{}

	if err := readYAML(filepath.Join(dir, "k8s.yaml"), &k8s); err != nil {
		return nil, err
	}

	if err := readYAML(filepath.Join(dir, "helm.yaml"), &helm); err != nil {
		return nil, err
	}

	if err := readYAML(filepath.Join(dir, "gitops.yaml"), &gitops); err != nil {
		return nil, err
	}

	if err := readYAML(filepath.Join(dir, "values.yaml"), &values); err != nil {
		return nil, err
	}

	return &Env{
		Name:   name,
		K8s:    k8s,
		Helm:   helm,
		GitOps: gitops,
		Values: values,
	}, nil
}
