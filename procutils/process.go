package procutils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hassanaziz0012/go-htop/types"
)

func ReadAllProcesses(processes *[]types.Process) {
	dirs, err := os.ReadDir("/proc")
	if err != nil {
		log.Fatal("failed to read /proc")
	}

	var wg sync.WaitGroup
	for _, dir := range dirs {
		if procId, err := strconv.Atoi(dir.Name()); err == nil {
			wg.Add(1)
			go func(processes *[]types.Process) {
				defer wg.Done()
				process, err := parseProcess(procId)
				if err != nil {
					return
				}
				*processes = append(*processes, process)
			}(processes)
		}
	}
	wg.Wait()
}

func readUptime() (string, error) {
	f, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", err
	}

	uptimeParsed, _ := strconv.ParseFloat(strings.Fields(string(f))[0], 64)
	uptimeInDuration := time.Duration(float64(uptimeParsed)) * time.Second

	hrs := int(uptimeInDuration.Hours())
	mins := int(uptimeInDuration.Minutes()) % 60
	secs := int(uptimeInDuration.Seconds()) % 60
	uptime := fmt.Sprintf("%02d:%02d:%02d", hrs, mins, secs)

	return uptime, nil
}

func readLoadAvgs() (string, error) {
	f, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return "", err
	}
	data := strings.Fields(string(f))
	if len(data) < 3 {
		return "", errors.New("unexpected format in /proc/loadavg")
	}
	return fmt.Sprintf("%s %s %s", data[0], data[1], data[2]), nil
}

func CreateMetadataFromProcesses(processes *[]types.Process) types.Metadata {
	uptime, err := readUptime()
	if err != nil {
		return types.Metadata{} // fail silently, and retry in next iteration.
	}
	loadAvgs, err := readLoadAvgs()
	if err != nil {
		return types.Metadata{} // fail silently, and retry in next iteration.
	}

	var threads int = 0
	var running int = 0
	for _, p := range *processes {
		threads += p.ThreadCount
		if p.Status == "R" {
			running++
		}
	}

	metadata := types.Metadata{
		Uptime:   uptime,
		LoadAvgs: loadAvgs,
		Tasks:    len(*processes),
		Threads:  threads,
		KThreads: 0, // TODO: idk where the kthreads are. build and connect this functionality.
		Running:  running,
	}

	return metadata
}

func CreateMemStatus() types.MemoryStatus {
	f, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return types.MemoryStatus{} // fail silently, retry in next iteration.
	}
	lines := strings.Split(string(f), "\n")

	parseLine := func(line string) int {
		parsed, _ := strconv.Atoi(strings.Fields(line)[1])
		return parsed
	}

	var memtotal, memfree, swaptotal, swapfree, buffers, cached, sReclaimable int

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "MemTotal"):
			parsed := parseLine(line)
			memtotal = parsed

		case strings.HasPrefix(line, "MemFree"):
			parsed := parseLine(line)
			memfree = parsed

		case strings.HasPrefix(line, "Buffers"):
			parsed := parseLine(line)
			buffers = parsed

		case strings.HasPrefix(line, "Cached"):
			parsed := parseLine(line)
			cached = parsed

		case strings.HasPrefix(line, "SReclaimable"):
			parsed := parseLine(line)
			sReclaimable = parsed

		case strings.HasPrefix(line, "SwapTotal"):
			parsed := parseLine(line)
			swaptotal = parsed

		case strings.HasPrefix(line, "SwapFree"):
			parsed := parseLine(line)
			swapfree = parsed
		}
	}

	return types.MemoryStatus{
		MemTotal: memtotal,
		MemUsed:  memtotal - memfree - buffers - cached - sReclaimable,
		SwpTotal: swaptotal,
		SwpUsed:  swaptotal - swapfree,
	}
}

func CreateCPUStatus() types.CPUStatus {
	getAllCoresSnapshot := func() []types.CPUCore {
		f, err := os.ReadFile("/proc/stat")
		if err != nil {
			return []types.CPUCore{} // fail silently, retry on next iteration.
		}
		lines := strings.SplitSeq(string(f), "\n")
		var cores []types.CPUCore
		for line := range lines {
			if !strings.HasPrefix(line, "cpu") {
				break
			}
			if strings.Fields(line)[0] == "cpu" {
				continue // skip the total and only calculate individual cores
			}

			core := getCpuCoreSnapshot(line)
			cores = append(cores, core)
		}
		return cores
	}

	prevCores := getAllCoresSnapshot()
	time.Sleep(time.Second * 1)
	curCores := getAllCoresSnapshot()

	for i := range prevCores {
		prev := &prevCores[i]
		cur := &curCores[i]

		if prev.Num != cur.Num {
			log.Fatal("core mismatch")
		}

		deltaTotal := cur.TotalCPU - prev.TotalCPU
		deltaIdle := cur.IdleCPU - prev.IdleCPU

		usage := float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100
		cur.UsagePercent = float64(usage)
	}

	return types.CPUStatus{
		Cores: curCores,
	}
}

func getCpuCoreSnapshot(line string) types.CPUCore {
	fields := strings.Fields(line)
	coreIndex, _ := strconv.Atoi(strings.Replace(fields[0], "cpu", "", 1))

	ticks := fields[1:]
	var totalCpuTime uint64
	var idleCpuTime uint64
	for _, tick := range ticks {
		parsedTick, err := strconv.ParseUint(tick, 10, 64)
		if err != nil {
			return types.CPUCore{}
		}
		totalCpuTime += parsedTick
	}

	idleTime, _ := strconv.ParseUint(ticks[3], 10, 64)
	ioWaitTime, _ := strconv.ParseUint(ticks[4], 10, 64)
	idleCpuTime = idleTime + ioWaitTime

	return types.CPUCore{
		Num:      coreIndex,
		TotalCPU: totalCpuTime,
		IdleCPU:  idleCpuTime,
	}
}

