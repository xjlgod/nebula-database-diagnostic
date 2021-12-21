package physical

import (
	"errors"
	"fmt"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/remote"
	"log"
	"strconv"
	"strings"
)

type PhyInfo struct {
	Process ProcessInfo
	Memory  MemoryInfo
	Disk    DiskInfo
	Swap    SwapInfo
	IO      IOInfo
	System  SystemInfo
	CPU     CPUInfo
}

type ProcessInfo struct {
	RunNumber  int
	WaitNumber int
}

type MemoryInfo struct { // kB
	MemTotal int
	MemFree  int
	MemBuff  int
	MemCache int
}

type DiskInfo struct { // kB
	DiskTotal     int
	DiskAvailable int
}

type SwapInfo struct { // kB
	SwapIn  int
	SwapOut int
}

type IOInfo struct { // kb
	BitIn  int
	BitOut int
}

type SystemInfo struct {
	InterruptCount     int
	ContextSwitchCount int
}

type CPUInfo struct { // percent
	UserUseTime   int
	SystemUseTime int
	IdleTime      int
	WaitPercent   int
}

func GetPhyInfo(conf config.SSHConfig) (PhyInfo, error) {
	info := PhyInfo{}

	c, err := remote.GetSSHClient(conf.Username, conf)
	if err != nil {
		return info, err
	}
	log.Printf("%+v", c)

	res, ok := c.Execute("vmstat 1 1")
	if !ok {
		return info, fmt.Errorf("exec got error: %v, with std error: %v", res.Err, errors.New(string(res.StdErr)))
	}

	fields := strings.Fields(string(res.StdOut))
	fields = fields[len(fields)-17:]

	process := ProcessInfo{}
	runNum, _ := strconv.Atoi(fields[0])
	process.RunNumber = runNum
	waitNum, _ := strconv.Atoi(fields[1])
	process.WaitNumber = waitNum
	info.Process = process

	memory := MemoryInfo{}
	memFree, _ := strconv.Atoi(fields[3])
	memory.MemFree = memFree
	memBuff, _ := strconv.Atoi(fields[4])
	memory.MemBuff = memBuff
	memCache, _ := strconv.Atoi(fields[5])
	memory.MemCache = memCache
	memory.MemTotal = memory.MemFree + memory.MemBuff + memory.MemCache
	info.Memory = memory

	disk, err := getDiskInfo(conf)
	info.Disk = disk

	swap := SwapInfo{}
	swapIn, _ := strconv.Atoi(fields[6])
	swap.SwapIn = swapIn
	swapOut, _ := strconv.Atoi(fields[7])
	swap.SwapOut = swapOut
	info.Swap = swap

	io := IOInfo{}
	bitIn, _ := strconv.Atoi(fields[8])
	io.BitIn = bitIn
	bitOut, _ := strconv.Atoi(fields[9])
	io.BitOut = bitOut
	info.IO = io

	system := SystemInfo{}
	systemIC, _ := strconv.Atoi(fields[10])
	system.InterruptCount = systemIC
	systemCC, _ := strconv.Atoi(fields[11])
	system.ContextSwitchCount = systemCC
	info.System = system

	cpu := CPUInfo{}
	cpuUU, _ := strconv.Atoi(fields[12])
	cpu.UserUseTime = cpuUU
	cpuSU, _ := strconv.Atoi(fields[13])
	cpu.SystemUseTime = cpuSU
	cpuIU, _ := strconv.Atoi(fields[14])
	cpu.IdleTime = cpuIU
	cpuWA, _ := strconv.Atoi(fields[15])
	cpu.WaitPercent = cpuWA
	info.CPU = cpu

	return info, err
}

func getDiskDetailInfo(conf config.SSHConfig) (map[string]DiskInfo, error) {
	c, err := remote.GetSSHClient(conf.Username, conf)
	if err != nil {
		return nil, err
	}
	log.Printf("%+v", c)

	res, ok := c.Execute("df -BK | grep -vE '^Filesystem|tmpfs|udev' | awk '{ print $1 \" \" $2 \" \" $4 }'")
	if !ok {
		return nil, fmt.Errorf("exec got error: %v, with std error: %v", res.Err, errors.New(string(res.StdErr)))
	}

	detailInfos := make(map[string]DiskInfo, 0)
	for _, s := range strings.Split(string(res.StdOut), "\n") {
		fields := strings.Fields(s)
		if len(fields) == 3 {
			info := DiskInfo{}
			f1, _ := strconv.Atoi(trimSuffix(strings.TrimSpace(fields[1]), "K"))
			info.DiskTotal = f1
			f2, _ := strconv.Atoi(trimSuffix(strings.TrimSpace(fields[2]), "K"))
			info.DiskAvailable = f2
			f0 := strings.TrimSpace(fields[0])
			detailInfos[f0] = info
		}
	}

	return detailInfos, nil
}

func getDiskInfo(conf config.SSHConfig) (DiskInfo, error) {
	detailInfo, err := getDiskDetailInfo(conf)
	if err != nil {
		return DiskInfo{}, err
	}

	diskTotal, diskAvailable := 0, 0
	for _, info := range detailInfo {
		diskTotal += info.DiskTotal
		diskAvailable += info.DiskAvailable
	}

	return DiskInfo{diskTotal, diskAvailable}, nil
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
