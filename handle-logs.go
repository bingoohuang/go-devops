package main

import (
	"encoding/json"
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type LogShowResult struct {
	Logger  string
	LogPath string
	Logs    []RpcResult
}

func HandleLogs(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)

	resultChan := make(chan *LogShowResult, len(devopsConf.Logs))

	var wg sync.WaitGroup
	for logger, log := range devopsConf.Logs {
		wg.Add(1)
		GoShowLog(&wg, logger, log, resultChan)
	}
	wg.Wait()
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

func GoShowLog(logsWg *sync.WaitGroup, logger string, log Log, results chan *LogShowResult) {
	go showLog(logsWg, logger, log, results)
}

func showLog(logsWg *sync.WaitGroup, logger string, log Log, results chan *LogShowResult) {
	defer logsWg.Done()

	resultChan := make(chan RpcResult, len(log.Machines))

	var wg sync.WaitGroup
	for _, logMachineName := range log.Machines {
		wg.Add(1)
		GoCallLogFileCommand(&wg, logMachineName, log, resultChan, "LogFileInfo", false, "", 0)
	}
	wg.Wait()
	close(resultChan)

	resultsMap := make(map[string]RpcResult)
	for commandResult := range resultChan {
		resultsMap[commandResult.GetMachineName()] = commandResult
	}

	results <- &LogShowResult{
		Logger:  logger,
		LogPath: log.Path,
		Logs:    createLogsResult(log, resultsMap),
	}
}

func createLogsResult(log Log, resultsMap map[string]RpcResult) []RpcResult {
	logs := make([]RpcResult, 0)
	for _, logMachineName := range log.Machines {
		machineName := findMachineName(logMachineName)
		result, ok := resultsMap[machineName]
		if ok {
			logs = append(logs, result)
		}
	}
	return logs
}

type LocateLogRsp struct {
	Err     string
	Results []*ShellResultCommandResult
}

func HandleLocateLogResult(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)
	rsp := &LocateLogRsp{Results: make([]*ShellResultCommandResult, 0)}

	qs := r.URL.Query()

	resultChan := make(chan RpcResult)
	for k, v := range qs {
		GoRpcExecuteTimeout(k, &ShellResultCommandArg{ShellKey: v[0]}, &ShellResultCommandExecute{}, 1*time.Second, resultChan)
	}

	for range qs {
		result := <-resultChan
		rsp.Results = append(rsp.Results, result.(*ShellResultCommandResult))
	}

	json.NewEncoder(w).Encode(rsp)
}

func HandleLocateLog(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]

	log, ok := devopsConf.Logs[loggerName]
	if !ok {
		return
	}

	logKey := vars["logKey"]
	preLines := vars["preLines"]
	lines := vars["lines"]

	const awkTmpl = `awk 'BEGIN{h=1;t=1;p=%s;f=0;m=%s}{if(f==0){if($0~/%s/){f=1;for(k in a)print a[k];print}else{a[t++]=$0;if(t-h>p)delete a[h++]}}else{print;if(++f==m)exit}}' %s`
	command := fmt.Sprintf(awkTmpl, preLines, lines, regexp.QuoteMeta(logKey), log.Path)

	executeCommand(log, command, w)
}

func executeCommand(log Log, command string, w http.ResponseWriter) {
	logMachinesNum := len(log.Machines)
	resultChan := make(chan RpcResult, logMachinesNum)
	for _, machine := range log.Machines {
		args := &CommandsArg{command, 3 * time.Minute}
		GoRpcExecuteTimeout(machine, args, &ShellCommandExecute{}, 3*time.Minute, resultChan)
	}
	resultsMap := make(map[string]RpcResult)
	for i := 0; i < logMachinesNum; i++ {
		result := <-resultChan
		resultsMap[result.GetMachineName()] = result
	}
	results := make([]RpcResult, 0)
	for _, machineName := range log.Machines {
		result := resultsMap[machineName]
		result.SetMachineName(machineName)
		results = append(results, result)
	}
	json.NewEncoder(w).Encode(results)
}
