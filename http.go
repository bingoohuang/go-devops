package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func StartHttpSever() {
	r := mux.NewRouter()

	r.HandleFunc(contextPath+"/iconfont.{extension}", HandleFont)
	r.HandleFunc(contextPath+"/favicon.ico", HandleFavicon)

	r.HandleFunc(contextPath+"/truncateLogFile/{loggerName}/{logMachine}", HandleTruncateLogFile)
	r.HandleFunc(contextPath+"/restartProcess/{loggerName}/{logMachine}", HandleRestartProcess)
	r.HandleFunc(contextPath+"/locateLog/{loggerName}/{timestampFrom}/{timestampTo}", gzipWrapper(HandleLocateLog))
	r.HandleFunc(contextPath+"/grepLog/{loggerName}/{grepText}", gzipWrapper(HandleGrepLog))
	r.HandleFunc(contextPath+"/tailLogFile/{loggerName}/{lines}", gzipWrapper(HandleTailLogFile))
	r.HandleFunc(contextPath+"/tailFLog/{loggerName}/{traceMobile}/{logSeq}", gzipWrapper(HandleTailFLog))
	r.HandleFunc(contextPath+"/machines", HandleMachines)
	r.HandleFunc(contextPath+"/logs", HandleLogs)
	r.HandleFunc(contextPath+"/saveConfig", HandleSaveConf)
	r.HandleFunc(contextPath+"/loadConfig", HandleLoadConf)
	r.HandleFunc(contextPath+"/", gzipWrapper(HandleHome))
	//r.HandleFunc(contextPath+"/", BasicAuth(gzipWrapper(HandleHome), []byte("bingoo"), []byte("bingoo")))

	http.Handle(contextPath+"/", r)

	fmt.Println("start to listen at ", httpPort)
	http.ListenAndServe(":"+httpPort, nil)
}
