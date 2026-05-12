package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all stored warp points",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configFile += "/.warprc"

		file, err := os.Open(configFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No warp points yet. Add one with 'wd add'.")
				return
			}
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
