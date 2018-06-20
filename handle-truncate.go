package main

import (
	"encoding/json"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleTruncateLogFile(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	logMachine := vars["logMachine"]

	log := devopsConf.Logs[loggerName]

	resultChan := make(chan *LogFileInfoResult, 1)
	CallLogFileCommand(nil, logMachine, log, resultChan, "TruncateLogFile", false, "", 0)

	result := <-resultChan
	json.NewEncoder(w).Encode(result)
}
