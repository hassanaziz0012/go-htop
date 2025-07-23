package procutils

import "github.com/hassanaziz0012/go-htop/types"

func AssignChildren(processes *[]types.Process) []*types.Process {
	var roots []*types.Process
	procmap := buildMap(processes)

	for i := range *processes {
		p := &(*processes)[i]
		parent, ok := procmap[p.PPID]
		if ok {
			parent.Children = append(parent.Children, p)
		} else {
			roots = append(roots, p)
		}
	}

	return roots
}

func buildMap(processes *[]types.Process) map[int]*types.Process {
	procmap := make(map[int]*types.Process, len(*processes))

	for i := range *processes {
		p := &(*processes)[i]
		procmap[p.PID] = p
	}
	return procmap
}
