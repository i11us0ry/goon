package public

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type HttpPar struct {
	Url 	string
	Timeout int
	Follow  bool
	Body    string
	Header  [][2]string
}

func NewHttpPar() *HttpPar{
	return &HttpPar{
		"",
		10,
		true,
		"",
		[][2]string{},
	}
}

var USER_AGENT = []string{
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0;",
	"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)",
	"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
	"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SE 2.X MetaSr 1.0; SE 2.X MetaSr 1.0; .NET CLR 2.0.50727; SE 2.X MetaSr 1.0)",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
}

func HttpDoGet2Body(par *HttpPar) (http.Header, []byte, int){
	/* 跳过https验证 */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Duration(par.Timeout) * time.Second,
	}
	if !par.Follow{
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return fmt.Errorf("first response")
		}
	}
	req, err := http.NewRequest("GET", par.Url, nil)
	if err != nil {
		return nil,nil,0
	}
	req.Header.Set("User-agent", USER_AGENT[rand.Intn(7)])

	if len(par.Header)!=0{
		for _,h := range(par.Header){
			hk, hv := h[0], h[1]
			req.Header.Set(hk, hv)
		}
	}

	resp, err := client.Do(req)
	if err!=nil{
		return nil,nil,0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil,nil,0
	}
	return resp.Header,body,resp.StatusCode
}

// 返回http请求code
func HttpDoGet2Code(url string,timeout int) int{
	/* 跳过https验证 */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{
		Transport: tr,
		Timeout: time.Duration(timeout) * time.Second,
	}
	resp, err := c.Get(url)
	if err != nil {
		return 50000
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func HttpDoPost2Body(par *HttpPar) (http.Header, []byte, int){
	/* 跳过https验证 */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Duration(par.Timeout) * time.Second,
	}
	if !par.Follow{
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return fmt.Errorf("first response")
		}
	}
	req, err := http.NewRequest("POST", par.Url, strings.NewReader(par.Body))
	if err != nil {
		return nil,nil,0
	}
	req.Header.Set("User-agent", USER_AGENT[rand.Intn(7)])
	if len(par.Header)!=0{
		for _,h := range(par.Header){
			hk, hv := h[0], h[1]
			req.Header.Set(hk, hv)
		}
	}

	resp, err := client.Do(req)
	if err!=nil{
		return nil,nil,0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil,nil,0
	}
	return resp.Header,body,resp.StatusCode
}