package main

import (
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"net/http"
)

func StartHttpSever() {
	r := mux.NewRouter()

	handleFunc(r, "/iconfont.{extension}", HandleFont, false)
	handleFunc(r, "/favicon.ico", HandleFavicon, false)

	handleFunc(r, "/truncateLogFile/{loggerName}/{logMachine}", HandleTruncateLogFile, false)
	handleFunc(r, "/restartProcess/{loggerName}/{logMachine}", HandleRestartProcess, false)
	handleFunc(r, "/locateLog/{loggerName}/{logKey}/{preLines}/{lines}", HandleLocateLog, true)
	handleFunc(r, "/tailLogFile/{loggerName}/{lines}", HandleTailLogFile, true)
	handleFunc(r, "/tailFLog/{loggerName}/{traceMobile}/{logSeq}", HandleTailFLog, true)
	handleFunc(r, "/machines", HandleMachines, false)
	handleFunc(r, "/logs", HandleLogs, false)
	handleFunc(r, "/saveConfig", HandleSaveConf, false)
	handleFunc(r, "/loadConfig", HandleLoadConf, false)

	handleFunc(r, "/", HandleHome, false)

	http.Handle(*contextPath+"/", r)

	fmt.Println("start to listen at ", httpPort)
	http.ListenAndServe(":"+httpPort, nil)
}

func handleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), requiredGzip bool) {
	wrap := go_utils.DumpRequest(f)
	wrap = go_utils.MustAuth(wrap, authParam)

	if requiredGzip {
		wrap = go_utils.GzipHandlerFunc(wrap)
	}

	r.HandleFunc(*contextPath+path, wrap)
}
