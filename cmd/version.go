package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-latest"
)

var (
	checkLatest bool
	// These values are replaced by GoReleaser using -ldflags.
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

const (
	OWNER = "AryanSharma9917"
	REPO  = "Codewise-CLI"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "outputs the cli version",
	Run: func(cmd *cobra.Command, args []string) {

		if checkLatest {
			checkForNewVersion()
		} else {
			fmt.Printf("codewise %s\ncommit: %s\nbuilt: %s\n", Version, Commit, Date)
		}

	},
}

func checkForNewVersion() {

	githubTag := &latest.GithubTag{
		Owner:             OWNER,
		Repository:        REPO,
		FixVersionStrFunc: latest.DeleteFrontV(),
	}

	if Version == "dev" {
		fmt.Println("Development build; update checks require a released version")
		return
	}

	res, err := latest.Check(githubTag, Version)

	if err != nil {
		fmt.Println("Unable to check for latest version. Check your internet connection")
		return
	}

	if res.Outdated {
		fmt.Printf("The latest version of codewise-cli is %s.\nPlease update to the latest version by running go get -u github.com/aryansharma9917/codewise-cli@latest", res.Current)
		return
	}

	fmt.Println("You are using the latest version of codewise-cli")

}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Flags for the version command
	versionCmd.Flags().BoolVarP(&checkLatest, "latest", "l", false, "Check if the latest version is installed")
}
