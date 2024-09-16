//go:build windows
// +build windows

package tui

import (
	"os"

	"golang.org/x/sys/windows"
)

func enableVirtualTerminalProcessing() error {
	handle := windows.Handle(os.Stdout.Fd())

	var mode uint32
	err := windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return err
	}

	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	err = windows.SetConsoleMode(handle, mode)
	if err != nil {
		return err
	}

	return nil
}

func putTerminalIntoRawMode() (interface{}, error) {
	return nil, nil
}

func restoreTerminal(termios interface{}) {
}
