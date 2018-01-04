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
	Logs    []*LogFileInfoResult
}

func HandleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	logsNum := len(devopsConf.Logs)
	resultChan := make(chan LogShowResult, logsNum)

	for logger, log := range devopsConf.Logs {
		go showLog(logger, log, resultChan)
	}

	resultsMap := make(map[string]*LogShowResult)
	for i := 0; i < logsNum; i++ {
		result := <-resultChan
		resultsMap[result.Logger] = &result
	}

	results := make([]*LogShowResult, 0)
	for _, logger := range loggers {
		results = append(results, resultsMap[logger])
	}

	json.NewEncoder(w).Encode(results)
}

func showLog(logger string, log Log, results chan LogShowResult) {
	machinesNum := len(log.Machines)
	resultChan := make(chan LogFileInfoResult, machinesNum)
	for _, machineName := range log.Machines {
		go CallLogFileCommand(machineName, log, resultChan,
			"LogFileInfo", false, "", 0)
	}

	resultsMap := make(map[string]*LogFileInfoResult)
	for i := 0; i < machinesNum; i++ {
		commandResult := <-resultChan
		resultsMap[commandResult.MachineName] = &commandResult
	}

	logs := createLogsResult(log, resultsMap)

	results <- LogShowResult{
		Logger:  logger,
		LogPath: log.Path,
		Logs:    logs,
	}
}

func createLogsResult(log Log, resultsMap map[string]*LogFileInfoResult) []*LogFileInfoResult {
	logs := make([]*LogFileInfoResult, 0)
	for _, machineName := range log.Machines {
		result, ok := resultsMap[machineName]
		if ok {
			logs = append(logs, result)
		}
	}
	return logs
}

func HandleLocateLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]

	log, ok := devopsConf.Logs[loggerName]
	if !ok {
		return
	}

	timestampFrom := vars["timestampFrom"]
	timestampTo := vars["timestampTo"]
	size := len(timestampFrom)

	// awk 'substr($0,1,23)>="2017-12-17 15:31:54.587" && substr($0,1,23)<="2017-12-17 15:31:54.588"' < demo.log
	const awkTmpl = `awk 'substr($0,1,%d)>="%s" && substr($0,1,%d)<="%s"' < %s`
	command := fmt.Sprintf(awkTmpl, size, timestampFrom, size, timestampTo, log.Path)

	logMachinesNum := len(log.Machines)
	resultChan := make(chan CommandsResult, logMachinesNum)

	for _, machine := range log.Machines {
		go TimeoutCallShellCommand(machine, command, resultChan)
	}

	resultsMap := make(map[string]*CommandsResult)
	for i := 0; i < logMachinesNum; i++ {
		result := <-resultChan
		resultsMap[result.MachineName] = &result
	}

	results := make([]*CommandsResult, 0)
	for _, machineName := range log.Machines {
		results = append(results, resultsMap[machineName])
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
	case <-time.After(3 * time.Minute):
		resultChan <- CommandsResult{
			Error:       "timeout",
			MachineName: machineName,
		}
	}
}

func DialAndCallShellCommand(machine Machine, commands string) CommandsResult {
	conn, err := net.DialTimeout("tcp", machine.IP+":"+rpcPort, 3*time.Second)
	if err != nil {
		return CommandsResult{Error: err.Error()}
	}

	client := rpc.NewClient(conn)
	defer client.Close()

	return CallShellCommand(client, commands)
}

func CallShellCommand(client *rpc.Client, commands string) CommandsResult {
	args := &CommandsArg{commands, 3 * time.Minute}
	var reply CommandsResult

	err := client.Call("ShellCommand.Execute", args, &reply)
	if err != nil {
		return CommandsResult{Error: err.Error()}
	}

	return reply
}
