package main

import "net/http"

type LogRotatesResult struct {
	Logs        string
	Size        string
	MachineName string
	Error       string
	CostTime    string
}

func HandleLogRotates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

}
