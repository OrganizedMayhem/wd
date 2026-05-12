package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <point>",
	Short: "Open the warp point in the default file explorer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pointToWarp := args[0]

		configFile, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		configFile += "/.warprc"

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
				path := parts[1]
				cmd := openDir(path)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					os.Exit(1)
				}
				return
			}
		}

		fmt.Fprintf(os.Stderr, "Warp point '%s' not found.\n", pointToWarp)
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
