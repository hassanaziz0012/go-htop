package processes

import (
	"strconv"
	"strings"

	"github.com/hassanaziz0012/go-htop/types"
)

func filterProcesses(table *types.TableWithConfig) {
	var filtered []types.Process

	query := strings.ToLower(table.Filter)

	for _, p := range *table.AllProcesses {
		if strings.Contains(strconv.Itoa(p.PID), query) ||
			strings.Contains(strings.ToLower(p.Command), query) ||
			strings.Contains(strings.ToLower(p.User), query) {
			filtered = append(filtered, p)
		}
	}

	table.Processes = &filtered
}
