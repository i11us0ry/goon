package public

import (
	"fmt"
	"goon3/lib/gonmap"
	"net"
	"strings"
	"time"
)

/*
获取services
@host		ip:port
*/
func host2Service(host string) string{
	// 获取tcp指纹
	service := gonmap.GetTcpBanner(host, gonmap.New(), 5 * time.Second)
	//fmt.Println(service,host)
	if service!=nil{
		return service.TcpFinger.Service
	}
	return ""
}

var bruteHost = &BruteHostType{}

func checkPortOpen(host string) bool{
	conn,err := net.DialTimeout("tcp",host,10 * time.Second)
	if err != nil{
		return false
	}  else {
		defer conn.Close()
		return true
	}
}

func GetPortInfo(hosts []string,thread int) *BruteHostType {
	gonmap.Init(10,10 * time.Second)
	input := make(chan string,  len(hosts))
	result := make(chan string, len(hosts))
	defer close(input)

	for _,host := range(hosts){
		input<-host
	}

	for i:=0;i<thread;i++{
		go work(input,result)
	}

	for i:=0;i<cap(result);i++{
		select {
		case host,ok:= <-result:
			if !ok {
				close(result)
				break;
			} else {
				if host!=""{
					Success.Println("[PortInfo] " + host)
				}
			}
		case <- time.After(time.Duration(5) * time.Second):
			if i==cap(result){
				close(result)
			} else {
				continue
			}
		}
	}
	fmt.Println()
	return bruteHost
}

func work(input,result chan string){
	for {
		task,ok := <-input
		if !ok{
			return
		}
		s := host2Service(task)
		if s == ""{
			result <- ""
		} else {
			if strings.Contains(s, "ftp"){
				bruteHost.Ftp.Host = append(bruteHost.Ftp.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s, "ssh"){
				bruteHost.Ssh.Host = append(bruteHost.Ssh.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s, "redis") {
				bruteHost.Redis.Host = append(bruteHost.Redis.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s, "mysql") {
				bruteHost.Mysql.Host = append(bruteHost.Mysql.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s, "ms-sql-") {
				bruteHost.Mssql.Host = append(bruteHost.Mssql.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s, "postgresql") {
				bruteHost.Postgres.Host = append(bruteHost.Postgres.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s,"netbios-") {
				bruteHost.NetBios.Host = append(bruteHost.NetBios.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s,"microsoft-") {
				bruteHost.Smb.Host = append(bruteHost.Smb.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s,"ms-wbt-server"){
				bruteHost.Rdp.Host = append(bruteHost.Rdp.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s,"rdp"){
				bruteHost.Rdp.Host = append(bruteHost.Rdp.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			} else if strings.Contains(s,"telnet"){
				bruteHost.Telnet.Host = append(bruteHost.Telnet.Host,task)
				result <- fmt.Sprintf("Host:%-30v Service:%v",task,s)
			}
		}
	}
}