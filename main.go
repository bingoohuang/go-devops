package main

func main() {
	if randomLogGen {
		createRandomLog()
		return
	}

	err := GoDevOpServer()
	FatalIfErr(err)

	StartHttpSever()
}
