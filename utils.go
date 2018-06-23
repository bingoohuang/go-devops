package main

import (
	"log"
	"strings"
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
