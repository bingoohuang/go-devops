package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Machines map[string]Machine
	Logs     map[string]Log
}

type Machine struct {
	IP string
}

type Log struct {
	Machines []string
	Path     string
	Process  string
}

var config Config

func init() {
	configFileArg := flag.String("config", "config.toml", "config file path")
	flag.Parse()

	configPath := *configFileArg
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return
	}

	_, err := toml.DecodeFile(configPath, &config)
	FatalIfErr(err)
}
