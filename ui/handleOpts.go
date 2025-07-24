package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/hassanaziz0012/go-htop/ui/opts"
	"github.com/hassanaziz0012/go-htop/ui/state"
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
	configureOptsShortcuts(state)
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
			opts.HandleQuitOpt(state)
		}
	})
}

func configureOptsShortcuts(state *state.AppState) {
	state.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			opts.HandleHelpOpt(state)
			return nil
		case tcell.KeyF2:
			return nil
		case tcell.KeyF3:
			opts.HandleSearchOpt(state)
			return nil
		case tcell.KeyF4:
			opts.HandleFilterOpt(state)
			return nil
		case tcell.KeyF5:
			opts.HandleTreeOpt(state)
			return nil
		case tcell.KeyF6:
			opts.HandleSortByOpt(state)
			return nil
		case tcell.KeyF7:
			opts.DecNice(state)
			return nil
		case tcell.KeyF8:
			opts.IncNice(state)
			return nil
		case tcell.KeyF9:
			opts.HandleKillOpt(state)
			return nil
		case tcell.KeyF10, tcell.KeyEscape:
			opts.HandleQuitOpt(state)
			return nil
		}

		if event.Rune() == 'q' {
			opts.HandleQuitOpt(state)
			return nil
		}

		return event
	})
}
