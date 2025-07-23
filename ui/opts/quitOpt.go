package opts

import (
	"context"

	"github.com/rivo/tview"
)

func HandleQuitOpt(app *tview.Application, pages *tview.Pages, cancel context.CancelFunc) {
	name, _ := pages.GetFrontPage()
	if name != "main" {
		pages.RemovePage(name)
		return
	}
	app.Stop()
	cancel()
}
