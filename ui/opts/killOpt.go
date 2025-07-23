package opts

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/hassanaziz0012/go-htop/ui/state"
	"github.com/hassanaziz0012/go-htop/utils"
	"github.com/rivo/tview"
)

var KillSignalsKeys = []string{
	"SIGHUP",
	"SIGINT",
	"SIGQUIT",
	"SIGILL",
	"SIGTRAP",
	"SIGABRT",
	"SIGBUS",
	"SIGFPE",
	"SIGKILL",
	"SIGUSR1",
	"SIGSEGV",
	"SIGUSR2",
	"SIGPIPE",
	"SIGALRM",
	"SIGTERM",
}

var KillSignalsMap = map[string]syscall.Signal{
	"SIGHUP":  syscall.SIGHUP,
	"SIGINT":  syscall.SIGINT,
	"SIGQUIT": syscall.SIGQUIT,
	"SIGILL":  syscall.SIGILL,
	"SIGTRAP": syscall.SIGTRAP,
	"SIGABRT": syscall.SIGABRT,
	"SIGBUS":  syscall.SIGBUS,
	"SIGFPE":  syscall.SIGFPE,
	"SIGKILL": syscall.SIGKILL,
	"SIGUSR1": syscall.SIGUSR1,
	"SIGSEGV": syscall.SIGSEGV,
	"SIGUSR2": syscall.SIGUSR2,
	"SIGPIPE": syscall.SIGPIPE,
	"SIGALRM": syscall.SIGALRM,
	"SIGTERM": syscall.SIGTERM,
}

var KillSignals = utils.NewOrderedMap(KillSignalsKeys, KillSignalsMap)

func HandleKillOpt(state *state.AppState) {
	process := getProc(state.ProcessesTable)

	killSignalsTable := tview.NewTable().SetSelectable(true, false)

	killSignalsTable.SetSelectedFunc(func(row, column int) {
		cell := killSignalsTable.GetCell(row, column)
		for k, v := range KillSignals.Data {
			if strings.Contains(cell.Text, k) {
				syscall.Kill(process.PID, v)
				utils.WriteLog(fmt.Sprintf("Killed process %d with signal %s", process.PID, v.String()))
			}
		}
	})

	var i int
	for _, k := range KillSignals.Keys {
		cell := tview.NewTableCell(fmt.Sprintf("%d) %s", i, k))
		killSignalsTable.SetCell(i, 0, cell)
		i++
	}

	headerText := fmt.Sprintf("Kill process? %d - %s", process.PID, process.Command)
	modal := tview.NewFrame(killSignalsTable).
		AddText(headerText, true, tview.AlignLeft, tcell.ColorDefault)

	state.Pages.AddPage("modal", modal, true, true)
}
