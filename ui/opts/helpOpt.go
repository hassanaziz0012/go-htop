package opts

import (
	"strings"

	"github.com/hassanaziz0012/go-htop/ui/state"
	"github.com/rivo/tview"
)

func HandleHelpOpt(state *state.AppState) {
	helpText := tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)
	layout := tview.NewFlex().AddItem(helpText, 60, 0, true)

	frame := tview.NewFrame(layout)
	setHelpText(helpText)
	state.Pages.AddPage("help", frame, true, true)
}

func setHelpText(t *tview.TextView) {
	var b strings.Builder
	b.WriteString("[blue]go-htop[-] - A simple, htop clone that I made as a learning project to improve my skills in Go and learn some new stuff about how operating systems work under the hood.")
	b.WriteString("\n")
	b.WriteString("\n")

	b.WriteString("I used htop as inspiration and basically just copied all its features, implemented them in my own way in Go.")
	b.WriteString("\n")
	b.WriteString("\n")

	b.WriteString("[blue]On Niceness[-]\n")
	b.WriteString("You can only increase the niceness of processes you own. And you can only decrease niceness of any process (including your own) if you are a ROOT (sudo) user. If you don't have the necessary permissions, Nice+ and Nice- will silently fail.")

	b.WriteString("[green]Press any key to return[-]")

	t.SetText(b.String())
}
