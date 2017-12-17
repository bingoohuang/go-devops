package main

import (
	"net/http"
	"net/rpc"
	"encoding/json"
	"time"
	"net"
)

func HandleMachines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	machineInfos := make([]MachineCommandResult, 0)
	for _, machine := range config.Machines {
		machineInfo := timeoutCall(machine)
		machineInfos = append(machineInfos, machineInfo)
	}

	json.NewEncoder(w).Encode(machineInfos)
}

func timeoutCall(machine Machine) MachineCommandResult {
	c := make(chan MachineCommandResult, 1)
	go func() { c <- DialAndCallMachineInfo(machine) }()
	select {
	case result := <-c:
		return result
	case <-time.After(3 * time.Second):
		return MachineCommandResult{
			Status: "NA",
			Error:  "timeout",
		}
	}
}

func DialAndCallMachineInfo(machine Machine) MachineCommandResult {
	conn, err := net.DialTimeout("tcp", machine.IP+":6979", 3*time.Second)
	if err != nil {
		return MachineCommandResult{
			Status: "NA",
			Error:  err.Error(),
		}
	}

	client := rpc.NewClient(conn)

	//client, err := rpc.DialHTTP("tcp", machine.IP+":6979")
	//if err != nil {
	//	return MachineCommandResult{
	//		Status: "NA",
	//		Error:  err.Error(),
	//	}
	//}
	defer client.Close()

	return CallMachineInfo(client)
}

func CallMachineInfo(client *rpc.Client) MachineCommandResult {
	args := &MachineCommandArg{}
	var reply MachineCommandResult
	err := client.Call("MachineCommand.MachineInfo", args, &reply)
	if (err != nil) {
		return MachineCommandResult{
			Status: "NA",
			Error:  err.Error(),
		}
	}

	return reply
}
