//go:build !windows
// +build !windows

package tui

import (
	"os"
	"syscall"
	"unsafe"
)

func enableVirtualTerminalProcessing() error {
	return nil
}

func putTerminalIntoRawMode() (*syscall.Termios, error) {
	fd := os.Stdin.Fd()
	termios := &syscall.Termios{}

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(termios)), 0, 0, 0); err != 0 {
		return nil, err
	}

	raw := *termios
	raw.Lflag &^= syscall.ICANON | syscall.ECHO

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&raw)), 0, 0, 0); err != 0 {
		return nil, err
	}

	return termios, nil
}

func restoreTerminal(termios *syscall.Termios) {
	fd := os.Stdin.Fd()
	syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(termios)), 0, 0, 0)
}
