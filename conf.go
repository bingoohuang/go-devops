package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

type DevopsConf struct {
	Machines       map[string]Machine
	Logs           map[string]Log
	Processes      map[string]Process
	Logrotates     map[string]LogRotate
	MessageTargets map[string]MessageTargetConf

	BlackcatThreshold    BlackcatThreshold
	BlackcatExLogs       map[string]BlackcatExLogConf
	BlackcatProcesses    map[string]BlackcatProcessConf
	BlackcatHttpCheckers map[string]BlackcatHttpChecker
	Misc                 MiscConf
}

type MiscConf struct {
	RedisServer string // redis server addr, eg: 127.0.0.1:6379, localhost:6388/0, password2/localhost:6388/0
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

func loadConfig() {
	meta, err := toml.DecodeFile(configFile, &devopsConf)
	if err != nil {
		log.Println("DecodeFile error:", err)
		return
	}

	redisServer = ParseServerItem(devopsConf.Misc.RedisServer)

	parseConfig(&meta)
}

func parseConfig(meta *toml.MetaData) {
	loadCrons()
	loadBlackcatCrons()
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
