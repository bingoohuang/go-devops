package main

import (
	"net"
	"net/rpc"
)

func GoDevOpServer() error {
	rpc.Register(new(ShellCommand))
	rpc.Register(new(MachineCommand))
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+rpcPort)
	FatalIfErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	go serve(listener)

	return nil
}

func serve(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
}
