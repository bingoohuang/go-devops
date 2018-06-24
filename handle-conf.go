package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/go-utils"
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

	go_utils.HeadContentTypeJson(w)
	json.NewEncoder(w).Encode(struct{ OK, Conf string }{"OK", string(conf)})
}

func HandleSaveConf(w http.ResponseWriter, r *http.Request) {
	go_utils.HeadContentTypeJson(w)
	config := r.FormValue("config")

	ioutil.WriteFile(configFile, []byte(config), 0644)
	meta, err := toml.Decode(config, &devopsConf)
	if err != nil {
		json.NewEncoder(w).Encode(struct{ OK, Msg string }{"ERROR", err.Error()})
		return
	}

	parseConfig(&meta)
	json.NewEncoder(w).Encode(struct{ OK, Msg string }{"OK", "OK"})
}
