//go:build linux

package cmd

import "os/exec"

func openDir(path string) *exec.Cmd {
	return exec.Command("xdg-open", path)
}
