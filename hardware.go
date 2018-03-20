package main

import (
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"time"
)

type HardwareInfo struct {
	OS                       string // like darwin
	TotalMemory              uint64 // like 8G
	HumanizedTotalMemory     string
	AvailableMemory          uint64
	HumanizedAvailableMemory string
	MemoryUsedPercent        float64
	Cores                    int32 // number of cores
	Hostname                 string
	TotalDisk                uint64
	HumanizedTotalDisk       string
	FreeDisk                 uint64
	HumanizedFreeDisk        string
	DiskUsedPercent          float64
	Ips                      []string
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
	MachineName string
	MachineInfo HardwareInfo
	CostTime    string
}

type MachineCommand int

func (t *MachineCommand) MachineInfo(args *MachineCommandArg, result *MachineCommandResult) error {
	start := time.Now()

	hardwareInfo := GetHardwareInfo()
	elapsed := time.Since(start)
	result.Error = ""
	result.MachineInfo = hardwareInfo
	result.CostTime = elapsed.String()
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
		OS:                       runtime.GOOS,
		TotalMemory:              vmStat.Total,
		HumanizedTotalMemory:     humanize.IBytes(vmStat.Total),
		AvailableMemory:          vmStat.Available,
		HumanizedAvailableMemory: humanize.IBytes(vmStat.Available),
		MemoryUsedPercent:        vmStat.UsedPercent,
		Cores:                    cores,
		Hostname:                 hostStat.Hostname,
		TotalDisk:                diskStat.Total,
		HumanizedTotalDisk:       humanize.IBytes(diskStat.Total),
		FreeDisk:                 diskStat.Free,
		HumanizedFreeDisk:        humanize.IBytes(diskStat.Free),
		DiskUsedPercent:          diskStat.UsedPercent,
		Ips:                      go_utils.GetLocalIps(),
	}
}
