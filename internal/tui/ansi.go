package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getTerminalSize() (int, int, error) {
	// Put the terminal into raw mode
	originalTermios, err := putTerminalIntoRawMode()
	if err != nil {
		return 0, 0, err
	}
	defer restoreTerminal(originalTermios)

	// Move cursor to bottom-right corner
	fmt.Print("\033[999C\033[999B")

	// Request cursor position
	fmt.Print("\033[6n")

	// Read the response from stdin
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('R')
	if err != nil {
		return 0, 0, err
	}

	// The response should be something like "\033[24;80R"
	if !strings.HasPrefix(response, "\033[") || !strings.HasSuffix(response, "R") {
		return 0, 0, fmt.Errorf("unexpected response: %s", response)
	}

	// Parse the response
	var rows, cols int
	_, err = fmt.Sscanf(response, "\033[%d;%dR", &rows, &cols)
	if err != nil {
		return 0, 0, err
	}

	return cols, rows, nil
}
