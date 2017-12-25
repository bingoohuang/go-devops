package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleTailLogFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	log := devopsConf.Logs[loggerName]

	logs := make([]LogFileInfoResult, 0)
	machinesNum := len(log.Machines)

	resultChan := make(chan LogFileInfoResult, machinesNum)
	for _, machine := range log.Machines {
		go TimeoutCallLogFileCommand(machine, log, resultChan, "TailLogFile", false)
	}

	for i := 0; i < machinesNum; i++ {
		commandResult := <-resultChan
		logs = append(logs, commandResult)
	}

	json.NewEncoder(w).Encode(logs)
}
