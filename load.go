package main

import "github.com/shirou/gopsutil/load"

func Load() *load.AvgStat {
	stat, err := load.Avg()
	dealwithErr(err)

	return stat
}
