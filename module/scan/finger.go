package scan

import (
	"fmt"
	"goon3/public"
	"regexp"
	"runtime"
)

type CheckDatas struct {
	Body    []byte
	Headers string
}

func FingerScan(urls []string){
	runtime.GOMAXPROCS(runtime.NumCPU())

	input := make(chan string, len(urls))
	result := make(chan string, len(urls))
	defer close(input)
	/* 将要扫描的host放到甬道中 */
	go func(){
		for _, url := range(urls){
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
		go fingerWork(input,result)
	}
	public.Out(result)
}

/* 从channel获取host和plugin，然后读取plugin内容 */
func fingerWork(input chan string,result chan string){
	for{
		task,ok := <-input
		find := false
		if !ok {
			return
		}
		/*
		进行 http get 请求，将response的header和body和rule匹配
		如果rule中url不为空，则重新请求添加url路由后再次请求
		*/
		var httpPar = public.NewHttpPar()
		httpPar.Url = task
		httpPar.Timeout = Par.Timeout
		httpPar.Follow = Par.Follow
		httpPar.Header = [][2]string{{"cookie","rememberMe=1;"}}
		header, body,_ := public.HttpDoGet2Body(httpPar)
		//fmt.Println(string(body))
		for _,rule := range(public.RuleDatas){
			if rule.Url!="" {
				if string(rule.Url[0])!="/"{
					url := task+"/"+rule.Url
					httpPar.Url = url
					header2, body2,_ := public.HttpDoGet2Body(httpPar)
					str := InfoCheck(url,fmt.Sprintf("%s",header2), string(body2), rule)
					if str != ""{
						result <- "[finger] " + str
						find = true
						break
					}
				} else {
					url := task+rule.Url
					httpPar.Url = url
					header2, body2,_ := public.HttpDoGet2Body(httpPar)
					str := InfoCheck(url,fmt.Sprintf("%s",header2), string(body2), rule)
					if str != ""{
						result <- "[finger] " + str
						find = true
						break
					}
				}
			} else {
				str := InfoCheck(task,fmt.Sprintf("%s",header), string(body), rule)
				if str != ""{
					result <- "[finger] " + str
					find = true
					break
				}
			}
		}
		if find == false{
			result <- ""
		}
	}
}

func InfoCheck(Url,Headers,Body string,rule public.RuleDataType) string {
	var matched bool
	var infoname []string
	//fmt.Println("============",rule,Body)
	if rule.Type == "code" {
		matched, _ = regexp.MatchString(rule.Rule, Body)
	} else {
		matched, _ = regexp.MatchString(rule.Rule, Headers)
	}
	if matched == true {
		infoname = append(infoname, rule.Name)
		return fmt.Sprintf("%-30v %v",Url,rule.Name)
	}
	return ""
}