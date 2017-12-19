package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func StartHttpSever() {
	r := mux.NewRouter()
	r.HandleFunc(contextPath+"/log/{logger}/{timestampFrom}/{timestampTo}", HandleLogs)
	r.HandleFunc(contextPath+"/machines", HandleMachines)
	r.HandleFunc(contextPath+"/", HandleHome)
	http.Handle(contextPath+"/", r)

	sport := strconv.Itoa(port)
	fmt.Println("start to listen at ", sport)
	http.ListenAndServe(":"+sport, nil)
}
