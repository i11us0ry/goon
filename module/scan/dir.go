package scan

import (
	"fmt"
	"goon3/public"
	"net/http"
	"regexp"
	"runtime"
	"strings"
)

func DirScan(urls []string){
	runtime.GOMAXPROCS(runtime.NumCPU())
	input := make(chan string, len(urls))
	result := make(chan string, len(urls))
	defer close(input)

	/* 将要扫描的host放到甬道中 */
	go func(){
		for _, url := range(urls){
			input <- url+ Par.DirInfo.Dir
		}
	}()
	thread := 10
	if len(urls) < Par.Thread {
		thread = len(urls)
	} else {
		thread = Par.Thread
	}
	/* 请求方式 */
	for i := 0; i< thread; i++{
		go scanDir(input,result)
	}

	public.Out(result,Par.Ofile)
}

func scanDir(input chan string,result chan string){
	for {
		task,ok := <-input
		if !ok{
			return
		}
		//fmt.Println(task)
		if find := strings.Contains(task, "http"); find {
			getDir(task,result)
		} else {
			result<-""
		}
	}
}

func getDir(url string,result chan string){
	/* 跳过https验证 */
	var httpPar = public.NewHttpPar()
	httpPar.Url = url
	httpPar.Timeout = Par.Timeout
	httpPar.Follow = Par.Follow
	httpPar.Body = Par.DirInfo.Body

	for _, h:= range(Par.DirInfo.Header) {
		kv := strings.Split(h,":")
		httpPar.Header = append(httpPar.Header,[2]string{kv[0],kv[1]})
	}

	var header http.Header
	var body []byte
	var code int

	if Par.DirInfo.Mode == "get" {
		header, body, code = public.HttpDoGet2Body(httpPar)
	} else {
		header, body, code = public.HttpDoPost2Body(httpPar)
	}

	//fmt.Println(httpPar)
	//fmt.Println(Par.DirInfo.Code,Par.DirInfo.RHeader,Par.DirInfo.RBody)
	//matched, _ = regexp.MatchString(rule.Rule, Headers)
	if code == 0 {
		result<-""
	} else if code == Par.DirInfo.Code{
		/* 只看code */
		if Par.DirInfo.RHeader=="" && Par.DirInfo.RBody==""{
			result<-url
			/* 同时看code,body,header */
		} else if Par.DirInfo.RHeader!="" && Par.DirInfo.RBody!=""{
			if findbody,_ := regexp.MatchString(Par.DirInfo.RBody, string(body));findbody{
				if find, _ := regexp.MatchString(Par.DirInfo.RHeader, fmt.Sprintf("%s",header)); find{
					result<-url
				} else {
					result<-""
				}
			}
		} else {
			/* 判断code和header */
			if Par.DirInfo.RHeader!=""{
				if find, _ := regexp.MatchString(Par.DirInfo.RHeader, fmt.Sprintf("%s",header)); find{
					result<-url
				} else {
					result<-""
				}
				/* 判断code和body */
			} else {
				if find, _ := regexp.MatchString(Par.DirInfo.RBody, string(body)); find{
					result<-url
				} else {
					result<-""
				}
			}
		}
	} else {
		/* code不等 */
		result<-""
	}
}