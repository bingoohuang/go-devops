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
	startHttp    bool
	rpcPort      string
	devMode      bool
	configFile   string
	randomLogGen bool
	hostname     string

	machineNames []string
	loggers      []string
)

type DevopsConf struct {
	Machines   map[string]Machine
	Logs       map[string]Log
	Processes  map[string]Process
	Logrotates map[string]LogRotate
}

type LogRotate struct {
	Machines   []string
	Files      []string
	Crons      []string
	Type       string
	Parameters string
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
}

var devopsConf DevopsConf

func init() {
	contextPathArg := flag.String("contextPath", "", "context path")
	httpPortArg := flag.Int("httpPort", 6879, "Port to serve.")
	startHttpArg := flag.Bool("startHttp", true, "Port to serve.")
	rpcPortArg := flag.Int("rpcPort", 6979, "Port to serve.")
	devModeArg := flag.Bool("devMode", false, "devMode(disable js/css minify)")
	configFileArg := flag.String("config", "config.toml", "config file path")
	directCmdsArg := flag.String("directCmds", "", "direct Cmds")
	randomLogGenArg := flag.Bool("randomLogGen", false, "random log generator to aaa.log")
	versionArg := flag.Bool("v", false, "print version")

	flag.Parse()

	if *versionArg {
		fmt.Println("Version 0.0.5")
		os.Exit(0)
	}

	if *directCmdsArg != "" {
		ExecuteCommands(*directCmdsArg, 3*time.Second)
		os.Exit(0)
	}

	contextPath = *contextPathArg
	httpPort = strconv.Itoa(*httpPortArg)
	startHttp = *startHttpArg
	rpcPort = strconv.Itoa(*rpcPortArg)
	devMode = *devModeArg
	configFile = *configFileArg
	randomLogGen = *randomLogGenArg

	ostStat, _ := host.Info()
	hostname = ostStat.Hostname

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return
	}

	loadConfig()
}

func loadConfig() {
	meta, err := toml.DecodeFile(configFile, &devopsConf)
	if err != nil {
		fmt.Println("DecodeFile error:", err)
		return
	}

	parseConfig(&meta)
}

func parseConfig(meta *toml.MetaData) {
	loadCrons()
	parseMeta(meta)
}

func parseMeta(meta *toml.MetaData) {
	machineNames = make([]string, 0)
	loggers = make([]string, 0)
	for _, key := range meta.Keys() {
		if len(key) != 2 {
			continue
		}

		if key[0] == "machines" {
			machineNames = append(machineNames, key[1])
		} else if key[0] == "logs" {
			loggers = append(loggers, key[1])
		}
	}
}
