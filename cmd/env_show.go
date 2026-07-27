package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/aryansharma9917/codewise-cli/pkg/env"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var envShowFormat string

type environmentView struct {
	Name   string           `json:"name" yaml:"name"`
	K8s    env.K8sConfig    `json:"k8s" yaml:"k8s"`
	Helm   env.HelmConfig   `json:"helm" yaml:"helm"`
	GitOps env.GitOpsConfig `json:"gitops" yaml:"gitops"`
	Values env.ValuesConfig `json:"values" yaml:"values"`
}

var envShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show the complete effective environment configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		loaded, err := env.LoadEnv(args[0])
		if err != nil {
			return err
		}
		view := environmentView{
			Name: loaded.Name, K8s: loaded.K8s, Helm: loaded.Helm,
			GitOps: loaded.GitOps, Values: loaded.Values,
		}
		var output []byte
		switch envShowFormat {
		case "yaml":
			output, err = yaml.Marshal(view)
		case "json":
			output, err = json.MarshalIndent(view, "", "  ")
		default:
			return fmt.Errorf("unsupported format %q: use yaml or json", envShowFormat)
		}
		if err != nil {
			return fmt.Errorf("encode environment: %w", err)
		}
		fmt.Println(string(output))
		return nil
	},
}

func init() {
	envShowCmd.Flags().StringVar(&envShowFormat, "format", "yaml", "Output format: yaml or json")
	envCmd.AddCommand(envShowCmd)
}
