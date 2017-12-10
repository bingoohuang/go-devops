package main

import (
	"os"
)

func main() {
	if os.Args[1] == "-server" {
		StartServer()
	} else if os.Args[1] == "-client" {
		CallShellCommandService(os.Args[2])
	}

}
