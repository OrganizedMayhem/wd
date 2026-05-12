package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.10.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wd-go",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("wd-go version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
