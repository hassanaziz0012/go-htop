package state

import (
	"context"
	"os/user"

	"github.com/hassanaziz0012/go-htop/types"
	"github.com/rivo/tview"
)

type Layout int

const (
	TreeLayout Layout = iota
	ListLayout
)

type AppState struct {
	User            *user.User
	App             *tview.Application
	Pages           *tview.Pages
	ProcessesTable  *types.TableWithConfig
	OptsTable       *tview.Table
	RootLayout      *tview.Flex
	Layout          Layout
	ProcessesByRoot []*types.Process
	Cancel          context.CancelFunc
}

func NewAppState(currentUser *user.User, app *tview.Application, pages *tview.Pages, ptable *types.TableWithConfig, optsTable *tview.Table, rootLayout *tview.Flex, layout Layout, cancel context.CancelFunc) *AppState {
	return &AppState{
		currentUser,
		app,
		pages,
		ptable,
		optsTable,
		rootLayout,
		layout,
		nil,
		cancel,
	}
}
