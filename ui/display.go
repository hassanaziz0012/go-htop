package ui

import (
	"context"
	"fmt"
	"os/user"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/hassanaziz0012/go-htop/types"
	"github.com/hassanaziz0012/go-htop/ui/processes"
	"github.com/hassanaziz0012/go-htop/ui/state"
	"github.com/hassanaziz0012/go-htop/utils"
	"github.com/rivo/tview"
)

func RunApplication(plist chan []types.Process, metadataChan chan types.Metadata, memStatusChan chan types.MemoryStatus, cpuStatusChan chan types.CPUStatus, cancel context.CancelFunc) {
	app := tview.NewApplication()
	table := processes.CreateProcessesTable()
	optsTable := createOptsTable()

	text := tview.NewTextView().
		ScrollToBeginning().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorLightSkyBlue)

	meminfo := tview.NewTextView().
		ScrollToBeginning().
		SetDynamicColors(true)

	cpuinfo := tview.NewTextView().
		ScrollToBeginning().
		SetDynamicColors(true)

	memCpuLayout := tview.NewFlex().
		SetDirection(0).
		AddItem(cpuinfo, 4, 1, false).
		AddItem(meminfo, 2, 1, false)

	overviewLayout := tview.NewFlex().
		SetDirection(1).
		AddItem(memCpuLayout, 0, 1, false).
		AddItem(text, 0, 1, false)

	rootLayout := tview.NewFlex().
		SetDirection(0).
		AddItem(overviewLayout, 7, 0, false).
		AddItem(table.Table, 0, 1, true).
		AddItem(optsTable, 1, 0, false)

	pages := tview.NewPages().
		AddPage("main", rootLayout, true, true)

	currentUser, err := user.Current()
	if err != nil {
		panic("failed to get current user")
	}
	appstate := state.NewAppState(
		currentUser,
		app,
		pages,
		table,
		optsTable,
		rootLayout,
		state.ListLayout,
		cancel,
	)

	ConfigureOpts(appstate)

	go func() {
		for {
			select {
			case ps := <-plist:
				app.QueueUpdateDraw(func() {
					allCopy := make([]types.Process, len(ps))
					copy(allCopy, ps)
					table.AllProcesses = &allCopy
					table.Processes = &ps

					processes.RenderProcesses(appstate)
				})
			case metadata := <-metadataChan:
				app.QueueUpdateDraw(func() {
					drawMetadata(text, &metadata)
				})
			case memstatus := <-memStatusChan:
				drawMemInfo(meminfo, &memstatus)

				app.QueueUpdateDraw(func() {
					meminfo.Clear()
				})
			case cpuStatus := <-cpuStatusChan:
				drawCPUCores(cpuinfo, &cpuStatus)

				app.QueueUpdateDraw(func() {
					cpuinfo.Clear()
				})
			}
		}
	}()

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		utils.WriteLog(fmt.Sprintf("%v", err))
	}
}

func createOptsTable() *tview.Table {
	table := tview.NewTable().SetSelectable(false, true).SetFixed(1, 10)

	optsCols := []string{
		"F1",
		"Help",
		"F2",
		"Setup",
		"F3",
		"Search",
		"F4",
		"Filter",
		"F5",
		"Tree",
		"F6",
		"SortBy",
		"F7",
		"Nice -",
		"F8",
		"Nice +",
		"F9",
		"Kill",
		"F10",
		"Quit",
	}

	for i, col := range optsCols {
		cell := tview.NewTableCell(col)

		if i%2 == 1 {
			cell.SetBackgroundColor(tcell.ColorLightSkyBlue)
		}

		table.SetCell(0, i, cell)
	}

	return table
}

func drawMetadata(text *tview.TextView, md *types.Metadata) {
	msg := fmt.Sprintf("Tasks: [green]%d[-], [green]%d[-] thr, [green]%d[-] kthr; [green]%d[-] running\nLoad average: [white]%s[-]\nUptime: [green]%s[-]", md.Tasks, md.Threads, md.KThreads, md.Running, md.LoadAvgs, md.Uptime)
	text.SetText(msg)
}

func drawMemInfo(text *tview.TextView, memstat *types.MemoryStatus) {
	text.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		barchar := "|"

		barwidth := max(width-20, 20)

		memPercentage := int((memstat.MemUsed * 100) / memstat.MemTotal)
		swpPercentage := int((memstat.SwpUsed * 100) / memstat.SwpTotal)

		memBar := drawBar(memPercentage, barwidth, barchar)
		swpBar := drawBar(swpPercentage, barwidth, barchar)

		mem := fmt.Sprintf("[blue]Mem[-][%s [gray]%.2fG/%.2fG[-]]", memBar, utils.KBtoGB(float64(memstat.MemUsed)), utils.KBtoGB(float64(memstat.MemTotal)))
		swp := fmt.Sprintf("[blue]Swp[-][%s [gray]%.2fG/%.2fG[-]]", swpBar, utils.KBtoGB(float64(memstat.SwpUsed)), utils.KBtoGB(float64(memstat.SwpTotal)))
		text.SetText(fmt.Sprintf("%s\n%s", mem, swp))

		return x, y, width, height
	})
}

func drawCPUCores(text *tview.TextView, cpustat *types.CPUStatus) {
	text.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		barchar := "|"

		var msg string
		for _, core := range cpustat.Cores {
			var barwidth int
			if core.UsagePercent > 10 {
				barwidth = max(width-13, 13)
			} else {
				barwidth = max(width-12, 12)
			}
			bar := drawBar(int(core.UsagePercent), barwidth, barchar)
			coreMsg := fmt.Sprintf("[blue]%d[-][%s [gray]%.2f%%[-]]", core.Num, bar, core.UsagePercent)
			msg += coreMsg + "\n"
		}

		text.SetText(msg)

		return x, y, width, height
	})
}

func drawBar(usedPercent int, width int, char string) string {
	usedBars := (usedPercent * width) / 100
	emptyBars := width - usedBars

	return fmt.Sprintf("[green]%s[transparent]%s[-]", strings.Repeat(char, usedBars), strings.Repeat(" ", emptyBars))
}
