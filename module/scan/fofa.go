package scan

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goon3/public"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

/* 用户转换检测用户信息成功json */
type checkUser struct {
	Email           string      `json:"email"`
	Username        string      `json:"username"`
	Fcoin           int         `json:"fcoin"`
	Isvip           bool        `json:"isvip"`
	Vip_level       int32       `json:"vip_level"`
	Is_verified     bool        `json:"is_verified"`
	Avatar          string      `json:"avatar"`
	Message         string       `json:"message"`
	Fofacli_ver     string      `json:"fofacli_ver"`
	Fofa_server     bool        `json:"fofa_server"`
}

/* 用户转换单个field返回的数据json */
type resultJson1 struct {
	Error       bool        `json:"error"`
	Mode        string      `json:"mode"`
	Page        int         `json:"page"`
	Query       string      `json:"query"`
	Results  [] string      `json:"results"`
	Size        int64       `json:"size"`
}

/* 用户转换多个field返回的数据json */
type resultJson2 struct {
	Error       bool        `json:"error"`
	Mode        string      `json:"mode"`
	Page        int         `json:"page"`
	Query       string      `json:"query"`
	Results  [][] string    `json:"results"`
	Size        int64       `json:"size"`
}

/* 返回的资产 */
type resultData struct {
	Results		[]string
}

/* 返回错误 */
type resultErr struct {
	Errmsg      string   `json:"errmsg"`
	Error       bool     `json:"error"`
}


var (
	Max = 0							// 当前用户最大获取资产数
	Level = ""                      // 当前用户会员等级
	GetNum = 0 					 	// 成功获取资产次数
	returnSize = 0  				// 本次请求资产的总量
	resulterr = &resultErr{}		// 解析错误信息
	resultjson1 = &resultJson1{}	// fields为单个时
	resultjson2 = &resultJson2{}	// fields为多个时
	resultdata = &resultData{}		// 解析fofa返回数据
	checkuser = &checkUser{}		// 解析用户会员信息
)


func FofaScan(words []string) {

	/* 获取用户信息 */
	if user := getUserInfo();!user{
		fmt.Println("")
		public.Error.Println("vip info is err! please check your email and key!")
		os.Exit(1)
	}

	if Par.FofaInfo.Num < 1{
		/* 公开版 */
		public.Warning.Println("num不能小于0！")
		os.Exit(1)
	} else {
		/* 普通模式 */
		for _,key := range(words){
			getFofa(key)
		}
	}

	/* 对结果去重 */
	result := []string{}
	m := make(map[string]bool) //map的值不重要
	for _, v := range resultdata.Results {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	/* 输出获取到的资产 */
	for _,data := range(result){
		public.Success.Println(data)
		public.FileWrite(Par.Ofile,(string(data)+"\n"))
	}
	fmt.Println("")
}



/*
获取fofa资产，正常模式
@key	keywords
*/
func getFofa(key string){
	asserts := getAssets(key)
	/* 判断fields是否为多数 */
	if find := strings.Contains(Par.FofaInfo.Fields, ","); find {
		str2json2(asserts,resultjson2)
	} else {
		str2json1(asserts,resultjson1)
	}
}

/*
获取fofa资产
@key	keyword
return  请求到的body
*/
func getAssets(key string) string{
	/* 对请求语法进行base64编码 */
	kbase64 := base64.StdEncoding.EncodeToString([]byte(key))
	time.Sleep((1/2)*time.Second)
	/* 生成url */
	Url := fmt.Sprintf("https://fofa.info/api/v1/search/all?email=%s&key=%s&qbase64=%s&size=%d&fields=%s",
		Par.FofaInfo.Email, Par.FofaInfo.Key,kbase64, Par.FofaInfo.Num, Par.FofaInfo.Fields)
	/* Get请求 */
	rsp,err:= http.Get(Url)
	if err != nil{
		public.Error.Println(err)
		return ""
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil{
		public.Error.Println(err)
		return ""
	}
	fmt.Println("数据请求成功，正在尝试请求下一条……")
	if strings.Contains(string(bytes), "errmsg"){
		str2json3(string(bytes),resulterr)
	} else {
		return string(bytes)
	}
	return ""
}

/*
格式化用户信息
@msg	fofa返回的用户数据
@User	用户信息结构体
*/
func str2json(msg string,User checkUser) bool{
	if err := json.Unmarshal([]byte(msg), &User); err == nil {
		switch User.Vip_level {
		case 1:
			Max = 100
			Level = "普通会员"
			break
		case 2:
			Max = 10000
			Level = "高级会员"
			break
		case 3:
			Max = 10000
			Level = "企业会员"
			break
		}
	} else{
		public.Error.Println(err)
		return false
	}
	return true
}

/*
格式化fofa返回的资产数据
@msg				fofa返回的资产数据
@resultJson1		当field为单个时的结构体
*/
func str2json1(msg string,resultJson *resultJson1){
	if err := json.Unmarshal([]byte(msg), &resultJson); err == nil {
		if resultJson.Error{
			public.Error.Println(resulterr.Errmsg)
		} else {
			for _,value := range(resultJson.Results){
				if Par.FofaInfo.Fields=="host" && Par.Web == true && !strings.Contains(value, "https"){
					value = "http://"+value
				}
				resultdata.Results = append(resultdata.Results,value)
			}
		}
	} else {
		public.Error.Println(err)
	}
}

/*
格式化fofa返回的资产数据
@msg				fofa返回的资产数据
@resultJson2		当field为多个时的结构体
*/
func str2json2(msg string,resultJson *resultJson2){
	if err := json.Unmarshal([]byte(msg), &resultJson); err == nil {
		if resultJson.Error{
			public.Error.Println(resulterr.Errmsg)
		} else {
			/* 获取host位置 */
			str := "ip,host,title"
			index := 0
			str2 := strings.SplitN(str,",",-1)
			for i,v :=range(str2){
				if v == "host"{
					index = i
				}
			}
			for _,value := range(resultJson.Results){
				r := ""
				for j,v := range(value){
					//r = r+v+"		"
					if Par.Web == true && j==index && !strings.Contains(v, "https"){
						v = "http://"+v
					}
					r = fmt.Sprintf("%s%-30s",r,v)
				}
				resultdata.Results = append(resultdata.Results,r)
			}
		}
	} else {
		public.Error.Println(err)
	}
}

/*
格式化fofa返回的资产数据
@msg			fofa返回的资产数据
@resultErr		错误数据结构体
*/
func str2json3(msg string,resultErr *resultErr){
	if err := json.Unmarshal([]byte(msg), &resultErr); err == nil {
		if resultErr.Error{
			public.Error.Println(resultErr.Errmsg)
		}
	} else{
		public.Error.Println(err)
	}
}

/* 获取用户信息 */
func getUserInfo() bool{
	Url := fmt.Sprintf("https://fofa.info/api/v1/info/my?email=%s&key=%s", Par.FofaInfo.Email, Par.FofaInfo.Key)
	rsp,err:= http.Get(Url)
	if err != nil{
		public.Error.Println(err)
		return false
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		public.Error.Println(err)
		return false
	}
	if strings.Contains(string(bytes), "errmsg"){
		public.Error.Println(err)
		return false
	} else{
		if !str2json(string(bytes), *checkuser){
			return false
		}
	}
	return true
}
