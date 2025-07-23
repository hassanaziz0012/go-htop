package opts

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/hassanaziz0012/go-htop/ui/state"
	"github.com/hassanaziz0012/go-htop/utils"
	"github.com/rivo/tview"
)

var filterPageName = "search"

func HandleFilterOpt(state *state.AppState) {
	controls := createFilterControls(state)
	input := tview.NewInputField().
		SetLabel("Filter: ").
		SetText(state.ProcessesTable.Filter).
		SetChangedFunc(func(text string) {
			state.ProcessesTable.Filter = text

			optsTable := state.OptsTable
			for i := 0; i < optsTable.GetColumnCount(); i++ {
				cell := optsTable.GetCell(0, i)
				utils.WriteLog(cell.Text)
				if strings.ToLower(cell.Text) == "filter" {
					if text == "" {
						cell.SetText("Filter")
					} else {
						cell.SetText("FILTER")
					}
				}
			}
		})

	filterBar := tview.NewFlex().
		AddItem(controls, 60, 0, false).
		AddItem(input, 0, 1, true)

	layout := tview.NewFlex().
		SetDirection(0).
		AddItem(state.RootLayout, 0, 1, false).
		AddItem(filterBar, 1, 0, true)

	page := state.Pages.AddPage(filterPageName, layout, true, true)

	state.App.SetRoot(page, true)
}

func resetFilterOptName(optsTable *tview.Table) {
	for i := 0; i < optsTable.GetColumnCount(); i++ {
		cell := optsTable.GetCell(0, i)
		utils.WriteLog(cell.Text)
		if strings.ToLower(cell.Text) == "filter" {
			cell.SetText("Filter")
		}
	}
}

func createFilterControls(appState *state.AppState) *tview.Flex {
	DoneHandler := NewButtonHandler("Done").
		SetHandler(func(state *state.AppState) {
			state.Pages.RemovePage(filterPageName)
		})
	CancelHandler := NewButtonHandler("Cancel").
		SetHandler(func(state *state.AppState) {
			state.ProcessesTable.Filter = ""
			resetFilterOptName(state.OptsTable)
			state.Pages.RemovePage(filterPageName)
		})

	controls := tview.NewFlex()
	optsKeys := []string{
		"Enter",
		"Esc",
	}
	opts := map[string]*ButtonHandler{
		"Enter": DoneHandler,
		"Esc":   CancelHandler,
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
				v.Handler(appState)
			}).
			SetStyle(btnStyle)
		controls.AddItem(val, 0, 1, false)
	}

	return controls
}
