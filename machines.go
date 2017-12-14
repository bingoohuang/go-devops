package main

import (
	"net/http"
	"net/rpc"
	"encoding/json"
)

func HandleMachines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	machineInfos := make([]MachineCommandResult, 0)
	for _, machine := range config.Machines {
		client, err := rpc.DialHTTP("tcp", machine.IP+":6979")
		FatalIfErr(err)
		defer client.Close()

		machineInfo := CallMachineInfo(client)
		machineInfos = append(machineInfos, machineInfo)
	}

	json.NewEncoder(w).Encode(machineInfos)
}

func CallMachineInfo(client *rpc.Client) MachineCommandResult {
	// Synchronous call
	args := &MachineCommandArg{}
	var reply MachineCommandResult
	err := client.Call("MachineCommand.MachineInfo", args, &reply)
	FatalIfErr(err)

	return reply
}
