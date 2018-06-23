package main

import (
	"encoding/json"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleRestartProcess(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)
	vars := mux.Vars(r)
	loggerName := vars["loggerName"]
	logMachine := vars["logMachine"]

	log := devopsConf.Logs[loggerName]

	resultChan := make(chan RpcResult, 1)
	CallLogFileCommand(nil, logMachine, log, resultChan, "RestartProcess", true, "", 0)

	result := <-resultChan
	json.NewEncoder(w).Encode(result)
}
