package cmd

import (
	"fmt"
	"strings"

	"github.com/aryansharma9917/codewise-cli/pkg/doctor"
	"github.com/spf13/cobra"
)

var (
	doctorCluster bool
	doctorContext string
	doctorRequire string
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check local tools and Kubernetes connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		required := []string{}
		if doctorRequire != "" {
			required = strings.Split(doctorRequire, ",")
		}
		results, err := doctor.Run(doctor.Options{
			Cluster: doctorCluster,
			Context: doctorContext,
			Require: required,
		})

		fmt.Println("Codewise Doctor")
		fmt.Println("---------------")
		for _, result := range results {
			fmt.Printf("%-8s %-10s %s\n", result.Status, result.Name, result.Message)
		}
		return err
	},
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorCluster, "cluster", false, "Check Kubernetes API connectivity")
	doctorCmd.Flags().StringVar(&doctorContext, "context", "", "Kubernetes context for the cluster check")
	doctorCmd.Flags().StringVar(&doctorRequire, "require", "", "Comma-separated tools that must pass (docker,kubectl,helm,git)")
	rootCmd.AddCommand(doctorCmd)
}
