package main

import (
	"errors"
	"net"
	"net/rpc"
	"time"
)

type RpcResult interface {
	GetMachineName() string
	SetMachineName(machineName string)
	GetError() string
	SetError(err error)
}

type RpcCallable interface {
	CreateResult(err error) RpcResult
	CommandName() string
}

func RpcCall(machineName, addr string, arg interface{}, callable RpcCallable) {
	resultChan := make(chan RpcResult)
	RpcCallTimeout(machineName, addr, "Execute", arg, callable, 3*time.Second, resultChan)
}

func RpcCallTimeout(machineName, addr, funcName string, arg interface{}, callable RpcCallable, executionTimeout time.Duration, resultChan chan RpcResult) {
	machine := devopsConf.Machines[machineName]
	if addr == "" {
		addr = machine.IP + ":" + rpcPort
	}

	c := make(chan RpcResult)
	go func() {
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err != nil {
			c <- callable.CreateResult(err)
			return
		}

		client := rpc.NewClient(conn)
		defer client.Close()

		reply := callable.CreateResult(nil)
		err = client.Call(callable.CommandName()+"."+funcName, arg, reply)
		if err != nil {
			reply.SetError(err)
		}

		c <- reply
	}()

	select {
	case result := <-c:
		result.SetMachineName(machineName)
		resultChan <- result
	case <-time.After(executionTimeout):
		result := callable.CreateResult(errors.New("timeout in 3 seconds"))
		result.SetMachineName(machineName)
		resultChan <- result
	}
}
