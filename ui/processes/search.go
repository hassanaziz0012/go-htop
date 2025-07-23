package processes

import (
	"strconv"
	"strings"

	"github.com/hassanaziz0012/go-htop/types"
)

func searchProcesses(table *types.TableWithConfig) {
	query := strings.ToLower(table.Search)

	for i := range *table.Processes {
		p := &(*table.Processes)[i]
		if strings.Contains(strconv.Itoa(p.PID), query) ||
			strings.Contains(strings.ToLower(p.Command), query) ||
			strings.Contains(strings.ToLower(p.User), query) {
			p.Highlight = true
		} else {
			p.Highlight = false
		}
	}
}
