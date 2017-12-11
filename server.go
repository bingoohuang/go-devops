package main

import (
	"net/rpc"
	"net"
	"net/http"
)

func GoServer() error {
	arith := new(Arith)
	rpc.Register(arith)
	shellCommand := new(ShellCommand)
	rpc.Register(shellCommand)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		return e
	}
	go http.Serve(l, nil)

	return nil
}
