//go:build windows

package cmd

import (
	"os/exec"
	"syscall"
)

// openDir launches the folder through `start`. The command line is built by
// hand so the path stays inside double quotes: Go only quotes arguments that
// contain spaces, and cmd.exe would otherwise treat characters like & or | in
// a path as command separators. Windows paths cannot contain double quotes.
func openDir(path string) *exec.Cmd {
	command := exec.Command("cmd.exe")
	command.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: `cmd.exe /d /c start "" "` + path + `"`,
	}
	return command
}
