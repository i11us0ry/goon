package scan

import (
	"fmt"
	"goon3/public"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"net/http"
	"io/ioutil"
	"crypto/tls"
	"time"
)

func TitleScan(urls []string){

	runtime.GOMAXPROCS(runtime.NumCPU())
	input := make(chan string, len(urls))
	result := make(chan string, len(urls))

	defer close(input)

	/* 将要扫描的host放到甬道中 */
	go func() {
		for _, host := range(urls){
			input <- host
		}
	}()
	thread := 10
	if len(urls) < Par.Thread {
		thread = len(urls)
	} else {
		thread = Par.Thread
	}
	/* 开启扫描  */
	for i := 0; i< thread; i++{
		go scanTitle(input,result)
	}
	public.Out(result,Par.Ofile)
}

func scanTitle(input chan string,result chan string){
	for {
		task,ok := <-input
		if !ok{
			return
		}
		if find := strings.Contains(task, "http"); !find {
			task = "http://"+task
		}
		title := getTitle(task)
		if title!=""{
			title = strings.TrimSpace(title)
			title = strings.Replace(title,"\\n","",-1)
			result<-fmt.Sprintf("[title] %-30v %v",task,title)
		} else {
			result<-("")
		}
	}
}

func getTitle(url string) string{
	rsp := getHttpRequest(url)
	title := ""
	if rsp != nil{
		c, b := rsp.StatusCode, "0"
		if c == 401{
			return fmt.Sprintf("code:%-5vlen:%-10v401 unauthorized",c,b)
		}
		/* 获取body */
		body, err:= ioutil.ReadAll(rsp.Body)
		b = strconv.Itoa(len(body))
		if err != nil{
			return ""
		} else if  string(body) == "" {
			return fmt.Sprintf("code:%-5vlen:%-10vNone",c,b)
		}
		/* 正常请求 */
		title = checkTitle(string(body))
		if title!=""{
			/* 转码 */
			return fmt.Sprintf("code:%-5vlen:%-10v%v",c,b,public.TitletoUtf8(string(body),title))
		} else {
			/* 检查首页是否跳转 */
			patterns := [8]string{
				`(?is)<meta[\s]*http-equiv[\s]*=[\s]*["|']refresh["|'][\s]*content[\s]*=[\s]*["|'].*?;[\s]*url[\s]*=[\s]*(.*?)[\s]*["|']`,
				`(?is)window.location[\s]*=[\s]*["|'](.*?)["|'][\s]*[;|"]`,
				`(?is)window.location.href[\s]*=[\s]*["|'](.*?)["|'][\s]*[;|"]`,
				`(?is)window.location.replace[\s]*\(["|'](.*?)["|']\)[\s]*[;|"]`,
				`(?is)window.navigate[\s]*\(["|'](.*?)["|']\)`,
				`(?is)location.href[\s]*=[\s]*["|'](.*?)["|']`,
				`(?is)parent.location[\s]*=[\s]*["|'](.*?)["|']`,
				`(?is)window.open\(["|'](.*?)["|'][,|\)]`,
			}
			for _,pattern := range(patterns){
				reg := regexp.MustCompile(pattern)
				jumps := reg.FindAllStringSubmatch(string(body), 1)
				if len(jumps)>0{
					jump := jumps[0][1]
					/* 处理跳转的url */
					if find:=strings.Contains(jump,"://");!find{
						if index:=strings.Index(jump,"/");index!=0{
							jump = "/"+jump
						}
						url = url + jump
					} else {
						url = jump
					}
					/* 获取第二次rsp */
					rsp = getHttpRequest(url)
					if rsp != nil{
						defer rsp.Body.Close()
						body, err:= ioutil.ReadAll(rsp.Body)
						if err != nil{
							return ""
						} else if  string(body) == "" {
							return fmt.Sprintf("code:%-5vlen:%-10vNone",rsp.StatusCode,strconv.Itoa(len(body)))
						}
						title = checkTitle(string(body))
						break
					}
				}
			}
			if title!=""{
				return fmt.Sprintf("code:%-5vlen:%-10v%v",c,b,title)
			} else {
				return fmt.Sprintf("code:%-5vlen:%-10vNone",c,b)
			}
		}
	} else {
		/* 第一次rsp为nil */
		return ""
	}
}

func getHttpRequest(url string) *http.Response{
	/* 跳过https验证 */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	/* 超时请求 */
	c := &http.Client{
		Transport: tr,
		Timeout: time.Duration(Par.Timeout) * time.Second,
	}
	/* 是否重定向 */
	if !Par.Follow{
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return fmt.Errorf("first response")
		}
	}
	rsp, err := c.Get(url)
	if err != nil{
		return nil
	}
	return rsp
}

/*
获取网站标题
@content	网站内容
return		title
*/
func checkTitle(content string) string{
	//var title string
	/* 直接匹配<title>标题</title> */
	reg := regexp.MustCompile(`(?is)<title[^>]*>(.*?)<\/title>`)
	titles := reg.FindAllStringSubmatch(content, -1)
	if len(titles)>0{
		/* 存在多个title，有的title为空 */
		for i:=0;i<len(titles);i++{
			if titles[i][1]!="" && titles[i][1]!=" "{
				return titles[i][1]
			}
		}
	}
	/* 匹配静态js */
	reg = regexp.MustCompile(`(?is)document.title[\s]*=[\s]*['|"](.*?)['|"];`)
	titles = reg.FindAllStringSubmatch(content, 1)
	if len(titles)>0{
		return titles[0][1]
	}
	return ""
}
