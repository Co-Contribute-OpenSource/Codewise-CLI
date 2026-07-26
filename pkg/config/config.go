package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDirName   = ".codewise"
	ConfigFileName  = "config.yaml"
	codewiseHomeEnv = "CODEWISE_HOME"
)

type Config struct {
	Version string `yaml:"version"`
	User    struct {
		Name string `yaml:"name"`
	} `yaml:"user"`

	Defaults struct {
		AppName   string `yaml:"app_name"`
		Image     string `yaml:"image"`
		ImageTag  string `yaml:"image_tag"`
		RepoURL   string `yaml:"repo_url"`
		Namespace string `yaml:"namespace"`
		Context   string `yaml:"context"`
		Branch    string `yaml:"branch"`
	} `yaml:"defaults"`
}

var DefaultConfig = []byte(`version: v1
user:
  name: aryan
defaults:
  app_name: myapp
  image: codewise
  image_tag: latest
  repo_url: https://github.com/example/repo
  namespace: default
  context: ""
  branch: main
`)

func configDir() (string, error) {
	if home := os.Getenv(codewiseHomeEnv); home != "" {
		return home, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve user home: %w", err)
	}
	return filepath.Join(home, ConfigDirName), nil
}

func InitConfig() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(dir, ConfigFileName)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	if _, err := os.Stat(configPath); err == nil {
		return configPath, fmt.Errorf("config already exists")
	}

	if err := os.WriteFile(configPath, DefaultConfig, 0644); err != nil {
		return "", err
	}

	return configPath, nil
}

func ReadConfig() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(dir, ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("config file not found")
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
