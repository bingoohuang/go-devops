package main

func main() {
	startMain()
}

func startMain() {
	err := StartDevOpServer()
	FatalIfErr(err)

	if appConfig.StartHttp {
		go StartHttpSever()
	}

	select {} // block forever
}
