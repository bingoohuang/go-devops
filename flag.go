package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"os"
	"strconv"
)

var (
	contextPath  string
	httpPort     string
	rpcPort      string
	devMode      bool
	configFile   string
	randomLogGen bool
)

type DevopsConf struct {
	Machines  map[string]Machine
	Logs      map[string]Log
	Processes map[string]Process
}

type Machine struct {
	IP string
}

type Log struct {
	Machines []string
	Path     string
	Process  string
}

type Process struct {
	Home  string
	Ps    string
	Kill  string
	Start string
	Conf  string
}

var devopsConf DevopsConf

func init() {
	contextPathArg := flag.String("contextPath", "", "context path")
	httpPortArg := flag.Int("httpPort", 6879, "Port to serve.")
	rpcPortArg := flag.Int("rpcPort", 6979, "Port to serve.")
	devModeArg := flag.Bool("devMode", false, "devMode(disable js/css minify)")
	configFileArg := flag.String("config", "config.toml", "config file path")
	randomLogGenArg := flag.Bool("randomLogGen", false, "random log generator to aaa.log")

	flag.Parse()

	contextPath = *contextPathArg
	httpPort = strconv.Itoa(*httpPortArg)
	rpcPort = strconv.Itoa(*rpcPortArg)
	devMode = *devModeArg
	configFile = *configFileArg
	randomLogGen = *randomLogGenArg

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return
	}

	_, err := toml.DecodeFile(configFile, &devopsConf)
	FatalIfErr(err)
}
