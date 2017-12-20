package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

// http://127.0.0.1:6879/log/yoga-system/2015-07-07%2011:23:33
func FindHandleLogsBetweenTimestamps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	logger := vars["logger"]

	log, ok := config.Logs[logger]
	if !ok {
		return
	}

	timestampFrom := vars["timestampFrom"]
	timestampTo := vars["timestampTo"]
	size := len(timestampFrom)

	// awk 'substr($0,1,23)>="2017-12-17 15:31:54.587" && substr($0,1,23)<="2017-12-17 15:31:54.588"' < demo.log
	const awkTmpl = `awk 'substr($0,1,%d)>="%s" && substr($0,1,%d)<="%s"' < %s`
	command := fmt.Sprintf(awkTmpl, size, timestampFrom, size, timestampTo, log.Path)

	results := make([]CommandsResult, 0)
	for _, machine := range log.Machines {
		result := timeoutCallShellCommand(config.Machines[machine], command)
		results = append(results, result)
	}

	json.NewEncoder(w).Encode(results)
}

func timeoutCallShellCommand(machine Machine, commands string) CommandsResult {
	c := make(chan CommandsResult, 1)
	go func() { c <- DialAndCallShellCommand(machine, commands) }()
	select {
	case result := <-c:
		return result
	case <-time.After(1 * time.Second):
		return CommandsResult{
			Error: "timeout",
		}
	}
}

func DialAndCallShellCommand(machine Machine, commands string) CommandsResult {
	conn, err := net.DialTimeout("tcp", machine.IP+":"+rpcPort, 1*time.Second)
	if err != nil {
		return CommandsResult{
			Error: err.Error(),
		}
	}

	client := rpc.NewClient(conn)
	defer client.Close()

	return CallShellCommand(client, commands)
}

func CallShellCommand(client *rpc.Client, commands string) CommandsResult {
	args := &CommandsArg{commands, 100 * time.Millisecond}
	var reply CommandsResult

	err := client.Call("ShellCommand.Execute", args, &reply)
	if err != nil {
		return CommandsResult{
			Error: err.Error(),
		}
	}

	return reply
}
