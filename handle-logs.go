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

type LogShowResult struct {
	Logger  string
	LogPath string
	Logs    []LogFileInfoResult
}

func HandleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	logsNum := len(devopsConf.Logs)
	resultChan := make(chan LogShowResult, logsNum)

	for logger, log := range devopsConf.Logs {
		go showLog(logger, log, resultChan)
	}

	results := make([]LogShowResult, 0)
	for i := 0; i < logsNum; i++ {
		result := <-resultChan
		results = append(results, result)
	}

	json.NewEncoder(w).Encode(results)
}

func showLog(logger string, log Log, results chan LogShowResult) {
	logs := make([]LogFileInfoResult, 0)

	machinesNum := len(log.Machines)

	resultChan := make(chan LogFileInfoResult, machinesNum)
	for _, machine := range log.Machines {
		go TimeoutCallLogFileCommand(machine, log, resultChan, "LogFileInfo", false, "")
	}

	for i := 0; i < machinesNum; i++ {
		commandResult := <-resultChan

		logs = append(logs, commandResult)
	}

	results <- LogShowResult{
		Logger:  logger,
		LogPath: log.Path,
		Logs:    logs,
	}
}

func FindHandleLogsBetweenTimestamps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	logger := vars["logger"]

	log, ok := devopsConf.Logs[logger]
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
	logMachinesNum := len(log.Machines)
	resultChan := make(chan CommandsResult, logMachinesNum)

	for _, machine := range log.Machines {
		go TimeoutCallShellCommand(machine, command, resultChan)
	}

	for i := 0; i < logMachinesNum; i++ {
		result := <-resultChan
		results = append(results, result)
	}

	json.NewEncoder(w).Encode(results)
}

func TimeoutCallShellCommand(machineName, commands string, resultChan chan CommandsResult) {
	machine := devopsConf.Machines[machineName]
	c := make(chan CommandsResult, 1)
	go func() { c <- DialAndCallShellCommand(machine, commands) }()
	select {
	case result := <-c:
		result.MachineName = machineName
		resultChan <- result
	case <-time.After(1 * time.Second):
		resultChan <- CommandsResult{
			Error:       "timeout",
			MachineName: machineName,
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
	args := &CommandsArg{commands, 500 * time.Millisecond}
	var reply CommandsResult

	err := client.Call("ShellCommand.Execute", args, &reply)
	if err != nil {
		return CommandsResult{
			Error: err.Error(),
		}
	}

	return reply
}
