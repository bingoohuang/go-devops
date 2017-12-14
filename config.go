package main

import (
	"os"
	"github.com/BurntSushi/toml"
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

func ReadConfig() Config {
	configPath := "config.toml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config := Config{}
	_, err := toml.DecodeFile(configPath, &config)
	FatalIfErr(err)

	return config
}
