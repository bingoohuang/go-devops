package main

import (
	"net"
	"net/rpc"
)

func StartDevOpServer() error {
	_ = rpc.Register(new(ShellCommand))
	_ = rpc.Register(new(ShellResultCommand))
	_ = rpc.Register(new(LogFileCommand))
	_ = rpc.Register(new(MachineCommand))
	_ = rpc.Register(new(CronCommand))
	_ = rpc.Register(new(ExLogCommand))
	_ = rpc.Register(new(AgentCommand))
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
