// Package cmd implements the command line interface for wd-go.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wd",
	Short: "Warp to custom directories in terminal",
	Long:  `wd (warp directory) is a tool that lets you jump to custom directories in the terminal`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Printf("wd-go version %s\n", version)
			return
		}

		if len(args) == 0 {
			cmd.Help()
			return
		}

		pointToWarp := args[0]

		configFile, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		configFile += "$HOME/.warprc"

		file, err := os.Open(configFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintln(os.Stderr, "No warp points yet.")
				os.Exit(1)
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && parts[0] == pointToWarp {
				fmt.Println(parts[1])
				return
			}
		}

		fmt.Fprintf(os.Stderr, "Warp point '%s' not found.\n", pointToWarp)
		os.Exit(1)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "print version")
}
