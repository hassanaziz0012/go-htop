package opts

import (
	"github.com/hassanaziz0012/go-htop/ui/state"
)

func HandleQuitOpt(state *state.AppState) {
	name, _ := state.Pages.GetFrontPage()
	if name != "main" {
		state.Pages.RemovePage(name)
		return
	}
	state.App.Stop()
	state.Cancel()
}
