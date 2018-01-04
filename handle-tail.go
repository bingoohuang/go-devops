package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func HandleTailLogFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	lines := vars["lines"]
	log := devopsConf.Logs[loggerName]

	machinesNum := len(log.Machines)

	resultChan := make(chan LogFileInfoResult, machinesNum)
	for _, machineName := range log.Machines {
		go CallLogFileCommand(machineName, log, resultChan,
			"TailLogFile", false, "-"+lines, 0)
	}

	resultsMap := make(map[string]*LogFileInfoResult)
	for i := 0; i < machinesNum; i++ {
		commandResult := <-resultChan
		resultsMap[commandResult.MachineName] = &commandResult
	}

	logs := createLogsResult(log, resultsMap)

	json.NewEncoder(w).Encode(logs)
}

func HandleTailFLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	logSeq, _ := vars["logSeq"]

	machineLogSeqMap := parseMachineSeqs(logSeq)
	log := devopsConf.Logs[loggerName]
	machinesNum := len(log.Machines)

	resultChan := make(chan LogFileInfoResult, machinesNum)
	for _, machineName := range log.Machines {
		seq := findSeq(machineLogSeqMap, machineName)
		go CallLogFileCommand(machineName, log, resultChan,
			"TailFLog", false, "", seq)
	}

	resultsMap := make(map[string]*LogFileInfoResult)
	newSeqMap := make(map[string]int)
	for i := 0; i < machinesNum; i++ {
		commandResult := <-resultChan
		resultsMap[commandResult.MachineName] = &commandResult
		newSeqMap[commandResult.MachineName] = commandResult.TailNextSeq
	}

	logs := make([]*LogFileInfoResult, 0)
	for _, machineName := range machineNames {
		logs = append(logs, resultsMap[machineName])
	}

	json.NewEncoder(w).Encode(
		struct {
			Results   []*LogFileInfoResult
			NewLogSeq string
		}{
			Results:   logs,
			NewLogSeq: createMachineSeqs(newSeqMap),
		})
}

func findSeq(machineLogSeqMap map[string]int, machineName string) int {
	machineLogSeq, ok := machineLogSeqMap[machineName]
	if ok {
		return machineLogSeq
	}
	return -1
}

func createMachineSeqs(newSeqMap map[string]int) string {
	newLogSeq := ""
	for key, value := range newSeqMap {
		if newLogSeq != "" {
			newLogSeq += ","
		}
		newLogSeq += key + "|" + strconv.Itoa(value)
	}

	return newLogSeq
}

func parseMachineSeqs(logSeq string) map[string]int {
	machineLogSeqMap := make(map[string]int)
	if logSeq != "init" {
		ss := strings.Split(logSeq, ",")
		for _, pair := range ss {
			z := strings.Split(pair, "|")
			machineLogSeqMap[z[0]], _ = strconv.Atoi(z[1])
		}
	}

	return machineLogSeqMap
}
