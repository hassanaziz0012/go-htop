package opts

import (
	"github.com/hassanaziz0012/go-htop/ui/state"
	"github.com/hassanaziz0012/go-htop/utils"
)

func HandleTreeOpt(appstate *state.AppState) {
	if appstate.Layout == state.ListLayout {
		utils.WriteLog("starting tree layout")
		appstate.Layout = state.TreeLayout
	} else {
		utils.WriteLog("stopping tree layout")
		appstate.Layout = state.ListLayout
	}
}
