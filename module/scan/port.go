package scan

import (
	"goon3/public"
	"net"
	"os"
	"runtime"
	"time"
)

func PortScan(ips []string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ilen, plen := len(ips), len(Par.Port)
	if ilen == 0 || plen == 0 {
		public.Error.Println("portscan cant find ip or port!")
		os.Exit(1)
	}
	input := make(chan string, len(ips)*len(Par.Port))
	result := make(chan string, len(ips)*len(Par.Port))
	defer close(input)

	/* 将要扫描的IP和端口放到甬道中 */
	for _, ip := range ips {
		for _, port := range Par.Port {
			input <- ip + ":" + port
		}
	}
	thread := 10
	if len(ips)*len(Par.Port) < Par.Thread {
		thread = len(ips) * len(Par.Port)
	} else {
		thread = Par.Thread
	}
	/* 启动扫描 */
	if public.InputValue.WebPtr == true {
		for i := 0; i < thread; i++ {
			go scanWeb(input, result)
		}
	} else {
		for i := 0; i < thread; i++ {
			go scanPort(input, result)
		}
	}
	public.Out(result)
}

/*
扫描端口、简单的TCP连接
@intput	ip:port
@result	扫描结果
*/
func scanPort(input chan string, result chan string) {
	for {
		task, ok := <-input
		if !ok {
			return
		}
		_, err := net.DialTimeout("tcp", task, time.Duration(Par.Timeout)*time.Second)
		if err != nil {
			result <- ""
		} else {
			result <- "[port] " + task
		}
	}
}

/*
将要扫描的HOST分为http、和https
@intput	ip:port
@result	扫描结果
*/
func scanWeb(input chan string, result chan string) {
	for {
		task, ok := <-input
		if !ok {
			return
		}
		url := "http://" + task
		isHttp := getWeb(url)
		if isHttp {
			result <- "[url] " + url
			Par.Url = append(Par.Url, url)
		} else {
			urls := "https://" + task
			isHttps := getWeb(urls)
			if isHttps {
				result <- "[url] " + urls
				Par.Url = append(Par.Url, urls)
			} else {
				result <- ""
			}
		}
	}
}

/*
扫描web
@url		http://ip:port
@result		扫描结果
*/
func getWeb(url string) bool {
	respCode := public.HttpDoGet2Code(url, Par.Timeout)
	for _, code := range Par.WebCode {
		if code == respCode {
			return true
		}
	}
	return false
}
