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
	mergeScripts := go_utils.MergeJs(MustAsset, go_utils.FilterAssetNames(AssetNames(), ".js"))
	js := go_utils.MinifyJs(mergeScripts, devMode)
	index = strings.Replace(index, "${contextPath}", contextPath, -1)
	index = strings.Replace(index, "/*.SCRIPT*/", js, 1)

	if ok {
		resultChan := make(chan RpcResult)
		GoRpcExecuteTimeout(machineName, &AgentCommandArg{Processes: make(map[string][]string), Topn: 0}, &AgentCommandExecute{}, 3*time.Second, resultChan)
		result := <-resultChan
		r := result.(*AgentCommandResult)
		index = buildAgentView(index, "", r)
	} else {
		index = strings.Replace(index, `<Error/>`, machineName+`'s Agent Info Not Available!`, -1)
	}

	w.Write([]byte(index))
}
