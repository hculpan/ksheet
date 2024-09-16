package sheet

import "strconv"

const (
	CELL_TYPE_INT = iota
	CELL_TYPE_NULL
	CELL_TYPE_STRING
	CELL_TYPE_UNCHANGED
)

var CellDataNull CellData = CellData{data: nil, dataType: CELL_TYPE_NULL}

type CellDataType int

type CellData struct {
	data     any
	dataType CellDataType
}

func NewCellData(dataType CellDataType, data any) *CellData {
	result := &CellData{
		data:     data,
		dataType: dataType,
	}

	return result
}

func (c *CellData) DataDisplay() string {
	switch c.dataType {
	case CELL_TYPE_INT:
		return strconv.Itoa(c.data.(int))
	case CELL_TYPE_STRING:
		return c.data.(string)
	default:
		return ""
	}
}

func (c *CellData) Data() any {
	return c.data
}

func (c *CellData) DataAsInt() int {
	return c.data.(int)
}

func (c *CellData) DataType() CellDataType {
	return c.dataType
}
