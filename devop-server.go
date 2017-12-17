package main

import (
	"net/rpc"
	"net"
)

func GoDevOpServer() error {
	rpc.Register(new(ShellCommand))
	rpc.Register(new(MachineCommand))
	//rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":6979")
	//if e != nil {
	//	return e
	//}
	//go http.Serve(l, nil)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":6979")
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
