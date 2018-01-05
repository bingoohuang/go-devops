package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"net/http"
	"os"
)

func HandleLoadConf(w http.ResponseWriter, r *http.Request) {
	var conf []byte

	_, err := os.Stat(configFile)
	if err == nil {
		conf, _ = ioutil.ReadFile(configFile)
	} else {
		conf = []byte("")
	}

	result := struct {
		OK   string
		Conf string
	}{
		OK:   "OK",
		Conf: string(conf),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}

func HandleSaveConf(w http.ResponseWriter, r *http.Request) {
	config := r.FormValue("config")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ioutil.WriteFile(configFile, []byte(config), 0644)
	meta, err := toml.Decode(config, &devopsConf)
	if err != nil {
		json.NewEncoder(w).Encode(struct {
			OK  string
			Msg string
		}{
			OK:  "ERROR",
			Msg: err.Error(),
		})
		return
	}

	parseMetas(&meta)

	json.NewEncoder(w).Encode(struct {
		OK  string
		Msg string
	}{
		OK:  "OK",
		Msg: "OK",
	})
}
