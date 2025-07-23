package processes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/hassanaziz0012/go-htop/types"
	"github.com/hassanaziz0012/go-htop/utils"
	"github.com/rivo/tview"
)

func HandleSort(cell *tview.TableCell, table *types.TableWithConfig, swapDirection bool) {
	removeSortSymbol(cell)
	sortSymbol := getSortSymbol(table, swapDirection)
	table.Sort = cell.Text

	for k := range types.ColsMap {
		table.Table.GetCell(0, k).SetBackgroundColor(tcell.ColorGreen)
		removeSortSymbol(cell)
	}

	cell.Text += fmt.Sprintf("  %s", sortSymbol)
	cell.SetBackgroundColor(tcell.ColorLightSkyBlue)

	utils.WriteLog(cell.Text)
}

func getSortSymbol(table *types.TableWithConfig, swapDirection bool) string {
	var sortSymbol string

	if !swapDirection {
		switch table.SortDirection {
		case types.SortAsc:
			return "v"
		case types.SortDsc:
			return "^"
		}
	}

	if table.SortDirection == types.SortAsc {
		table.SortDirection = types.SortDsc
		sortSymbol = "^"
	} else {
		table.SortDirection = types.SortAsc
		sortSymbol = "v"
	}
	return sortSymbol
}

func removeSortSymbol(cell *tview.TableCell) {
	cell.Text = strings.ReplaceAll(cell.Text, "v", "")
	cell.Text = strings.ReplaceAll(cell.Text, "^", "")
	cell.Text = strings.TrimSpace(cell.Text)
}

func sortProcesses(table *types.TableWithConfig, processes *[]types.Process) {
	ascending := table.SortDirection == types.SortAsc

	sort.Slice(*processes, func(i, j int) bool {
		a := (*processes)[i]
		b := (*processes)[j]

		switch table.Sort {
		case "PID":
			if ascending {
				return a.PID < b.PID
			}
			return a.PID > b.PID
		case "U":
			if ascending {
				return strings.ToLower(a.User) < strings.ToLower(b.User)
			}
			return strings.ToLower(a.User) > strings.ToLower(b.User)
		case "PI":
			if ascending {
				return a.Priority < b.Priority
			}
			return a.Priority > b.Priority
		case "NI":
			if ascending {
				return a.Nice < b.Nice
			}
			return a.Nice > b.Nice
		case "S":
			return true
		case "CPU":
			if ascending {
				return a.CPU < b.CPU
			}
			return a.CPU > b.CPU
		case "VMEM":
			if ascending {
				return a.VirtualMemory < b.VirtualMemory
			}
			return a.VirtualMemory > b.VirtualMemory
		case "SMEM":
			if ascending {
				return a.SharedMemory < b.SharedMemory
			}
			return a.SharedMemory > b.SharedMemory
		case "MEM":
			if ascending {
				return a.Memory < b.Memory
			}
			return a.Memory > b.Memory
		case "TIME":
			if ascending {
				return a.Time < b.Time
			}
			return a.Time > b.Time
		case "CMD":
			if ascending {
				return strings.ToLower(a.Command) < strings.ToLower(b.Command)
			}
			return strings.ToLower(a.Command) > strings.ToLower(b.Command)
		default:
			return true
		}
	})
}
