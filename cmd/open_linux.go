//go:build linux

package cmd

import "github.com/spf13/cobra"

func openDir(cmd *cobra.Command, path string) error {
	return runOpener(cmd, "xdg-open", path)
}
