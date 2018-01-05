package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

type LogShowResult struct {
	Logger  string
	LogPath string
	Logs    []*LogFileInfoResult
}

func HandleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resultChan := make(chan *LogShowResult, len(devopsConf.Logs))

	var wg sync.WaitGroup
	for logger, log := range devopsConf.Logs {
		wg.Add(1)
		go showLog(wg, logger, log, resultChan)
	}
	wg.Done()
	close(resultChan)

	resultsMap := make(map[string]*LogShowResult)
	for result := range resultChan {
		resultsMap[result.Logger] = result
	}

	results := make([]*LogShowResult, 0)
	for _, logger := range loggers {
		results = append(results, resultsMap[logger])
	}

	json.NewEncoder(w).Encode(results)
}

func showLog(logsWg sync.WaitGroup, logger string, log Log, results chan *LogShowResult) {
	defer logsWg.Done()

	resultChan := make(chan *LogFileInfoResult, len(log.Machines))

	var wg sync.WaitGroup
	for _, logMachineName := range log.Machines {
		wg.Add(1)
		go CallLogFileCommand(&wg, logMachineName, log, resultChan,
			"LogFileInfo", false, "", 0)
	}
	wg.Wait()
	close(resultChan)

	resultsMap := make(map[string]*LogFileInfoResult)
	for commandResult := range resultChan {
		resultsMap[commandResult.MachineName] = commandResult
	}

	results <- &LogShowResult{
		Logger:  logger,
		LogPath: log.Path,
		Logs:    createLogsResult(log, resultsMap),
	}
}

func createLogsResult(log Log, resultsMap map[string]*LogFileInfoResult) []*LogFileInfoResult {
	logs := make([]*LogFileInfoResult, 0)
	for _, logMachineName := range log.Machines {
		machineName := findMachineName(logMachineName)
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
