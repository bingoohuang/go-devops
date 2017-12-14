package main

import (
	"net/rpc"
	"net"
	"net/http"
)

func GoServer() error {
	rpc.Register(new(ShellCommand))
	rpc.Register(new(MachineCommand))
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		return e
	}
	go http.Serve(l, nil)

	return nil
}
