package main

import (
	"log"
	"net"
	"net/rpc"
	"strings"
	"time"
)

func FatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func DialAndCall(address string, callFun func(client *rpc.Client) error) error {
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return err
	}

	client := rpc.NewClient(conn)
	defer client.Close()

	return callFun(client)
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
