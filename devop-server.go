package main

import (
	"net"
	"net/rpc"
)

func StartDevOpServer() error {
	rpc.Register(new(ShellCommand))
	rpc.Register(new(LogFileCommand))
	rpc.Register(new(MachineCommand))
	rpc.Register(new(CronCommand))
	rpc.Register(new(ExLogCommand))
	rpc.Register(new(AgentCommand))
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
