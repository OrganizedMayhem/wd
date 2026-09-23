package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// completeWarpPoint offers stored warp point names, described by their paths,
// as completions for a command's first argument.
func completeWarpPoint(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	store, err := newStore()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	points, err := store.Load()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []cobra.Completion
	for _, point := range points {
		if strings.HasPrefix(point.Name, toComplete) {
			completions = append(completions, cobra.CompletionWithDesc(point.Name, point.Path))
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
