package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleRestartProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	logMachine := vars["logMachine"]

	log := devopsConf.Logs[loggerName]

	resultChan := make(chan LogFileInfoResult, 1)
	TimeoutCallLogFileCommand(logMachine, log, resultChan, "RestartProcess", true, "", 0)

	result := <-resultChan
	json.NewEncoder(w).Encode(result)
}
