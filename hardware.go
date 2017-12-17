package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"time"
)

type HardwareInfo struct {
	OS                string // like darwin
	MemoryTotal       uint64 // like 8G
	FriendMemoryTotal string
	FreeMemory        uint64
	FriendFreeMemory  string
	MemoryUsedPercent float64
	Cores             int32 // number of cores
	Hostname          string
	TotalDisk         uint64
	FriendTotalDisk   string
	FreeDisk          uint64
	FriendFreeDisk    string
	DiskUsedPercent   float64
}

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

type MachineCommandArg struct {
}

type MachineCommandResult struct {
	Error       string
	Name        string
	MachineInfo HardwareInfo
	CostMillis  int64
}

type MachineCommand int

func (t *MachineCommand) MachineInfo(args *MachineCommandArg, result *MachineCommandResult) error {
	start := time.Now()

	hardwareInfo := GetHardwareInfo()
	elapsed := time.Since(start)
	result.Error = ""
	result.MachineInfo = hardwareInfo
	result.CostMillis = elapsed.Nanoseconds() / 1e6
	return nil
}

func GetHardwareInfo() HardwareInfo {
	// memory
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)

	// cpu - get CPU number of cores and speed
	cpuStats, err := cpu.Info()
	dealwithErr(err)

	var cores int32 = 0
	for _, cpuState := range cpuStats {
		cores += cpuState.Cores
	}

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	dealwithErr(err)

	diskStat, err := disk.Usage("/")
	dealwithErr(err)

	return HardwareInfo{
		OS:                runtime.GOOS,
		MemoryTotal:       vmStat.Total,
		FriendMemoryTotal: humanize.IBytes(vmStat.Total),
		FreeMemory:        vmStat.Free,
		FriendFreeMemory:  humanize.IBytes(vmStat.Free),
		MemoryUsedPercent: vmStat.UsedPercent,
		Cores:             cores,
		Hostname:          hostStat.Hostname,
		TotalDisk:         diskStat.Total,
		FriendTotalDisk:   humanize.IBytes(diskStat.Total),
		FreeDisk:          diskStat.Free,
		FriendFreeDisk:    humanize.IBytes(diskStat.Free),
		DiskUsedPercent:   diskStat.UsedPercent,
	}
}
