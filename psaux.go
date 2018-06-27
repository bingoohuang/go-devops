package main

import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type PsAuxItem struct {
	User    string
	Pid     string
	Ppid    string
	Cpu     string
	Mem     string
	Vsz     string
	Rss     string
	Tty     string
	Stat    string
	Start   string
	Time    string
	Command string
}

func PsAuxAll() []*PsAuxItem {
	return PsAux(psaux + `|sed '1d'|sort -nrk 2,2`)
}

var psaux = `ps axo user,pid,ppid,pcpu,pmem,vsz,rss,tname,stat,start,time,args`

func PsAuxGrep(keywords ...string) []*PsAuxItem {
	shellCmd := psaux + `|sed '1d'`

	for _, keyword := range keywords {
		if keyword != "" {
			shellCmd += `|grep ` + keyword
		}
	}

	if len(keywords) > 0 {
		shellCmd += `|grep -v grep`
	}

	return PsAux(shellCmd)
}

func PsAuxTop(n int) []*PsAuxItem {
	return PsAux(psaux + `|sed '1d'|sort -nrk 4,4|head -n ` + strconv.Itoa(n))
}

var BlankRegex = regexp.MustCompile(`\s+`)

func PsAux(shellCmd string) []*PsAuxItem {
	items := ExecuteBash(shellCmd, func(line string) interface{} {
		// USER  PID  %CPU %MEM  VSZ RSS TT  STAT STARTED TIME COMMAND
		f := BlankRegex.Split(line, 12)
		return &PsAuxItem{
			User: f[0], Pid: f[1], Ppid: f[2], Cpu: f[3], Mem: f[4], Vsz: f[5], Rss: f[6],
			Tty: f[7], Stat: f[8], Start: f[9], Time: f[10], Command: f[11]}
	})

	var auxItems []*PsAuxItem
	for _, item := range items {
		auxItems = append(auxItems, item.(*PsAuxItem))
	}

	return auxItems
}

// 执行Shell脚本，返回行解析对象数组
func ExecuteBash(shellScripts string, lineFunc func(line string) interface{}) []interface{} {
	cmd := exec.Command("bash", "-c", shellScripts)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}

	cmd.Start()
	defer cmd.Process.Kill()
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	var items []interface{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}

		items = append(items, lineFunc(strings.TrimSpace(line)))
	}

	return items
}
