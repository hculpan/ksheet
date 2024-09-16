package main

import (
	"fmt"
	"ksheet/internal/sheet"
	"ksheet/internal/tui"
)

func main() {
	s := sheet.NewSheet()
	/*	s.SetDataByCellAddress("A1", sheet.CELL_TYPE_INT, 10)
		s.SetDataByCellAddress("B1", sheet.CELL_TYPE_STRING, "hello")
		s.SetDataByCellAddress("C1", sheet.CELL_TYPE_INT, 110)*/
	err := tui.Start(*s)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
