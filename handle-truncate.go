package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleTruncateLogFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	logMachine := vars["logMachine"]

	log := devopsConf.Logs[loggerName]

	resultChan := make(chan LogFileInfoResult, 1)
	CallLogFileCommand(logMachine, log, resultChan, "TruncateLogFile", false, "", 0)

	result := <-resultChan
	json.NewEncoder(w).Encode(result)
}