func parseProcess(procId int) (types.Process, error) {
	pagesize, _ := getPageSize()

	// sampling requires a 1s delay, so i run this in a goroutine so it doesn't delay the rest of the app.
	cpuChan := make(chan float64)
	go getCpuUsage(procId, cpuChan)

	stat, err := readProcStatFile(procId)
	if err != nil {
		return types.Process{}, err
	}

	status, err := readProcStatusFile(procId)
	if err != nil {
		return types.Process{}, err
	}
	statusData := string(status)

	p := types.Process{}
	p.PID = procId

	p.PPID, err = getProcessPPID(statusData)
	if err != nil {
		return types.Process{}, err
	}

	p.Command, err = getProcessCmd(procId)
	if err != nil {
		return types.Process{}, err
	}

	p.Status = getProcessState(statusData)
	p.User = getProcessUsername(statusData)
	p.Time, _ = getProcessTime(procId, string(stat))

	p.ThreadCount, _ = getThreadCount(statusData)
	p.Memory, _ = getUsedMem(statusData)
	p.VirtualMemory, _ = getVirtualMem(statusData)

	p.SharedMemory, err = getSharedMemory(procId, pagesize)
	if err != nil {
		return types.Process{}, err
	}

	p.Priority, _ = getPriority(string(stat))
	p.Nice, _ = getNice(string(stat))

	cpu := <-cpuChan
	p.CPU = cpu

	return p, nil
}

func getProcessPPID(statusData string) (int, error) {
	lines := strings.Split(statusData, "\n")
	ppid, err := strconv.Atoi(strings.Fields(lines[7-1])[1])
	return ppid, err
}

func getProcessCmd(procId int) (string, error) {
	command, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", procId))
	if err != nil {
		return "", fmt.Errorf("failed to read process cmd: %d", procId)
	}
	return string(command), nil
}

func getProcessState(statusData string) string {
	lines := strings.Split(statusData, "\n")
	state := strings.Fields(lines[3-1])[1]
	return state
}

func getProcessUsername(statusData string) string {
	lines := strings.Split(statusData, "\n")
	uid, _ := strconv.Atoi(strings.Fields(lines[9-1])[1])
	username, _ := getUsernameFromUid(uid)
	return username
}

func getProcessTime(procId int, statData string) (time.Duration, error) {
	proctime, err := getProcTime(procId, statData)
	if err != nil {
		return 0, err
	}
	clockTicks, _ := getClockTicks()
	proctimeInSecs := proctime / clockTicks
	ptime := time.Duration(proctimeInSecs) * time.Second
	return ptime, nil
}

func getThreadCount(statusData string) (int, error) {
	re := regexp.MustCompile(`Threads:\s+(\d+)`)
	threads := re.FindString(statusData)
	threadCount, err := strconv.Atoi(strings.TrimSpace(strings.Split(threads, ":")[1]))
	if err != nil {
		return 0, errors.New("failed to parse thread count")
	}
	return threadCount, nil
}

func getPriority(statData string) (int, error) {
	fields := strings.Fields(statData)
	priority := fields[18-1] // 18th field
	parsed, err := strconv.Atoi(strings.TrimSpace(priority))
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func getNice(statData string) (int, error) {
	fields := strings.Fields(statData)
	nice := fields[19-1] // 19th field
	parsed, err := strconv.Atoi(strings.TrimSpace(nice))
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func readProcStatFile(procId int) ([]byte, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", procId))
	if err != nil {
		return nil, fmt.Errorf("failed to read stat file: %d - %w", procId, err)
	}
	return data, nil
}

func readProcStatusFile(procId int) ([]byte, error) {
	status, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", procId))
	if err != nil {
		return nil, errors.New("failed to read process status: " + strconv.Itoa(procId))
	}
	return status, err
}

func getProcTime(procId int, statData string) (int, error) {
	fields := strings.Fields(statData)
	if len(fields) < 15 {
		return 0, fmt.Errorf("unexpected format in /proc/%d/stat", procId)
	}
	utime, _ := strconv.Atoi(fields[14-1]) // get the 14th value
	stime, _ := strconv.Atoi(fields[15-1]) // get the 15th value
	procTime := utime + stime
	return procTime, nil
}

func getClockTicks() (int, error) {
	out, err := exec.Command("getconf", "CLK_TCK").CombinedOutput()
	if err != nil {
		return 0, errors.New("failed to get CLK_TCK")
	}
	ticks, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, errors.New("failed to parse clock ticks")
	}
	return ticks, nil
}

func getPageSize() (int, error) {
	out, err := exec.Command("getconf", "PAGESIZE").CombinedOutput()
	if err != nil {
		return 0, errors.New("failed to get PAGESIZE")
	}
	psize, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, errors.New("failed to parse page size")
	}
	return psize, nil
}

func getUsernameFromUid(uid int) (string, error) {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return "", errors.New("failed to read /etc/passwd")
	}
	users := strings.Lines(string(data))
	for user := range users {
		userData := strings.Split(user, ":")
		userId := userData[2]
		if userId == strconv.Itoa(uid) {
			username := userData[0]
			return username, nil
		}
	}
	return "", nil
}
