package brute

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goon3/public"
	"runtime"
	"strings"
	"time"
)

func Mysql(hosts []string) {
	/* CPU核 */
	runtime.GOMAXPROCS(runtime.NumCPU())
	input := make(chan BruteInfo,  len(hosts))
	result := make(chan string, len(hosts))
	defer close(input)
	user, pass := []string{},[]string{}
	if len(Par.User)>0{
		user = Par.User
	} else {
		user = Par.Brute.Mssql.Info.User
	}
	if len(Par.Pass)>0{
		pass = Par.Pass
	} else {
		pass = Par.Brute.Mssql.Info.Pass
	}
	CreateInfo(&input,hosts,user,pass)
	thread := 10
	if len(hosts) < Par.Thread {
		thread = len(hosts)
	} else {
		thread = Par.Thread
	}
	/* 开启扫描 */
	for i:=0;i<thread;i++{
		go mysqlWork(input,result)
	}
	/* 输出 */
	public.Out(result,Par.Ofile)
}

func mysqlWork(input chan BruteInfo,result chan string){
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
			go mysqlConn(ctx, brute_pool, brute_result)
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

func mysqlConn(ctx context.Context,Pool chan BrutePool,scan chan string) {
	for{
		Pool, ok := <-Pool
		if !ok{
			return
		}
		Host, User, Pass := Pool.Host, Pool.User, Pool.Pass
		Pass = strings.Replace(Pass, "{{user}}", User, -1)
		dataSourceName := fmt.Sprintf("%v:%v@tcp(%v)/mysql?charset=utf8&timeout=%v", User, Pass, Host, time.Duration(Par.Timeout)*time.Second)
		db, err := sql.Open("mysql", dataSourceName)
		if err == nil {
			/* SetConnMaxLifetime 连接池里面的连接最大存活时长，默认值为0,表示不限制。 */
			db.SetConnMaxLifetime(5 * time.Second)
			/* SetConnMaxIdleTime 连接池里面的连接最大空闲时长 */
			db.SetConnMaxIdleTime(5 * time.Second)
			/* SetMaxOpenConns 设置与数据库的最大打开连接数，服务器cpu核心数 * 2 + 服务器有效磁盘数 */
			db.SetMaxOpenConns(runtime.NumCPU()*2)
			/* SetMaxIdleConns 设置空闲连接池中的最大连接数,小于MaxOpenConns */
			db.SetMaxIdleConns(runtime.NumCPU())
			defer db.Close()
			err  = db.Ping()
			if err == nil {
				result := fmt.Sprintf("[Mysql] %-40v %v:%-10v", Host, User, Pass)
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