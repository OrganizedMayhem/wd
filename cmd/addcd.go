package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var addcdCmd = &cobra.Command{
	Use:   "addcd <path> [point]",
	Short: "Adds a path to your warp points",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("resolve %s: %w", args[0], err)
		}
		point := filepath.Base(path)
		if len(args) > 1 {
			point = args[1]
		}

		return addPoint(cmd, point, path)
	},
}

func init() {
	rootCmd.AddCommand(addcdCmd)
}
