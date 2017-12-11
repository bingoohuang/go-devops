package main

import (
	"net/rpc"
	"log"
	"fmt"
)

func CallService(client *rpc.Client, a, b int) {
	// Synchronous call
	args := &Args{a, b}
	var reply int
	err := client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
}
