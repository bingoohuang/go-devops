package main

import (
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"net/http"
)

func StartHttpSever() {
	r := mux.NewRouter()

	r.HandleFunc(contextPath+"/iconfont.{extension}", HandleFont)
	r.HandleFunc(contextPath+"/favicon.ico", HandleFavicon)

	r.HandleFunc(contextPath+"/truncateLogFile/{loggerName}/{logMachine}", HandleTruncateLogFile)
	r.HandleFunc(contextPath+"/restartProcess/{loggerName}/{logMachine}", HandleRestartProcess)
	r.HandleFunc(contextPath+"/locateLog/{loggerName}/{timestampFrom}/{timestampTo}", go_utils.GzipHandlerFunc(HandleLocateLog))
	r.HandleFunc(contextPath+"/grepLog/{loggerName}/{grepText}", go_utils.GzipHandlerFunc(HandleGrepLog))
	r.HandleFunc(contextPath+"/tailLogFile/{loggerName}/{lines}", go_utils.GzipHandlerFunc(HandleTailLogFile))
	r.HandleFunc(contextPath+"/tailFLog/{loggerName}/{traceMobile}/{logSeq}", go_utils.GzipHandlerFunc(HandleTailFLog))
	r.HandleFunc(contextPath+"/machines", HandleMachines)
	r.HandleFunc(contextPath+"/logs", HandleLogs)
	r.HandleFunc(contextPath+"/saveConfig", HandleSaveConf)
	r.HandleFunc(contextPath+"/loadConfig", HandleLoadConf)
	r.HandleFunc(contextPath+"/", go_utils.GzipHandlerFunc(HandleHome))

	http.Handle(contextPath+"/", r)

	fmt.Println("start to listen at ", httpPort)
	http.ListenAndServe(":"+httpPort, nil)
}
