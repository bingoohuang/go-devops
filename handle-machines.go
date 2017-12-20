package main

import (
	"encoding/json"
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
		go TimeoutCallMachineInfo(machine, machineName, resultChan)
	}

	for i := 0; i < size; i++ {
		result := <-resultChan
		results = append(results, result)
	}

	json.NewEncoder(w).Encode(results)
}

func TimeoutCallMachineInfo(machine Machine, machineName string, resultChan chan MachineCommandResult) {
	c := make(chan MachineCommandResult, 1)
	go func() { c <- DialAndCall(machine, CallMachineInfo, &MachineCommandArg{}).(MachineCommandResult) }()
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

func CallMachineInfo(client *rpc.Client, args interface{}) interface{} {
	var reply MachineCommandResult
	err := client.Call("MachineCommand.MachineInfo", args, &reply)
	if err != nil {
		return MachineCommandResult{
			Error: err.Error(),
		}
	}

	return reply
}
