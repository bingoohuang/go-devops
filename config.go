package main

import (
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
}

var config Config

func init() {
	configPath := "config.toml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return
	}

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	_, err := toml.DecodeFile(configPath, &config)
	FatalIfErr(err)
}
