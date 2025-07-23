package ui

import (
	"context"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/hassanaziz0012/go-htop/ui/opts"
	"github.com/hassanaziz0012/go-htop/ui/state"
	"github.com/rivo/tview"
)

var GlobalOpts = map[string]string{
	"F1":  "Help",
	"F2":  "Setup",
	"F3":  "Search",
	"F4":  "Filter",
	"F5":  "Tree",
	"F6":  "SortBy",
	"F7":  "Nice -",
	"F8":  "Nice +",
	"F9":  "Kill",
	"F10": "Quit",
}

func ConfigureOpts(state *state.AppState) {
	setOptionSelectFunc(state)
	configureOptsShortcuts(state.App, state.Pages, state.Cancel)
}

func setOptionSelectFunc(state *state.AppState) {
	state.OptsTable.SetSelectionChangedFunc(func(row, column int) {
		cell := state.OptsTable.GetCell(row, column)

		switch strings.ToLower(cell.Text) {
		case "help":
			opts.HandleHelpOpt(state)
		case "setup":
			return
		case "search":
			opts.HandleSearchOpt(state)
		case "filter":
			opts.HandleFilterOpt(state)
		case "tree":
			opts.HandleTreeOpt(state)
		case "sortby":
			opts.HandleSortByOpt(state)
		case "nice -":
			opts.DecNice(state)
		case "nice +":
			opts.IncNice(state)
		case "kill":
			opts.HandleKillOpt(state)
		case "quit":
			opts.HandleQuitOpt(state.App, state.Pages, state.Cancel)
		}
	})
}

func configureOptsShortcuts(app *tview.Application, pages *tview.Pages, cancel context.CancelFunc) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyF10 || event.Key() == tcell.KeyEscape {
			opts.HandleQuitOpt(app, pages, cancel)
			return nil
		}
		if event.Key() == tcell.KeyF9 {
			return nil
		}
		return event
	})
}
