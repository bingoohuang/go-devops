package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

func HandleRestartProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	logMachine := vars["logMachine"]

	log := devopsConf.Logs[loggerName]

	resultChan := make(chan *LogFileInfoResult, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	CallLogFileCommand(&wg, logMachine, log, resultChan, "RestartProcess", true, "", 0)
	wg.Wait()
	close(resultChan)

	result := <-resultChan
	json.NewEncoder(w).Encode(result)
}
