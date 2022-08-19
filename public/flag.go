package public

import (
	"flag"
)

func Flag() {
	var svar string
	Banner()
	// 通用参数
	flag.StringVar(&InputValue.ModePtr,"mode", "all", "运行模式如:webscan、brute、title、fofa、mysql、mssql等，默认为all")
	flag.StringVar(&InputValue.IpsPtr,"ip","","如:127.0.0.1、127.0.0.1/24、127.0.0.1-255，IP段支持/8-/31任意CIDR")
	flag.StringVar(&InputValue.IfilePtr,"ifile", "", "输入文件，可以是url或ip或ip段，IP段支持/8-/31任意CIDR")
	flag.StringVar(&InputValue.OfilePtr,"ofile", "", "输出文件，默认为result.txt")
	flag.IntVar(&InputValue.ThreadPtr,"thread", 0, "thread,默认从配置文件读取")
	flag.IntVar(&InputValue.TimePtr,"time", 0, "timeout,默认从配置文件读取")
	flag.BoolVar(&InputValue.HelpPtr,"help", false, "mode详解")
	flag.StringVar(&InputValue.UrlPtr,"url", "", "url，必须有http")
	flag.BoolVar(&InputValue.NoPingPtr,"np",false,"np，不进行icmp探活")

	// 端口扫描参数
	flag.StringVar(&InputValue.PortPtr,"port","","扫描端口如:80,443-445,8000-9000")
	flag.BoolVar(&InputValue.WebPtr,"web",false,"port和fofa host输出格式如:http://127.0.0.1:80")

	// dir扫描参数
	flag.StringVar(&InputValue.DirPtr,"dir", "", "dir fuzz请求的路径如:/login.jsp，适用于对批量url进行单个dir探测，支持post发包，支持正则匹配，可探测简单poc")
	flag.StringVar(&InputValue.DModePtr,"dmode", "get", "dir fuzz请求方式:get或post")
	flag.StringVar(&InputValue.HeaderPtr,"header", "", "dir请求包的header如:Content-Type:Application/json")
	flag.StringVar(&InputValue.BodyPtr,"body", "", "dir请求包的body如:{\"name\":\"username\"}")
	flag.IntVar(&InputValue.RCodePtr,"code", 200, "dir返回包code如:200、302")
	flag.StringVar(&InputValue.RHeaderPtr,"rheader", "", "dir返回包header如:rememberMe，支持正则")
	flag.StringVar(&InputValue.RBodyPtr,"rbody", "", "dir返回包body如:root:x:0:0，支持正则")

	// fofa
	flag.StringVar(&InputValue.KeyPtr,"key", "", "fofa查询语句如:domain='fofa.so'")
	flag.IntVar(&InputValue.NumPtr,"num", 0, "fofa请求数量如:100、10000")
	flag.StringVar(&InputValue.FieldsPtr,"fields", "", "fofa返回类型如:ip,host")

	// brute
	flag.StringVar(&InputValue.UserFilePtr,"ufile", "", "用户字典，默认从配置文件读取")
	flag.StringVar(&InputValue.PassFilePtr,"pfile", "", "密码字典，默认从配置文件读取")
	flag.StringVar(&InputValue.UserPtr,"user", "", "用户，默认从配置文件读取")
	flag.StringVar(&InputValue.PassPtr,"pass", "", "密码，默认从配置文件读取")

	flag.StringVar(&svar, "svar", "bar", "a string var") // 对变量取址
	flag.Parse()
}