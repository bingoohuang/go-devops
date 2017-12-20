package main

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

func FatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func DialAndCall(machine Machine,
	callFun func(client *rpc.Client, args interface{}) interface{},
	args interface{}) interface{} {
	conn, err := net.DialTimeout("tcp", machine.IP+":"+rpcPort, 1*time.Second)
	if err != nil {
		return MachineCommandResult{
			Error: err.Error(),
		}
	}

	client := rpc.NewClient(conn)
	defer client.Close()

	return callFun(client, args)
}
