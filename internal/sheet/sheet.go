package sheet

import (
	"strconv"
	"unicode"
)

const (
	MAX_COLS = 999
	MAX_ROWS = 9999
)

type Sheet struct {
	data map[int]*CellData
}

func NewSheet() *Sheet {
	result := &Sheet{
		data: map[int]*CellData{},
	}

	for i := range MAX_ROWS {
		result.SetData(i, 0, CELL_TYPE_INT, i)
	}

	for i := range MAX_COLS {
		result.SetData(0, i, CELL_TYPE_STRING, result.GetColumnLabel(i))
	}

	return result
}

func (s *Sheet) GetColumnLabel(col int) string {
	var result byte = byte(col + 64)
	return string(result)
}

func (s *Sheet) SetDataByCellAddress(cell string, dataType CellDataType, value any) error {
	col, row, err := cellToColRow(cell)
	if err != nil {
		return err
	}

	s.SetData(col, row, dataType, value)
	return nil
}

func (s *Sheet) SetData(col, row int, dataType CellDataType, value any) {
	cell, ok := s.data[col*MAX_COLS+row]

	if dataType == CELL_TYPE_NULL && ok {
		delete(s.data, col*MAX_COLS+row)
		return
	} else if !ok {
		cell = NewCellData(dataType, value)
		s.data[col*MAX_COLS+row] = cell
	}

	switch v := value.(type) {
	case int:
		if dataType == CELL_TYPE_STRING {
			value = strconv.Itoa(v)
		}
	case string:
		if dataType == CELL_TYPE_INT {
			n, err := strconv.Atoi(v)
			if err != nil {
				value = -1
			} else {
				value = n
			}
		}
	}

	cell.data = value
	cell.dataType = dataType
}

func (s *Sheet) HasData(col, row int) bool {
	return s.GetData(col, row) != &CellDataNull && s.GetData(col, row).data != nil
}

func (s *Sheet) RemoveData(col, row int) {
	s.SetData(col, row, CELL_TYPE_NULL, "")
}

func (s *Sheet) SetDataString(col int, row int, data string) error {
	if data == "" {
		s.RemoveData(col, row)
	} else if data[0] == '=' {

	} else {
		n, err := strconv.Atoi(data)
		if err != nil {
			s.SetData(col, row, CELL_TYPE_STRING, data)
		} else {
			s.SetData(col, row, CELL_TYPE_INT, n)
		}
	}

	return nil
}

func (s *Sheet) GetData(col, row int) *CellData {
	if result, ok := s.data[col*MAX_COLS+row]; !ok {
		return &CellDataNull
	} else {
		return result
	}
}

func (s *Sheet) GetDataByCellAddress(cell string) (*CellData, error) {
	col, row, err := cellToColRow(cell)
	if err != nil {
		return nil, err
	}

	return s.GetData(col, row), nil
}

func cellToColRow(cell string) (int, int, error) {
	letters := ""
	digits := ""

	// Separate letters and digits
	for _, char := range cell {
		if unicode.IsLetter(char) {
			letters += string(char)
		} else if unicode.IsDigit(char) {
			digits += string(char)
		}
	}

	// Convert the letters to a column number
	col := 0
	for i := 0; i < len(letters); i++ {
		col = col*26 + int(letters[i]-'A'+1)
	}

	// Convert the digits to a row number
	row, err := strconv.Atoi(digits)
	if err != nil {
		return 0, 0, err
	}

	return col, row, nil
}
