package procutils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getSharedMemory(procId int, pageSize int) (int, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/statm", procId))
	if err != nil {
		return 0, fmt.Errorf("failed to read statm: %d", procId)
	}
	fields := strings.Fields(string(data))
	sharedMem, _ := strconv.Atoi(fields[2])
	return sharedMem * pageSize / 1024, nil // conversion to KB
}

func getUsedMem(statusData string) (int, error) {
	re := regexp.MustCompile(`VmRSS:\s+(\d+)\s(\S+)`)
	memUsed := re.FindString(statusData)
	vals := strings.Split(memUsed, ":")

	if len(vals) < 2 {
		return 0, errors.New("failed to parse MemUsed")
	}

	mem, err := strconv.Atoi(strings.Split(strings.TrimSpace(vals[1]), " ")[0])
	if err != nil {
		return 0, errors.New("failed to parse memory used")
	}
	return mem, nil
}

func getVirtualMem(statusData string) (int, error) {
	lines := strings.Split(statusData, "\n")
	line := lines[18] // 19th line
	vmem, err := strconv.Atoi(strings.Fields(line)[1])
	if err != nil {
		return 0, err
	}
	return vmem, nil
}
