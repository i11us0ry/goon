package scan

import (
	"goon3/public"
)

var Par = &public.AutoParType{}

func Init(par **public.AutoParType) {
	Par = *par
}

type Info struct {
	Ip string
	Port string
}

type Dic struct {
	User		[]string
	Pass    	[]string
}

var dic = &Dic{}

/* 待扫描资产详情 */
type ScanInfo struct {
	Host 		string
	Users		[]string
	Passs    	[]string
	Timeout 	int
}