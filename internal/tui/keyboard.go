package tui

import (
	"bufio"
	"fmt"
	"ksheet/internal/sheet"
	"os"
)

func getInput() error {
	originalTermios, err := putTerminalIntoRawMode()
	if err != nil {
		return err
	}
	defer restoreTerminal(originalTermios)

	reader := bufio.NewReader(os.Stdin)

	for {
		key, _ := reader.ReadByte()

		// If the key is an escape character, handle special keys
		if key == 27 {
			handleArrowKeys(reader)
		} else if key == 13 || key == 10 {
			handleReturn()
		} else if key == 8 || key == 127 || key == 126 {
			handleDeleteBackspace()
		} else if key >= 1 && key <= 26 {
			updateStatusBar(fmt.Sprintf("key = %d", key))
			switch key + 64 {
			case 'C':
				break
			case 'S':
				if maxCol == 0 && maxRow == 0 {
					updateStatusBar("Error: No data to save")
					break
				}
				updateStatusBar("Saving file...")
			}
		} else {
			startEditMode(string(rune(key)))
		}
	}

	return nil
}

func handleDeleteBackspace() {
	if len(editBuffer) > 0 {
		editBuffer = editBuffer[:len(editBuffer)-1]
		drawEditBar()
	} else {
		currSheet.RemoveData(activeCol, activeRow)
		drawCell(activeCol, activeRow, true)
	}
}

func handleReturn() {
	endEditMode()
	drawCell(activeCol, activeRow, false)
	activeRow++
	drawCell(activeCol, activeRow, true)
	drawTopBar()
}

func handleArrowKeys(reader *bufio.Reader) {
	prevActiveCol := activeCol
	prevActiveRow := activeRow
	// Escape sequence
	seq, _ := reader.Peek(2)
	if len(seq) == 2 && seq[0] == '[' {
		reader.Discard(2)
		switch seq[1] {
		case 'A':
			endEditMode()
			if activeRow > 1 {
				activeRow--
			}
		case 'B':
			endEditMode()
			if activeRow < sheet.MAX_ROWS {
				activeRow++
			}
		case 'C':
			endEditMode()
			if activeCol < sheet.MAX_COLS {
				activeCol++
			}
		case 'D':
			endEditMode()
			if activeCol > 1 {
				activeCol--
			}
		default:
			fmt.Printf("Unknown escape sequence: %v\n", seq)
		}
	}
	drawCell(prevActiveCol, prevActiveRow, false)
	drawCell(activeCol, activeRow, true)
	drawTopBar()

}
