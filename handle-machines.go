package main

import (
	"encoding/json"
	"github.com/bingoohuang/go-utils"
	"net/http"
	"net/rpc"
	"time"
)

func HandleMachines(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)

	size := len(devopsConf.Machines)
	resultChan := make(chan MachineCommandResult, size)
	for machineName, machine := range devopsConf.Machines {
		go TimeoutCallMachineInfo(machine, machineName, resultChan)
	}

	resultsMap := make(map[string]*MachineCommandResult)
	for i := 0; i < size; i++ {
		result := <-resultChan
		resultsMap[result.MachineName] = &result
	}

	results := make([]*MachineCommandResult, 0)
	for _, machineName := range machineNames {
		results = append(results, resultsMap[machineName])
	}

	json.NewEncoder(w).Encode(results)
}

func TimeoutCallMachineInfo(machine Machine, machineName string, resultChan chan MachineCommandResult) {
	c := make(chan MachineCommandResult, 1)
	reply := MachineCommandResult{
		MachineName: machineName,
	}

	go func() {
		err := DialAndCall(machine.IP+":"+rpcPort, func(client *rpc.Client) error {
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
