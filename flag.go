package main

import (
	"flag"
	"github.com/bingoohuang/go-utils"
	"github.com/shirou/gopsutil/host"
	"log"
	"os"
	"strconv"
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

	authParam go_utils.MustAuthParam

	redisServer *RedisServer
)

func init() {
	flag.StringVar(&contextPath, "contextPath", "", "context path")
	httpPortArg := flag.Int("httpPort", 6879, "Port to serve.")
	flag.BoolVar(&startHttp, "startHttp", true, "Port to serve.")
	rpcPortArg := flag.Int("rpcPort", 6979, "Port to serve.")
	flag.BoolVar(&devMode, "devMode", false, "devMode(disable js/css minify)")
	flag.StringVar(&configFile, "config", "config.toml", "config file path")
	versionArg := flag.Bool("v", false, "print version")

	go_utils.PrepareMustAuthFlag(&authParam)

	flag.Parse()

	if *versionArg {
		log.Println("Version 0.1.1")
		os.Exit(0)
	}

	httpPort = strconv.Itoa(*httpPortArg)
	rpcPort = strconv.Itoa(*rpcPortArg)
	hostname = Hostname()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return
	}

	loadConfig()
	startBizCron()
}

func Hostname() string {
	ostStat, _ := host.Info()
	return ostStat.Hostname
}
