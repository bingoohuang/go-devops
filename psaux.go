package main

import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode"
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
	// refer: go tool dist list -json
	if runtime.GOOS == "darwin" {
		return PsAux(psauxNoHeading)
	}

	return PsAux(psauxNoHeading + `--sort=-pid --forest`)
}

var psaux = `ps axo user,pid,ppid,pcpu,pmem,vsz,rss,tname,stat,start,time,args `
var psauxNoHeading = CreatePsAuxNoHeading()

func CreatePsAuxNoHeading() string {
	if runtime.GOOS == "darwin" {
		return psaux + `|sed '1d'`
	}
	return psaux + ` --no-heading `
}

func PsAuxGrep(keywords ...string) []*PsAuxItem {
	shellCmd := psauxNoHeading

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
	return PsAux(psauxNoHeading + `--sort=-pcpu|head -n ` + strconv.Itoa(n))
}

var BlankRegex = regexp.MustCompile(`\s+`)

var PasAuxHeadingStartedEndIndex int
var PasAuxHeadingHasTty bool

type PsAuxHeading struct {
}

// USER       PID  PPID %CPU %MEM    VSZ   RSS TTY      STAT  STARTED     TIME COMMAND
// root         1     0  0.0  0.0  19356   384 ?        Ss     Oct 26 00:00:16 /sbin/init
// root     14857 12361  0.0  0.2 100012  3956 ?        Ss   14:22:50 00:00:00  \_ sshd: app [priv]
func (t *PsAuxHeading) Lining(line string) {
	started := "STARTED"
	index := strings.Index(line, started)
	PasAuxHeadingStartedEndIndex = index + len(started)
	PasAuxHeadingHasTty = strings.Index(line, "TTY") >= 0
}

func init() {
	ExecuteBash(psaux+`|head -n 1`, &PsAuxHeading{})
}

type PsAuxLiner struct {
	AuxItems []*PsAuxItem
}

func NewPsAuxLiner() *PsAuxLiner {
	return &PsAuxLiner{AuxItems: make([]*PsAuxItem, 0)}
}

func (t *PsAuxLiner) Lining(line string) {
	offset := strings.IndexFunc(line[PasAuxHeadingStartedEndIndex:], func(r rune) bool {
		return unicode.IsSpace(r)
	})

	n := 9
	if PasAuxHeadingHasTty {
		n = 10
	}
	f := BlankRegex.Split(line[0:PasAuxHeadingStartedEndIndex+offset], n)
	otherPart := line[PasAuxHeadingStartedEndIndex+offset+1:]
	otherPart = strings.TrimSpace(otherPart)
	spaceIndex := strings.IndexFunc(otherPart, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	time := otherPart[0:spaceIndex]
	command := otherPart[spaceIndex+1:]

	var item *PsAuxItem = nil
	if PasAuxHeadingHasTty {
		item = &PsAuxItem{
			User: f[0], Pid: f[1], Ppid: f[2], Cpu: f[3], Mem: f[4], Vsz: f[5], Rss: f[6],
			Tty: f[7], Stat: f[8], Start: f[9], Time: time, Command: command}

	} else {
		item = &PsAuxItem{
			User: f[0], Pid: f[1], Ppid: f[2], Cpu: f[3], Mem: f[4], Vsz: f[5], Rss: f[6],
			Tty: "", Stat: f[7], Start: f[8], Time: time, Command: command}
	}
	t.AuxItems = append(t.AuxItems, item)
}

func PsAux(shellCmd string) []*PsAuxItem {
	liner := NewPsAuxLiner()
	ExecuteBash(shellCmd, liner)

	return liner.AuxItems
}

type BashOutputLiner interface {
	Lining(line string)
}

// 执行Shell脚本，返回行解析对象数组
func ExecuteBash(shellScripts string, liner BashOutputLiner) error {
	cmd := exec.Command("bash", "-c", shellScripts)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}

	cmd.Start()
	defer cmd.Process.Kill()
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		liner.Lining(line)
	}

	return nil
}
