// Package cmd implements the command-line interface for wd.
package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "wd",
	Short:             "Warp to custom directories in terminal",
	Long:              `wd (warp directory) is a tool that lets you jump to custom directories in the terminal`,
	Args:              cobra.ArbitraryArgs,
	SilenceErrors:     true,
	SilenceUsage:      true,
	ValidArgsFunction: completeWarpPoint,
	RunE: func(cmd *cobra.Command, args []string) error {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Fprintf(cmd.OutOrStdout(), "wd version %s\n", version)
			return nil
		}
		if len(args) == 0 {
			return cmd.Help()
		}

		path, err := getWarpPoint(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), path)
		return nil
	},
}

func getWarpPoint(name string) (string, error) {
	store, err := newStore()
	if err != nil {
		return "", err
	}
	path, err := store.Get(name)
	if errors.Is(err, fs.ErrNotExist) {
		return "", errors.New("no warp points yet")
	}
	return path, err
}

// passthroughArgs lists every first argument the shell wrappers must hand
// straight to the binary instead of treating as a warp point: all subcommands
// (including cobra's help and completion commands), their aliases, the root
// flags, and the hidden requests shell completion scripts make, which cobra
// only registers when they are the command being run.
func passthroughArgs() []string {
	args := []string{"--help", "-h", "--version", "-v", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd}
	for _, command := range rootCmd.Commands() {
		args = append(args, command.Name())
		args = append(args, command.Aliases...)
	}
	return args
}

// addPoint stores a warp point after checking that the shell wrapper will be
// able to jump to it.
func addPoint(cmd *cobra.Command, point, path string) error {
	// The PowerShell wrapper matches passthrough names without regard to case.
	isReserved := slices.ContainsFunc(passthroughArgs(), func(arg string) bool {
		return strings.EqualFold(arg, point)
	})
	if strings.HasPrefix(point, "-") || isReserved {
		return fmt.Errorf("warp point name %q is reserved by a wd command; choose another name", point)
	}

	store, err := newStore()
	if err != nil {
		return err
	}
	if err := store.Put(point, path); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Added warp point '%s' to '%s'\n", point, path)
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "print version")
}
