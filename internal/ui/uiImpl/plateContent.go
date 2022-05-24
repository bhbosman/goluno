package uiImpl

import (
	"github.com/bhbosman/gocomms/connectionManager/model"
	"github.com/rivo/tview"
	"strconv"
)

type connectionPlateContent struct {
	Grid []model.LineData
}

func newConnectionPlateContent(Grid []model.LineData) *connectionPlateContent {
	return &connectionPlateContent{
		Grid: Grid,
	}
}

func (self *connectionPlateContent) GetCell(row, column int) *tview.TableCell {
	switch column {
	case 0:
		switch row {
		case 0:
			return tview.NewTableCell("*")
		default:

			return tview.NewTableCell("")
		}
	case 1:
		switch row {
		case 0:
			return tview.NewTableCell("Name")
		default:
			return tview.NewTableCell(self.Grid[row-1].InValue.Name)
		}
	case 2:
		switch row {
		case 0:
			return tview.NewTableCell("In(Other)")
		default:
			return tview.NewTableCell(strconv.Itoa(self.Grid[row-1].InValue.OtherMsgCount)).
				SetAlign(tview.AlignRight)
		}
	case 3:
		switch row {
		case 0:
			return tview.NewTableCell("In(RWS)")
		default:
			return tview.NewTableCell(strconv.Itoa(self.Grid[row-1].InValue.RwsBytesIn)).
				SetAlign(tview.AlignRight)
		}
	case 4:
		switch row {
		case 0:
			return tview.NewTableCell("In(Bytes)")
		default:
			return tview.NewTableCell(strconv.Itoa(self.Grid[row-1].InValue.RwsBytesIn)).
				SetAlign(tview.AlignRight)
		}
	case 5:
		switch row {
		case 0:
			return tview.NewTableCell("Out(Bytes)")
		default:
			return tview.NewTableCell(strconv.Itoa(self.Grid[row-1].InValue.RwsBytesOut)).
				SetAlign(tview.AlignRight)
		}
	case 6:
		switch row {
		case 0:
			return tview.NewTableCell("Name")
		default:
			return tview.NewTableCell(self.Grid[row-1].OutValue.Name)
		}
	case 7:
		switch row {
		case 0:
			return tview.NewTableCell("In(Other)")
		default:
			return tview.NewTableCell(strconv.Itoa(self.Grid[row-1].OutValue.OtherMsgCount)).
				SetAlign(tview.AlignRight)
		}
	case 8:
		switch row {
		case 0:
			return tview.NewTableCell("In(RWS)")
		default:
			return tview.NewTableCell(strconv.Itoa(self.Grid[row-1].OutValue.RwsBytesIn)).
				SetAlign(tview.AlignRight)
		}
	case 9:
		switch row {
		case 0:
			return tview.NewTableCell("In(Bytes)")
		default:
			return tview.NewTableCell(strconv.Itoa(self.Grid[row-1].OutValue.RwsBytesIn)).
				SetAlign(tview.AlignRight)
		}
	case 10:
		switch row {
		case 0:
			return tview.NewTableCell("Out(Bytes)")
		default:
			return tview.NewTableCell(strconv.Itoa(self.Grid[row-1].OutValue.RwsBytesOut)).
				SetAlign(tview.AlignRight)
		}

	}
	return tview.NewTableCell("")
}

func (self *connectionPlateContent) GetRowCount() int {
	return len(self.Grid) + 1
}

func (self *connectionPlateContent) GetColumnCount() int {
	return 11
}

func (self *connectionPlateContent) SetCell(row, column int, cell *tview.TableCell) {
}

func (self *connectionPlateContent) RemoveRow(row int) {
}

func (self *connectionPlateContent) RemoveColumn(column int) {
}

func (self *connectionPlateContent) InsertRow(row int) {
}

func (self *connectionPlateContent) InsertColumn(column int) {
}

func (self *connectionPlateContent) Clear() {
}
