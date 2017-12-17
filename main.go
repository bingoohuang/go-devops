package main

var config Config

func main() {
	config = ReadConfig()

	err := GoDevOpServer()
	FatalIfErr(err)

	StartHttpSever()
}
