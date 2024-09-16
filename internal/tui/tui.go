package tui

import (
	"fmt"
	"ksheet/internal/sheet"
	"os"
	"os/signal"
)

const (
	HEADER_ROWS = 4
)

var termHeight int
var termWidth int
var currSheet sheet.Sheet
var activeRow int = 1
var activeCol int = 1
var maxRow int = 0
var maxCol int = 0
var numCols int
var editMode bool = false
var editBuffer string = ""
var menuBar string = "  Ctl+S:Save  Ctl+A:Save as  Ctl+L:Load"
var statusBar string

func Start(s sheet.Sheet) error {
	if err := enableVirtualTerminalProcessing(); err != nil {
		return err
	}
	currSheet = s

	var err error
	termWidth, termHeight, err = getTerminalSize()
	if err != nil {
		return err
	}

	numCols = (termWidth - 4) / 9
	if numCols > 26 {
		numCols = 26
	}

	if err := drawSpreadsheet(); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("\nExiting...\n")
		os.Exit(0)
	}()

	getInput()

	return nil
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func drawCell(col int, row int, highlighted bool) {
	cell := currSheet.GetData(col, row)
	text := cell.DataDisplay()
	if cell.DataType() == sheet.CELL_TYPE_INT && len(text) > 9 {
		text = text[:8] + "X"
	}
	drawCellData(col, row, text, highlighted)
}

func drawCellData(col int, row int, data string, highlighted bool) {
	var howManyCells int = 1
	/*	if len(data) > 9 {
		howManyCells := len(data) / 9
		canDisplayAll := true
		for i := range howManyCells {
			if currSheet.HasData(col+i+1, row) || (activeCol == col+i+1 && activeRow == row) {
				canDisplayAll = false
				break
			}
		}

		if !canDisplayAll {
			data = data[:9]
		}
	}*/
	if len(data) > 9 {
		data = data[:9]
	}

	colPos := (col*(howManyCells*9) - 4)
	goTo(colPos, row)
	if highlighted {
		fmt.Printf("\033[%d;%dH", row+HEADER_ROWS, colPos)
		fmt.Printf("\033[30m\033[47m%9s\033[0m", data)
	} else {
		fmt.Printf("\033[%d;%dH%9s", row+HEADER_ROWS, colPos, data)
	}
}

func goHome() {
	goTo(1, 1)
}

func goReturn(row int) {
	goTo(1, row)
}

func goTo(col, row int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func colToLetters(col int) string {
	letters := ""
	for col > 0 {
		col--
		letters = string(rune('A'+col%26)) + letters
		col /= 26
	}
	return letters
}

func colRowToCell(col int, row int) string {
	letters := colToLetters(col)
	return fmt.Sprintf("%s%d", letters, row)
}

func getSelectedCellText() string {
	return colRowToCell(activeCol, activeRow)
}

func drawTopBar() {
	selectedText := getSelectedCellText() + " : " + currSheet.GetData(activeCol, activeRow).DataDisplay()
	goTo(1, 1)
	fmt.Printf("\033[30m\033[47m%s%*s\033[0m\n", selectedText, termWidth-len(selectedText), "")
}

func updateStatusBar(msg string) {
	statusBar = msg
	drawStatusBar()
}

func drawStatusBar() {
	goTo(1, termHeight-1)
	fmt.Printf("\033[30m\033[47m%s%*s\033[0m\n", statusBar, termWidth-len(statusBar), "")
}

func drawMenuBar() {
	goTo(1, 2)
	fmt.Printf("\033[30m\033[47m%s%*s\033[0m\n", menuBar, termWidth-len(menuBar), "")
}

func drawEditBar() {
	goTo(1, 3)
	fmt.Printf("%s%*s", editBuffer, termWidth-len(editBuffer), "")
	if !editMode {
		fmt.Print("\033[?25l") // Hide cursor
	} else {
		fmt.Print("\033[?25h") // Show cursor
		goTo(len(editBuffer)+1, 3)
	}
}

func drawColHeader() {
	hdr := []rune(fmt.Sprintf("\033[30m\033[47m%*s\033[0m", termWidth, ""))
	var c byte = 0
	for idx := 16; idx < len(hdr)-7 && c < byte(numCols); idx += 9 {
		hdr[idx+2] = rune(c + 65)
		c++
	}
	goTo(1, 4)
	fmt.Println(string(hdr))
}
func drawSpreadsheet() error {
	clearScreen()

	goHome()

	// Draw header
	drawTopBar()
	drawMenuBar()
	drawEditBar()

	drawColHeader()

	// Draw sheet grid
	for row := 1; row < termHeight-HEADER_ROWS; row++ {
		goReturn(row + HEADER_ROWS)
		for col := range numCols {
			if col == 0 {
				fmt.Printf("\033[30m\033[47m%4d\033[0m", row)
			} else {
				drawCell(col, row, activeCol == col && activeRow == row)
			}
		}
	}

	drawStatusBar()

	return nil
}

func startEditMode(buff string) {
	if activeCol > maxCol {
		maxCol = activeCol
	}

	if activeRow > maxRow {
		maxRow = activeRow
	}

	editMode = true
	editBuffer += buff
	drawEditBar()
}

func endEditMode() {
	if editMode && editBuffer != "" {
		currSheet.SetDataString(activeCol, activeRow, editBuffer)
		editMode = false
		editBuffer = ""
	}
	drawEditBar()
}
