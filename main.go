package main

import (
	"os"
	"net/rpc"
	"log"
	"time"
	"fmt"
)

var config Config

func main() {
	config = ReadConfig()
	fmt.Println(config)

	err := GoServer()
	if err == nil {
		GoHttpSever()
	}

	if len(os.Args) > 1 && os.Args[1] == "-client" {
		serverAddress := "127.0.0.1"
		client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}

		defer client.Close()

		CallShellCommandService(client, os.Args[2])
	} else {
		for {
			time.Sleep(10 * time.Second)
		}
	}

}
