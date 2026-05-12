package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove points warping to nonexistent directories",
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
				fmt.Println("No warp points to clean.")
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
			if len(parts) == 2 {
				path := parts[1]
				if _, err := os.Stat(path); err == nil {
					lines = append(lines, line)
				}
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

		fmt.Println("Cleaned warp points.")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
