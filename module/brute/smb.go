package brute

import (
	"context"
	"fmt"
	"goon3/public"
	"runtime"
	"strings"
	"time"
	"github.com/stacktitan/smb/smb"
)

func Smb(hosts []string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	user, pass := []string{},[]string{}
	if len(Par.User)>0{
		user = Par.User
	} else {
		user = Par.Brute.Smb.Info.User
	}
	if len(Par.Pass)>0{
		pass = Par.Pass
	} else {
		pass = Par.Brute.Smb.Info.Pass
	}
	input := make(chan BruteInfo,  len(hosts))
	result := make(chan string, len(hosts))
	defer close(input)

	/* 添加扫描数据 */
	go func() {
		CreateInfo(&input,hosts,user,pass)
	}()
	thread := 10
	if len(hosts) < Par.Thread {
		thread = len(hosts)
	} else {
		thread = Par.Thread
	}
	/* 开启扫描 */
	for i:=0;i< thread;i++{
		go smbWork(input,result)
	}

	/* 输出 */
	public.Out(result,Par.Ofile)
}

func smbWork(input chan BruteInfo,result chan string){
	for {
		task, ok := <-input
		if !ok {
			return
		}
		pool := BrutePool{}
		brute_result := make(chan string, 1)
		brute_pool := make(chan BrutePool, len(task.Users)*len(task.Passs))
		defer close(brute_pool)
		defer close(brute_result)
		// 由上下文统一管理
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Par.Timeout+1)*time.Second)
		defer cancel()
		// 新建资源池
		for _, user := range (task.Users) {
			for _, pass := range (task.Passs) {
				pool.Host = task.Host
				pool.User = user
				pool.Pass = pass
				brute_pool <- pool
			}
		}
		// 启动爆破
		for i := 0; i < Par.Thread; i++ {
			go SmblConn(ctx, brute_pool, brute_result)
		}
		// 等待回显
		select {
		case <-ctx.Done():
			result <- ""
		case flag := <-brute_result:
			result <- flag
			cancel()
		}
		time.Sleep(1*time.Second)
	}
}

func SmblConn(ctx context.Context,Pool chan BrutePool,scan chan string) {
	for{
		Pool, ok := <-Pool
		if !ok{
			return
		}
		Hosts, User, Pass := Pool.Host, Pool.User, Pool.Pass
		Host := strings.SplitN(Hosts,":",-1)
		Pass = strings.Replace(Pass, "{{user}}", User, -1)
		options := smb.Options{
			Host:        Host[0],
			Port:        445,
			User:        User,
			Password:    Pass,
			Domain:      "",
			Workstation: "",
		}
		session, err := smb.NewSession(options, false)
		if err == nil {
			session.Close()
			if session.IsAuthenticated {
				result := fmt.Sprintf("[SMB] %-40v %v:%-20v", Host[0], User, Pass)
				select {
				case <- ctx.Done():
					return
				default:
					scan <- result
				}
			}
		}
	}
}
