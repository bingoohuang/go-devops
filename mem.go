package main

import "github.com/shirou/gopsutil/mem"

func Memory() *mem.VirtualMemoryStat {
	// memory
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)
	return vmStat
}
