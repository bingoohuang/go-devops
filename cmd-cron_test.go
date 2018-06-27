package main

import (
	"fmt"
	"github.com/metakeule/fmtdate"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func assertEquals(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func TestCopyTruncate(t *testing.T) {
	d1 := []byte("a23456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789b")
	file1 := "log-2018-02-07.log"
	ioutil.WriteFile(file1, d1, 0644)

	var cronCommand CronCommand
	var result CronCommandResult
	arg := &CronCommandArg{
		Files:      []string{file1},
		Type:       "CopyTruncate",
		Parameters: "maxSize=50",
	}
	cronCommand.Execute(arg, &result)
	bytes, _ := ioutil.ReadFile(file1)

	assertEquals(t, "1234567890123456789012345678901234567890123456789b", string(bytes), "")
	os.Remove(file1)
	os.Remove(file1 + ".tmp")
}

func TestDelete(t *testing.T) {
	ioutil.WriteFile("log-2018-02-07.log", []byte("a"), 0644)
	os.MkdirAll("logs-xxxyyy", 0777)

	var cronCommand CronCommand
	var result CronCommandResult
	arg := &CronCommandArg{
		Files: []string{"log-2018-02-07.log", "logs-xxxyyy"},
		Type:  "Delete",
	}
	cronCommand.Execute(arg, &result)

	_, e := os.Stat("log-2018-02-07.log")
	assertEquals(t, os.IsNotExist(e), true, "")
	_, e = os.Stat("logs-xxxyyy")
	assertEquals(t, os.IsNotExist(e), true, "")
}

func TestDeleleOlds(t *testing.T) {
	os.MkdirAll("logs", 0777)
	file1 := fmtdate.Format("logs/log-YYYY-MM-DD.log", time.Now())
	ioutil.WriteFile(file1, []byte("a"), 0644)

	file2 := fmtdate.Format("logs/log-YYYY-MM-DD.log", time.Now().AddDate(0, 0, -3))
	ioutil.WriteFile(file2, []byte("a"), 0644)

	var cronCommand CronCommand
	var result CronCommandResult
	arg := &CronCommandArg{
		Files:      []string{"./logs"},
		Type:       "DeleteOlds",
		Parameters: "days=3,pattern=log-YYYY-MM-DD.log",
	}
	cronCommand.Execute(arg, &result)

	_, e := os.Stat(file1)
	assertEquals(t, os.IsNotExist(e), false, "")

	_, e = os.Stat(file2)
	assertEquals(t, os.IsNotExist(e), true, "")

	os.RemoveAll("logs")
}
