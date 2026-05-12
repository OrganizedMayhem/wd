package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addcdCmd = &cobra.Command{
	Use:   "addcd <path> [point]",
	Short: "Adds a path to your warp points",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		var point string
		if len(args) > 1 {
			point = args[1]
		} else {
			point = filepath.Base(path)
		}

		configFile, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configFile += "/.warprc"

		f, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		if _, err := f.WriteString(fmt.Sprintf("%s:%s\n", point, path)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Added warp point '%s' to '%s'\n", point, path)
	},
}

func init() {
	rootCmd.AddCommand(addcdCmd)
}
