package main

import (
	"reflect"
	"testing"
)

func TestPsAuxLiner(t *testing.T) {
	liner := NewPsAuxLiner()
	heading := &PsAuxHeading{}
	heading.Lining("USER       PID  PPID %CPU %MEM    VSZ   RSS TTY      STAT  STARTED     TIME COMMAND")
	liner.Lining("root         1     0  0.0  0.0  19356   384 ?        Ss     Oct 26 00:00:16 /sbin/init")
	if !reflect.DeepEqual(liner.AuxItems[0], &PsAuxItem{
		User: "root", Pid: "1", Ppid: "0", Cpu: "0.0", Mem: "0.0", Vsz: "19356", Rss: "384", Tty: "?", Stat: "Ss", Start: "Oct 26", Time: "00:00:16", Command: "/sbin/init"}) {
		t.Error("PsAuxLiner Failed")
	}

	liner.Lining(`root     14857 12361  0.0  0.2 100012  3956 ?        Ss   14:22:50 00:00:00  \_ sshd: app [priv]`)
	if !reflect.DeepEqual(liner.AuxItems[1], &PsAuxItem{
		User: "root", Pid: "14857", Ppid: "12361", Cpu: "0.0", Mem: "0.2", Vsz: "100012", Rss: "3956", Tty: "?", Stat: "Ss", Start: "14:22:50", Time: "00:00:00", Command: ` \_ sshd: app [priv]`}) {
		t.Error("PsAuxLiner Failed")
	}

	// USER       PID  PPID %CPU %MEM    VSZ   RSS TTY      STAT  STARTED     TIME COMMAND
	// root         1     0  0.0  0.0  19356  1084 ?        Ss     Apr 23 00:00:01 /sbin/init
	// root         2     0  0.0  0.0      0     0 ?        S      Apr 23 00:00:00 [kthreadd]
	// root         3     2  0.0  0.0      0     0 ?        S      Apr 23 00:00:00 [migration/0]
}

func TestPsAuxAll(t *testing.T) {
	ExecuteBash(psaux+`|head -n 1`, &PsAuxHeading{})
	all := PsAuxAll()

	if len(all) == 0 {
		t.Error("PsAuxAll Failed")
	}
}
