package main

import "flag"

var (
	contextPath string
	port        int
	devMode     bool
)

func init() {
	contextPathArg := flag.String("contextPath", "", "context path")
	portArg := flag.Int("port", 6879, "Port to serve.")
	devModeArg := flag.Bool("devMode", false, "devMode(disable js/css minify)")

	flag.Parse()

	contextPath = *contextPathArg
	port = *portArg
	devMode = *devModeArg
}
