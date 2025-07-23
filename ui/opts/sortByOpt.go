package opts

import (
	"github.com/gdamore/tcell/v2"
	"github.com/hassanaziz0012/go-htop/types"
	"github.com/hassanaziz0012/go-htop/ui/processes"
	"github.com/hassanaziz0012/go-htop/ui/state"
	"github.com/hassanaziz0012/go-htop/utils"
	"github.com/rivo/tview"
)

func HandleSortByOpt(state *state.AppState) {
	table := tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorLightSkyBlue))

	table.SetSelectedFunc(func(row, column int) {
		ptableCell := state.ProcessesTable.Table.GetCell(0, row-1)
		utils.WriteLog("Sorting by: " + ptableCell.Text)
		processes.HandleSort(ptableCell, state.ProcessesTable, true)
		state.App.SetRoot(state.Pages, true)
	})

	headerCell := tview.NewTableCell("Sort by").SetBackgroundColor(tcell.ColorGreen).SetSelectable(false)
	table.SetCell(0, 0, headerCell)

	for colIndex, colName := range types.ColsMap {
		cell := tview.NewTableCell(colName)
		table.SetCell(colIndex+1, 0, cell)
	}

	offsetBox := tview.NewBox()

	subLayout := tview.NewFlex().
		SetDirection(0).
		AddItem(offsetBox, 6, 0, false).
		AddItem(table, 0, 1, true)

	// add table to root layout
	layout := tview.NewFlex().
		AddItem(subLayout, 12, 0, true).
		AddItem(state.RootLayout, 0, 1, false)

	state.App.SetRoot(layout, true)
}
