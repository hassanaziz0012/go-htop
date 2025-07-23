package processes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/hassanaziz0012/go-htop/procutils"
	"github.com/hassanaziz0012/go-htop/types"
	"github.com/hassanaziz0012/go-htop/ui/state"
	"github.com/rivo/tview"
)

func CreateProcessesTable() *types.TableWithConfig {
	table := tview.NewTable().SetSelectable(true, false).SetFixed(1, 0)

	tableWithConfig := types.TableWithConfig{
		Table: table,
	}

	addTableHeaders(&tableWithConfig)

	handleRowSelected(&tableWithConfig)

	return &tableWithConfig
}

func addTableHeaders(table *types.TableWithConfig) {
	// add table headers
	for colIndex, colName := range types.ColsMap {
		cell := tview.NewTableCell(colName).
			SetBackgroundColor(tcell.ColorGreen).
			SetSelectable(false)
		table.Table.SetCell(0, colIndex, cell)

		if table.Sort == cell.Text {
			HandleSort(cell, table, false)
		}
	}
}

func clearTable(table *types.TableWithConfig) {
	table.Table.Clear()
	addTableHeaders(table)
}

func RenderProcesses(appstate *state.AppState) {
	if appstate.Layout == state.ListLayout {
		RenderProcessesTable(appstate)
	} else {
		roots := procutils.AssignChildren(appstate.ProcessesTable.Processes)
		appstate.ProcessesByRoot = roots

		RenderProcessesTree(appstate)
	}
}

func RenderProcessesTable(state *state.AppState) *tview.Table {
	table := state.ProcessesTable

	if table.Filter != "" {
		filterProcesses(table)
	}
	if table.Search != "" {
		searchProcesses(table)
	}
	sortProcesses(table, table.Processes)

	curRow, curCol := table.Table.GetSelection()
	clearTable(table)
	table.Table.Select(curRow, curCol)
	curRow, _ = table.Table.GetSelection()

	for i, p := range *table.Processes {
		row := i + 1
		renderProcessToTable(state, p, row)
	}

	return table.Table
}

func RenderProcessesTree(state *state.AppState) *tview.Table {
	clearTable(state.ProcessesTable)

	row := 1 // row 0 is for table headers

	var render func(p *types.Process, prefix string, isLast bool)

	render = func(p *types.Process, prefix string, isLast bool) {
		var connector string
		if prefix == "" {
			connector = "" // root node
		} else if isLast {
			connector = "└── "
		} else {
			connector = "├── "
		}

		renderProcessToTable(state, *p, row)

		cmdCell := state.ProcessesTable.Table.GetCell(row, 10)
		cmdCell.SetText(fmt.Sprintf("[blue]%s%s[-] %s", prefix, connector, p.Command))

		row++

		var newPrefix string
		if isLast {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}

		for i, child := range p.Children {
			last := i == len(p.Children)-1
			render(child, newPrefix, last)
		}
	}

	for i, root := range state.ProcessesByRoot {
		last := i == len(state.ProcessesByRoot)-1
		render(root, "", last)
	}

	return state.ProcessesTable.Table
}

func renderProcessToTable(state *state.AppState, p types.Process, row int) {
	for colIndex, colName := range types.ColsMap {
		cell := tview.NewTableCell("")

		var value string
		switch colName {
		case "PID":
			value = strconv.Itoa(p.PID)
		case "U":
			value = p.User
			if p.User != state.User.Username {
				cell.SetTextColor(tcell.ColorGray)
			}
		case "PI":
			value = strconv.Itoa(p.Priority)
		case "NI":
			value = strconv.Itoa(p.Nice)
		case "S":
			value = p.Status
			switch value {
			case "S":
				cell.SetTextColor(tcell.ColorGray)
			case "R":
				cell.SetTextColor(tcell.ColorGreen)
			case "D":
				cell.SetTextColor(tcell.ColorRed)
			}
		case "CPU":
			parsed := strconv.FormatFloat(p.CPU, 'f', 2, 64)
			value = parsed

			if p.CPU == 0.0 {
				cell.SetTextColor(tcell.ColorGray)
			}
		case "VMEM":
			value = strconv.Itoa(p.VirtualMemory)
		case "SMEM":
			value = strconv.Itoa(p.SharedMemory)
		case "MEM":
			value = strconv.Itoa(p.Memory)
		case "TIME":
			value = parseProcessTime(p.Time)
		case "CMD":
			value = p.Command
		}

		cell.SetText(value)

		if p.Highlight {
			cell.SetBackgroundColor(tcell.ColorYellow)
		} else {
			cell.SetBackgroundColor(tcell.ColorBlack)
		}

		state.ProcessesTable.Table.SetCell(row, colIndex, cell)
	}
}

func handleRowSelected(table *types.TableWithConfig) {
	table.Table.SetSelectionChangedFunc(func(row, column int) {
		if row == 0 {
			cell := table.Table.GetCell(row, column)
			HandleSort(cell, table, true)
		}
	})
}

func parseProcessTime(d time.Duration) string {
	totalSecs := d.Seconds()
	mins := int(totalSecs) / 60
	secs := int(totalSecs) % 60
	centisecs := int(d.Milliseconds()/10) % 100

	return fmt.Sprintf("%02d:%02d.%02d", mins, secs, centisecs)
}
