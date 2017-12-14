package main

import "log"

func FatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
