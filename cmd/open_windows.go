//go:build windows

package cmd

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/spf13/cobra"
)

var procShellExecuteW = syscall.NewLazyDLL("shell32.dll").NewProc("ShellExecuteW")

// openDir hands the folder straight to the Windows shell. Going through
// cmd.exe would expand %VAR% pairs in the path, and explorer.exe exits with
// status 1 even when it succeeds.
func openDir(_ *cobra.Command, path string) error {
	verb, err := syscall.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	file, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	const swShowNormal = 1
	ret, _, _ := procShellExecuteW.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), 0, 0, swShowNormal)
	// ShellExecuteW returns a value greater than 32 on success.
	if ret <= 32 {
		return fmt.Errorf("ShellExecute failed with code %d", ret)
	}
	return nil
}
