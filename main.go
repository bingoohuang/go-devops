package main

func main() {
	startMain()
}

func startMain() {
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
