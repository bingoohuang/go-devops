package main

import (
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var machineInfoMap = make(map[string]*AgentCommandResult)

func HandleMachineInfo(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeHtml(w)

	vars := mux.Vars(r)
	machineName := vars["machineName"]
	result, ok := machineInfoMap[machineName]

	index := string(MustAsset("res/viewagent.html"))
	if ok {
		index = buildAgentView(index, "", result)
	} else {
		index = strings.Replace(index, `<Error/>`, machineName+`'s Agent Info Not Available!`, -1)
	}

	w.Write([]byte(index))
}
