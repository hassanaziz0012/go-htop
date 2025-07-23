package opts

import (
	"strconv"

	"github.com/hassanaziz0012/go-htop/types"
)

func getProc(ptable *types.TableWithConfig) types.Process {
	prow, _ := ptable.Table.GetSelection()
	pid, _ := strconv.Atoi(ptable.Table.GetCell(prow, 0).Text)

	var process types.Process
	for _, p := range *ptable.Processes {
		if pid == p.PID {
			process = p
			break
		}
	}
	return process
}
