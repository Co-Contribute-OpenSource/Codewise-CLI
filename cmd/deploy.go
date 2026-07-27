package cmd

import (
	"github.com/aryansharma9917/codewise-cli/pkg/deploy"
	"github.com/spf13/cobra"
)

var (
	deployEnv string
	dryRun    bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deployment operations",
}

var deployRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if deployEnv == "" {
			return ExitError("please provide --env")
		}
		return deploy.Run(deployEnv, dryRun)
	},
}

var deployPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Preview deployment execution plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		if deployEnv == "" {
			return ExitError("please provide --env")
		}
		return deploy.Plan(deployEnv)
	},
}

var deployExplainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explain strategy selection, execution, and recovery",
	RunE: func(cmd *cobra.Command, args []string) error {
		if deployEnv == "" {
			return ExitError("please provide --env")
		}
		return deploy.Explain(deployEnv)
	},
}

func init() {

	deployRunCmd.Flags().StringVar(&deployEnv, "env", "", "Environment to deploy")
	deployRunCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview deployment")

	deployPlanCmd.Flags().StringVar(&deployEnv, "env", "", "Environment to plan")
	deployExplainCmd.Flags().StringVar(&deployEnv, "env", "", "Environment to explain")

	deployCmd.AddCommand(deployRunCmd)
	deployCmd.AddCommand(deployPlanCmd)
	deployCmd.AddCommand(deployExplainCmd)

	rootCmd.AddCommand(deployCmd)
}
