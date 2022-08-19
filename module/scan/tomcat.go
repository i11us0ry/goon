package scan

import (
	"fmt"
	"goon3/public"
	"runtime"
	"strings"
	"encoding/base64"
)

func TomcatScan(urls []string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(Par.User)>0{
		dic.User = Par.User
	} else {
		dic.User = Par.Brute.Tomcat.Info.User
	}
	if len(Par.Pass)>0{
		dic.Pass = Par.Pass
	} else {
		dic.Pass = Par.Brute.Tomcat.Info.Pass
	}

	input := make(chan string, len(urls))
	result := make(chan string, len(urls))
	defer close(input)

	/* 将要扫描的host放到甬道中 */
	go func(){
		for _, url := range(urls){
			if !strings.Contains(url,"/manager/html"){
				url = url + "/manager/html"
			}
			input <- url
		}
	}()
	thread := 10
	if len(urls) < Par.Thread {
		thread = len(urls)
	} else {
		thread = Par.Thread
	}
	/* 启动扫描 */
	for i:=0; i< thread; i++{
		go tomcatWork(input,result)
	}
	public.Out(result,Par.Ofile)
}

func tomcatWork(input chan string,result chan string){
	for {
		task,ok := <-input
		find := false
		if !ok{
			return
		}
		for _, user := range(dic.User){
			for _, pass := range(dic.Pass){
				if flag := startBrute(task, user, pass);flag!=""{
					result<-flag
					find = true
				}
			}
			if find{
				break
			}
		}
		if !find{
			result <- ""
		}
	}
}

func startBrute(url, user, pass string) string{

	pass = strings.Replace(pass, "{{user}}", user, -1)
	//fmt.Println(url,user,pass)
	var httpPar = public.NewHttpPar()
	httpPar.Url = url
	httpPar.Timeout = Par.Timeout
	httpPar.Follow = Par.Follow

	strbytes := []byte(fmt.Sprintf("%v:%v",user, pass))
	encoded := base64.StdEncoding.EncodeToString(strbytes)

	login := fmt.Sprintf("Basic %v", encoded)
	httpPar.Header = [][2]string{{"Authorization",login}}
	_,_,code := public.HttpDoGet2Body(httpPar)
	if code == 200 {
		return fmt.Sprintf("[tomcat] %-30v %v:%v",url, user, pass)
	}
	return ""
}