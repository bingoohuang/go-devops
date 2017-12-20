package main

import (
	"encoding/json"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

func HandleMachines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	results := make([]MachineCommandResult, 0)
	size := len(config.Machines)
	resultChan := make(chan MachineCommandResult, size)
	for machineName, machine := range config.Machines {
		go timeoutCallMachineInfo(machine, machineName, resultChan)
	}

	for i := 0; i < size; i++ {
		result := <-resultChan
		results = append(results, result)
	}

	json.NewEncoder(w).Encode(results)
}

func timeoutCallMachineInfo(machine Machine, machineName string, resultChan chan MachineCommandResult) {
	c := make(chan MachineCommandResult, 1)
	go func() { c <- DialAndCallMachineInfo(machine) }()
	select {
	case result := <-c:
		result.Name = machineName
		resultChan <- result
	case <-time.After(1 * time.Second):
		resultChan <- MachineCommandResult{
			Name:  machineName,
			Error: "timeout",
		}
	}
}

func DialAndCallMachineInfo(machine Machine) MachineCommandResult {
	conn, err := net.DialTimeout("tcp", machine.IP+":6979", 1*time.Second)
	if err != nil {
		return MachineCommandResult{
			Error: err.Error(),
		}
	}

	client := rpc.NewClient(conn)
	defer client.Close()

	return CallMachineInfo(client)
}

func CallMachineInfo(client *rpc.Client) MachineCommandResult {
	args := &MachineCommandArg{}
	var reply MachineCommandResult
	err := client.Call("MachineCommand.MachineInfo", args, &reply)
	if err != nil {
		return MachineCommandResult{
			Error: err.Error(),
		}
	}

	return reply
}
