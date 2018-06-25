package main

import (
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func HandleMachineInfo(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeHtml(w)

	vars := mux.Vars(r)
	machineName := vars["machineName"]

	_, ok := devopsConf.Machines[machineName]
	index := string(MustAsset("res/viewagent.html"))
	if ok {
		resultChan := make(chan RpcResult)
		go RpcExecuteTimeout(machineName, &AgentCommandArg{Processes: make(map[string][]string), Topn: 0}, &AgentCommandExeucte{}, 3*time.Second, resultChan)
		result := <-resultChan
		r := result.(*AgentCommandResult)
		index = buildAgentView(index, "", r)
	} else {
		index = strings.Replace(index, `<Error/>`, machineName+`'s Agent Info Not Available!`, -1)
	}

	w.Write([]byte(index))
}
