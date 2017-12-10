package main

import (
	"net/rpc"
	"net"
	"net/http"
	"log"
)

func StartServer() {
	arith := new(Arith)
	rpc.Register(arith)
	shellCommand := new(ShellCommand)
	rpc.Register(shellCommand)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
