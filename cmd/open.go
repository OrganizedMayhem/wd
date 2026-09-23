package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <point>",
	Short: "Open the warp point in the default file explorer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := getWarpPoint(args[0])
		if err != nil {
			return err
		}
		if err := openDir(cmd, path); err != nil {
			return fmt.Errorf("open %s: %w", path, err)
		}
		return nil
	},
}

// runOpener runs an external opener with its output sent to the command's
// output streams.
func runOpener(cmd *cobra.Command, name string, args ...string) error {
	open := exec.Command(name, args...)
	open.Stdout = cmd.OutOrStdout()
	open.Stderr = cmd.ErrOrStderr()
	return open.Run()
}

func init() {
	rootCmd.AddCommand(openCmd)
}
