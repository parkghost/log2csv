// +build windows

package term

import (
	"os"
	"syscall"
)

func IsANSI(f *os.File) bool {
	return false
}

// IsTerminal returns false on Windows.
func IsTerminal(f *os.File) bool {
	ft, _ := syscall.GetFileType(syscall.Handle(f.Fd()))
	return ft == syscall.FILE_TYPE_CHAR
}
