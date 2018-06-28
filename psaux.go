package main

import (
	"regexp"
	"strconv"
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
	return PsAux(prefix + `--sort=-pid --forest`)
}

func PsAuxTop(n int) []*PsAuxItem {
	return PsAux(prefix + `--sort=-pcpu|head -n ` + strconv.Itoa(n))
}

var prefix = `ps axo lstart,user,pid,ppid,pcpu,pmem,vsz,rss,tname,stat,time,args --no-heading `

// STARTED USER       PID  PPID %CPU %MEM    VSZ   RSS TTY      STAT     TIME COMMAND
// 2018-06-01 10:14:43     haproxy 1750 1 0.2 0.0 14716 3800 ? Ss 01:25:08 sbin/haproxy -f conf/app-haproxy.conf
// 2018-06-01 10:13:16     root 1253 1 0.0 0.0 4064 548 tty6 Ss+ 00:00:00 /sbin/mingetty /dev/tty6
// 2018-06-01 10:13:16     root 1251 1 0.0 0.0 4064 548 tty5 Ss+ 00:00:00 /sbin/mingetty /dev/tty5

var BlankRegex = regexp.MustCompile(`\s+`)

var fixedLtime = `|awk '{c="date -d\""$1 FS $2 FS $3 FS $4 FS $5"\" +\047%Y-%m-%d %H:%M:%S\047"; c|getline d; close(c); $1=$2=$3=$4=$5=""; printf "%s\n",d$0 }'`

func PsAux(shellCmd string) []*PsAuxItem {
	auxItems := make([]*PsAuxItem, 0)

	shell := shellCmd + fixedLtime
	ExecuteBash(shell, func(line string) bool {
		f := BlankRegex.Split(line, 13)
		auxItems = append(auxItems, &PsAuxItem{
			User: f[2], Pid: f[3], Ppid: f[4], Cpu: f[5], Mem: f[6], Vsz: f[7], Rss: f[8],
			Tty: f[9], Stat: f[10], Start: f[0] + ` ` + f[1], Time: f[11], Command: f[12]})
		return true
	})

	return auxItems
}
