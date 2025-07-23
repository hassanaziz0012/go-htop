package opts

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/hassanaziz0012/go-htop/ui/state"
	"github.com/hassanaziz0012/go-htop/utils"
	"github.com/rivo/tview"
)

var searchPageName = "search"

func HandleSearchOpt(state *state.AppState) {
	state.ProcessesTable.Search = ""

	controls := createSearchControls(state)
	input := tview.NewInputField().
		SetLabel("Search: ").
		SetText("").
		SetChangedFunc(func(text string) {
			state.ProcessesTable.Search = text
		})

	searchBar := tview.NewFlex().
		AddItem(controls, 60, 0, false).
		AddItem(input, 0, 1, true)

	layout := tview.NewFlex().
		SetDirection(0).
		AddItem(state.RootLayout, 0, 1, false).
		AddItem(searchBar, 1, 0, true)

	page := state.Pages.AddPage(searchPageName, layout, true, true)

	state.App.SetRoot(page, true)
}

func createSearchControls(appstate *state.AppState) *tview.Flex {
	NextHandler := NewButtonHandler("Next").
		SetHandler(func(state *state.AppState) {
			processes := state.ProcessesTable.Processes
			table := state.ProcessesTable.Table
			currentRow, _ := table.GetSelection()

			for i := currentRow + 1; i < table.GetRowCount(); i++ {
				pid, _ := strconv.Atoi(table.GetCell(i, 0).Text)

				for _, p := range *processes {
					if p.PID == pid {
						if p.Highlight {
							table.Select(i, 0)
							utils.WriteLog("search -> row " + strconv.Itoa(i) + p.User)
							return
						}
					}
				}
			}

			// loop back to first selection if nothing found
			table.Select(currentRow, 0)
		})
	PrevHandler := NewButtonHandler("Prev").
		SetHandler(func(state *state.AppState) {
			processes := state.ProcessesTable.Processes
			table := state.ProcessesTable.Table
			currentRow, _ := table.GetSelection()

			for i := currentRow - 1; i >= 0; i-- {
				pid, _ := strconv.Atoi(table.GetCell(i, 0).Text)

				for _, p := range *processes {
					if p.PID == pid {
						if p.Highlight {
							table.Select(i, 0)
							utils.WriteLog("search -> row " + strconv.Itoa(i) + p.User)
							return
						}
					}
				}
			}

			// loop back to first selection if nothing found
			table.Select(currentRow, 0)
		})
	CancelHandler := NewButtonHandler("Cancel").
		SetHandler(func(state *state.AppState) {
			state.ProcessesTable.Search = ""
			state.Pages.RemovePage(searchPageName)
		})

	controls := tview.NewFlex()
	optsKeys := []string{
		"N",
		"P",
		"Esc",
	}
	opts := map[string]*ButtonHandler{
		"N":   NextHandler,
		"P":   PrevHandler,
		"Esc": CancelHandler,
	}

	btnStyle := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorLightSkyBlue)

	for _, key := range optsKeys {
		label := tview.NewTextView().SetText(key)
		controls.AddItem(label, 0, 1, false)

		v := opts[key]

		val := tview.NewButton(v.Label).
			SetSelectedFunc(func() {
				v.Handler(appstate)
			}).
			SetStyle(btnStyle)
		controls.AddItem(val, 0, 1, false)
	}

	return controls
}
