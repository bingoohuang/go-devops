package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/go-utils"
	"github.com/shirou/gopsutil/host"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	httpPort string
	rpcPort  string
	hostname string

	machineNames []string
	loggers      []string

	redisServer *RedisServer
)

type AppConfig struct {
	ContextPath string
	HttpPort    int
	StartHttp   bool
	RpcPort     int
	ConfigFile  string
	DevMode     bool

	EncryptKey  string
	CookieName  string
	RedirectUri string
	LocalUrl    string
	ForceLogin  bool
}

var appConfig AppConfig

var authParam go_utils.MustAuthParam

func init() {
	var configFile string

	versionArg := flag.Bool("v", false, "print version")
	flag.StringVar(&configFile, "configFile", "appConfig.toml", "app config file path")

	flag.Parse()
	if *versionArg {
		log.Println("Version 0.1.1")
		os.Exit(0)
	}

	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		if _, err := toml.Decode(`
ContextPath = ""
HttpPort  = 6879
StartHttp = false
RpcPort = 6979
DevMode = false
ConfigFile = "config.toml"
`, &appConfig); err != nil {
			log.Panic("config file decode error", err.Error())
		}
	} else {
		if _, err := toml.DecodeFile(configFile, &appConfig); err != nil {
			log.Panic("config file decode error", err.Error())
		}
	}

	if appConfig.ContextPath != "" && strings.Index(appConfig.ContextPath, "/") < 0 {
		appConfig.ContextPath = "/" + appConfig.ContextPath
	}

	authParam = go_utils.MustAuthParam{
		EncryptKey:  appConfig.EncryptKey,
		CookieName:  appConfig.CookieName,
		RedirectUri: appConfig.RedirectUri,
		LocalUrl:    appConfig.LocalUrl,
		ForceLogin:  appConfig.ForceLogin,
	}

	httpPort = strconv.Itoa(appConfig.HttpPort)
	rpcPort = strconv.Itoa(appConfig.RpcPort)
	hostname = Hostname()

	if _, err := os.Stat(appConfig.ConfigFile); os.IsNotExist(err) {
		return
	}

	loadConfig()
	startBizCron()
}

func Hostname() string {
	ostStat, _ := host.Info()
	return ostStat.Hostname
}
