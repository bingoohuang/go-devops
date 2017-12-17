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
	for machineName, machine := range config.Machines {
		result := timeoutCallMachineInfo(machine)
		result.Name = machineName
		results = append(results, result)
	}

	json.NewEncoder(w).Encode(results)
}

func timeoutCallMachineInfo(machine Machine) MachineCommandResult {
	c := make(chan MachineCommandResult, 1)
	go func() { c <- DialAndCallMachineInfo(machine) }()
	select {
	case result := <-c:
		return result
	case <-time.After(1 * time.Second):
		return MachineCommandResult{
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
