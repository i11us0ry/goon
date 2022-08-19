package brute

import (
	"context"
	"fmt"
	"goon3/public"
	"net"
	"runtime"
	"strings"
	"time"
)


func Redis(hosts []string) {
	/* CPU */
	runtime.GOMAXPROCS(runtime.NumCPU())
	user, pass := []string{},[]string{}
	if len(Par.User)>0{
		user = Par.User
	} else {
		user = Par.Brute.Redis.Info.User
	}
	if len(Par.Pass)>0{
		pass = Par.Pass
	} else {
		pass = Par.Brute.Redis.Info.Pass
	}
	/* 建立chan，传输待扫描数据和扫描结果 */
	input := make(chan BruteInfo, len(hosts))
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
	/* 启动扫描任务 */
	for i := 0; i< thread; i++{
		go redisWork(input,result)
	}

	/* 输出 */
	public.Out(result,Par.Ofile)
}

func redisWork(input chan BruteInfo,result chan string){
	for {
		task,ok := <-input
		if !ok{
			return
		}
		scan := make(chan string)
		ctx, cancel := context.WithCancel(context.Background())
		go redisUnauth(ctx,task,scan)
		select {
		case flag,ok:= <-scan:
			if flag == "" || !ok{
				result<-""
			} else {
				/* 爆破成功 */
				result <- flag
			}
		case <- time.After(time.Duration(Par.Timeout) * time.Duration(len(Par.Pass)+1) * time.Second):
			result<-""
		}
		cancel()
		time.Sleep(1*time.Second)
	}
}

/*
未授权访问
@info	    host信息
return:
string
*/
func redisUnauth(ctx context.Context,info BruteInfo,scan chan string) {
	select {
	case <- ctx.Done():
		return
	default:
		/* 建立连接 */
		conn, err := net.DialTimeout("tcp", info.Host, time.Duration(Par.Timeout)*time.Second)
		if err != nil {
			scan <- ""
		}
		/* 设置读写超时 */
		err = conn.SetReadDeadline(time.Now().Add(time.Duration(Par.Timeout)*time.Second))
		if err!= nil{
			scan <- ""
		}
		err = conn.SetWriteDeadline(time.Now().Add(time.Duration(Par.Timeout)*time.Second))
		if err!= nil{
			scan <- ""
		}
		defer conn.Close()
		/* 发送info命令查看是否存在未授权 */
		_, err = conn.Write([]byte("info\r\n"))
		if err != nil {
			scan <- ""
		}
		/* 读取redis答复 */
		reply, err := readreply(conn)
		if err != nil || reply==""{
			scan <- ""
		}
		/* 如果relpy存在redis_version则说明存在未授权 */
		if strings.Contains(reply, "redis_version") {
			result := fmt.Sprintf("[Redis] %-40s unauthorized", info.Host)
			/* 输出成功 */
			scan <- result
		} else if strings.Contains(reply, "Authentication required") {
			/* 如果relpy存在Authentication required则说明需要输入密码（确认资产是否属于redis） */
			for _,pass := range(info.Passs){
				/* 输入认证密码 */
				_, err = conn.Write([]byte(fmt.Sprintf("auth %s\r\n", pass)))
				if err != nil {
					scan <- ""
				}
				reply, err := readreply(conn)
				if err != nil {
					scan <- ""
				}
				if strings.Contains(reply, "+OK") {
					result := fmt.Sprintf("[Redis] %-40s %s", info.Host, pass)
					scan <- result
				} else if strings.Contains(reply, "password") {
					continue
				} else {
					scan <- ""
				}
			}
		} else {
			scan <- ""
		}
	}
}


/*
读取答复信息
@conn		redis连接
return
result
err
*/
func readreply(conn net.Conn) (result string, err error) {
	buf := make([]byte, 4096)
	for {
		count, err := conn.Read(buf)
		if err != nil{
			break
		}
		result += string(buf[0:count])
		if count < 4096 {
			break
		}
	}
	return result, err
}
