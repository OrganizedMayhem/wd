package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:       "init <shell>",
	Short:     "Generate shell wrapper to allow changing directories",
	ValidArgs: []string{"bash", "zsh", "powershell"},
	Args:      cobra.ExactArgs(1),
	Long: `Outputs a shell script that wraps the wd binary.
Because a child process cannot change the parent shell's directory, you must use a shell wrapper to get the full "warp" functionality.

Add the following to your .bashrc or .zshrc:
eval "$(wd init bash)"

Add the following to your PowerShell profile ($PROFILE):
Invoke-Expression (@(wd init powershell) -join "` + "`n" + `")
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := args[0]
		passthrough := passthroughArgs()
		var script string
		switch shell {
		case "bash", "zsh":
			script = strings.Replace(posixScript, passthroughPlaceholder, strings.Join(passthrough, "|"), 1)
		case "powershell", "pwsh":
			script = strings.Replace(powershellScript, passthroughPlaceholder, "'"+strings.Join(passthrough, "', '")+"'", 1)
		default:
			return fmt.Errorf("unsupported shell %q; supported shells are bash, zsh, and powershell", shell)
		}
		fmt.Fprintln(cmd.OutOrStdout(), script)
		return nil
	},
}

// passthroughPlaceholder is replaced with the arguments the wrapper forwards
// to the binary unchanged, so the list never drifts from the registered commands.
const passthroughPlaceholder = "@PASSTHROUGH@"

const posixScript = `wd_cd() {
    local target_path
    target_path=$(command wd "$@")
    local exit_code=$?

    if [ $exit_code -eq 0 ] && [ -n "$target_path" ]; then
        cd "$target_path"
    elif [ -n "$target_path" ]; then
        echo "$target_path"
        return $exit_code
    else
        return $exit_code
    fi
}

wd() {
    if [ $# -eq 0 ]; then
        command wd
        return
    fi

    case "$1" in
        @PASSTHROUGH@)
            command wd "$@"
            ;;
        *)
            wd_cd "$@"
            ;;
    esac
}`

const powershellScript = `function wd {
    $wdBin = (Get-Command wd -CommandType Application -ErrorAction Stop | Select-Object -First 1).Source

    if ($args.Count -eq 0) {
        & $wdBin
        return
    }

    switch ($args[0]) {
        { $_ -in @PASSTHROUGH@ } {
            & $wdBin @args
            return
        }
    }

    $targetPath = & $wdBin @args
    $exitCode = $LASTEXITCODE

    if ($exitCode -eq 0 -and $targetPath) {
        Set-Location -LiteralPath $targetPath
    } elseif ($targetPath) {
        Write-Output $targetPath
        $global:LASTEXITCODE = $exitCode
    } else {
        $global:LASTEXITCODE = $exitCode
    }
}`

func init() {
	rootCmd.AddCommand(initCmd)
}
