package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [point]",
	Short: "Adds the current working directory to your warp points",
	Long: `Adds the current working directory to your warp points.
If no point is specified, the current directory's name will be used.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Add a warp point
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var point string
		if len(args) > 0 {
			point = args[0]
		} else {
			point = filepath.Base(pwd)
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

		if _, err := f.WriteString(fmt.Sprintf("%s:%s\n", point, pwd)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Added warp point '%s' to '%s'\n", point, pwd)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
