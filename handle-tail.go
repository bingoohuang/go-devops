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

	logs := make([]LogFileInfoResult, 0)
	machinesNum := len(log.Machines)

	resultChan := make(chan LogFileInfoResult, machinesNum)
	for _, machine := range log.Machines {
		go TimeoutCallLogFileCommand(machine, log, resultChan, "TailLogFile", false, "-"+lines, 0)
	}

	for i := 0; i < machinesNum; i++ {
		commandResult := <-resultChan
		logs = append(logs, commandResult)
	}

	json.NewEncoder(w).Encode(logs)
}

func HandleTailFLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	logSeq, _ := vars["logSeq"]

	machineLogSeqMap := make(map[string]int)
	if logSeq != "init" {
		ss := strings.Split(logSeq, ",")
		for _, pair := range ss {
			z := strings.Split(pair, "|")
			machineLogSeqMap[z[0]], _ = strconv.Atoi(z[1])
		}
	}

	log := devopsConf.Logs[loggerName]

	logs := make([]LogFileInfoResult, 0)
	machinesNum := len(log.Machines)

	resultChan := make(chan LogFileInfoResult, machinesNum)
	for _, machine := range log.Machines {
		machineLogSeq, ok := machineLogSeqMap[machine]
		seq := -1
		if ok {
			seq = machineLogSeq
		}
		go TimeoutCallLogFileCommand(machine, log, resultChan, "TailFLogFile", false, "", seq)
	}

	newSeqMap := make(map[string]int)
	for i := 0; i < machinesNum; i++ {
		commandResult := <-resultChan
		logs = append(logs, commandResult)
		newSeqMap[commandResult.MachineName] = commandResult.TailNextSeq
	}

	newLogSeq := ""
	for key, value := range newSeqMap {
		if newLogSeq != "" {
			newLogSeq += ","
		}
		newLogSeq += key + "|" + strconv.Itoa(value)
	}

	json.NewEncoder(w).Encode(
		struct {
			Results   []LogFileInfoResult
			NewLogSeq string
		}{
			Results:   logs,
			NewLogSeq: newLogSeq,
		})
}
