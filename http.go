package main

import (
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"net/http"
)

func StartHttpSever() {
	r := mux.NewRouter()

	handleFunc(r, "/iconfont.{extension}", HandleFont, false, false)
	handleFunc(r, "/favicon.ico", HandleFavicon, false, false)

	handleFunc(r, "/truncateLogFile/{loggerName}/{logMachine}", HandleTruncateLogFile, false, true)
	handleFunc(r, "/restartProcess/{loggerName}/{logMachine}", HandleRestartProcess, false, true)
	handleFunc(r, "/locateLog/{loggerName}/{logKey}/{preLines}/{lines}", HandleLocateLog, true, true)
	handleFunc(r, "/tailLogFile/{loggerName}/{lines}", HandleTailLogFile, true, true)
	handleFunc(r, "/tailFLog/{loggerName}/{traceMobile}/{logSeq}", HandleTailFLog, true, true)
	handleFunc(r, "/machines", HandleMachines, false, true)
	handleFunc(r, "/logs", HandleLogs, false, true)
	handleFunc(r, "/saveConfig", HandleSaveConf, false, true)
	handleFunc(r, "/loadConfig", HandleLoadConf, false, true)

	handleFunc(r, "/", serveWelcome, false, false)
	handleFunc(r, "/home", HandleHome, true, true)

	http.Handle(contextPath+"/", r)

	fmt.Println("start to listen at ", httpPort)
	http.ListenAndServe(":"+httpPort, nil)
}

func handleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), requiredGzip, requiredBasicAuth bool) {
	wrap := go_utils.DumpRequest(f)
	if requiredBasicAuth && authBasic {
		wrap = go_utils.RandomPoemBasicAuth(wrap)
	}

	if requiredGzip {
		wrap = go_utils.GzipHandlerFunc(wrap)
	}

	r.HandleFunc(contextPath+path, wrap)
}

func serveWelcome(w http.ResponseWriter, r *http.Request) {
	if !authBasic {
		// fmt.Println("Redirect to", contextPath+"/home")
		// http.Redirect(w, r, contextPath+"/home", 301)
		HandleHome(w, r)
	} else {
		welcome := MustAsset("res/welcome.html")
		go_utils.ServeWelcome(w, welcome, contextPath)
	}
}
