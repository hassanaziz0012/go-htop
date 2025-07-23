package main

import (
	"context"
	"time"

	"github.com/hassanaziz0012/go-htop/procutils"
	"github.com/hassanaziz0012/go-htop/types"
	"github.com/hassanaziz0012/go-htop/ui"
	"github.com/hassanaziz0012/go-htop/utils"
)

var RefreshRate time.Duration = time.Second * 1

func main() {
	plist := make(chan []types.Process)
	metadataChan := make(chan types.Metadata)
	memStatusChan := make(chan types.MemoryStatus)
	cpuStatusChan := make(chan types.CPUStatus)

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		ticker := time.NewTicker(RefreshRate)
		for {
			select {
			case <-ticker.C:
				var processes []types.Process = []types.Process{}
				procutils.ReadAllProcesses(&processes)
				plist <- processes
				metadataChan <- procutils.CreateMetadataFromProcesses(&processes)
				memStatusChan <- procutils.CreateMemStatus()
				cpuStatusChan <- procutils.CreateCPUStatus()
			case <-ctx.Done():
				utils.WriteLog("cancelled by user")
				return
			}
		}
	}(ctx)

	ui.RunApplication(plist, metadataChan, memStatusChan, cpuStatusChan, cancel)
}
