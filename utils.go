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

func DialAndCall(machine Machine, callFun func(client *rpc.Client) error) error {
	conn, err := net.DialTimeout("tcp", machine.IP+":"+rpcPort, 1*time.Second)
	if err != nil {
		return err
	}

	client := rpc.NewClient(conn)
	defer client.Close()

	return callFun(client)
}
