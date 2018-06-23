package main

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os/exec"
)

type Tailer interface {
	line(line string)
	error(err error)
}

func Tailf(logFile string, tailer Tailer, stop chan bool) {
	expanded, _ := homedir.Expand(logFile)
	cmd := exec.Command("bash", "-c", "tail -F "+expanded)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		tailer.error(err)
		return
	}

	if err := cmd.Start(); err != nil {
		tailer.error(err)
		return
	}
	defer cmd.Wait()

	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				tailer.error(err)
				fmt.Println("read err", err.Error())
				stop <- true
				return
			}

			tailer.line(line)
		}
	}()

	<-stop
	stdout.Close()

	fmt.Println("Tailf stopped", expanded)
}

/*
type LogTailer struct {
}

func (t *LogTailer) line(line string) {
	fmt.Print(line)
}
func (t *LogTailer) error(err error) {
	fmt.Println(err)
}

func main() {
	stop := make(chan bool)
	go Tailf("./a.log", &LogTailer{}, stop)
	time.Sleep(10 * time.Second)
	stop <- true
}

*/
