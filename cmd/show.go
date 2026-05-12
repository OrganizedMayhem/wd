package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [point]",
	Short: "Print path to given warp point",
	Args:  cobra.MaximumNArgs(1),
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
				fmt.Println("No warp points yet.")
				return
			}
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		if len(args) > 0 {
			pointToShow := args[0]
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 && parts[0] == pointToShow {
					fmt.Println(parts[1])
					return
				}
			}
			fmt.Printf("Warp point '%s' not found.\n", pointToShow)
		} else {
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			var points []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 && parts[1] == pwd {
					points = append(points, parts[0])
				}
			}

			if len(points) > 0 {
				fmt.Printf("Warp points for current directory: %s\n", strings.Join(points, ", "))
			} else {
				fmt.Println("No warp points for current directory.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
