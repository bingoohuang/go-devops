package main

import (
	"encoding/json"
	"github.com/bingoohuang/go-utils"
	"net/http"
	"time"
)

func HandleMachines(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)

	size := len(devopsConf.Machines)
	resultChan := make(chan RpcResult, size)
	for machineName := range devopsConf.Machines {
		GoRpcExecuteTimeout(machineName, &MachineCommandArg{}, &MachineCommandExecute{}, 3*time.Second, resultChan)
	}

	resultsMap := make(map[string]RpcResult)
	for i := 0; i < size; i++ {
		result := <-resultChan
		resultsMap[result.GetMachineName()] = result
	}

	results := make([]RpcResult, 0)
	for _, machineName := range machineNames {
		results = append(results, resultsMap[machineName])
	}

	json.NewEncoder(w).Encode(results)
}
