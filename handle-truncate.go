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

	resultChan := make(chan RpcResult, 1)
	CallLogFileCommand(logMachine, log, resultChan, "TruncateLogFile", false, "", 0)

	result := <-resultChan
	_ = json.NewEncoder(w).Encode(result)
}
