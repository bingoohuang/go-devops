package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/shirou/gopsutil/host"
	"os"
	"strconv"
	"time"
)

var (
	contextPath  string
	httpPort     string
	rpcPort      string
	devMode      bool
	configFile   string
	randomLogGen bool
	hostname     string
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
	directCmdsArg := flag.String("directCmds", "", "direct Cmds")
	randomLogGenArg := flag.Bool("randomLogGen", false, "random log generator to aaa.log")
	versionArg := flag.Bool("v", false, "print version")

	flag.Parse()

	if *versionArg {
		fmt.Println("Version 0.0.3")
		os.Exit(0)
	}

	if *directCmdsArg != "" {
		ExecuteCommands(*directCmdsArg, 3*time.Second)
		os.Exit(0)
	}

	contextPath = *contextPathArg
	httpPort = strconv.Itoa(*httpPortArg)
	rpcPort = strconv.Itoa(*rpcPortArg)
	devMode = *devModeArg
	configFile = *configFileArg
	randomLogGen = *randomLogGenArg

	ostStat, _ := host.Info()
	hostname = ostStat.Hostname

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return
	}

	_, err := toml.DecodeFile(configFile, &devopsConf)
	FatalIfErr(err)
}
