package main

import "github.com/shirou/gopsutil/disk"

func Disk() []*disk.UsageStat {
	result := make([]*disk.UsageStat, 0)
	ret, err := disk.Partitions(false)
	dealwithErr(err)
	if err != nil {
		return result
	}

	empty := disk.PartitionStat{}

	for _, partion := range ret {
		if partion != empty {
			diskStat, err := disk.Usage(partion.Mountpoint)
			dealwithErr(err)
			if err != nil {
				break
			}

			result = append(result, diskStat)
		}
	}

	return result
}
