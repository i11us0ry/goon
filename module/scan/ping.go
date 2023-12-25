package scan

import (
	"bytes"
	"goon3/public"
	"os/exec"
	"runtime"
	"strings"
)

func Ping(ips []string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	input := make(chan string, len(ips))
	result := make(chan string, len(ips))
	defer close(input)
	for _, ip := range ips {
		input <- ip
	}
	thread := 10
	if len(ips) < Par.Thread {
		thread = len(ips)
	} else {
		thread = Par.Thread
	}
	for i := 0; i < thread; i++ {
		go pingWork(input, result)
	}
	public.Out(result)
}

func pingWork(input, result chan string) {
	for {
		task, ok := <-input
		if !ok {
			return
		}
		alive := PingScan(task)
		if alive {
			result <- "[ping] " + task
			Par.IpAlive = append(Par.IpAlive, task)
		} else {
			result <- ""
		}
	}
}

func PingScan(ip string) bool {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("cmd", "/c", "ping -n 1 -w 1 "+ip+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	case "darwin":
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -W 1 "+ip+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	default: //linux
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -w 1 "+ip+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	}
	outinfo := bytes.Buffer{}
	command.Stdout = &outinfo
	err := command.Start()
	if err != nil {
		return false
	}
	if err = command.Wait(); err != nil {
		return false
	} else {
		if strings.Contains(outinfo.String(), "true") && strings.Count(outinfo.String(), ip) > 2 {
			return true
		} else {
			return false
		}
	}
}
