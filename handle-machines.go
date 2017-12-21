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
	reply := MachineCommandResult{
		Name: machineName,
	}

	go func() {
		err := DialAndCall(machine, func(client *rpc.Client) error {
			return client.Call("MachineCommand.MachineInfo", &MachineCommandArg{}, &reply)
		})
		if err != nil {
			reply.Error = err.Error()
		}
		c <- reply
	}()
	select {
	case result := <-c:
		resultChan <- result
	case <-time.After(1 * time.Second):
		reply.Error = "timeout"
		resultChan <- reply
	}
}
