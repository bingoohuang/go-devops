package main

import (
	"github.com/dustin/go-humanize"
	"log"
	"strconv"
	"strings"
	"time"
)

func FatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func findMachineName(logMachineName string) string {
	index := strings.Index(logMachineName, ":")
	if index > 0 {
		return strings.TrimSpace(logMachineName[0:index])
	}

	return logMachineName
}

func parseMachineNameAndAddress(logMachineName string) (string, string, string) {
	machineName := logMachineName
	machinePort := rpcPort
	errorMsg := ""

	index := strings.Index(logMachineName, ":")
	if index > 0 {
		machineName = strings.TrimSpace(logMachineName[0:index])
		machinePort = strings.TrimSpace(logMachineName[index+1:])
	}

	machine, ok := devopsConf.Machines[machineName]
	if !ok {
		errorMsg = machineName + " is unknown"
	}

	return machineName, machine.IP + ":" + machinePort, errorMsg
}

func HumanizedKib(kib string) string {
	u, e := strconv.ParseUint(kib, 10, 64)
	if e != nil {
		return kib + "KiB"
	}
	return strings.Replace(humanize.IBytes(u*1024), " ", "", 1)
}

func IsDurationAgo(maybeTs string, duration time.Duration) bool {
	if len(maybeTs) < 19 {
		return false
	}

	ts, e := time.ParseInLocation("2006-01-02 15:04:05", maybeTs[0:19], time.Local)

	now := time.Now()
	return e == nil && now.Sub(ts).Nanoseconds() > duration.Nanoseconds()
}
