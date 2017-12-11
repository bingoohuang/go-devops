package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)

func GoHttpSever() {
	r := mux.NewRouter()
	r.HandleFunc("/log/{logger}/{timestamp}", HandleLog)
	http.Handle("/", r)

	go http.ListenAndServe(":6879", nil)
}

// http://127.0.0.1:6879/log/yoga-system/2015-07-07%2011:23:33
func HandleLog(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	logger := vars["logger"]
	timestamp := vars["timestamp"]
	fmt.Println("logger:", logger, " timestamp:", timestamp)

	command := `sed -n "/2015-07-07 11:23:33/,/2015-07-07 15:00:33/p" file.log`
	w.Write([]byte(command))

}
