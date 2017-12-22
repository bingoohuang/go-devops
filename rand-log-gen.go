package main

import (
	"time"
	"os"
	"math/rand"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createLog() {
	f, err := os.OpenFile("aa.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check(err)

	for {
		f.WriteString(time.Now().String() + " " + RandStringBytesMaskImpr(1000) + "\n")
		f.Sync()
		time.Sleep(1 * time.Second)
	}

}

// 用掩码进行替换
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
