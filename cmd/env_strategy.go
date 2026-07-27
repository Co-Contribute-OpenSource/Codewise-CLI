package cmd

import (
	"fmt"

	"github.com/aryansharma9917/codewise-cli/pkg/env"
	"github.com/spf13/cobra"
)

var envStrategyCmd = &cobra.Command{
	Use:   "strategy <name> <auto|helm|kubectl|gitops>",
	Short: "Set the deployment strategy for an environment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := env.SetStrategy(args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("Environment %q strategy set to %s\n", args[0], args[1])
		return nil
	},
}

func init() {
	envCmd.AddCommand(envStrategyCmd)
}
