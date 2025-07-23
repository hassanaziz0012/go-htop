package opts

import (
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/hassanaziz0012/go-htop/ui/state"
)

func IncNice(state *state.AppState) {
	p := getProc(state.ProcessesTable)

	if p.User != state.User.Username {
		return // users can only increase niceness of processes they own
	}

	pid := strconv.Itoa(p.PID)
	newNice := strconv.Itoa(p.Nice + 1)
	_, err := exec.Command("renice", "-n", newNice, "-p", pid).CombinedOutput()

	if err != nil {
		log.Fatal("increment error:", err)
	}
}

func DecNice(state *state.AppState) {
	p := getProc(state.ProcessesTable)
	if os.Geteuid() != 0 {
		return // only root users can decrement niceness, so we gracefully fail all other users.
	}

	pid := strconv.Itoa(p.PID)
	newNice := strconv.Itoa(p.Nice - 1)
	_, err := exec.Command("renice", "-n", newNice, "-p", pid).CombinedOutput()
	if err != nil {
		log.Fatal("decrement error:", err)
	}
}
