package brute

import (
	"context"
	"fmt"
	"goon3/lib/telnet"
	"goon3/public"
	"runtime"
	"strings"
	"time"
)

func Telnet(hosts []string){
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(Par.User)>0{
		dic.User = Par.User
	} else {
		dic.User = Par.Brute.Telnet.Info.User
	}
	if len(Par.Pass)>0{
		dic.Pass = Par.Pass
	} else {
		dic.Pass = Par.Brute.Telnet.Info.Pass
	}
	input := make(chan string,  len(hosts))
	result := make(chan string, len(hosts))
	defer close(input)
	for _, host := range(hosts){
		input <- host
	}
	thread := 10
	if len(hosts) < Par.Thread {
		thread = len(hosts)
	} else {
		thread = Par.Thread
	}
	/* 开启扫描 */
	for i:=0;i< thread;i++{
		go telnetWork(input,result)
	}
	/* 输出 */
	public.Out(result,Par.Ofile)
}


func telnetWork(input chan string,result chan string){
	for {
		task, ok := <-input
		if !ok {
			return
		}
		find := false
		// 由上下文统一管理
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		brute_result := make(chan string ,len(dic.User)*len(dic.Pass))
		defer close(brute_result)
		// 启动爆破
		for _, user := range(dic.User){
			for _, pass := range(dic.Pass){
				go telnetConn(ctx, task, user, pass, brute_result)
				select {
				case flag := <- brute_result:
					if flag != ""{
						result <- flag
						cancel()
						find = true
						break
					}
				case <- time.After(time.Duration(Par.Timeout+4) * time.Second):
					result <- ""
					cancel()
					find = true
					break
				}
			}
			if find{
				break
			}
		}
		if !find{
			result <- ""
		}
		time.Sleep(1*time.Second)
	}
}

func telnetConn(ctx context.Context, Host, User, Pass string,brute_result chan string){
	Pass = strings.Replace(Pass, "{{user}}", User, -1)
	t := new(telnet.TelnetClient)
	t.Host = Host
	t.UserName = User
	t.Password = Pass
	t.Time = Par.Timeout
	if flag := t.Conn();flag {
		select {
		case <- ctx.Done():
			return
		default:
			brute_result <- fmt.Sprintf("[Telnet] %-30v %v:%v",t.Host,t.UserName,t.Password)
		}
	} else {
		select {
		case <- ctx.Done():
			return
		default:
			brute_result <- ""
		}
	}
}
