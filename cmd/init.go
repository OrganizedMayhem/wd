package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Generate shell wrapper to allow changing directories",
	Long: `Outputs a shell script that wraps the wd-go binary.
Because a child process cannot change the parent shell's directory, you must use a shell wrapper to get the full "warp" functionality.

Add the following to your .bashrc or .zshrc:
eval "$(wd-go init bash)"
`,
	ValidArgs: []string{"bash", "zsh"},
	Run: func(cmd *cobra.Command, args []string) {
		script := `wd() {
    if [ $# -eq 0 ]; then
        wd-go
        return
    fi

    case "$1" in
        add|addcd|clean|list|ls|open|path|rm|show|version|init|help|--help|-h|--version|-v)
            wd-go "$@"
            ;;
        *)
            local target_path
            target_path=$(wd-go "$@")
            local exit_code=$?
            
            if [ $exit_code -eq 0 ] && [ -n "$target_path" ]; then
                cd "$target_path"
            elif [ -n "$target_path" ]; then
                echo "$target_path"
                return $exit_code
            else
                return $exit_code
            fi
            ;;
    esac
}`
		fmt.Println(script)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
