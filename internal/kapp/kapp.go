package kapp

import "ksheet/internal/sheet"

var Kapp KSheets

const (
	MAIN_WIDTH  int = 2700
	MAIN_HEIGHT int = 1535
)

type KSheets struct {
	Sheet       sheet.Sheet
	DisplayRows int
	DisplayCols int
}

func init() {
	Kapp = KSheets{
		Sheet: *sheet.NewSheet(),
	}
}
