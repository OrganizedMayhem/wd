//go:build darwin

package cmd

import "os/exec"

func openDir(path string) *exec.Cmd {
	return exec.Command("open", path)
}
