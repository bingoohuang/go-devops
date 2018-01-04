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
```bash
#!/bin/bash
START=$(date +%s)

awk 'substr($0,1,19)>="2018-01-04 16:07:18" && substr($0,1,19)<="2018-01-04 16:07:19"' < src.log > cut.log

END=$(date +%s)
DIFF=$(( $END - $START ))
echo "It took $DIFF seconds"

```