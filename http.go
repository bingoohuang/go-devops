package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func StartHttpSever() {
	r := mux.NewRouter()
	r.HandleFunc("/log/{logger}/{timestampFrom}/{timestampTo}", HandleLogs)
	r.HandleFunc("/machines", HandleMachines)
	http.Handle("/", r)

	http.ListenAndServe(":6879", nil)
}
