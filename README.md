# go-devops
devops based on go

## Snapshots
Home page:
![image](https://user-images.githubusercontent.com/1940588/34563611-30ef73d0-f18e-11e7-8318-27f259577246.png)
Tail page:
![image](https://user-images.githubusercontent.com/1940588/34563629-42f2ab74-f18e-11e7-8b39-b48c707db69a.png)
Config page:
![image](https://user-images.githubusercontent.com/1940588/34563633-49f3958c-f18e-11e7-8b64-871c4d1b40f8.png)

## Utils
1. [online png to ico](https://cloudconvert.com/png-to-ico)
2. [online png to svg](https://image.online-convert.com/convert-to-svg)

## ps aux output meaning
```
$ ps aux  
USER       PID  %CPU %MEM  VSZ RSS     TTY   STAT START   TIME COMMAND
timothy  29217  0.0  0.0 11916 4560 pts/21   S+   08:15   0:00 pine  
root     29505  0.0  0.0 38196 2728 ?        Ss   Mar07   0:00 sshd: can [priv]   
can      29529  0.0  0.0 38332 1904 ?        S    Mar07   0:00 sshd: can@notty  
```

```
ps aux|head -1 
# Returns: USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND

#ps aux with ppid
ps axo user,pid,ppid,pcpu,pmem,vsz,rss,tname,stat,start,time,args|head -1 
# Returns: USER       PID  PPID %CPU %MEM    VSZ   RSS TTY      STAT  STARTED     TIME COMMAND
```

    USER = user owning the process
    PID = process ID of the process
    %CPU = It is the CPU time used divided by the time the process has been running.
    %MEM = ratio of the process’s resident set size to the physical memory on the machine
    VSZ = virtual memory usage of entire process (in KiB)
    RSS = resident set size, the non-swapped physical memory that a task has used (in KiB)
    TTY = controlling tty (terminal)
    STAT = multi-character process state
    START = starting time or date of the process
    TIME = cumulative CPU time
    COMMAND = command with all its arguments

Process state codes:

    R running or runnable (on run queue)
    D uninterruptible sleep (usually IO)
    S interruptible sleep (waiting for an event to complete)
    Z defunct/zombie, terminated but not reaped by its parent
    T stopped, either by a job control signal or because it is being traced

Some extra modifiers:

    < high-priority (not nice to other users)
    N low-priority (nice to other users)
    L has pages locked into memory (for real-time and custom IO)
    s is a session leader
    l is multi-threaded (using CLONE_THREAD, like NPTL pthreads do)
    + is in the foreground process group

## Some scripts

# Count cost time of command
```bash
#!/bin/bash
START=$(date +%s)

awk 'substr($0,1,19)>="2018-01-04 16:07:18" && substr($0,1,19)<="2018-01-04 16:07:19"' < src.log > cut.log

END=$(date +%s)
DIFF=$(( $END - $START ))
echo "It took $DIFF seconds"

```

## SSH from host1 to host2 without a password on linux
```bash
[yogaapp@host1 ~]$ ssh-keygen -t rsa 
[yogaapp@host1 ~]$ ssh-copy-id -i .ssh/id_rsa.pub cpapp@host2
```

## Startup JAVA program
Go program os.Getenv("PATH") print PATH: /usr/local/bin:/bin:/usr/bin, So we need link java executable to /usr/bin or /usr/local/bin.

```bash
[root@hitest.app.02 ~]# which java
/opt/jdk1.8.0_20/bin/java
[root@hitest.app.02 ~]# ln -s /opt/jdk1.8.0_20/bin/java /usr/bin/java
[root@hitest.app.02 ~]# ls -l /usr/bin/java
lrwxrwxrwx 1 root root 25 Jan 10 10:56 /usr/bin/java -> /opt/jdk1.8.0_20/bin/java
```

## ls -l command output explanation
[From](https://superuser.com/questions/171858/how-do-i-interpret-the-results-of-the-ls-l-command)
<pre>

      +-permissions that apply to the owner
      |
      |     +-permissions that apply to all other users
      |     |
      |     |  +-number of hard links
      |     |  |
      |     |  |             +-size      +-last modification date and time
     _|_   _|_ |            _|__ ________|_______
    drwxr-xr-x 2 ataka root 4096 2008-11-04 16:58 ataka
        ___      _____ ____                       _____
         |         |    |                           |
         |         |    |                           +-name of file or directory
         |         |    |
         |         |    +-the group that the group permissions applies to
         |         |
         |         +-owner
         |
         +-permissions that apply to users who are members of the group
         
</pre>

## Linux中查看各文件夹大小命令
```bash
# -h或–human-readable 以K，M，G为单位，提高信息的可读性。
du -h --max-depth=1
```
## [Limit log file size.](https://www.bluedelta.nl/linux/limit-log-file-size/)
```bash
#!/bin/bash
logfile=$1

# truncate a log-file at 50000 lines.
if [ -f ${logfile} ]
then
    ls -lh ${logfile}
    tail -50000 $logfile > $logfile.tmp
    cat $logfile.tmp > $logfile
    rm $logfile.tmp
    ls -lh ${logfile}
else
    echo "File $logfile not found"
fi
```

## Search log context lines by keyword
```awk
#!/bin/awk -f
BEGIN {
    phead = 1 # head of previous lines queue
    ptail = 1 # size of previous lines queue
    pmax = 3 # max size of previous lines queue
    found = 0 # found lines from keyword
    max =  2 # max lines from keyword
}
{
    if (found == 0) {
      if ($0 ~ /7/) {
        found = 1
        for(key in parr) print parr[key]
        print
      } else {
        parr[ptail++] = $0
        if(ptail - phead > pmax) delete parr[phead++]
      }
    } else {
      print
      if(++found == max) exit
    }
}
END {
}

``` 


##  Show top five CPU consuming processes with `ps`
```bash
# n for numeric sorting, r for reverse order and k 2,2 for the second column.
ps aux | sort -nrk 3,3 | head -n 5

# If you want something that's truly 'top'esq with constant updates, use watch
watch "ps aux | sort -nrk 3,3 | head -n 5"
```

## Cron Expression
```
0 0 8-20 * * ?

```
####[CRON Expression Format](https://godoc.org/github.com/robfig/cron)

A cron expression represents a set of times, using 6 space-separated fields.

Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

####Predefined schedules 

You may use one of several pre-defined schedules in place of a cron expression.

Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *

####Intervals
You may also schedule a job to execute at fixed intervals, starting at the time it's added or cron is run. This is supported by formatting the cron spec like this:

`@every <duration>`

where "duration" is a string accepted by [time.ParseDuration](http://golang.org/pkg/time/#ParseDuration).

`A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`

For example, "@every 1h30m10s" would indicate a schedule that activates after 1 hour, 30 minutes, 10 seconds, and then every interval after that.