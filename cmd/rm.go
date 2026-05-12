package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <point>",
	Short: "Removes the given warp point",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pointToRemove := args[0]

		configFile, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configFile += "/.warprc"

		file, err := os.Open(configFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No warp points to remove.")
				return
			}
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && parts[0] != pointToRemove {
				lines = append(lines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		output := strings.Join(lines, "\n")
		err = os.WriteFile(configFile, []byte(output), 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Removed warp point '%s'\n", pointToRemove)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
