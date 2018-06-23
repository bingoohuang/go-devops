package main

import (
	"flag"
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/shirou/gopsutil/host"
	"os"
	"strconv"
)

var (
	contextPath  *string
	httpPort     string
	startHttp    *bool
	rpcPort      string
	devMode      *bool
	configFile   *string
	randomLogGen bool
	hostname     string

	machineNames []string
	loggers      []string

	authParam go_utils.MustAuthParam
)

func init() {
	contextPath = flag.String("contextPath", "", "context path")
	httpPortArg := flag.Int("httpPort", 6879, "Port to serve.")
	startHttp = flag.Bool("startHttp", true, "Port to serve.")
	rpcPortArg := flag.Int("rpcPort", 6979, "Port to serve.")
	devMode = flag.Bool("devMode", false, "devMode(disable js/css minify)")
	configFile = flag.String("config", "config.toml", "config file path")
	versionArg := flag.Bool("v", false, "print version")

	go_utils.PrepareMustAuthFlag(&authParam)

	flag.Parse()

	if *versionArg {
		fmt.Println("Version 0.1.0")
		os.Exit(0)
	}

	httpPort = strconv.Itoa(*httpPortArg)
	rpcPort = strconv.Itoa(*rpcPortArg)
	ostStat, _ := host.Info()
	hostname = ostStat.Hostname

	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		return
	}

	loadConfig()
}
