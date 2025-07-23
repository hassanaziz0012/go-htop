package procutils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func parseCpuWithPS(procId int) error {
	id := strconv.Itoa(procId)
	out, err := exec.Command("ps", "-p", id, "-o", "%cpu").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to read cpu usage. %d, %w", procId, err)
	}
	fmt.Println(string(out))
	return nil
}

func getCpuSnapshot(procId int) (procTime int, totalCpuTime int, err error) {
	stat, _ := readProcStatFile(procId)
	procTime, _ = getProcTime(procId, string(stat))

	systemStatFile, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read /proc/stat: %w", err)
	}
	lines := strings.Split(string(systemStatFile), "\n")
	if len(lines) == 0 {
		return 0, 0, fmt.Errorf("unexpected format in /proc/stat")
	}
	cpuFields := strings.Fields(lines[0])
	for _, val := range cpuFields[1:] {
		num, err := strconv.Atoi(val)
		if err != nil {
			continue // skip bad values
		}
		totalCpuTime += num
	}

	return procTime, totalCpuTime, nil
}

func getCpuUsage(procId int, ch chan float64) {
	defer close(ch)
	procTime1, totalCpuTime1, _ := getCpuSnapshot(procId)
	time.Sleep(time.Second * 1)
	procTime2, totalCpuTime2, _ := getCpuSnapshot(procId)

	procDelta := procTime2 - procTime1
	totalDelta := totalCpuTime2 - totalCpuTime1

	cpuUsage := 100 * float64(procDelta) / float64(totalDelta)
	numCores := float64(runtime.NumCPU())

	ch <- cpuUsage * numCores
}
