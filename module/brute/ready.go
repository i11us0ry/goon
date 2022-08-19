package brute

import (
	"goon3/public"
)

var Par = &public.AutoParType{}

type Dic struct {
	User		[]string
	Pass    	[]string
}

var dic = &Dic{}

/* 待扫描资产详情 */
type BruteInfo struct {
	Host 		string
	Port 		string
	Users		[]string
	Passs    	[]string
	User		string
	Pass		string
	Timeout 	int
}
// 资源池
type BrutePool struct {
	Host 		string
	Port 		string
	User		string
	Pass		string
	Timeout 	int
}

/* 生成info */
func CreateInfo(input *chan BruteInfo,Host,User,Pass []string){
	var info BruteInfo
	for _, host := range(Host){
		info.Host = host
		info.Users = User
		info.Passs = Pass
		*input <- info
	}
}

func Init(par **public.AutoParType) {
	Par = *par
}
