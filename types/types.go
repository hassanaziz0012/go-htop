package types

import (
	"time"

	"github.com/rivo/tview"
)

type Process struct {
	PID           int
	PPID          int
	User          string
	Status        string
	CPU           float64
	VirtualMemory int
	SharedMemory  int
	Memory        int
	Time          time.Duration
	ThreadCount   int
	Priority      int
	Nice          int
	Command       string
	Children      []*Process
	Highlight     bool
}

type Metadata struct {
	Uptime   string
	LoadAvgs string
	Tasks    int
	Threads  int
	KThreads int
	Running  int
}

type MemoryStatus struct {
	MemTotal int
	MemUsed  int
	SwpTotal int
	SwpUsed  int
}

type CPUStatus struct {
	Cores []CPUCore
}

type CPUCore struct {
	Num          int
	UsagePercent float64
	TotalCPU     uint64
	IdleCPU      uint64
}

type SortDirection int

const (
	SortAsc SortDirection = iota
	SortDsc
)

type TableWithConfig struct {
	Table        *tview.Table
	AllProcesses *[]Process
	Processes    *[]Process
	Search       string
	Filter       string
	Sort         string
	SortDirection
}
