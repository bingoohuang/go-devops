package main

import (
	"errors"
	"github.com/hako/durafmt"
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

func RpcExecute(machineName string, arg interface{}, callable RpcCallable) {
	machine := devopsConf.Machines[machineName]
	addr := machine.IP + ":" + rpcPort
	RpcAddrExecute(machineName, addr, arg, callable)
}

func RpcAddrExecute(machineName, addr string, arg interface{}, callable RpcCallable) {
	RpcAddrCall(machineName, addr, "Execute", arg, callable)
}

func RpcCall(machineName, funcName string, arg interface{}, callable RpcCallable) {
	machine := devopsConf.Machines[machineName]
	addr := machine.IP + ":" + rpcPort
	RpcAddrCall(machineName, addr, funcName, arg, callable)
}

func RpcAddrCall(machineName, addr, funcName string, arg interface{}, callable RpcCallable) {
	machine := devopsConf.Machines[machineName]
	if addr == "" {
		addr = machine.IP + ":" + rpcPort
	}

	go func() {
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err != nil {
			return
		}

		client := rpc.NewClient(conn)
		defer client.Close()

		reply := callable.CreateResult(nil)
		client.Call(callable.CommandName()+"."+funcName, arg, reply)
	}()
}

func RpcExecuteTimeout(machineName string, arg interface{}, callable RpcCallable, timeout time.Duration, resultChan chan RpcResult) {
	machine := devopsConf.Machines[machineName]
	addr := machine.IP + ":" + rpcPort
	RpcAddrExecuteTimeout(machineName, addr, arg, callable, timeout, resultChan)
}

func RpcAddrExecuteTimeout(machineName, addr string, arg interface{}, callable RpcCallable, timeout time.Duration, resultChan chan RpcResult) {
	RpcAddrCallTimeout(machineName, addr, "Execute", arg, callable, timeout, resultChan)
}

func RpcCallTimeout(machineName, funcName string, arg interface{}, callable RpcCallable, timeout time.Duration, resultChan chan RpcResult) {
	machine := devopsConf.Machines[machineName]
	addr := machine.IP + ":" + rpcPort
	RpcAddrCallTimeout(machineName, addr, funcName, arg, callable, timeout, resultChan)
}

func RpcAddrCallTimeout(machineName, addr, funcName string, arg interface{}, callable RpcCallable, timeout time.Duration, resultChan chan RpcResult) {
	c := make(chan RpcResult)
	go func() {
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err != nil {
			result := callable.CreateResult(err)
			result.SetMachineName(machineName)
			c <- result
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
	case <-time.After(timeout):
		fmt := durafmt.ParseShort(timeout)
		result := callable.CreateResult(errors.New("timeout in " + fmt.String() + " to call " + machineName + " " + callable.CommandName() + "." + funcName + "@" + addr))
		result.SetMachineName(machineName)
		resultChan <- result
	}
}
