package main

import (
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

type AgentCommandArg struct {
	Processes map[string][]string
}

type DiskUsage struct {
	Path        string
	Fstype      string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

type AgentCommandResult struct {
	Load1  float64
	Load5  float64
	Load15 float64

	MemTotal       uint64
	MemAvailable   uint64
	MemUsed        uint64
	MemUsedPercent float64

	DiskUsages []DiskUsage

	Cores    int32 // number of cores
	Hostname string

	Processes map[string]PsAuxItem
	Top       []PsAuxItem

	MachineName string
	Error       string
}

var Cores int32 // number of cores
var Hostname string

func init() {
	// cpu - get CPU number of cores and speed
	cpuStats, _ := cpu.Info()
	var cores int32 = 0
	for _, cpuState := range cpuStats {
		cores += cpuState.Cores
	}

	Cores = cores

	hostStat, _ := host.Info()
	Hostname = hostStat.Hostname
}

type AgentCommand int

type AgentCommandExeucte struct {
}

func (t *AgentCommandResult) GetMachineName() string {
	return t.MachineName
}

func (t *AgentCommandResult) SetMachineName(machineName string) {
	t.MachineName = machineName
}

func (t *AgentCommandResult) GetError() string {
	return t.Error
}

func (t *AgentCommandResult) SetError(err error) {
	if err != nil {
		t.Error += err.Error()
	}
}

func (t *AgentCommandExeucte) CreateResult(err error) RpcResult {
	result := &AgentCommandResult{}
	result.SetError(err)
	return result
}

func (t *AgentCommandExeucte) CommandName() string {
	return "AgentCommand"
}

func (t *AgentCommand) Execute(a *AgentCommandArg, r *AgentCommandResult) error {
	load := Load()
	r.Load1 = load.Load1
	r.Load5 = load.Load5
	r.Load15 = load.Load15

	memory := Memory()
	r.MemTotal = memory.Total
	r.MemAvailable = memory.Available
	r.MemUsed = memory.Used
	r.MemUsedPercent = memory.UsedPercent

	r.Cores = Cores
	r.Hostname = Hostname

	disks := Disk()
	r.DiskUsages = make([]DiskUsage, len(disks))
	for i, d := range disks {
		r.DiskUsages[i] = DiskUsage{d.Path, d.Fstype, d.Total, d.Free, d.Used, d.UsedPercent}
	}

	processes := make(map[string]PsAuxItem)
	for k, v := range a.Processes {
		grep := PsAuxGrep(v...)
		if len(grep) > 0 {
			processes[k] = *grep[0]
		} else {
			r.SetError(errors.New("Process " + k + " not found!"))
		}
	}

	r.Processes = processes
	top := PsAuxTop(10)
	r.Top = make([]PsAuxItem, 0)
	for _, i := range top {
		r.Top = append(r.Top, *i)
	}

	return nil
}
