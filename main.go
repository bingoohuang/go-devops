package main

func main() {
	err := GoDevOpServer()
	FatalIfErr(err)

	StartHttpSever()
}
