package main

func main() {
	if randomLogGen {
		createRandomLog()
		return
	}

	err := StartDevOpServer()
	FatalIfErr(err)

	if *startHttp {
		go StartHttpSever()
	}

	select {} // block forever
}
