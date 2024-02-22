package check

import (
	"fmt"
	"goon3/public"
	"os"
	"strings"
)

func CheckInput() {
	fmt.Println()
	public.Info.Println("checking input……")
	fmt.Println()
	CheckPar()
}

// 一些公共的参数par
func CheckPar() {
	/*
		端口扫描和fofa导出host是否导出为web格式
		- 如:baidu.com --> http://baidu.com
	*/

	if public.InputValue.OfilePtr == "" {
		public.InputValue.OfilePtr = public.GetCurrentDir() + "/result.txt"
	}

	// thread和timeout
	if public.InputValue.TimePtr != 0 {
		public.ConfValue.Timeout = public.InputValue.TimePtr
	}

	/*
		如果输入ufile，则读取ufile的user到public.ConfValue.User
		如果输入user，则将user添加到public.ConfValue.User
	*/
	if public.InputValue.UserFilePtr != "" {
		users := public.FileReadByline(public.InputValue.UserFilePtr)
		if len(users) > 0 {
			public.ConfValue.User = users
		}
	} else if public.InputValue.UserPtr != "" {
		public.ConfValue.User = []string{}
		public.ConfValue.User = append(public.ConfValue.User, public.InputValue.UserPtr)
	}
	/*
		如果输入pfile，则读取pfile的password到public.ConfValue.Pass
		如果输入pass，则将pass添加到public.ConfValue.Pass
	*/
	if public.InputValue.PassFilePtr != "" {
		passs := public.FileReadByline(public.InputValue.PassFilePtr)
		if len(passs) > 0 {
			public.ConfValue.Pass = passs
		}
	} else if public.InputValue.PassPtr != "" {
		public.ConfValue.Pass = []string{}
		public.ConfValue.Pass = append(public.ConfValue.User, public.InputValue.PassPtr)
	}
	/*
		如果输入port指定扫描端口，则将端口进行整理后放到public.ConfValue.Port
		反则将配置文件中的端口整理后放到public.ConfValue.Port
		整理和输入规则:80,81-85 --> [80,81,82,83,84,85]
	*/
	if public.InputValue.PortPtr != "" {
		ports := GetPort(public.InputValue.PortPtr)
		if len(ports) > 0 {
			public.ConfValue.Port = ports
		}
	} else {
		tempPort := public.ConfValue.Port
		public.ConfValue.Port = []string{}
		for _, port := range tempPort {
			ports := GetPort(port)
			for _, port1 := range ports {
				public.ConfValue.Port = append(public.ConfValue.Port, port1)
			}
		}
	}
	/*
		处理ip、url和ifile
		如果输入ip指定扫描ip，则将ip进行整理后放到public.ConfValue.ip
		如果输入url，判断url格式正确后放到public.ConfValue.Url
		如果输入ifile，则要分类讨论：
			- 如果mode是web方面的，则ifile内容应该为url
			- 如果mode是爆破方面的，则ifile内容应该是ip
			- 如果mode是all，则对ifile内容进行分类，ip放到public.ConfValue.ip，url放到public.ConfValue.Url
			- 如果mode是fofa，则ifile内容放到
	*/
	if public.InputValue.IpsPtr != "" {
		// 如果需要探活，则判断是否是ip段，并检查ip段的第一个和最后一个地址
		if !public.InputValue.NoPingPtr {
			if find := strings.Count(public.InputValue.IpsPtr, "/"); find == 1 {
				if !CheckSub(public.InputValue.IpsPtr) {
					// fmt.Println(public.InputValue.IpsPtr)
					public.Warning.Println("CheckSub cant find any ip 1!")
					os.Exit(0)
				}
			} else {
				ips := GetIp(public.InputValue.IpsPtr)
				if len(ips) > 0 {
					public.ConfValue.Ip = ips
				}
			}
		} else {
			ips := GetIp(public.InputValue.IpsPtr)
			if len(ips) > 0 {
				public.ConfValue.Ip = ips
			}
		}
	} else if public.InputValue.UrlPtr != "" {
		if strings.Contains(public.InputValue.UrlPtr, "http") {
			public.ConfValue.Url = append(public.ConfValue.Url, public.InputValue.UrlPtr)
		} else {
			public.Error.Println("url is err!")
			os.Exit(1)
		}
	} else if public.InputValue.IfilePtr != "" {
		lines := public.FileReadByline(public.InputValue.IfilePtr)
		if len(lines) == 0 {
			public.Error.Println("ifile is err!")
			os.Exit(1)
		} else {
			/*
				如果mode是web方面的，则ifile内容应该为url、或domain
				如果mode是爆破方面的，则ifile内容应该是ip
			*/
			if public.InputValue.ModePtr == "title" || public.InputValue.ModePtr == "dfuzz" || public.InputValue.ModePtr == "finger" || public.InputValue.ModePtr == "tomcat" {
				for _, line := range lines {
					if strings.Contains(line, "http") {
						public.ConfValue.Url = append(public.ConfValue.Url, line)
					} else {
						// 如果目标不带http则加上http
						line = "http://" + line
						public.ConfValue.Url = append(public.ConfValue.Url, line)
					}
				}
			} else if public.InputValue.ModePtr == "brute" || public.InputValue.ModePtr == "webscan" || public.InputValue.ModePtr == "ftp" || public.InputValue.ModePtr == "mysql" || public.InputValue.ModePtr == "mssql" || public.InputValue.ModePtr == "redis" || public.InputValue.ModePtr == "ssh" || public.InputValue.ModePtr == "ms17010" || public.InputValue.ModePtr == "smb" || public.InputValue.ModePtr == "postgres" || public.InputValue.ModePtr == "ip" || public.InputValue.ModePtr == "port" || public.InputValue.ModePtr == "netbios" || public.InputValue.ModePtr == "rdp" || public.InputValue.ModePtr == "telnet" {
				for _, line := range lines {
					// 这里做一步判断，如果要爆破的资产文件中已经指定了端口，如文件中为127.0.0.1:21，则不需要通过GetIp验证ip是否合法
					if strings.Contains(line, ":") {
						public.ConfValue.Ip = append(public.ConfValue.Ip, line)
					} else {
						if !public.InputValue.NoPingPtr {
							if find := strings.Count(line, "/"); find == 1 {
								if !CheckSub(line) {
									continue
								}
							} else {
								public.ConfValue.Ip = append(public.ConfValue.Ip, line)
							}
						} else {
							ips := GetIp(line)
							if len(ips) > 0 {
								for _, ip := range ips {
									public.ConfValue.Ip = append(public.ConfValue.Ip, ip)
								}
							}
						}
					}
				}
			} else if public.InputValue.ModePtr == "all" {
				for _, line := range lines {
					if strings.Contains(line, "http") {
						public.ConfValue.Url = append(public.ConfValue.Url, line)
					} else if strings.Contains(line, ":") {
						line = "http://" + line
						public.ConfValue.Url = append(public.ConfValue.Url, line)
					} else {
						if !public.InputValue.NoPingPtr {
							if find := strings.Count(line, "/"); find == 1 {
								if !CheckSub(line) {
									continue
								}
							} else {
								public.ConfValue.Ip = append(public.ConfValue.Ip, line)
							}
						} else {
							ips := GetIp(line)
							if len(ips) > 0 {
								for _, ip := range ips {
									public.ConfValue.Ip = append(public.ConfValue.Ip, ip)
								}
							}
						}
					}
				}
			} else if public.InputValue.ModePtr == "fofa" {
				public.ConfValue.FofaWord = lines
			} else {
				fmt.Println(public.Banner)
				public.OutMode()
				os.Exit(1)
			}
		}
	} else {
		if public.InputValue.ModePtr != "fofa" {
			public.Error.Println("you should input -ip or -url or -ifile or -mode!")
			public.OutMode()
			os.Exit(1)
		}
	}
	/*
		dir的一些参数
		- dir：要请求的dir路径
		- body：response body包含内容
		- header：response header包含内容
		- code：response code
	*/
	if public.InputValue.DirPtr != "" {
		public.ConfValue.DirInfo.Dir = public.InputValue.DirPtr
	}
	if public.InputValue.BodyPtr != "" {
		public.ConfValue.DirInfo.Body = public.InputValue.BodyPtr
	}
	if public.InputValue.HeaderPtr != "" {
		public.ConfValue.DirInfo.Header = append(public.ConfValue.DirInfo.Header, public.InputValue.HeaderPtr)
	}
	if public.InputValue.RCodePtr != 200 {
		public.ConfValue.DirInfo.Code = public.InputValue.RCodePtr
	}
	if public.InputValue.RBodyPtr != "" {
		public.ConfValue.DirInfo.RBody = public.InputValue.RBodyPtr
	}
	if public.InputValue.RHeaderPtr != "" {
		public.ConfValue.DirInfo.RHeader = public.InputValue.RHeaderPtr
	}
	if public.InputValue.DModePtr != "" {
		public.ConfValue.DirInfo.Mode = public.InputValue.DModePtr
	}
	/*
		fofa的一些参数
		- num: api返回数据量
		- key：搜索语法
		- fields：fields
	*/
	if public.InputValue.NumPtr != 0 {
		public.ConfValue.FofaInfo.Num = public.InputValue.NumPtr
	}
	if public.InputValue.KeyPtr != "" {
		public.ConfValue.FofaWord = append(public.ConfValue.FofaWord, public.InputValue.KeyPtr)
	}
	if public.InputValue.FieldsPtr != "" {
		public.ConfValue.FofaInfo.Fields = public.InputValue.FieldsPtr
	}
}

func CheckSub(ips string) bool {
	ipsub := getSubNet(ips)
	if ipsub != nil && len(ipsub) != 0 {
		for _, v := range ipsub {
			//public.ConfValue.IpAlive = append(public.ConfValue.IpAlive,v)
			ips := GetIp(v)
			if len(ips) > 0 {
				for _, ip := range ips {
					public.ConfValue.Ip = append(public.ConfValue.Ip, ip)
				}
			}
		}
		return true
	}
	return false
}
